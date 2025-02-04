// Copyright 2022-2025 The sacloud/iaas-service-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package setup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/accessor"
	"github.com/sacloud/iaas-api-go/types"
)

// MaxRetryCountExceededError リトライ最大数超過エラー
type MaxRetryCountExceededError error

// CreateFunc リソース作成関数
type CreateFunc func(ctx context.Context, zone string) (accessor.ID, error)

// ProvisionBeforeUpFunc リソース作成後、起動前のプロビジョニング関数
//
// リソース作成後に起動が行われないリソース(VPCルータなど)向け。
// 必要であればこの中でリソース起動処理を行う。
type ProvisionBeforeUpFunc func(ctx context.Context, zone string, id types.ID, target interface{}) error

// DeleteFunc リソース削除関数。
//
// リソース作成時のコピー待ちの間にリソースのAvailabilityがFailedになった場合に利用される。
type DeleteFunc func(ctx context.Context, zone string, id types.ID) error

// ReadFunc リソース起動待ちなどで利用するリソースのRead用Func
type ReadFunc func(ctx context.Context, zone string, id types.ID) (interface{}, error)

// RetryableSetup リソース作成時にコピー待ちや起動待ちが必要なリソースのビルダー。
//
// リソースのビルドの際、必要に応じてリトライ(リソースの削除&再作成)を行う。
type RetryableSetup struct {
	// Create リソース作成用関数
	Create CreateFunc
	// IsWaitForCopy コピー待ちを行うか
	IsWaitForCopy bool
	// IsWaitForUp 起動待ちを行うか
	IsWaitForUp bool
	// ProvisionBeforeUp リソース起動前のプロビジョニング関数
	ProvisionBeforeUp ProvisionBeforeUpFunc
	// Delete リソース削除用関数
	Delete DeleteFunc
	// Read リソース起動待ち関数
	Read ReadFunc

	// Options .
	Options *Options
}

// Setup リソースのビルドを行う。必要に応じてリトライ(リソースの削除&再作成)を行う。
func (r *RetryableSetup) Setup(ctx context.Context, zone string) (interface{}, error) {
	if (r.IsWaitForCopy || r.IsWaitForUp) && r.Read == nil {
		return nil, errors.New("failed: Read is required when IsWaitForCopy or IsWaitForUp is true")
	}

	r.init()

	var created interface{}
	for r.Options.RetryCount+1 > 0 {
		r.Options.RetryCount--

		// リソース作成
		target, err := r.createResource(ctx, zone)
		if err != nil {
			return nil, err
		}
		id := target.GetID()

		// コピー待ち
		if r.IsWaitForCopy {
			// コピー待ち、Failedになった場合はリソース削除
			state, err := r.waitForCopyWithCleanup(ctx, zone, id)
			if err != nil {
				return state, err
			}
			if state != nil {
				created = state
			}
		} else {
			created = target
		}

		// 起動前の設定など
		if err := r.provisionBeforeUp(ctx, zone, id, created); err != nil {
			return created, err
		}

		// 起動待ち
		if err := r.waitForUp(ctx, zone, id, created); err != nil {
			return created, err
		}

		if created != nil {
			break
		}
	}

	if created == nil {
		return nil, MaxRetryCountExceededError(fmt.Errorf("max retry count exceeded"))
	}
	return created, nil
}

func (r *RetryableSetup) init() {
	if r.Options == nil {
		r.Options = &Options{}
	}
	r.Options.Init()
}

func (r *RetryableSetup) createResource(ctx context.Context, zone string) (accessor.ID, error) {
	if r.Create == nil {
		return nil, fmt.Errorf("create func is required")
	}
	return r.Create(ctx, zone)
}

func (r *RetryableSetup) waitForCopyWithCleanup(ctx context.Context, zone string, id types.ID) (interface{}, error) {
	waiter := &iaas.StatePollingWaiter{
		ReadFunc: func() (interface{}, error) {
			return r.Read(ctx, zone, id)
		},
		TargetAvailability: []types.EAvailability{
			types.Availabilities.Available,
			types.Availabilities.Failed,
		},
		PendingAvailability: []types.EAvailability{
			types.Availabilities.Unknown,
			types.Availabilities.Migrating,
			types.Availabilities.Uploading,
			types.Availabilities.Transferring,
			types.Availabilities.Discontinued,
		},
		Interval: r.Options.PollingInterval,
	}

	// wait
	compChan, progressChan, errChan := waiter.WaitForStateAsync(ctx)
	var state interface{}
	var err error

loop:
	for {
		select {
		case v := <-compChan:
			state = v
			break loop
		case v := <-progressChan:
			state = v
		case e := <-errChan:
			err = e
			break loop
		}
	}

	if state != nil {
		// Availabilityを持ち、Failedになっていた場合はリソースを削除してリトライ
		if f, ok := state.(accessor.Availability); ok && f != nil {
			if f.GetAvailability().IsFailed() {
				// FailedになったばかりだとDelete APIが失敗する(コピー進行中など)場合があるため、
				// 任意の回数リトライ&待機を行う
				for i := 0; i < r.Options.DeleteRetryCount; i++ {
					time.Sleep(r.Options.DeleteRetryInterval)
					if err = r.Delete(ctx, zone, id); err == nil {
						break
					}
				}

				return nil, nil
			}
		}

		return state, nil
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *RetryableSetup) provisionBeforeUp(ctx context.Context, zone string, id types.ID, created interface{}) error {
	if r.ProvisionBeforeUp != nil && created != nil {
		var err error
		for i := 0; i < r.Options.ProvisioningRetryCount; i++ {
			if err = r.ProvisionBeforeUp(ctx, zone, id, created); err == nil {
				break
			}
			time.Sleep(r.Options.ProvisioningRetryInterval)
		}
		return err
	}
	return nil
}

func (r *RetryableSetup) waitForUp(ctx context.Context, zone string, id types.ID, created interface{}) error {
	if r.IsWaitForUp && created != nil {
		waiter := &iaas.StatePollingWaiter{
			ReadFunc: func() (interface{}, error) {
				return r.Read(ctx, zone, id)
			},
			TargetAvailability: []types.EAvailability{
				types.Availabilities.Available,
			},
			PendingAvailability: []types.EAvailability{
				types.Availabilities.Unknown,
				types.Availabilities.Migrating,
				types.Availabilities.Uploading,
				types.Availabilities.Transferring,
				types.Availabilities.Discontinued,
			},
			TargetInstanceStatus: []types.EServerInstanceStatus{
				types.ServerInstanceStatuses.Up,
			},
			PendingInstanceStatus: []types.EServerInstanceStatus{
				types.ServerInstanceStatuses.Unknown,
				types.ServerInstanceStatuses.Cleaning,
				types.ServerInstanceStatuses.Down,
			},
			Interval: r.Options.PollingInterval,
		}
		_, err := waiter.WaitForState(ctx)
		return err
	}
	return nil
}

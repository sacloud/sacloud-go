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

package builder

import (
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/types"
)

// Builder SIMのセットアップを行う
type Builder struct {
	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID
	ICCID       string
	PassCode    string

	Activate bool
	IMEI     string
	Carrier  []*iaas.SIMNetworkOperatorConfig

	Client *APIClient
}

// Validate 値の検証
func (b *Builder) Validate(ctx context.Context) error {
	if b.ICCID == "" {
		return fmt.Errorf("iccid is required")
	}
	if len(b.Carrier) == 0 {
		return fmt.Errorf("carrier is required")
	}
	return nil
}

// Build SIMの作成
func (b *Builder) Build(ctx context.Context) (*iaas.SIM, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	sim, err := b.Client.SIM.Create(ctx, &iaas.SIMCreateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
		ICCID:       b.ICCID,
		PassCode:    b.PassCode,
	})
	if err != nil {
		return nil, err
	}

	if err := b.Client.SIM.SetNetworkOperator(ctx, sim.ID, b.Carrier); err != nil {
		return sim, err
	}

	if b.Activate {
		if err := b.Client.SIM.Activate(ctx, sim.ID); err != nil {
			return sim, err
		}
	}

	if b.IMEI != "" {
		if err := b.Client.SIM.IMEILock(ctx, sim.ID, &iaas.SIMIMEILockRequest{IMEI: b.IMEI}); err != nil {
			return sim, err
		}
	}

	// reload
	refreshed, err := query.FindSIMByID(ctx, b.Client.SIM, sim.ID)
	if err != nil {
		return sim, err
	}
	return refreshed, nil
}

// Update SIMの更新
func (b *Builder) Update(ctx context.Context, id types.ID) (*iaas.SIM, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	sim, err := query.FindSIMByID(ctx, b.Client.SIM, id)
	if err != nil {
		return nil, err
	}

	_, err = b.Client.SIM.Update(ctx, id, &iaas.SIMUpdateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
	})
	if err != nil {
		return nil, err
	}

	if err := b.Client.SIM.SetNetworkOperator(ctx, sim.ID, b.Carrier); err != nil {
		return nil, err
	}

	if !b.Activate && sim.Info.Activated {
		if err := b.Client.SIM.Deactivate(ctx, sim.ID); err != nil {
			return nil, err
		}
	}
	if b.Activate && !sim.Info.Activated {
		if err := b.Client.SIM.Activate(ctx, sim.ID); err != nil {
			return nil, err
		}
	}

	// Unlock -> ロックされている、かつIMEIが空か前と変わっている場合
	if sim.Info.IMEILock && (b.IMEI == "" || b.IMEI != sim.Info.IMEI) {
		if err := b.Client.SIM.IMEIUnlock(ctx, sim.ID); err != nil {
			return nil, err
		}
	}
	// Lock -> IMEIが変わっている場合は上でアンロックされた状態になっているはず
	if b.IMEI != "" && b.IMEI != sim.Info.IMEI {
		if err := b.Client.SIM.IMEILock(ctx, sim.ID, &iaas.SIMIMEILockRequest{IMEI: b.IMEI}); err != nil {
			return nil, err
		}
	}

	// reload
	sim, err = query.FindSIMByID(ctx, b.Client.SIM, id)
	if err != nil {
		return nil, err
	}
	return sim, nil
}

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

package bridge

import (
	"context"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-service-go/serviceutil"
)

func (s *Service) Delete(req *DeleteRequest) error {
	return s.DeleteWithContext(context.Background(), req)
}

func (s *Service) DeleteWithContext(ctx context.Context, req *DeleteRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	if req.WaitForRelease {
		opt := query.CheckReferencedOption{
			Timeout: time.Duration(req.WaitForReleaseTimeout) * time.Second,
			Tick:    time.Duration(req.WaitForReleaseTick) * time.Second,
		}
		if err := cleanup.DeleteBridge(ctx, s.caller, req.Zone, req.Zones, req.ID, opt); err != nil {
			return serviceutil.HandleNotFoundError(err, !req.FailIfNotFound)
		}
	} else {
		client := iaas.NewBridgeOp(s.caller)
		if err := client.Delete(ctx, req.Zone, req.ID); err != nil {
			return serviceutil.HandleNotFoundError(err, !req.FailIfNotFound)
		}
	}
	return nil
}

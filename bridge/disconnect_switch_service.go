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
	"fmt"

	"github.com/sacloud/iaas-api-go"
)

func (s *Service) DisconnectSwitch(req *DisconnectSwitchRequest) error {
	return s.DisconnectSwitchWithContext(context.Background(), req)
}

func (s *Service) DisconnectSwitchWithContext(ctx context.Context, req *DisconnectSwitchRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return err
	}

	bridgeOp := iaas.NewBridgeOp(s.caller)
	switchOp := iaas.NewSwitchOp(s.caller)

	bridge, err := bridgeOp.Read(ctx, req.Zone, req.ID)
	if err != nil {
		return err
	}
	if bridge.SwitchInZone == nil {
		return fmt.Errorf("target bridge[%s] is not connected any switches", req.ID)
	}

	return switchOp.DisconnectFromBridge(ctx, req.Zone, bridge.SwitchInZone.ID)
}

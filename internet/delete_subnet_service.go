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

package internet

import (
	"context"

	"github.com/sacloud/iaas-api-go"
)

func (s *Service) DeleteSubnet(req *DeleteSubnetRequest) error {
	return s.DeleteSubnetWithContext(context.Background(), req)
}

func (s *Service) DeleteSubnetWithContext(ctx context.Context, req *DeleteSubnetRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	client := iaas.NewInternetOp(s.caller)
	return client.DeleteSubnet(ctx, req.Zone, req.ID, req.SubnetID)
}

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
	"fmt"

	"github.com/sacloud/iaas-api-go"
)

func (s *Service) ReadIPv6(req *ReadIPv6Request) (*iaas.IPv6Net, error) {
	return s.ReadIPv6WithContext(context.Background(), req)
}

func (s *Service) ReadIPv6WithContext(ctx context.Context, req *ReadIPv6Request) (*iaas.IPv6Net, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	internetOp := iaas.NewInternetOp(s.caller)
	current, err := internetOp.Read(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, fmt.Errorf("reading the internet resource[%s] failed: %s", req.ID, err)
	}
	if len(current.Switch.IPv6Nets) == 0 {
		return nil, iaas.NewNoResultsError()
	}

	ipv6Op := iaas.NewIPv6NetOp(s.caller)
	return ipv6Op.Read(ctx, req.Zone, current.Switch.IPv6Nets[0].ID)
}

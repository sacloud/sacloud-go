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

package database

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-service-go/serviceutil"
)

func (s *Service) MonitorDisk(req *MonitorDiskRequest) ([]*iaas.MonitorDiskValue, error) {
	return s.MonitorDiskWithContext(context.Background(), req)
}

func (s *Service) MonitorDiskWithContext(ctx context.Context, req *MonitorDiskRequest) ([]*iaas.MonitorDiskValue, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	client := iaas.NewDatabaseOp(s.caller)
	cond, err := serviceutil.MonitorCondition(req.Start, req.End)
	if err != nil {
		return nil, err
	}

	values, err := client.MonitorDisk(ctx, req.Zone, req.ID, cond)
	if err != nil {
		return nil, err
	}
	return values.Values, nil
}

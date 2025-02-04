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

package mobilegateway

import (
	"context"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/require"
)

func TestMobileGatewayService_convertCreateRequest(t *testing.T) {
	ctx := context.Background()
	name := testutil.ResourceName("mobile-gateway-service-create")
	zone := testutil.TestZone()
	caller := testutil.SingletonAPICaller()

	// setup
	swOp := iaas.NewSwitchOp(caller)
	sw, err := swOp.Create(ctx, zone, &iaas.SwitchCreateRequest{Name: name})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		swOp.Delete(ctx, zone, sw.ID) //nolint
	}()

	// test
	cases := []struct {
		in     *CreateRequest
		expect *ApplyRequest
	}{
		{
			in: &CreateRequest{
				Zone:        zone,
				Name:        name,
				Description: "description",
				Tags:        types.Tags{"tag1", "tag2"},
				PrivateInterface: &PrivateInterfaceSetting{
					SwitchID:       sw.ID,
					IPAddress:      "192.168.0.1",
					NetworkMaskLen: 24,
				},
				StaticRoutes: []*iaas.MobileGatewayStaticRoute{
					{
						Prefix:  "192.168.1.0/24",
						NextHop: "192.168.0.2",
					},
				},
				InternetConnectionEnabled:       true,
				InterDeviceCommunicationEnabled: true,
				DNS: &DNSSetting{
					DNS1: "8.8.8.8",
					DNS2: "8.8.4.4",
				},
				SIMs: nil,
				TrafficConfig: &TrafficConfig{
					TrafficQuotaInMB:     10,
					BandWidthLimitInKbps: 128,
					EmailNotifyEnabled:   true,
					AutoTrafficShaping:   true,
				},
				NoWait:          false,
				BootAfterCreate: true,
			},
			expect: &ApplyRequest{
				Zone:        zone,
				Name:        name,
				Description: "description",
				Tags:        types.Tags{"tag1", "tag2"},
				PrivateInterface: &PrivateInterfaceSetting{
					SwitchID:       sw.ID,
					IPAddress:      "192.168.0.1",
					NetworkMaskLen: 24,
				},
				StaticRoutes: []*iaas.MobileGatewayStaticRoute{
					{
						Prefix:  "192.168.1.0/24",
						NextHop: "192.168.0.2",
					},
				},
				InternetConnectionEnabled:       true,
				InterDeviceCommunicationEnabled: true,
				DNS: &DNSSetting{
					DNS1: "8.8.8.8",
					DNS2: "8.8.4.4",
				},
				SIMs: nil,
				TrafficConfig: &TrafficConfig{
					TrafficQuotaInMB:     10,
					BandWidthLimitInKbps: 128,
					EmailNotifyEnabled:   true,
					AutoTrafficShaping:   true,
				},
				NoWait:          false,
				BootAfterCreate: true,
			},
		},
	}

	for _, tc := range cases {
		require.EqualValues(t, tc.expect, tc.in.ApplyRequest())
	}
}

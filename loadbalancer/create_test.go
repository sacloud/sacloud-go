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

package loadbalancer

import (
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancerService_convertCreateRequest(t *testing.T) {
	name := testutil.ResourceName("load-balancer-service")
	zone := testutil.TestZone()

	cases := []struct {
		in     *CreateRequest
		expect *ApplyRequest
	}{
		{
			in: &CreateRequest{
				Zone:           zone,
				Name:           name,
				Description:    "desc",
				Tags:           types.Tags{"tag1", "tag2"},
				SwitchID:       102,
				PlanID:         types.LoadBalancerPlans.Standard,
				VRID:           10,
				IPAddresses:    []string{"192.168.0.101", "192.168.0.102"},
				NetworkMaskLen: 24,
				DefaultRoute:   "192.168.0.1",
				VirtualIPAddresses: []*iaas.LoadBalancerVirtualIPAddress{
					{
						VirtualIPAddress: "192.168.0.201",
						Port:             80,
						DelayLoop:        10,
						SorryServer:      "192.168.0.99",
						Description:      "desc",
						Servers: []*iaas.LoadBalancerServer{
							{
								IPAddress: "192.168.0.202",
								Port:      80,
								Enabled:   true,
								HealthCheck: &iaas.LoadBalancerServerHealthCheck{
									Protocol:     types.LoadBalancerHealthCheckProtocols.HTTP,
									Path:         "/",
									ResponseCode: 200,
								},
							},
							{
								IPAddress: "192.168.0.203",
								Port:      80,
								Enabled:   true,
								HealthCheck: &iaas.LoadBalancerServerHealthCheck{
									Protocol:     types.LoadBalancerHealthCheckProtocols.HTTP,
									Path:         "/",
									ResponseCode: 200,
								},
							},
						},
					},
				},
				NoWait: true,
			},
			expect: &ApplyRequest{
				Zone:           zone,
				Name:           name,
				Description:    "desc",
				Tags:           types.Tags{"tag1", "tag2"},
				SwitchID:       102,
				PlanID:         types.LoadBalancerPlans.Standard,
				VRID:           10,
				IPAddresses:    []string{"192.168.0.101", "192.168.0.102"},
				NetworkMaskLen: 24,
				DefaultRoute:   "192.168.0.1",
				VirtualIPAddresses: []*iaas.LoadBalancerVirtualIPAddress{
					{
						VirtualIPAddress: "192.168.0.201",
						Port:             80,
						DelayLoop:        10,
						SorryServer:      "192.168.0.99",
						Description:      "desc",
						Servers: []*iaas.LoadBalancerServer{
							{
								IPAddress: "192.168.0.202",
								Port:      80,
								Enabled:   true,
								HealthCheck: &iaas.LoadBalancerServerHealthCheck{
									Protocol:     types.LoadBalancerHealthCheckProtocols.HTTP,
									Path:         "/",
									ResponseCode: 200,
								},
							},
							{
								IPAddress: "192.168.0.203",
								Port:      80,
								Enabled:   true,
								HealthCheck: &iaas.LoadBalancerServerHealthCheck{
									Protocol:     types.LoadBalancerHealthCheckProtocols.HTTP,
									Path:         "/",
									ResponseCode: 200,
								},
							},
						},
					},
				},
				NoWait: true,
			},
		},
	}

	for _, tc := range cases {
		require.EqualValues(t, tc.expect, tc.in.ApplyRequest())
	}
}

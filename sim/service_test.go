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

package sim

import (
	"os"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/pointer"
)

func TestSIMService_CRUD(t *testing.T) {
	svc := New(testutil.SingletonAPICaller())
	name := testutil.ResourceName("sim-service")

	iccid := os.Getenv("SAKURACLOUD_SIM_ICCID")
	passcode := os.Getenv("SAKURACLOUD_SIM_PASSCODE")

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		PreCheck:           testutil.PreCheckEnvsFunc("SAKURACLOUD_SIM_ICCID", "SAKURACLOUD_SIM_PASSCODE"),
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Setup:              nil,
		IgnoreStartupWait:  true,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Create(&CreateRequest{
					Name:     name,
					ICCID:    iccid,
					PassCode: passcode,
					Carriers: []*iaas.SIMNetworkOperatorConfig{
						{Allow: true, Name: types.SIMOperators.SoftBank.String()},
					},
				})
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, v interface{}) error {
				sim := v.(*iaas.SIM)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, sim.Info, "SIM.Info"),
				)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Read(&ReadRequest{ID: ctx.ID})
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, v interface{}) error {
				sim := v.(*iaas.SIM)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, sim.Info, "SIM.Info"),
				)
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					return svc.Find(&FindRequest{Names: []string{name}})
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, v interface{}) error {
					values := v.([]*iaas.SIM)
					return testutil.DoAsserts(
						testutil.AssertLenFunc(t, values, 1, "SIMs"),
						testutil.AssertNotNilFunc(t, values[0].Info, "SIMs[0].Info"),
					)
				},
				SkipExtractID: true,
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
					return svc.Update(&UpdateRequest{
						ID:   ctx.ID,
						Name: pointer.NewString(name + "-upd"),
					})
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, v interface{}) error {
					sim := v.(*iaas.SIM)
					return testutil.DoAsserts(
						testutil.AssertNotNilFunc(t, sim.Info, "SIM.Info"),
					)
				},
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) error {
				return svc.Delete(&DeleteRequest{ID: ctx.ID})
			},
		},
		Cleanup: nil,
	})
}

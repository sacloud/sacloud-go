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

package containerregistry

import (
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/containerregistry/builder"
	"github.com/sacloud/packages-go/pointer"
)

func TestContainerRegistryService_CRUD(t *testing.T) {
	svc := New(testutil.SingletonAPICaller())
	name := testutil.ResourceName("container-registry")

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		PreCheck:           nil,
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Setup:              nil,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Create(&CreateRequest{
					Name:           name,
					Description:    "test",
					Tags:           types.Tags{"tag1", "tag2"},
					AccessLevel:    types.ContainerRegistryAccessLevels.ReadOnly,
					VirtualDomain:  name + ".usacloud.jp",
					SubDomainLabel: name,
					Users: []*builder.User{
						{
							UserName:   "user1",
							Password:   "password1",
							Permission: types.ContainerRegistryPermissions.ReadOnly,
						},
						{
							UserName:   "user2",
							Password:   "password2",
							Permission: types.ContainerRegistryPermissions.ReadOnly,
						},
					},
				})
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Read(&ReadRequest{ID: ctx.ID})
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
					return svc.Update(&UpdateRequest{
						ID:            ctx.ID,
						Name:          pointer.NewString(name + "-upd"),
						Description:   pointer.NewString("test-upd"),
						Tags:          &types.Tags{"tag1-upd", "tag2-upd"},
						AccessLevel:   &types.ContainerRegistryAccessLevels.ReadOnly,
						VirtualDomain: pointer.NewString(name + "-upd.usacloud.jp"),
						Users: &[]*builder.User{
							{
								UserName:   "user1",
								Password:   "password1",
								Permission: types.ContainerRegistryPermissions.ReadWrite,
							},
							{
								UserName:   "user3",
								Password:   "password3",
								Permission: types.ContainerRegistryPermissions.ReadOnly,
							},
						},
					})
				},
			},
		},
		Shutdown: nil,
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) error {
				return svc.Delete(&DeleteRequest{ID: ctx.ID})
			},
		},
		Cleanup: nil,
	})
}

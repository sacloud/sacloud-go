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

package localrouter

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

// Builder ローカルルータの構築を行う
type Builder struct {
	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID

	Switch       *iaas.LocalRouterSwitch
	Interface    *iaas.LocalRouterInterface
	Peers        []*iaas.LocalRouterPeer
	StaticRoutes []*iaas.LocalRouterStaticRoute

	SettingsHash string

	Client *APIClient
}

// Validate 設定値の検証
func (b *Builder) Validate(_ context.Context) error {
	return nil
}

// Build ローカルルータの作成や設定をまとめて行う
func (b *Builder) Build(ctx context.Context) (*iaas.LocalRouter, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	localRouter, err := b.Client.LocalRouter.Create(ctx, &iaas.LocalRouterCreateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
	})
	if err != nil {
		return nil, err
	}

	if b.hasNetworkSettings() {
		lr, err := b.Client.LocalRouter.UpdateSettings(ctx, localRouter.ID, &iaas.LocalRouterUpdateSettingsRequest{
			Switch:       b.Switch,
			Interface:    b.Interface,
			StaticRoutes: b.StaticRoutes,
			SettingsHash: b.SettingsHash,
		})
		if err != nil {
			return localRouter, err
		}
		localRouter = lr

		if len(b.Peers) > 0 {
			lr, err := b.Client.LocalRouter.UpdateSettings(ctx, localRouter.ID, &iaas.LocalRouterUpdateSettingsRequest{
				Switch:       localRouter.Switch,
				Interface:    localRouter.Interface,
				StaticRoutes: localRouter.StaticRoutes,
				Peers:        b.Peers,
				SettingsHash: localRouter.SettingsHash,
			})
			if err != nil {
				return localRouter, err
			}
			localRouter = lr
		}
	}

	return localRouter, nil
}

func (b *Builder) hasNetworkSettings() bool {
	return b.Interface != nil && b.Switch != nil &&
		b.Interface.NetworkMaskLen > 0 &&
		b.Interface.VirtualIPAddress != "" &&
		len(b.Interface.IPAddress) > 0 &&
		b.Switch.Code != ""
}

// Update ローカルルータの更新
func (b *Builder) Update(ctx context.Context, id types.ID) (*iaas.LocalRouter, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	// check Internet is exists
	_, err := b.Client.LocalRouter.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	localRouter, err := b.Client.LocalRouter.Update(ctx, id, &iaas.LocalRouterUpdateRequest{
		Switch:       b.Switch,
		Interface:    b.Interface,
		Peers:        b.Peers,
		StaticRoutes: b.StaticRoutes,
		SettingsHash: b.SettingsHash,
		Name:         b.Name,
		Description:  b.Description,
		Tags:         b.Tags,
		IconID:       b.IconID,
	})
	if err != nil {
		return nil, err
	}

	return localRouter, nil
}

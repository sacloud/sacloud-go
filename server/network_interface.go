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

package server

import (
	"fmt"

	"github.com/sacloud/iaas-api-go/types"
	serverBuilder "github.com/sacloud/iaas-service-go/server/builder"
	"github.com/sacloud/packages-go/validate"
)

type NetworkInterface struct {
	Upstream       string // スイッチID or "disconnected"(切断) or "shared"(共有セグメント) 省略時は"disconnected"
	PacketFilterID types.ID
	UserIPAddress  string `validate:"omitempty,ipv4"`
}

func (s *NetworkInterface) Validate() error {
	if err := validate.New().Struct(s); err != nil {
		return err
	}
	if s.Upstream == "" || s.Upstream == "shared" || s.Upstream == "disconnected" {
		return nil
	}
	if types.StringID(s.Upstream).IsEmpty() {
		return fmt.Errorf(`upstream require to be "shared" or "disconnected" or SwitchID. value:%v`, s.Upstream)
	}
	return nil
}

func (s *NetworkInterface) NICSettingHolder() serverBuilder.NICSettingHolder {
	switch s.Upstream {
	case "shared":
		return &serverBuilder.SharedNICSetting{PacketFilterID: s.PacketFilterID}
	case "", "disconnected":
		return &serverBuilder.DisconnectedNICSetting{}
	default:
		return &serverBuilder.ConnectedNICSetting{
			SwitchID:         types.StringID(s.Upstream),
			DisplayIPAddress: s.UserIPAddress,
			PacketFilterID:   s.PacketFilterID,
		}
	}
}

func (s *NetworkInterface) AdditionalNICSettingHolder() serverBuilder.AdditionalNICSettingHolder {
	switch s.Upstream {
	case "", "disconnected":
		return &serverBuilder.DisconnectedNICSetting{}
	default:
		return &serverBuilder.ConnectedNICSetting{
			SwitchID:         types.StringID(s.Upstream),
			DisplayIPAddress: s.UserIPAddress,
			PacketFilterID:   s.PacketFilterID,
		}
	}
}

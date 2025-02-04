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
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/validate"
)

type SetDNSRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-" validate:"required"`

	DNS1 string `validate:"required,ipv4"`
	DNS2 string `validate:"required,ipv4"`
}

func (req *SetDNSRequest) Validate() error {
	return validate.New().Struct(req)
}

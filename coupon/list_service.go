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

package coupon

import (
	"context"

	"github.com/sacloud/iaas-api-go"
)

func (s *Service) List() ([]*iaas.Coupon, error) {
	return s.ListWithContext(context.Background())
}

func (s *Service) ListWithContext(ctx context.Context) ([]*iaas.Coupon, error) {
	authOp := iaas.NewAuthStatusOp(s.caller)
	couponOp := iaas.NewCouponOp(s.caller)

	account, err := authOp.Read(ctx)
	if err != nil {
		return nil, err
	}
	searched, err := couponOp.Find(ctx, account.AccountID)
	if err != nil {
		return nil, err
	}
	return searched.Coupons, nil
}

package verifyint64_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifyint64"
)

func TestRequiresIfTrue(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.Bool
		validator validator.Int64
		expErrors int
	}

	testCases := map[string]testCase{
		"base-true": {
			in: types.BoolValue(true),
			validator: verifyint64.RequiresIfTrue(path.Expressions{
				path.MatchRoot("base-true"),
			}...),
			expErrors: 0,
		},
		"base-false": {
			in: types.BoolValue(false),
			validator: verifyint64.RequiresIfTrue(path.Expressions{
				path.MatchRoot("base-true"),
			}...),
			expErrors: 0,
		},
		"depend-false": {
			in: types.BoolValue(true),
			validator: verifyint64.RequiresIfTrue(path.Expressions{
				path.MatchRoot("base-false"),
			}...),
			expErrors: 1,
		},
		"depend-true": {
			in: types.BoolValue(true),
			validator: verifyint64.RequiresIfTrue(path.Expressions{
				path.MatchRoot("base-true"),
			}...),
			expErrors: 0,
		},
		"is-null": {
			in: types.BoolNull(),
			validator: verifyint64.RequiresIfTrue(
				path.Expressions{
					path.MatchRoot("base-true"),
				}...,
			),
			expErrors: 0,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validator.Int64Request{
				ConfigValue: boolToInt64(test.in),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"base-true":  schema.BoolAttribute{},
							"base-false": schema.BoolAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"base-true":  tftypes.Bool,
							"base-false": tftypes.Bool,
						},
					}, map[string]tftypes.Value{
						"base-true":  tftypes.NewValue(tftypes.Bool, true),
						"base-false": tftypes.NewValue(tftypes.Bool, false),
					}),
				},
			}
			res := validator.Int64Response{}
			test.validator.ValidateInt64(context.TODO(), req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatalf("expected %d error(s), got none", test.expErrors)
			}

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}

func boolToInt64(b types.Bool) types.Int64 {
	if b.ValueBool() {
		return types.Int64Value(1)
	}
	return types.Int64Value(0)
}

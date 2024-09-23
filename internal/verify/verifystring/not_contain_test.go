package verifystring_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifystring"
)

func TestContainValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.String
		validator validator.String
		expErrors int
	}

	testCases := map[string]testCase{
		"contains-substring": {
			in:        types.StringValue("hello123"),
			validator: verifystring.NotContain("hello"),
			expErrors: 1,
		},
		"does-not-contain-substring": {
			in:        types.StringValue("goodbye"),
			validator: verifystring.NotContain("hello"),
			expErrors: 0,
		},
		"confirm-other-attr": {
			in:        types.StringValue("gooyehi99993"),
			validator: verifystring.NotContain(path.MatchRoot("does-not-contain-substring").String()),
			expErrors: 0,
		},
		"skip-validation-on-null": {
			in:        types.StringNull(),
			validator: verifystring.NotContain("hello"),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in:        types.StringUnknown(),
			validator: verifystring.NotContain("hello"),
			expErrors: 0,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				ConfigValue: test.in,
			}
			res := validator.StringResponse{}
			test.validator.ValidateString(context.TODO(), req, &res)

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

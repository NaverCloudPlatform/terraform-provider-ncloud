package verifystring

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type notContainValidator struct {
	target string
}

func NotContain(target string) validator.String {
	return notContainValidator{
		target: target,
	}
}

func (v notContainValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Value must not contain the string '%s'.", v.target)
}

func (v notContainValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v notContainValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if strings.Contains(req.ConfigValue.ValueString(), v.target) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Value",
			fmt.Sprintf("Value must not contain the string '%s'.", v.target),
		)
	}
}

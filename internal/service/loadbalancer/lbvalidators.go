package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Int32 = idleTimeoutValidator{}

type idleTimeoutValidator struct {
	ty string
}

func (v idleTimeoutValidator) Description(_ context.Context) string {
	return fmt.Sprintf("idleTimeout is not supported for the %s type load balancer", v.ty)
}

func (v idleTimeoutValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v idleTimeoutValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		// If code block does not exist, config is valid.
		return
	}

	typePath := req.Path.ParentPath().AtName("type")

	var m types.String

	diags := req.Config.GetAttribute(ctx, typePath, &m)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if m.IsNull() || m.IsUnknown() {
		// Only validate if mode value is known.
		return
	}

	if m.ValueString() == v.ty {
		resp.Diagnostics.AddAttributeError(
			typePath,
			"Unsupported attribute combination (network type and idleTimeout)",
			v.Description(ctx),
		)
	}
}

func checkUnsupportedIdleTimeout(ty string) validator.Int32 {
	return idleTimeoutValidator{
		ty: ty,
	}
}

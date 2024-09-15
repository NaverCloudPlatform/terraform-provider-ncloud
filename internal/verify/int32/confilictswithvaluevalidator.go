package verify

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Int32 = conflictsWithValueValidator{}

type conflictsWithValueValidator struct {
	ty    string
	value string
}

func (v conflictsWithValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("conflicts with %s value in %s attritube ", v.value, v.ty)
}

func (v conflictsWithValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v conflictsWithValueValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		// If code block does not exist, config is valid.
		return
	}

	typePath := req.Path.ParentPath().AtName(v.ty)

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

	if m.ValueString() == v.value {
		resp.Diagnostics.AddAttributeError(
			typePath,
			"Unsupported attribute combination (network type and idleTimeout)",
			v.Description(ctx),
		)
	}
}

func ConflictsWithVaule(ty string, value string) validator.Int32 {
	return conflictsWithValueValidator{
		ty:    ty,
		value: value,
	}
}

package verify

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Int32 = conflictsWithValueValidator[attr.Value]{}

type conflictsWithValueValidator[V attr.Value] struct {
	expr  path.Expression
	value V
}

func (v conflictsWithValueValidator[V]) Description(_ context.Context) string {
	return fmt.Sprintf("conflicts with attribute `%s` value %s", v.expr, v.value)
}

func (v conflictsWithValueValidator[V]) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v conflictsWithValueValidator[V]) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		// If code block does not exist, config is valid.
		return
	}

	var m V

	targetPaths, diags := req.Config.PathMatches(ctx, v.expr)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(targetPaths) > 1 {
		resp.Diagnostics.AddError(
			"Ambiguous expression was set to validator",
			"The path traversed by the expression turned out multiple paths, can't determine which path to validate")
		return
	}
	p := targetPaths[0]

	diags = req.Config.GetAttribute(ctx, p, &m)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if m.IsNull() || m.IsUnknown() {
		// Only validate if mode value is known.
		return
	}

	if m.Equal(v.value) {
		resp.Diagnostics.AddAttributeError(
			p,
			fmt.Sprintf("Conflicts `%s` attritube value %s", v.expr, m.String()),
			v.Description(ctx),
		)
	}
}

func ConflictsWithValue[V attr.Value](expr path.Expression, value V) validator.Int32 {
	return conflictsWithValueValidator[V]{
		expr:  expr,
		value: value,
	}
}

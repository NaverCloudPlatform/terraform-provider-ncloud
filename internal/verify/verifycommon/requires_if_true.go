package verifycommon

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var (
	_ validator.String = RequiresIfTrueValidator{}
	_ validator.Bool   = RequiresIfTrueValidator{}
	_ validator.Int64  = RequiresIfTrueValidator{}
)

type RequiresIfTrueValidator struct {
	PathExpressions path.Expressions
}

type RequiresIfTrueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type RequiresIfTrueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (rt RequiresIfTrueValidator) Description(ctx context.Context) string {
	return rt.MarkdownDescription(ctx)
}

func (rt RequiresIfTrueValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Requires that if a boolean attribute is true, the dependent attribute %q must be set.", rt.PathExpressions)
}

func (rt RequiresIfTrueValidator) Validate(ctx context.Context, req RequiresIfTrueValidatorRequest, res *RequiresIfTrueValidatorResponse) {
	if req.ConfigValue.IsNull() {
		return
	}

	// Merge the path expressions
	expressions := req.PathExpression.MergeExpressions(rt.PathExpressions...)

	for _, expression := range expressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)
		res.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(req.Path) {
				continue
			}

			var mpVal attr.Value
			diags := req.Config.GetAttribute(ctx, mp, &mpVal)
			res.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attributes have a known value
			if mpVal.IsUnknown() {
				return
			}

			// Check if the attribute value is `true`
			if mpVal.IsNull() || mpVal.String() == "false" || mpVal.String() == "0" {
				res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.Path,
					fmt.Sprintf("Attribute %q must be set to `true` when %q is specified", mp, req.Path),
				))
			}
		}
	}
}

func (rt RequiresIfTrueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := RequiresIfTrueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RequiresIfTrueValidatorResponse{}

	rt.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (rt RequiresIfTrueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := RequiresIfTrueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RequiresIfTrueValidatorResponse{}

	rt.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (rt RequiresIfTrueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := RequiresIfTrueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RequiresIfTrueValidatorResponse{}

	rt.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

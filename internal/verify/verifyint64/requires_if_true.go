package verifyint64

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifycommon"
)

func RequiresIfTrue(expressions ...path.Expression) validator.Int64 {
	return verifycommon.RequiresIfTrueValidator{
		PathExpressions: expressions,
	}
}

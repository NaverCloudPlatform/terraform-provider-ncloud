package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EmptyStringToNull converts a Framework empty string to null.
// This ensures consistency where "" and null are equivalent.
func EmptyStringToNull(s types.String) types.String {
	if s == types.StringValue("") {
		return types.StringNull()
	}

	return s
}

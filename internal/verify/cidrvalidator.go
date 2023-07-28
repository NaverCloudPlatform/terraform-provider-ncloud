package verify

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = cidrBlockValidator{}

// cidrBlock validates that a string Attribute's following the cidr block format
type cidrBlockValidator struct {
}

// CidrBlock returns an validator which ensures that string follows the cidr block format
func CidrBlockValidator() []validator.String {
	return []validator.String{
		cidrBlockValidator{},
	}
}

// Description describes the validation in plain text formatting.
func (validator cidrBlockValidator) Description(_ context.Context) string {
	return "string must follow CIDR notation like \"192.0.2.0/24\" as defined in RFC4632 and RFC4291"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator cidrBlockValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v cidrBlockValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if err := ValidateCIDRBlock(value); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"CIDRBlock Type Validation Error",
			err.Error(),
		))
		return
	}
}

// ValidateCIDRBlock validates that the specified CIDR block is valid:
// - The CIDR block parses to an IP address and network
// - The CIDR block is the CIDR block for the network
func ValidateCIDRBlock(cidr string) error {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("%q is not a valid CIDR block: %w", cidr, err)
	}

	if !CIDRBlocksEqual(cidr, ipnet.String()) {
		return fmt.Errorf("%q is not a valid CIDR block; did you mean %q?", cidr, ipnet)
	}

	return nil
}

// CIDRBlocksEqual returns whether or not two CIDR blocks are equal:
// - Both CIDR blocks parse to an IP address and network
// - The string representation of the IP addresses are equal
// - The string representation of the networks are equal
// This function is especially useful for IPv6 CIDR blocks which have multiple valid representations.
func CIDRBlocksEqual(cidr1, cidr2 string) bool {
	ip1, ipnet1, err := net.ParseCIDR(cidr1)
	if err != nil {
		return false
	}
	ip2, ipnet2, err := net.ParseCIDR(cidr2)
	if err != nil {
		return false
	}

	return ip2.String() == ip1.String() && ipnet2.String() == ipnet1.String()
}

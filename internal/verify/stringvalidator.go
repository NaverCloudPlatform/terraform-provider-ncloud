package verify

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func InstanceNameValidator() []validator.String {
	return []validator.String{
		stringvalidator.LengthAtLeast(3),
		stringvalidator.LengthAtMost(30),
		stringvalidator.RegexMatches(
			regexp.MustCompile("^[a-z][a-z0-9-]*$"),
			"must only lowercase letters, numbers and special characters \"-\" are allowed and must start with an alphabetic character",
		),
		stringvalidator.RegexMatches(
			regexp.MustCompile("^.*[^-_]$"),
			"must end with an alphabetic character or number",
		),
	}
}

func ValidateEmptyStringElement(i []interface{}) error {
	for _, v := range i {
		if v == nil || v == "" {
			return fmt.Errorf("empty string element found")
		}
	}
	return nil
}

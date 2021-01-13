package ncloud

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func validateInstanceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) < 3 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be shorter than 3 characters", k))
	}

	if len(value) > 30 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 30 characters", k))
	}

	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%s can only lowercase letters, numbers and special characters \"-\" are allowed and must start with an alphabetic character", k))
	}

	if regexp.MustCompile(`.*(-|_)$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must end with an alphabetic character or number", k))
	}

	return
}

func validatePortRange(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[1-9][0-9]*([-][1-9][0-9]*|[0-9]*)`).MatchString(value) {
		errors = append(errors, fmt.Errorf("Invalid %s type format (eg. 1-65535, 22)", k))
		return
	}

	if !isValidPortRange(value) {
		errors = append(errors, fmt.Errorf("%s must be 1 to 65535", k))
		return
	}

	return
}

func isValidPortRange(value string) bool {
	ports := strings.Split(value, "-")

	if len(ports) == 2 {
		start, err := strconv.Atoi(ports[0])
		if err != nil {
			return false
		}

		end, err := strconv.Atoi(ports[1])
		if err != nil {
			return false
		}

		if start > 65535 || end > 65535 {
			return false
		}

		if start > end {
			return false
		}

		return true
	} else if len(ports) > 2 {
		return false
	}

	portNumber, err := strconv.Atoi(value)
	if err != nil {
		return false
	}

	return portNumber <= 65535
}

func validateOneResult(resultCount int) error {
	if resultCount < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}
	if resultCount > 1 {
		return fmt.Errorf("more than one found results(%d). please change search criteria and try again", resultCount)
	}
	return nil
}

func validateParseDuration(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	duration, err := time.ParseDuration(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as a duration: %s", k, err))
	}
	if duration < 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be greater than zero", k))
	}
	return
}

func ToDiagFunc(validator schema.SchemaValidateFunc) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		attr := path[len(path)-1].(cty.GetAttrStep)
		warnings, errors := validator(v, attr.Name)

		for _, w := range warnings {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       w,
				AttributePath: path,
			})
		}
		for _, err := range errors {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				AttributePath: path,
			})
		}
		return diags
	}
}

package ncloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func validateInternetLineTypeCode(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "PUBLC" && value != "GLBL" {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, "PUBLC", "GLBL"))
	}
	return
}

var serverNamePattern = regexp.MustCompile(`[(A-Z|a-z|0-9|\\-|\\*)]+`)

func validateServerName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validateStringLengthInRange(3, 30)
	if len(value) < 3 || len(value) > 30 {
		errors = append(errors, fmt.Errorf("must be a valid %q characters between 1 and 30", k))
	}

	// 알파벳, 숫자, 하이픈(-) 으로만 구성 가능하며, 마지막 문자는 하이픈(-)이 올 수 없다.
	if !serverNamePattern.MatchString(value) || strings.LastIndex(value, "-") == (len(value)-1) {
		errors = append(errors, fmt.Errorf("server name is composed of alphabets, numbers, hyphen (-) and wild card (*).<br> Hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input.<br> Maximum length is 63Bytes, and the minimum is 1Byte"))
	}
	return
}

func validateStringLengthInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if len(value) < min || len(value) > max {
			errors = append(errors, fmt.Errorf("must be a valid %q characters between %d and %d", k, min, max))
		}
		return
	}
}

func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

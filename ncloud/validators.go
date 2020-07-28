package ncloud

import (
	"fmt"
	"regexp"
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

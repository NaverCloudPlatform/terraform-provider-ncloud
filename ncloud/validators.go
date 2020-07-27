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

	if !regexp.MustCompile(`^[A-Za-z0-9-*]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q can only composed of alphabets, numbers, hyphen (-) and wild card (*)", k))
	}
	if !regexp.MustCompile(`.*[^\\-]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q with hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input", k))
	}

	return
}

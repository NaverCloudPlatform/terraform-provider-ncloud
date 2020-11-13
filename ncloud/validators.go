package ncloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

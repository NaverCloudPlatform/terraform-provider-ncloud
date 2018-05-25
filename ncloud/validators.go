package ncloud

import "fmt"

func validateInternetLineTypeCode(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "PUBLC" && value != "GLBL" {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, "PUBLC", "GLBL"))
	}
	return
}

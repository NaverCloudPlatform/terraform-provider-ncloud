package common

import "fmt"

// ErrorRequiredArgOnVpc return error for required on vpc
func ErrorRequiredArgOnVpc(name string) error {
	return fmt.Errorf("missing required argument: The argument \"%s\" is required on vpc", name)
}

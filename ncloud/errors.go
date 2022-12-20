package ncloud

import "fmt"

// NotSupportClassic return error for not support classic
func NotSupportClassic(name string) error {
	return fmt.Errorf("%s doesn't support classic", name)
}

// NotSupportVpc return error for not support vpc
func NotSupportVpc(name string) error {
	return fmt.Errorf("%s doesn't support vpc", name)
}

// ErrorRequiredArgOnVpc return error for required on vpc
func ErrorRequiredArgOnVpc(name string) error {
	return fmt.Errorf("missing required argument: The argument \"%s\" is required on vpc", name)
}

// ErrorRequiredArgOnClassic return error for required on classic
func ErrorRequiredArgOnClassic(name string) error {
	return fmt.Errorf("missing required argument: The argument \"%s\" is required on classic", name)
}

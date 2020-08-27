package ncloud

import "fmt"

//NotSupportClassic return error for not support classic
func NotSupportClassic(name string) error {
	return fmt.Errorf("%s doesn't support classic", name)
}

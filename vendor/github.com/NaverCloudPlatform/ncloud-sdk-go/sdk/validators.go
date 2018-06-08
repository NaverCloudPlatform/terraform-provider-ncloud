package sdk

import (
	"fmt"
	"strings"
)

var boolValueStrings = []string{"true", "false"}

func validateRequiredField(key string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if v == "" {
			return fmt.Errorf("%s field is required", key)
		}
	case int:
		if v == 0 {
			return fmt.Errorf("%s field is required", key)
		}
	}
	return nil
}

func validateIntegerInRange(key string, value interface{}, min, max int) error {
	v, ok := value.(int)
	if !ok {
		return fmt.Errorf("expected type of %s to be int", key)
	}
	if v < min {
		return fmt.Errorf("%q cannot be lower than %d: %d", key, min, value)
	}
	if v > max {
		return fmt.Errorf("%q cannot be higher than %d: %d", key, max, value)
	}
	return nil
}

func validateMultipleValue(key string, value int, multiple int) error {
	if int(value/multiple)*multiple != value {
		return fmt.Errorf("%s must be a multiple of %d", key, multiple)
	}
	return nil
}

func validateIncludeValues(key string, value string, includeValues []string) error {
	for _, included := range includeValues {
		if value == included {
			return nil
		}
	}
	return fmt.Errorf("%s should be %s", key, strings.Join(includeValues, " or "))
}

func validateIncludeValuesIgnoreCase(key string, value string, includeValues []string) error {
	for _, included := range includeValues {
		if strings.EqualFold(value, included) {
			return nil
		}
	}
	return fmt.Errorf("%s should be %s", key, strings.Join(includeValues, " or "))
}

func validateBoolValue(key string, value string) error {
	for _, included := range boolValueStrings {
		if value == included {
			return nil
		}
	}
	return fmt.Errorf("%s should be %s", key, strings.Join(boolValueStrings, " or "))
}

func validateStringMaxLen(key string, value interface{}, max int) error {
	return validateStringLenBetween(key, value, 0, max)
}

func validateStringLenBetween(key string, value interface{}, min, max int) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected type of %s to be string", key)
	}
	if len(v) < min || len(v) > max {
		return fmt.Errorf("expected length of %s to be in the range (%d - %d), got %s", key, min, max, v)
	}
	return nil
}

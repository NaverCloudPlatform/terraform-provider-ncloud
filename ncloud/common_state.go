package ncloud

import (
	"encoding/json"
	"reflect"
	"unicode"
)

// VpcCommonStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch a instances
func VpcCommonStateRefreshFunc(instance interface{}, err error, statusName string) (interface{}, string, error) {
	if err != nil {
		return nil, "", err
	}

	if instance == nil || reflect.ValueOf(instance).IsNil() {
		return instance, "TERMINATED", nil
	}

	b, err := json.Marshal(instance)
	if err != nil {
		return nil, "", err
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)

	a := []rune(statusName)
	a[0] = unicode.ToLower(a[0])
	statusName = string(a)

	u, ok := m[statusName].(map[string]interface{})
	if ok {
		return instance, u["code"].(string), nil
	}

	return instance, "", nil
}

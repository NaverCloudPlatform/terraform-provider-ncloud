package common

import (
	"encoding/json"
	"regexp"
)

// MarshalUnchecked return the JSON encoding of value
// It does not check for errors about marshaling failing, so be careful to use
func MarshalUnchecked(value interface{}) []byte {
	v, _ := json.Marshal(value)
	return v
}

func MarshalUncheckedString(value interface{}) string {
	return string(MarshalUnchecked(value))
}

func ReplaceNull(s string) string {
	re := regexp.MustCompile(`:<null>`)
	return re.ReplaceAllString(s, ":null")
}

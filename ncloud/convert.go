package ncloud

import (
	"encoding/base64"
	"strconv"
)

func String(n int) string {
	return strconv.Itoa(n)
}

func StringList(input []interface{}) []string {
	vs := make([]string, 0, len(input))
	for _, v := range input {
		vs = append(vs, v.(string))
	}
	return vs
}

func Base64Decode(sEncData string) string {
	v, err := base64.StdEncoding.DecodeString(sEncData)
	if err != nil {
		v = []byte(sEncData)
	}
	return string(v)
}

func Int(v interface{}) (int, error) {
	var value int
	switch v.(type) {
	case int:
		value = v.(int)
	case string:
		converted, err := strconv.Atoi(v.(string))
		if err != nil {
			return 0, err
		}
		value = converted
	}

	return value, nil

}

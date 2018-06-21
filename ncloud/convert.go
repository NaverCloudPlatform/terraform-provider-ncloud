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

func Base64Decode(data string) string {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	base64.StdEncoding.Decode(base64Text, []byte(data))
	return string(base64Text)
}

func Base64Encode(data string) string {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(base64Text, []byte(data))
	return string(base64Text)
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

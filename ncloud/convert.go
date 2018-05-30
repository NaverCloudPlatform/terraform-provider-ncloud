package ncloud

import "encoding/base64"

func StringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
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

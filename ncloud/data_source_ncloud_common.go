package ncloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"io/ioutil"
	"os"
)

// Generates a hash for the set hash function used by the ID
func dataResourceIdHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func writeToFile(filePath string, data interface{}) {
	err := os.Remove(filePath)
	if err != nil {
		return
	}

	if bs, err := json.MarshalIndent(data, "", "\t"); err == nil {
		str := string(bs)
		_ = ioutil.WriteFile(filePath, []byte(str), 777)
	}
}

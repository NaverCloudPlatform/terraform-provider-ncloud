package ncloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

// Generates a hash for the set hash function used by the ID
func dataResourceIdHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func writeToFile(filePath string, data interface{}) error {
	log.Printf("[INFO] WriteToFile FilePath: %s", filePath)

	if err := os.Remove(filePath); err != nil && os.IsNotExist(err) != true {
		return err
	}

	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	str := string(bs)
	return ioutil.WriteFile(filePath, []byte(str), 777)
}

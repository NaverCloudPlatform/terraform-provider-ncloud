package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Generates a hash for the set hash function used by the ID
func DataResourceIdHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", Hashcode(buf.String()))
}

func WriteToFile(filePath string, data interface{}) error {
	log.Printf("[INFO] WriteToFile FilePath: %s", filePath)

	if err := os.Remove(filePath); err != nil && os.IsExist(err) {
		return err
	}

	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	str := string(bs)
	return os.WriteFile(filePath, []byte(str), 0777)
}

package ncloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/hashcode"
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
	log.Printf("[INFO] WriteToFile FilaPath: %s", filePath)

	if err := os.Remove(filePath); err != nil {
		return err
	}

	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	str := string(bs)
	return ioutil.WriteFile(filePath, []byte(str), 777)
}

func validateOneResult(resultCount int) error {
	if resultCount < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}
	if resultCount > 1 {
		return fmt.Errorf("more than one found results. please change search criteria and try again")
	}
	return nil
}

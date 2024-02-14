package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

func WriteStringListToFile(path string, list types.List) error {
	var dataList []string

	for _, v := range list.Elements() {
		var data string
		if err := json.Unmarshal([]byte(v.String()), &data); err != nil {
			return err
		}
		dataList = append(dataList, data)
	}

	if err := WriteToFile(path, dataList); err != nil {
		return err
	}
	return nil
}

func WriteImageProductToFile(path string, images types.List) error {
	var imagesToJson []imageProductToJson

	for _, image := range images.Elements() {
		imageJson := imageProductToJson{}
		if err := json.Unmarshal([]byte(image.String()), &imageJson); err != nil {
			return err
		}
		imagesToJson = append(imagesToJson, imageJson)
	}

	if err := WriteToFile(path, imagesToJson); err != nil {
		return err
	}
	return nil
}

type imageProductToJson struct {
	ProductCode    string `json:"product_code"`
	GenerationCode string `json:"generation_code"`
	ProductName    string `json:"product_name"`
	ProductType    string `json:"product_type"`
	PlatformType   string `json:"platform_type"`
	OsInformation  string `json:"os_information"`
}

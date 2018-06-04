package ncloud

import (
	"encoding/json"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

var commonCodeSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"code": {
			Type: schema.TypeString,
		},
		"code_name": {
			Type: schema.TypeString,
		},
	},
}

func logErrorResponse(tag string, err error, args interface{}) {
	param, _ := json.Marshal(args)
	log.Printf("[DEBUG] %s error params=%s, err=%s", tag, param, err)
}

func logCommonResponse(tag string, args interface{}, commonResponse common.CommonResponse) {
	param, _ := json.Marshal(args)
	result := fmt.Sprintf("RequestID: %s, ReturnCode: %d, ReturnMessage: %s", commonResponse.RequestID, commonResponse.ReturnCode, commonResponse.ReturnMessage)
	log.Printf("[DEBUG] %s success params=%s, response=%s", tag, param, result)
}

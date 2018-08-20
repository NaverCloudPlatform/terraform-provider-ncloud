package ncloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudLoginKey() *schema.Resource {
	return &schema.Resource{
		Read:   resourceNcloudLoginKeyRead,
		Create: resourceNcloudLoginKeyCreate,
		Update: nil,
		Delete: resourceNcloudLoginKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(3, 30),
				Description:  "Key name to generate. If the generated key name exists, an error occurs.",
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudLoginKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	keyName := d.Get("key_name").(string)

	loginKey, err := getLoginKey(client, keyName)
	if err != nil {
		return err
	}

	if loginKey != nil {
		d.Set("fingerprint", loginKey.Fingerprint)
		d.Set("create_date", loginKey.CreateDate)
	}

	return nil
}

func resourceNcloudLoginKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	keyName := d.Get("key_name").(string)

	resp, err := client.server.V2Api.CreateLoginKey(&server.CreateLoginKeyRequest{KeyName: ncloud.String(keyName)})
	if err != nil {
		logErrorResponse("CreateLoginKey", err, keyName)
		return err
	}
	logCommonResponse("CreateLoginKey", keyName, GetCommonResponse(resp))

	d.SetId(keyName)
	d.Set("private_key", resp.PrivateKey)

	return resourceNcloudLoginKeyRead(d, meta)
}

func resourceNcloudLoginKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	keyName := d.Get("key_name").(string)
	return waitForDeleteLoginKey(client, keyName)
}

func getLoginKey(client *NcloudAPIClient, keyName string) (*server.LoginKey, error) {
	reqParams := &server.GetLoginKeyListRequest{
		KeyName: ncloud.String(keyName),
	}
	resp, err := client.server.V2Api.GetLoginKeyList(reqParams)
	if err != nil {
		logErrorResponse("GetLoginKeyList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetLoginKeyList", reqParams, GetCommonResponse(resp))

	if len(resp.LoginKeyList) > 0 {
		return resp.LoginKeyList[0], err
	}

	return nil, err
}

func waitForDeleteLoginKey(client *NcloudAPIClient, keyName string) error {

	c1 := make(chan error, 1)

	go func() {
		for {
			resp, err := client.server.V2Api.DeleteLoginKey(&server.DeleteLoginKeyRequest{KeyName: ncloud.String(keyName)})

			if err == nil || *resp.ReturnCode == "200" {
				c1 <- nil
				return
			}
			// ignore resp.ReturnCode == 10407
			logCommonResponse("DeleteLoginKey", keyName, GetCommonResponse(resp))

			log.Printf("[DEBUG] Wait to delete login key (%s)", keyName)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : Wait to delete login key (%s)", keyName)
	}
}

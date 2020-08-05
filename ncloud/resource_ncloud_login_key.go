package ncloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNcloudLoginKey() *schema.Resource {
	return &schema.Resource{
		Read:   resourceNcloudLoginKeyRead,
		Create: resourceNcloudLoginKeyCreate,
		Update: resourceNcloudLoginKeyUpdate,
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
				ValidateFunc: validation.StringLenBetween(3, 30),
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
		},
	}
}

func resourceNcloudLoginKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	keyName := d.Get("key_name").(string)

	loginKey, err := getLoginKey(client, keyName)
	if err != nil {
		return err
	}

	if loginKey != nil {
		d.Set("fingerprint", loginKey.Fingerprint)
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudLoginKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudLoginKeyRead(d, meta)
}

func resourceNcloudLoginKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	keyName := d.Get("key_name").(string)
	reqParams := &server.CreateLoginKeyRequest{KeyName: ncloud.String(keyName)}

	logCommonRequest("CreateLoginKey", reqParams)

	resp, err := client.server.V2Api.CreateLoginKey(reqParams)
	if err != nil {
		logErrorResponse("CreateLoginKey", err, keyName)
		return err
	}
	logCommonResponse("CreateLoginKey", GetCommonResponse(resp))

	d.SetId(keyName)
	d.Set("private_key", resp.PrivateKey)

	time.Sleep(time.Second * 1) // for internal Master / Slave DB sync

	return resourceNcloudLoginKeyRead(d, meta)
}

func resourceNcloudLoginKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	keyName := d.Get("key_name").(string)

	if err := deleteLoginKey(client, keyName); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func getLoginKeyList(client *NcloudAPIClient, keyName *string) (*server.GetLoginKeyListResponse, error) {
	reqParams := &server.GetLoginKeyListRequest{}
	if keyName != nil {
		reqParams.KeyName = keyName
	}

	logCommonRequest("GetLoginKeyList", reqParams)

	resp, err := client.server.V2Api.GetLoginKeyList(reqParams)
	if err != nil {
		logErrorResponse("GetLoginKeyList", err, reqParams)
		return nil, err
	}

	var totalRowsLog string
	if resp != nil {
		totalRowsLog = fmt.Sprintf("totalRows: %d", ncloud.Int32Value(resp.TotalRows))
	}
	logCommonResponse("GetLoginKeyList", GetCommonResponse(resp), totalRowsLog)
	return resp, nil
}

func getLoginKey(client *NcloudAPIClient, keyName string) (*server.LoginKey, error) {
	resp, err := getLoginKeyList(client, ncloud.String(keyName))
	if len(resp.LoginKeyList) > 0 {
		return resp.LoginKeyList[0], err
	}

	return nil, err
}

func deleteLoginKey(client *NcloudAPIClient, keyName string) error {
	reqParams := &server.DeleteLoginKeyRequest{KeyName: ncloud.String(keyName)}
	logCommonRequest("DeleteLoginKey", reqParams)

	resp, err := client.server.V2Api.DeleteLoginKey(reqParams)
	if err != nil {
		logErrorResponse("DeleteLoginKey", err, keyName)
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}
	logCommonResponse("DeleteLoginKey", commonResponse)

	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := getLoginKeyList(client, ncloud.String(keyName))
			if err != nil {
				return 0, "", err
			}

			if ncloud.Int32Value(resp.TotalRows) == 0 {
				return 0, "OK", err
			}

			return resp, "", nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting to delete LoginKey: %s", err)
	}

	return nil
}

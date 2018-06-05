package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
	"time"
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
	log.Println("[DEBUG] resourceNcloudLoginKeyRead")
	conn := meta.(*NcloudSdk).conn

	keyName := d.Get("key_name").(string)

	loginKey, err := getLoginKey(conn, keyName)
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
	log.Println("[DEBUG] resourceNcloudLoginKeyCreate")
	conn := meta.(*NcloudSdk).conn

	keyName := d.Get("key_name").(string)

	resp, err := conn.CreateLoginKey(keyName)
	if err != nil {
		logErrorResponse("CreateLoginKey", err, keyName)
		return err
	}
	logCommonResponse("CreateLoginKey", keyName, resp.CommonResponse)

	d.SetId(keyName)
	d.Set("private_key", resp.PrivateKey)

	return resourceNcloudLoginKeyRead(d, meta)
}

func resourceNcloudLoginKeyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudLoginKeyDelete")
	conn := meta.(*NcloudSdk).conn

	keyName := d.Get("key_name").(string)
	return waitForDeleteLoginKey(conn, keyName, DefaultStopTimeout)
}

func waitForDeleteLoginKey(conn *sdk.Conn, keyName string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultWaitForInterval
	}
	for {
		resp, err := conn.DeleteLoginKey(keyName)
		if err == nil || resp.ReturnCode == 200 {
			break
		}
		// resp.ReturnCode == 10407
		logCommonResponse("DeleteLoginKey", keyName, *resp)

		timeout = timeout - DefaultWaitForInterval
		if timeout <= 0 {
			return fmt.Errorf("error: Timeout: %d", timeout)
		}
		time.Sleep(DefaultWaitForInterval * time.Second)
	}
	return nil
}

func getLoginKey(conn *sdk.Conn, keyName string) (*sdk.LoginKey, error) {
	log.Println("[DEBUG] getLoginKey")
	reqParams := &sdk.RequestGetLoginKeyList{
		KeyName: keyName,
	}
	resp, err := conn.GetLoginKeyList(reqParams)
	if err != nil {
		logErrorResponse("GetLoginKeyList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetLoginKeyList", reqParams, resp.CommonResponse)

	if len(resp.LoginKeyList) > 0 {
		return &resp.LoginKeyList[0], err
	}

	return nil, err
}

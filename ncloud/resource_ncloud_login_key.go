package ncloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
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
	keyName := d.Get("key_name").(string)

	fingerprint, err := getFingerPrint(meta.(*ProviderConfig), &keyName)
	if err != nil {
		return err
	}

	if fingerprint != nil {
		d.Set("fingerprint", fingerprint)
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
	var privateKey *string
	var err error

	keyName := d.Get("key_name").(string)

	if meta.(*ProviderConfig).SupportVPC == true || meta.(*ProviderConfig).Site == "fin" {
		privateKey, err = createVpcLoginKey(meta.(*ProviderConfig), &keyName)
	} else {
		privateKey, err = createClassicLoginKey(meta.(*ProviderConfig), &keyName)
	}

	if err != nil {
		return err
	}

	d.SetId(keyName)
	d.Set("private_key", privateKey)

	time.Sleep(time.Second * 1) // for internal Master / Slave DB sync

	return resourceNcloudLoginKeyRead(d, meta)
}

func resourceNcloudLoginKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	keyName := d.Get("key_name").(string)

	if config.SupportVPC == true || config.Site == "fin" {
		if err := deleteVpcLoginKey(config.Client, &keyName); err != nil {
			return err
		}
	} else {
		if err := deleteClassicLoginKey(config.Client, &keyName); err != nil {
			return err
		}
	}

	d.SetId("")

	return nil
}

func getClassicFingerPrintList(client *NcloudAPIClient, keyName *string) ([]*string, error) {
	reqParams := &server.GetLoginKeyListRequest{}
	if keyName != nil {
		reqParams.KeyName = keyName
	}

	logCommonRequest("getClassicFingerPrintList", reqParams)
	resp, err := client.server.V2Api.GetLoginKeyList(reqParams)
	if err != nil {
		logErrorResponse("getClassicFingerPrintList", err, reqParams)
		return nil, err
	}
	logResponse("getClassicFingerPrintList", resp)

	keyList := make([]*string, 0, *resp.TotalRows)
	for _, v := range resp.LoginKeyList {
		keyList = append(keyList, v.Fingerprint)
	}

	return keyList, nil
}

func getVpcFingerPrintList(client *NcloudAPIClient, keyName *string) ([]*string, error) {
	reqParams := &vserver.GetLoginKeyListRequest{}
	if keyName != nil {
		reqParams.KeyName = keyName
	}

	logCommonRequest("getVpcFingerPrintList", reqParams)
	resp, err := client.vserver.V2Api.GetLoginKeyList(reqParams)
	if err != nil {
		logErrorResponse("getVpcFingerPrintList", err, reqParams)
		return nil, err
	}
	logResponse("getVpcFingerPrintList", resp)

	keyList := make([]*string, 0, *resp.TotalRows)
	for _, v := range resp.LoginKeyList {
		keyList = append(keyList, v.Fingerprint)
	}

	return keyList, nil
}

func getFingerPrint(config *ProviderConfig, keyName *string) (*string, error) {
	var fingerPrintList []*string
	var err error

	if config.SupportVPC == true {
		fingerPrintList, err = getVpcFingerPrintList(config.Client, keyName)
	} else {
		fingerPrintList, err = getClassicFingerPrintList(config.Client, keyName)
	}

	if err != nil {
		return nil, err
	}

	if len(fingerPrintList) > 0 {
		return fingerPrintList[0], nil
	}

	return nil, nil
}

func deleteClassicLoginKey(client *NcloudAPIClient, keyName *string) error {
	reqParams := &server.DeleteLoginKeyRequest{KeyName: keyName}

	logCommonRequest("deleteClassicLoginKey", reqParams)
	resp, err := client.server.V2Api.DeleteLoginKey(reqParams)
	if err != nil {
		logErrorResponse("deleteClassicLoginKey", err, keyName)
		return err
	}
	logCommonResponse("deleteClassicLoginKey", GetCommonResponse(resp))

	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := getClassicFingerPrintList(client, keyName)
			if err != nil {
				return 0, "", err
			}

			if len(resp) == 0 {
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

func deleteVpcLoginKey(client *NcloudAPIClient, keyName *string) error {
	reqParams := &vserver.DeleteLoginKeysRequest{KeyNameList: []*string{keyName}}

	logCommonRequest("deleteVpcLoginKey", reqParams)
	resp, err := client.vserver.V2Api.DeleteLoginKeys(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcLoginKey", err, keyName)
		return err
	}
	logCommonResponse("deleteVpcLoginKey", GetCommonResponse(resp))

	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := getVpcFingerPrintList(client, keyName)
			if err != nil {
				return 0, "", err
			}

			if len(resp) == 0 {
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

func createClassicLoginKey(config *ProviderConfig, keyName *string) (*string, error) {
	client := config.Client

	reqParams := &server.CreateLoginKeyRequest{KeyName: keyName}

	logCommonRequest("createClassicLoginKey", reqParams)
	resp, err := client.server.V2Api.CreateLoginKey(reqParams)
	if err != nil {
		logErrorResponse("createClassicLoginKey", err, keyName)
		return nil, err
	}
	logCommonResponse("createClassicLoginKey", GetCommonResponse(resp))

	return resp.PrivateKey, nil
}

func createVpcLoginKey(config *ProviderConfig, keyName *string) (*string, error) {
	client := config.Client

	reqParams := &vserver.CreateLoginKeyRequest{KeyName: keyName}

	logCommonRequest("createVpcLoginKey", reqParams)
	resp, err := client.vserver.V2Api.CreateLoginKey(reqParams)
	if err != nil {
		logErrorResponse("createVpcLoginKey", err, keyName)
		return nil, err
	}
	logCommonResponse("createVpcLoginKey", GetCommonResponse(resp))

	return resp.PrivateKey, nil
}

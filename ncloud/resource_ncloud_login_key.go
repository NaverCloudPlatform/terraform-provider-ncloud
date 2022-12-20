package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_login_key", resourceNcloudLoginKey())
}

func resourceNcloudLoginKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLoginKeyCreate,
		Read:   resourceNcloudLoginKeyRead,
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 30)),
				Description:      "Key name to generate. If the generated key name exists, an error occurs.",
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudLoginKeyRead(d *schema.ResourceData, meta interface{}) error {
	loginKey, err := getLoginKey(meta.(*ProviderConfig), d.Id())
	if err != nil {
		return err
	}

	if loginKey == nil {
		d.SetId("") // resource not found
		return nil
	}

	d.Set("key_name", loginKey.KeyName)
	d.Set("fingerprint", loginKey.Fingerprint)
	return nil
}

func resourceNcloudLoginKeyCreate(d *schema.ResourceData, meta interface{}) error {
	var privateKey *string
	var err error

	keyName := d.Get("key_name").(string)

	if meta.(*ProviderConfig).SupportVPC == true {
		privateKey, err = createVpcLoginKey(meta.(*ProviderConfig), &keyName)
	} else {
		privateKey, err = createClassicLoginKey(meta.(*ProviderConfig), &keyName)
	}

	if err != nil {
		return err
	}

	d.SetId(keyName)
	d.Set("private_key", strings.TrimSpace(*privateKey))

	time.Sleep(time.Second * 1) // for internal Master / Slave DB sync

	return resourceNcloudLoginKeyRead(d, meta)
}

func resourceNcloudLoginKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	keyName := d.Id()

	if config.SupportVPC == true {
		if err := deleteVpcLoginKey(config, keyName); err != nil {
			return err
		}
	} else {
		if err := deleteClassicLoginKey(config, keyName); err != nil {
			return err
		}
	}

	d.SetId("")

	return nil
}

func getLoginKey(config *ProviderConfig, keyName string) (*LoginKey, error) {
	if config.SupportVPC {
		return getVpcLoginKey(config, keyName)
	} else {
		return getClassicLoginKey(config, keyName)
	}
}

func getVpcLoginKey(config *ProviderConfig, keyName string) (*LoginKey, error) {
	resp, err := config.Client.vserver.V2Api.GetLoginKeyList(&vserver.GetLoginKeyListRequest{
		KeyName: ncloud.String(keyName),
	})

	if err != nil {
		return nil, err
	}

	if len(resp.LoginKeyList) < 1 {
		return nil, nil
	}

	l := resp.LoginKeyList[0]
	return &LoginKey{
		KeyName:     l.KeyName,
		Fingerprint: l.Fingerprint,
	}, nil
}

func getClassicLoginKey(config *ProviderConfig, keyName string) (*LoginKey, error) {
	resp, err := config.Client.server.V2Api.GetLoginKeyList(&server.GetLoginKeyListRequest{
		KeyName: ncloud.String(keyName),
	})

	if err != nil {
		return nil, err
	}

	if len(resp.LoginKeyList) < 1 {
		return nil, nil
	}

	l := resp.LoginKeyList[0]
	return &LoginKey{
		KeyName:     l.KeyName,
		Fingerprint: l.Fingerprint,
	}, nil
}

func deleteClassicLoginKey(config *ProviderConfig, keyName string) error {
	reqParams := &server.DeleteLoginKeyRequest{KeyName: ncloud.String(keyName)}

	logCommonRequest("deleteClassicLoginKey", reqParams)
	resp, err := config.Client.server.V2Api.DeleteLoginKey(reqParams)
	if err != nil {
		logErrorResponse("deleteClassicLoginKey", err, keyName)
		return err
	}
	logCommonResponse("deleteClassicLoginKey", GetCommonResponse(resp))

	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := getClassicLoginKey(config, keyName)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
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

func deleteVpcLoginKey(config *ProviderConfig, keyName string) error {
	reqParams := &vserver.DeleteLoginKeysRequest{KeyNameList: []*string{ncloud.String(keyName)}}

	logCommonRequest("deleteVpcLoginKey", reqParams)
	resp, err := config.Client.vserver.V2Api.DeleteLoginKeys(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcLoginKey", err, keyName)
		return err
	}
	logCommonResponse("deleteVpcLoginKey", GetCommonResponse(resp))

	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := getVpcLoginKey(config, keyName)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
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

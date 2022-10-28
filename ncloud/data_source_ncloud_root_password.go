package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_root_password", dataSourceNcloudRootPassword())
}

func dataSourceNcloudRootPassword() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRootPasswordRead,

		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"root_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceNcloudRootPasswordRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	rootPassword, err := getRootPassword(d, config)
	if err != nil {
		return err
	}

	d.SetId(d.Get("server_instance_no").(string))
	d.Set("root_password", rootPassword)

	return nil
}

func getRootPassword(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return getVpcRootPassword(d, config)
	} else {
		return getClassicRootPassword(d, config)
	}
}

func getClassicRootPassword(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	reqParams := &server.GetRootPasswordRequest{
		ServerInstanceNo: ncloud.String(d.Get("server_instance_no").(string)),
		PrivateKey:       ncloud.String(d.Get("private_key").(string)),
	}

	logCommonRequest("getClassicRootPassword", reqParams)
	resp, err := config.Client.server.V2Api.GetRootPassword(reqParams)
	if err != nil {
		logErrorResponse("getClassicRootPassword", err, reqParams)
		return nil, err
	}
	logCommonResponse("getClassicRootPassword", GetCommonResponse(resp))

	return resp.RootPassword, nil
}

func getVpcRootPassword(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	reqParams := &vserver.GetRootPasswordRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: ncloud.String(d.Get("server_instance_no").(string)),
		PrivateKey:       ncloud.String(d.Get("private_key").(string)),
	}

	logCommonRequest("getVpcRootPassword", reqParams)
	resp, err := config.Client.vserver.V2Api.GetRootPassword(reqParams)
	if err != nil {
		logErrorResponse("getVpcRootPassword", err, reqParams)
		return nil, err
	}
	logCommonResponse("getVpcRootPassword", GetCommonResponse(resp))

	return resp.RootPassword, nil
}

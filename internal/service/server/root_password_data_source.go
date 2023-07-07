package server

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudRootPassword() *schema.Resource {
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
	config := meta.(*conn.ProviderConfig)

	rootPassword, err := getRootPassword(d, config)
	if err != nil {
		return err
	}

	d.SetId(d.Get("server_instance_no").(string))
	d.Set("root_password", rootPassword)

	return nil
}

func getRootPassword(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return getVpcRootPassword(d, config)
	} else {
		return getClassicRootPassword(d, config)
	}
}

func getClassicRootPassword(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	reqParams := &server.GetRootPasswordRequest{
		ServerInstanceNo: ncloud.String(d.Get("server_instance_no").(string)),
		PrivateKey:       ncloud.String(d.Get("private_key").(string)),
	}

	LogCommonRequest("getClassicRootPassword", reqParams)
	resp, err := config.Client.Server.V2Api.GetRootPassword(reqParams)
	if err != nil {
		LogErrorResponse("getClassicRootPassword", err, reqParams)
		return nil, err
	}
	LogCommonResponse("getClassicRootPassword", GetCommonResponse(resp))

	return resp.RootPassword, nil
}

func getVpcRootPassword(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	reqParams := &vserver.GetRootPasswordRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: ncloud.String(d.Get("server_instance_no").(string)),
		PrivateKey:       ncloud.String(d.Get("private_key").(string)),
	}

	LogCommonRequest("getVpcRootPassword", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetRootPassword(reqParams)
	if err != nil {
		LogErrorResponse("getVpcRootPassword", err, reqParams)
		return nil, err
	}
	LogCommonResponse("getVpcRootPassword", GetCommonResponse(resp))

	return resp.RootPassword, nil
}

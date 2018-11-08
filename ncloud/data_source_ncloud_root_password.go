package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudRootPassword() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRootPasswordRead,

		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server instance number",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Server’s login key (auth key)",
			},
			"root_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Password of a root account",
			},
		},
	}
}

func dataSourceNcloudRootPasswordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	serverInstanceNo := d.Get("server_instance_no").(string)
	privateKey := d.Get("private_key").(string)
	reqParams := &server.GetRootPasswordRequest{
		ServerInstanceNo: ncloud.String(serverInstanceNo),
		PrivateKey:       ncloud.String(privateKey),
	}

	logCommonRequest("GetRootPassword", reqParams)
	resp, err := client.server.V2Api.GetRootPassword(reqParams)
	if err != nil {
		logErrorResponse("GetRootPassword", err, reqParams)
		return err
	}
	logCommonResponse("GetRootPassword", GetCommonResponse(resp))

	d.SetId(serverInstanceNo)
	d.Set("root_password", resp.RootPassword)

	return nil
}

package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

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
	conn := meta.(*NcloudSdk).conn

	serverInstanceNo := d.Get("server_instance_no").(string)
	privateKey := d.Get("private_key").(string)
	reqParams := &sdk.RequestGetRootPassword{
		ServerInstanceNo: serverInstanceNo,
		PrivateKey:       privateKey,
	}
	resp, err := conn.GetRootPassword(reqParams)
	if err != nil {
		logErrorResponse("GetRootPassword", err, reqParams)
		return err
	}
	logCommonResponse("GetRootPassword", reqParams, resp.CommonResponse)

	d.SetId(serverInstanceNo)
	d.Set("root_password", resp.RootPassword)

	return nil
}

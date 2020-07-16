package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudVpcsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceNcloudVpcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.GetVpcListRequest{
		RegionCode: regionCode,
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		reqParams.VpcStatusCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_no_list"); ok {
		reqParams.VpcNoList = expandStringInterfaceList(v.([]interface{}))
	}

	logCommonRequest("GetVpcList", reqParams)
	resp, err := client.vpc.V2Api.GetVpcList(reqParams)

	if err != nil {
		logErrorResponse("GetVpcList", err, reqParams)
		return err
	}
	logResponse("GetVpcList", resp)

	if resp == nil || len(resp.VpcList) == 0 {
		return fmt.Errorf("no matching VPC found")
	}

	vpcs := make([]string, 0)

	for _, vpc := range resp.VpcList {
		vpcs = append(vpcs, ncloud.StringValue(vpc.VpcNo))
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("ids", vpcs); err != nil {
		return fmt.Errorf("Error setting vpc ids: %s", err)
	}

	return nil
}

package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_network_acls", dataSourceNcloudNetworkAcls())
}

func dataSourceNcloudNetworkAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNetworkAclsRead,

		Schema: map[string]*schema.Schema{
			"network_acl_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of Network ACL ID to retrieve.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the field to filter by.",
			},
			"vpc_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC ID that you want to filter from.",
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

func dataSourceNcloudNetworkAclsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_network_acls`")
	}

	reqParams := &vpc.GetNetworkAclListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("network_acl_no_list"); ok {
		reqParams.NetworkAclNoList = expandStringInterfaceList(v.([]interface{}))
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NetworkAclName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNo = ncloud.String(v.(string))
	}

	logCommonRequest("GetNetworkAclList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNetworkAclList(reqParams)

	if err != nil {
		logErrorResponse("GetNetworkAclList", err, reqParams)
		return err
	}
	logResponse("GetNetworkAclList", resp)

	if resp == nil || len(resp.NetworkAclList) == 0 {
		return fmt.Errorf("no matching Network ACL found")
	}

	instances := make([]string, 0)

	for _, vpc := range resp.NetworkAclList {
		instances = append(instances, ncloud.StringValue(vpc.VpcNo))
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("ids", instances); err != nil {
		return fmt.Errorf("Error setting Network ACL ids: %s", err)
	}

	return nil
}

package vpc

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudNetworkAcls() *schema.Resource {
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
			"filter": DataSourceFiltersSchema(),

			"network_acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(ResourceNcloudNetworkACL()),
			},
		},
	}
}

func dataSourceNcloudNetworkAclsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_network_acls`")
	}

	reqParams := &vpc.GetNetworkAclListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("network_acl_no_list"); ok {
		// reqParams.NetworkAclNoList = []*string{ncloud.String(v.(string))}
		list := make([]*string, 0, len(v.([]interface{})))
		for _, v := range v.([]interface{}) {
			list = append(list, ncloud.String(v.(string)))
		}
		reqParams.NetworkAclNoList = list
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NetworkAclName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNo = ncloud.String(v.(string))
	}

	LogCommonRequest("GetNetworkAclList", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetNetworkAclList(reqParams)

	if err != nil {
		LogErrorResponse("GetNetworkAclList", err, reqParams)
		return err
	}
	LogResponse("GetNetworkAclList", resp)

	if resp == nil || len(resp.NetworkAclList) == 0 {
		return fmt.Errorf("no matching Network ACL found")
	}

	var resources []map[string]interface{}

	for _, r := range resp.NetworkAclList {
		instance := map[string]interface{}{
			"id":             *r.NetworkAclNo,
			"network_acl_no": *r.NetworkAclNo,
			"name":           *r.NetworkAclName,
			"description":    *r.NetworkAclDescription,
			"vpc_no":         *r.VpcNo,
			"is_default":     *r.IsDefault,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, ResourceNcloudNetworkACL().Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("network_acls", resources); err != nil {
		return fmt.Errorf("Error setting Network ACLs: %s", err)
	}

	return nil
}

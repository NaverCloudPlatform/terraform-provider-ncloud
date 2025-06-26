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

func DataSourceNcloudNetworkACLDenyAllowGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNetworkACLDenyAllowGroupsRead,

		Schema: map[string]*schema.Schema{
			"network_acl_deny_allow_group_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": DataSourceFiltersSchema(),
			"network_acl_deny_allow_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(ResourceNcloudNetworkACLDenyAllowGroup()),
			},
		},
	}
}

func dataSourceNcloudNetworkACLDenyAllowGroupsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	reqParams := &vpc.GetNetworkAclDenyAllowGroupListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("network_acl_deny_allow_group_no_list"); ok {
		reqParams.NetworkAclDenyAllowGroupNoList = ExpandStringInterfaceList(v.([]interface{}))
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NetworkAclDenyAllowGroupName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNo = ncloud.String(v.(string))
	}

	LogCommonRequest("GetNetworkAclDenyAllowGroupList", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetNetworkAclDenyAllowGroupList(reqParams)

	if err != nil {
		LogErrorResponse("GetNetworkAclDenyAllowGroupList", err, reqParams)
		return err
	}
	LogResponse("GetNetworkAclDenyAllowGroupList", resp)

	if resp == nil || len(resp.NetworkAclDenyAllowGroupList) == 0 {
		return fmt.Errorf("no matching NetworkAclDenyAllowGroup found")
	}

	var resources []map[string]interface{}

	for _, r := range resp.NetworkAclDenyAllowGroupList {
		m := map[string]interface{}{
			"id":                              *r.NetworkAclDenyAllowGroupNo,
			"network_acl_deny_allow_group_no": *r.NetworkAclDenyAllowGroupNo,
			"vpc_no":                          *r.VpcNo,
			"name":                            *r.NetworkAclDenyAllowGroupName,
			"description":                     *r.NetworkAclDenyAllowGroupDescription,
		}

		// only can get `ip_list` data from `getNetworkAclDenyAllowGroupDetail`
		if g, err := GetNetworkAclDenyAllowGroupDetail(config, *r.NetworkAclDenyAllowGroupNo); err != nil {
			return err
		} else {
			m["ip_list"] = g.IpList
		}

		resources = append(resources, m)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, ResourceNcloudNetworkACLDenyAllowGroup().Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("network_acl_deny_allow_groups", resources); err != nil {
		return fmt.Errorf("error setting NetworkAclDenyAllowGroups: %s", err)
	}

	return nil
}

package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_init_script", dataSourceNcloudInitScript())
}

func dataSourceNcloudInitScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudInitScriptRead,
		Schema: map[string]*schema.Schema{
			"init_script_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"LNX", "WND"}, false),
			},
			"filter": dataSourceFiltersSchema(),

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudInitScriptRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcInitScriptListFiltered(d, config)
	} else {
		return NotSupportClassic("data source `ncloud_init_script`")
	}

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getVpcInitScriptListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vserver.GetInitScriptListRequest{
		RegionCode:     &config.RegionCode,
		OsTypeCode:     StringPtrOrNil(d.GetOk("os_type")),
		InitScriptName: StringPtrOrNil(d.GetOk("name")),
	}

	if v, ok := d.GetOk("init_script_no"); ok {
		reqParams.InitScriptNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getVpcInitScriptList", reqParams)
	resp, err := config.Client.vserver.V2Api.GetInitScriptList(reqParams)

	if err != nil {
		logErrorResponse("getVpcInitScriptList", err, reqParams)
		return nil, err
	}
	logResponse("getVpcInitScriptList", resp)

	var resources []map[string]interface{}

	for _, r := range resp.InitScriptList {
		instance := map[string]interface{}{
			"id":             *r.InitScriptNo,
			"init_script_no": *r.InitScriptNo,
			"name":           *r.InitScriptName,
			"description":    *r.InitScriptDescription,
			"os_type":        *r.OsType,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudInitScript().Schema)
	}

	return resources, nil
}

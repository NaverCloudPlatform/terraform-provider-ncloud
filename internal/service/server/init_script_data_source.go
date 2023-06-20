package server

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func init() {
	provider.RegisterDataSource("ncloud_init_script", dataSourceNcloudInitScript())
}

func dataSourceNcloudInitScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudInitScriptRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"init_script_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"LNX", "WND"}, false)),
			},
			"filter": DataSourceFiltersSchema(),
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudInitScriptRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*provider.ProviderConfig)
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

	if err := ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getVpcInitScriptListFiltered(d *schema.ResourceData, config *provider.ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vserver.GetInitScriptListRequest{
		RegionCode:     &config.RegionCode,
		OsTypeCode:     StringPtrOrNil(d.GetOk("os_type")),
		InitScriptName: StringPtrOrNil(d.GetOk("name")),
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.InitScriptNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getVpcInitScriptList", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetInitScriptList(reqParams)

	if err != nil {
		LogErrorResponse("getVpcInitScriptList", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcInitScriptList", resp)

	var resources []map[string]interface{}

	for _, r := range resp.InitScriptList {
		instance := map[string]interface{}{
			"id":             *r.InitScriptNo,
			"init_script_no": *r.InitScriptNo,
			"name":           *r.InitScriptName,
			"description":    *r.InitScriptDescription,
			"os_type":        *r.OsType.Code,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudInitScript().Schema)
	}

	return resources, nil
}

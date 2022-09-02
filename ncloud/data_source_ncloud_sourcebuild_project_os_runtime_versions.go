package ncloud

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcebuild_project_os_runtime_versions", dataSourceNcloudSourceBuildRuntimeVersions())
}

func dataSourceNcloudSourceBuildRuntimeVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildRuntimeVersionsRead,
		Schema: map[string]*schema.Schema{
			"os_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"runtime_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": dataSourceFiltersSchema(),
			"runtime_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSourceBuildRuntimeVersionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	osIdParam := Int32PtrOrNil(d.GetOk("os_id"))
	osId := ncloud.IntString(int(ncloud.Int32Value(osIdParam)))
	runtimeIdParam := Int32PtrOrNil(d.GetOk("runtime_id"))
	runtimeId := ncloud.IntString(int(ncloud.Int32Value(runtimeIdParam)))

	logCommonRequest("GetRuntimeVersionEnv", "")
	resp, err := config.Client.sourcebuild.V1Api.GetRuntimeVersionEnv(ctx, osId, runtimeId)
	if err != nil {
		logErrorResponse("GetRuntimeVersionEnv", err, "")
		return diag.FromErr(err)
	}
	logResponse("GetRuntimeVersionEnv", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Version {
		runtime := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}

		resources = append(resources, runtime)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceBuildRuntimeVersions().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("runtime_versions", resources)

	return nil
}

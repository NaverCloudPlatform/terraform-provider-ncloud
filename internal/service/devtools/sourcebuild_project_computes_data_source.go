package devtools

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func init() {
	RegisterDataSource("ncloud_sourcebuild_project_computes", dataSourceNcloudSourceBuildComputes())
}

func dataSourceNcloudSourceBuildComputes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildComputesRead,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"computes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSourceBuildComputesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	LogCommonRequest("GetComputeEnv", "")
	resp, err := config.Client.Sourcebuild.V1Api.GetComputeEnv(ctx)
	if err != nil {
		LogErrorResponse("GetComputeEnv", err, "")
		return diag.FromErr(err)
	}
	LogResponse("GetComputeEnv", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Compute {
		compute := map[string]interface{}{
			"id":  *r.Id,
			"cpu": *r.Cpu,
			"mem": *r.Mem,
		}

		resources = append(resources, compute)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceBuildComputes().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("computes", resources)

	return nil
}

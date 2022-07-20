package ncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcebuild_compute", dataSourceNcloudSourceBuildCompute())
}

func dataSourceNcloudSourceBuildCompute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildComputeRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": dataSourceFiltersSchema(),
			"compute": {
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

func dataSourceNcloudSourceBuildComputeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	logCommonRequest("GetComputeEnv", "")
	resp, err := config.Client.sourcebuild.V1Api.GetComputeEnv(ctx)
	if err != nil {
		logErrorResponse("GetComputeEnv", err, "")
		return diag.FromErr(err)
	}
	logResponse("GetComputeEnv", resp)

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
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceBuildCompute().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("compute", resources)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return diag.FromErr(writeToFile(output.(string), resources))
	}

	return nil
}

package ncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcebuild_project_os", dataSourceNcloudSourceBuildOs())
}

func dataSourceNcloudSourceBuildOs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildOsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"os": {
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
						"archi": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSourceBuildOsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	logCommonRequest("GetOsEnv", "")
	resp, err := config.Client.sourcebuild.V1Api.GetOsEnv(ctx)
	if err != nil {
		logErrorResponse("GetOsEnv", err, "")
		return diag.FromErr(err)
	}
	logResponse("GetOsEnv", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Os {
		os := map[string]interface{}{
			"id":      *r.Id,
			"name":    *r.Name,
			"archi":   *r.Archi,
			"version": *r.Version,
		}

		resources = append(resources, os)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceBuildOs().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("os", resources)

	return nil
}

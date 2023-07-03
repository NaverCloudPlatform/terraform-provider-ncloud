package devtools

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourceBuildOs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildOsRead,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
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
	config := meta.(*conn.ProviderConfig)

	LogCommonRequest("GetOsEnv", "")
	resp, err := config.Client.Sourcebuild.V1Api.GetOsEnv(ctx)
	if err != nil {
		LogErrorResponse("GetOsEnv", err, "")
		return diag.FromErr(err)
	}
	LogResponse("GetOsEnv", resp)

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
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceBuildOs().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("os", resources)

	return nil
}

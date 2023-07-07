package devtools

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourceBuildDockerEngines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildDockerEnginesRead,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"docker_engines": {
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

func dataSourceNcloudSourceBuildDockerEnginesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	LogCommonRequest("GetDockerEnv", "")
	resp, err := config.Client.Sourcebuild.V1Api.GetDockerEnv(context.Background())
	if err != nil {
		LogErrorResponse("GetDockerEnv", err, "")
		return diag.FromErr(err)
	}
	LogResponse("GetDockerEnv", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Docker {
		docker := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}

		resources = append(resources, docker)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceBuildDockerEngines().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("docker_engines", resources)

	return nil
}

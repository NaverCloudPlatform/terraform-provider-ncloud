package devtools

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourceDeployProjectsContext() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceDeployProjectsReadContext,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"projects": {
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

func dataSourceNcloudSourceDeployProjectsReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	reqParams := make(map[string]interface{})

	reqParams["projectName"] = ncloud.StringValue(StringPtrOrNil(d.GetOk("name")))
	resp, err := config.Client.Vsourcedeploy.V1Api.GetProjects(ctx, reqParams)

	if err != nil {
		return diag.FromErr(err)
	}
	LogResponse("GetProjects", resp)

	resources := []map[string]interface{}{}
	for _, r := range resp.ProjectList {
		project := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}

		resources = append(resources, project)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceDeployProjectsContext().Schema)
	}
	d.SetId(config.RegionCode)
	d.Set("projects", resources)

	return nil
}

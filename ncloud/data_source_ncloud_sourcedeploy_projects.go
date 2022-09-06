package ncloud

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcedeploy_projects", dataSourceNcloudSourceDeployProjectsContext())
}

func dataSourceNcloudSourceDeployProjectsContext() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceDeployProjectsReadContext,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
	config := meta.(*ProviderConfig)
	reqParams := make(map[string]interface{})

	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_sourcedeploy_projects_context`"))
	}

	reqParams["projectName"] = ncloud.StringValue(StringPtrOrNil(d.GetOk("name")))
	resp, err := config.Client.vsourcedeploy.V1Api.GetProjects(ctx, reqParams)

	if err != nil {
		return diag.FromErr(err)
	}
	logResponse("GetProjects", resp)

	resources := []map[string]interface{}{}
	for _, r := range resp.ProjectList {
		project := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}

		resources = append(resources, project)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceDeployProjectsContext().Schema)
	}
	d.SetId(config.RegionCode)
	d.Set("projects", resources)

	return nil
}

package ncloud

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcebuild_projects", dataSourceNcloudSourceBuildProjects())
}

func dataSourceNcloudSourceBuildProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildProjectsRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
						"project_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSourceBuildProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	reqParams := make(map[string]interface{})
	reqParams["projectName"] = ncloud.StringValue(StringPtrOrNil(d.GetOk("project_name")))

	logCommonRequest("GetSourceBuildProjects", reqParams)
	resp, err := config.Client.sourcebuild.V1Api.GetProjects(ctx, reqParams)
	if err != nil {
		logErrorResponse("GetSourceBuildProjects", err, reqParams)
		return diag.FromErr(err)
	}
	logResponse("GetSourceBuildProjects", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Project {
		project := map[string]interface{}{
			"id":          *r.Id,
			"project_no":  *r.Id,
			"name":        *r.Name,
			"permission":  *r.Permission,
			"action_name": *r.ActionName,
		}

		resources = append(resources, project)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceBuildProjects().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("projects", resources)

	return nil
}

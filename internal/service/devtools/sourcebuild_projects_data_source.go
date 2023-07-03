package devtools

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourceBuildProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildProjectsRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
	config := meta.(*conn.ProviderConfig)

	reqParams := make(map[string]interface{})
	reqParams["projectName"] = ncloud.StringValue(StringPtrOrNil(d.GetOk("project_name")))

	LogCommonRequest("GetSourceBuildProjects", reqParams)
	resp, err := config.Client.Sourcebuild.V1Api.GetProjects(ctx, reqParams)
	if err != nil {
		LogErrorResponse("GetSourceBuildProjects", err, reqParams)
		return diag.FromErr(err)
	}
	LogResponse("GetSourceBuildProjects", resp)

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
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceBuildProjects().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("projects", resources)

	return nil
}

package devtools

import (
	"context"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcepipeline"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourcePipelineProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourcePipelineProjectsRead,
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

func dataSourceNcloudSourcePipelineProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	projects, err := getSourcePipelineProjects(ctx, config)
	if err != nil {
		LogErrorResponse("getSourcePipelineProjects", err, projects)
		return diag.FromErr(err)
	}
	LogResponse("getSourcePipelineProjects", projects)

	if projects == nil {
		d.SetId("")
		return nil
	}

	var resources []map[string]interface{}
	for _, project := range projects {
		mapping := map[string]interface{}{
			"id":   ncloud.Int32Value(project.Id),
			"name": ncloud.StringValue(project.Name),
		}
		resources = append(resources, mapping)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourcePipelineProjects().Schema)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("projects", resources)

	return nil
}

func getSourcePipelineProjects(ctx context.Context, config *conn.ProviderConfig) ([]*PipelineProjects, error) {
	resp, err := config.Client.Vsourcepipeline.V1Api.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	return convertVpcPipelineProjects(resp), nil
}

func convertVpcPipelineProjects(r *vsourcepipeline.GetProjectListResponse) []*PipelineProjects {
	projects := []*PipelineProjects{}

	for _, project := range r.ProjectList {
		pi := &PipelineProjects{
			Id:   project.Id,
			Name: project.Name,
		}

		projects = append(projects, pi)
	}

	return projects
}

type PipelineProjects struct {
	Id *int32 `json:"id,omitempty"`

	Name *string `json:"name,omitempty"`
}

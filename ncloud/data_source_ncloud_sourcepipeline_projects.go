package ncloud

import (
	"context"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcepipeline"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcepipeline"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcepipeline_projects", dataSourceNcloudSourcePipelineProjects())
}

func dataSourceNcloudSourcePipelineProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourcePipelineProjectsRead,
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

func dataSourceNcloudSourcePipelineProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	projects, err := getSourcePipelineProjects(ctx, config)
	if err != nil {
		logErrorResponse("getSourcePipelineProjects", err, projects)
		return diag.FromErr(err)
	}
	logResponse("getSourcePipelineProjects", projects)

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
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourcePipelineProjects().Schema)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("projects", resources)

	return nil
}

func getSourcePipelineProjects(ctx context.Context, config *ProviderConfig) ([]*PipelineProjects, error) {
	if config.SupportVPC {
		return getVpcSourcePipelineProjects(ctx, config)
	}
	return getClassicSourcePipelineProjects(ctx, config)
}

func getClassicSourcePipelineProjects(ctx context.Context, config *ProviderConfig) ([]*PipelineProjects, error) {
	resp, err := config.Client.sourcepipeline.V1Api.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	return convertClassicPipelineProjects(resp), nil
}

func getVpcSourcePipelineProjects(ctx context.Context, config *ProviderConfig) ([]*PipelineProjects, error) {
	resp, err := config.Client.vsourcepipeline.V1Api.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	return convertVpcPipelineProjects(resp), nil
}

func convertClassicPipelineProjects(r *sourcepipeline.GetProjectListResponse) []*PipelineProjects {
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

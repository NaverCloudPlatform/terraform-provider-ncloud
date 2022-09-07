package ncloud

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcepipeline_project", dataSourceNcloudSourcePipelineProject())
}

func dataSourceNcloudSourcePipelineProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourcePipelineProjectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"stage_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"scenario_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"target": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"repository_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"repository_branch": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"project_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"file": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"manifest": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"full_manifest": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"linked_tasks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"triggers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sourcecommit": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repository_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"branch": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSourcePipelineProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(*ncloud.IntString(d.Get("id").(int)))

	return resourceNcloudSourcePipelineRead(ctx, d, meta)
}

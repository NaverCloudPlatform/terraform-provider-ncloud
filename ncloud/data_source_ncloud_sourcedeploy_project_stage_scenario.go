package ncloud

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcedeploy_project_stage_scenario", dataSourceNcloudSourceDeployScenarioContext())
}

func dataSourceNcloudSourceDeployScenarioContext() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceDeployScenarioReadContext,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"stage_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
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
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"object_storage": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"object": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"source_build": {
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
							},
						},
						"rollback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deploy_command": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pre_deploy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"command": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"path": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"deploy_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"post_deploy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"command": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"load_balancer": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_target_group_no": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"load_balancer_target_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"delete_server": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"manifest": {
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
									"branch": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"path": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"canary_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"canary_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"analysis_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"prometheus": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"env": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"baseline": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"canary": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"metrics": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"success_criteria": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"query_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"weight": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"metric": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"filter": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"query": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"analysis_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"duration": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"delay": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"interval": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"step": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"pass_score": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"path": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_path": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"deploy_path": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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

func dataSourceNcloudSourceDeployScenarioReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(*ncloud.IntString(d.Get("id").(int)))
	return resourceNcloudSourceDeployScenarioRead(ctx, d, meta)
}

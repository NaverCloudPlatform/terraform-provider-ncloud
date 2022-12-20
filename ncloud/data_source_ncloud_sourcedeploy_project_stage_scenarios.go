package ncloud

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcedeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcedeploy_project_stage_scenarios", dataSourceNcloudSourceDeployscenariosContext())
}

func dataSourceNcloudSourceDeployscenariosContext() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceDeployScenariosReadContext,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"stage_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": dataSourceFiltersSchema(),
			"scenarios": {
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

func dataSourceNcloudSourceDeployScenariosReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_sourcedeploy_project_stage_scenarios`"))
	}

	projectId := ncloud.IntString(d.Get("project_id").(int))
	stageId := ncloud.IntString(d.Get("stage_id").(int))
	resp, err := GetScenarios(ctx, config, projectId, stageId)
	if err != nil {
		return diag.FromErr(err)
	}
	logResponse("GetScenarios", resp)

	resources := []map[string]interface{}{}
	for _, r := range resp.ScenarioList {
		project := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}

		resources = append(resources, project)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceDeployscenariosContext().Schema)
	}
	d.SetId(config.RegionCode)
	d.Set("scenarios", resources)

	return nil
}

func GetScenarios(ctx context.Context, config *ProviderConfig, projectId *string, stageId *string) (*vsourcedeploy.GetScenarioListResponse, error) {

	reqParams := make(map[string]interface{})
	logCommonRequest("GetScenarios", reqParams)
	resp, err := config.Client.vsourcedeploy.V1Api.GetScenarioes(ctx, projectId, stageId, reqParams)

	if err != nil {
		logErrorResponse("GetScenarios", err, "")
		return nil, err
	}
	logResponse("GetScenarios", resp)

	return resp, nil
}

package devtools

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcedeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourceDeployscenariosContext() *schema.Resource {
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
			"filter": DataSourceFiltersSchema(),
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
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_sourcedeploy_project_stage_scenarios`"))
	}

	projectId := ncloud.IntString(d.Get("project_id").(int))
	stageId := ncloud.IntString(d.Get("stage_id").(int))
	resp, err := GetScenarios(ctx, config, projectId, stageId)
	if err != nil {
		return diag.FromErr(err)
	}
	LogResponse("GetScenarios", resp)

	resources := []map[string]interface{}{}
	for _, r := range resp.ScenarioList {
		project := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}

		resources = append(resources, project)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceDeployscenariosContext().Schema)
	}
	d.SetId(config.RegionCode)
	d.Set("scenarios", resources)

	return nil
}

func GetScenarios(ctx context.Context, config *conn.ProviderConfig, projectId *string, stageId *string) (*vsourcedeploy.GetScenarioListResponse, error) {

	reqParams := make(map[string]interface{})
	LogCommonRequest("GetScenarios", reqParams)
	resp, err := config.Client.Vsourcedeploy.V1Api.GetScenarioes(ctx, projectId, stageId, reqParams)

	if err != nil {
		LogErrorResponse("GetScenarios", err, "")
		return nil, err
	}
	LogResponse("GetScenarios", resp)

	return resp, nil
}

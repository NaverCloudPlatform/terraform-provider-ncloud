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

func DataSourceNcloudSourceDeployStagesContext() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceDeployStagesReadContext,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": DataSourceFiltersSchema(),
			"stages": {
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

func dataSourceNcloudSourceDeployStagesReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	projectId := ncloud.IntString(d.Get("project_id").(int))
	resp, err := GetStages(ctx, config, projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	resources := []map[string]interface{}{}
	for _, r := range resp.StageList {
		stage := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}
		resources = append(resources, stage)
	}
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceDeployStagesContext().Schema)
	}
	d.SetId(config.RegionCode)
	d.Set("stages", resources)

	return nil
}

func GetStages(ctx context.Context, config *conn.ProviderConfig, projectId *string) (*vsourcedeploy.GetStageListResponse, error) {

	reqParams := make(map[string]interface{})
	LogCommonRequest("getStages", reqParams)
	resp, err := config.Client.Vsourcedeploy.V1Api.GetStages(ctx, projectId, reqParams)

	if err != nil {
		LogErrorResponse("getStages", err, "")
		return nil, err
	}
	LogResponse("getStages", resp)

	return resp, nil
}

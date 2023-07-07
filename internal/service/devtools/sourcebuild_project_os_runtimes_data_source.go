package devtools

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourceBuildRuntimes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourceBuildRuntimesRead,
		Schema: map[string]*schema.Schema{
			"os_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": DataSourceFiltersSchema(),
			"runtimes": {
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

func dataSourceNcloudSourceBuildRuntimesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	osIdParam := Int32PtrOrNil(d.GetOk("os_id"))
	osId := ncloud.IntString(int(ncloud.Int32Value(osIdParam)))

	LogCommonRequest("GetRuntimeEnv", "")
	resp, err := config.Client.Sourcebuild.V1Api.GetRuntimeEnv(context.Background(), osId)
	if err != nil {
		LogErrorResponse("GetRuntimeEnv", err, "")
		return diag.FromErr(err)
	}
	LogResponse("GetRuntimeEnv", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Runtime {
		runtime := map[string]interface{}{
			"id":   *r.Id,
			"name": *r.Name,
		}

		resources = append(resources, runtime)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceBuildRuntimes().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("runtimes", resources)

	return nil
}

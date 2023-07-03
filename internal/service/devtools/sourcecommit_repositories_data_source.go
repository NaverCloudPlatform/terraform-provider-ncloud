package devtools

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourceCommitRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSousrceNcloudSourceCommitRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": DataSourceFiltersSchema(),
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"repository_no": {
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

func dataSousrceNcloudSourceCommitRepositoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*conn.ProviderConfig)

	LogCommonRequest("GetSourceCommitRepositories", "")
	resp, err := GetRepositories(ctx, config)
	if err != nil {
		LogErrorResponse("GetSourceCommitRepositories", err, "")
		return diag.FromErr(err)
	}
	LogResponse("GetSourceCommitRepositories", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Repository {
		repo := map[string]interface{}{
			"id":            *r.Id,
			"repository_no": *r.Id,
			"name":          *r.Name,
			"permission":    *r.Permission,
			"action_name":   *r.ActionName,
		}

		resources = append(resources, repo)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSourceCommitRepository().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("repositories", resources)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return diag.FromErr(WriteToFile(output.(string), resources))
	}

	return nil
}

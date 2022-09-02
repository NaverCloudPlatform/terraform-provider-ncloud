package ncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcecommit_repositories", dataSourceNcloudSourceCommitRepositories())
}

func dataSourceNcloudSourceCommitRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSousrceNcloudSourceCommitRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": dataSourceFiltersSchema(),
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

	config := meta.(*ProviderConfig)

	logCommonRequest("GetSourceCommitRepositories", "")
	resp, err := getRepositories(ctx, config)
	if err != nil {
		logErrorResponse("GetSourceCommitRepositories", err, "")
		return diag.FromErr(err)
	}
	logResponse("GetSourceCommitRepositories", resp)

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
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSourceCommitRepository().Schema)
	}

	d.SetId(config.RegionCode)
	d.Set("repositories", resources)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return diag.FromErr(writeToFile(output.(string), resources))
	}

	return nil
}

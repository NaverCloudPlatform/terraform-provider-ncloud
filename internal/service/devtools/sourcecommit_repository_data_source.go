package devtools

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func init() {
	RegisterDataSource("ncloud_sourcecommit_repository", dataSourceNcloudSourceCommitRepository())
}

func dataSourceNcloudSourceCommitRepository() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSousrceNcloudSourceCommitRepositoryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_https_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_ssh_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_safer": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSousrceNcloudSourceCommitRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*ProviderConfig)

	name := d.Get("name").(string)

	LogCommonRequest("GetSourceCommitRepository", "")
	repository, err := getRepository(ctx, config, name)

	var diags diag.Diagnostics

	if err != nil {
		LogErrorResponse("GetSourceCommitRepository", err, "")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to search repository",
			Detail:   "Unable to search repository - detail",
		})
		return diags
	}

	if repository == nil {
		LogErrorResponse("GetSourceCommitRepository", err, "")
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "there is no such repository",
			Detail:   "there is no such repository - detail",
		})
		return diags
	}

	LogResponse("GetSourceCommitRepository", repository)
	d.SetId(strconv.Itoa(*repository.Id))
	d.Set("repository_no", strconv.Itoa(*repository.Id))
	d.Set("name", repository.Name)
	d.Set("description", repository.Description)
	d.Set("creator", repository.Created.User)
	d.Set("git_https_url", repository.Git.Https)
	d.Set("git_ssh_url", repository.Git.Ssh)
	d.Set("file_safer", repository.Linked.FileSafer)

	return nil
}

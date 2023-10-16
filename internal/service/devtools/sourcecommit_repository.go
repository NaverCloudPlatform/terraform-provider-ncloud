package devtools

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcecommit"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudSourceCommitRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudSourceCommitRepositoryCreate,
		ReadContext:   resourceNcloudSourceCommitRepositoryRead,
		UpdateContext: resourceNcloudSourceCommitRepositoryUpdate,
		DeleteContext: resourceNcloudSourceCommitRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultTimeout),
			Read:   schema.DefaultTimeout(conn.DefaultTimeout),
			Update: schema.DefaultTimeout(conn.DefaultTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 100)),
			},
			"repository_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 500)),
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
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNcloudSourceCommitRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	reqParams := &sourcecommit.CreateRepository{
		Name:        ncloud.String(d.Get("name").(string)),
		Description: StringPtrOrNil(d.GetOk("description")),
	}

	if fileSafer, ok := d.GetOk("file_safer"); ok {
		reqParams.Linked = &sourcecommit.CreateRepositoryLinked{
			FileSafer: BoolPtrOrNil(fileSafer, ok),
		}
	}

	LogCommonRequest("resourceNcloudSourceCommitRepositoryCreate", reqParams)
	resp, err := config.Client.Sourcecommit.V1Api.CreateRepository(ctx, reqParams)
	LogCommonResponse("resourceNcloudSourceCommitRepositoryCreate", GetCommonResponse(nil))
	var diags diag.Diagnostics

	if err != nil {
		LogErrorResponse("resourceNcloudSourceCommitRepositoryCreate", err, reqParams)

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Fail to create repository",
			Detail:   err.Error(),
		})
		return diags
	}

	LogResponse("resourceNcloudSourceCommitRepositoryCreate", resp)

	name := ncloud.StringValue(reqParams.Name)

	if err := waitForSourceCommitRepositoryActive(ctx, d, config, name); err != nil {

		diags := append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to search repository",
			Detail:   fmt.Sprintf("Unable to search repository - detail , name : (%s)", name),
		})
		return diags
	}

	return resourceNcloudSourceCommitRepositoryRead(ctx, d, meta)
}

func resourceNcloudSourceCommitRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	name := ncloud.String(d.Get("name").(string))
	id := ncloud.String(d.Id())

	repository, err := GetRepositoryById(ctx, config, *id)

	LogCommonRequest("resourceNcloudSourceCommitRepositoryRead", name)
	var diags diag.Diagnostics

	if err != nil {
		LogErrorResponse("resourceNcloudSourceCommitRepositoryRead", err, *name)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to search repository",
			Detail:   fmt.Sprintf("Unable to search repository - detail repository : %s", *name),
		})
		return diags
	}

	LogResponse("resourceNcloudSourceCommitRepositoryRead", repository)

	if repository == nil {
		d.SetId("")
		return nil
	}

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

func resourceNcloudSourceCommitRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	if d.HasChanges("description", "file_safer") {

		reqParams := &sourcecommit.ChangeRepository{
			Description: ncloud.String(d.Get("description").(string)),
		}

		reqParams.Linked = &sourcecommit.CreateRepositoryLinked{
			FileSafer: ncloud.Bool(d.Get("file_safer").(bool)),
		}

		id := ncloud.String(d.Id())

		LogCommonRequest("resourceNcloudSourceCommitRepositoryUpdate", reqParams)
		_, err := config.Client.Sourcecommit.V1Api.ChangeRepository(ctx, reqParams, id)

		if err != nil {
			LogErrorResponse("resourceNcloudSourceCommitRepositoryUpdate", err, *id)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudSourceCommitRepositoryUpdate", id)
	}

	return resourceNcloudSourceCommitRepositoryRead(ctx, d, meta)
}

func resourceNcloudSourceCommitRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	id := ncloud.String(d.Id())

	LogCommonRequest("resourceNcloudSourceCommitRepositoryDelete", *id)

	if _, err := config.Client.Sourcecommit.V1Api.DeleteRepository(ctx, id); err != nil {
		LogErrorResponse("resourceNcloudSourceCommitRepositoryDelete", err, *id)
		return diag.FromErr(err)
	}

	LogResponse("resourceNcloudSourceCommitRepositoryDelete", id)
	d.SetId("")
	return nil
}

func waitForSourceCommitRepositoryActive(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig, name string) error {

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"RESOLVE"},
		Refresh: func() (result interface{}, state string, err error) {
			repository, err := getRepository(ctx, config, name)
			if err != nil {
				return nil, "", fmt.Errorf("Repository response error , name : (%s) to become activating: %s", name, err)
			}
			if repository == nil {
				return name, "NULL", nil
			}

			if ncloud.StringValue(repository.Name) == name {
				d.SetId(strconv.Itoa(*repository.Id))
				return repository, "RESOLVE", nil
			}

			return nil, "PENDING", nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for SourceCommit Repository id : (%s) to become activating: %s", name, err)
	}
	return nil
}

func getRepository(ctx context.Context, config *conn.ProviderConfig, name string) (*sourcecommit.GetRepositoryDetailResponse, error) {

	LogCommonRequest("getRepository", name)
	resp, err := config.Client.Sourcecommit.V1Api.GetRepository(ctx, &name)

	if err != nil {
		LogErrorResponse("getRepository", err, name)
		return nil, err
	}
	LogResponse("getRepository", resp)

	return resp, nil
}

func GetRepositoryById(ctx context.Context, config *conn.ProviderConfig, id string) (*sourcecommit.GetRepositoryDetailResponse, error) {

	LogCommonRequest("getRepositoryById", id)
	resp, err := config.Client.Sourcecommit.V1Api.GetRepositoryById(ctx, &id)

	if err != nil {
		LogErrorResponse("getRepositoryById", err, id)
		return nil, err
	}
	LogResponse("getRepositoryById", resp)

	return resp, nil
}

func GetRepositories(ctx context.Context, config *conn.ProviderConfig) (*sourcecommit.GetRepositoryListResponse, error) {
	LogCommonRequest("getRepositories", "")
	resp, err := config.Client.Sourcecommit.V1Api.GetRepositories(ctx)
	if err != nil {
		LogErrorResponse("getRepositories", err, "")
		return nil, err
	}
	LogResponse("getRepositories", resp)

	return resp, nil
}

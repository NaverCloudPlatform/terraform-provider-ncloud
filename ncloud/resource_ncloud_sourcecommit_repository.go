package ncloud

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
)

func init() {
	RegisterResource("ncloud_sourcecommit_repository", resourceNcloudSourceCommitRepository())
}

func resourceNcloudSourceCommitRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudSourceCommitRepositoryCreate,
		ReadContext:   resourceNcloudSourceCommitRepositoryRead,
		UpdateContext: resourceNcloudSourceCommitRepositoryUpdate,
		DeleteContext: resourceNcloudSourceCommitRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultTimeout),
			Read:   schema.DefaultTimeout(DefaultTimeout),
			Update: schema.DefaultTimeout(DefaultTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(1, 100)),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 500)),
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_https": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_ssh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filesafer": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceNcloudSourceCommitRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	reqParams := &sourcecommit.CreateRepository{
		Name:        ncloud.String(d.Get("name").(string)),
		Description: StringPtrOrNil(d.GetOk("description")),
	}

	if fileSafer, ok := d.GetOk("filesafer"); ok {
		reqParams.Linked = &sourcecommit.CreateRepositoryLinked{
			FileSafer: BoolPtrOrNil(fileSafer, ok),
		}
	}

	logCommonRequest("resourceNcloudSourceCommitRepositoryCreate", reqParams)
	resp, err := config.Client.sourcecommit.V1Api.CreateRepository(ctx, reqParams)
	logCommonResponse("resourceNcloudSourceCommitRepositoryCreate", GetCommonResponse(nil))
	var diags diag.Diagnostics

	if err != nil {
		logErrorResponse("resourceNcloudSourceCommitRepositoryCreate", err, reqParams)

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Fail to create repository",
			Detail:   err.Error(),
		})
		return diags
	}

	logResponse("resourceNcloudSourceCommitRepositoryCreate", resp)

	name := ncloud.StringValue(reqParams.Name)

	if err := waitForSourceCommitRepositoryActive(ctx, d, config, name); err != nil {

		name := d.Get("name").(string)

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
	config := meta.(*ProviderConfig)
	id := ncloud.String(d.Id())

	repository, err := getRepositoryById(ctx, config, *id)

	logCommonRequest("resourceNcloudSourceCommitRepositoryRead", id)
	var diags diag.Diagnostics

	if err != nil {
		logErrorResponse("resourceNcloudSourceCommitRepositoryRead", err, *id)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to search repository",
			Detail:   fmt.Sprintf("Unable to search repository - detail repository id : %s", *id),
		})
		return diags
	}

	logResponse("resourceNcloudSourceCommitRepositoryRead", repository)

	if repository == nil {
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(int(*repository.Id)))
	d.Set("name", repository.Name)
	d.Set("description", repository.Description)
	d.Set("creator", repository.Created.User)
	d.Set("git_https", repository.Git.Https)
	d.Set("git_ssh", repository.Git.Ssh)
	d.Set("filesafer", repository.Linked.FileSafer)

	return nil
}

func resourceNcloudSourceCommitRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	reqParams := &sourcecommit.ChangeRepository{
		Description: StringPtrOrNil(d.GetOk("description")),
	}

	if fileSafer, ok := d.GetOk("filesafer"); ok {
		reqParams.Linked = &sourcecommit.CreateRepositoryLinked{
			FileSafer: BoolPtrOrNil(fileSafer, ok),
		}
	}

	id := ncloud.String(d.Id())

	logCommonRequest("resourceNcloudSourceCommitRepositoryUpdate", reqParams)
	_, err := config.Client.sourcecommit.V1Api.ChangeRepository(ctx, reqParams, id)

	if err != nil {
		logErrorResponse("resourceNcloudSourceCommitRepositoryUpdate", err, *id)
		return diag.FromErr(err)
	}

	logResponse("resourceNcloudSourceCommitRepositoryUpdate", id)
	return resourceNcloudSourceCommitRepositoryRead(ctx, d, meta)
}

func resourceNcloudSourceCommitRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	id := ncloud.String(d.Id())

	logCommonRequest("resourceNcloudSourceCommitRepositoryDelete", *id)

	if _, err := config.Client.sourcecommit.V1Api.DeleteRepository(ctx, id); err != nil {
		logErrorResponse("resourceNcloudSourceCommitRepositoryDelete", err, *id)
		return diag.FromErr(err)
	}

	logResponse("resourceNcloudSourceCommitRepositoryDelete", id)
	d.SetId("")
	return nil
}

func waitForSourceCommitRepositoryActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, name string) error {

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

			return nil, "PENDING", err
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

func getRepository(ctx context.Context, config *ProviderConfig, name string) (*sourcecommit.GetRepositoryDetailResponse, error) {

	logCommonRequest("getRepository", name)
	resp, err := config.Client.sourcecommit.V1Api.GetRepository(ctx, &name)

	if err != nil {
		logErrorResponse("getRepository", err, name)
		return nil, err
	}
	logResponse("getRepository", resp)

	return resp, nil
}

func getRepositoryById(ctx context.Context, config *ProviderConfig, id string) (*sourcecommit.GetRepositoryDetailResponse, error) {

	logCommonRequest("getRepositoryById", id)
	resp, err := config.Client.sourcecommit.V1Api.GetRepositoryById(ctx, &id)

	if err != nil {
		logErrorResponse("getRepositoryById", err, id)
		return nil, err
	}
	logResponse("getRepositoryById", resp)

	return resp, nil
}

func getRepositories(ctx context.Context, config *ProviderConfig) (*sourcecommit.GetRepositoryListResponse, error) {
	logCommonRequest("getRepositories", "")
	resp, err := config.Client.sourcecommit.V1Api.GetRepositories(ctx)
	if err != nil {
		logErrorResponse("getRepositories", err, "")
		return nil, err
	}
	logResponse("getRepositories", resp)

	return resp, nil
}

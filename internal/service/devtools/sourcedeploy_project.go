package devtools

import (
	"context"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcedeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudSourceDeployProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudSourceDeployProjectCreate,
		ReadContext:   resourceNcloudSourceDeployProjectRead,
		DeleteContext: resourceNcloudSourceDeployProjectDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultTimeout),
			Read:   schema.DefaultTimeout(conn.DefaultTimeout),
			Update: schema.DefaultTimeout(conn.DefaultTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(regexp.MustCompile(`^[^ !@#$%^&*()+\=\[\]{};':"\\|,.<>\/?]+$`), `Cannot contain special characters ( !@#$%^&*()+\=\[\]{};':"\\|,.<>\/?).`),
				)),
			},
		},
	}
}

func resourceNcloudSourceDeployProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	reqParams := &vsourcedeploy.CreateProject{
		Name: StringPtrOrNil(d.GetOk("name")),
	}

	LogCommonRequest("CreateSourceDeployProject", reqParams)
	resp, err := config.Client.Vsourcedeploy.V1Api.CreateProject(ctx, reqParams)
	if err != nil {
		LogErrorResponse("CreateSourceDeployProject", err, reqParams)
		return diag.FromErr(err)
	}
	LogResponse("CreateSourceDeployProject", resp)
	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(resp.Id))))

	return resourceNcloudSourceDeployProjectRead(ctx, d, meta)
}

func resourceNcloudSourceDeployProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	project, err := GetSourceDeployProjectByName(ctx, config, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if project == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(project.Id))))
	d.Set("name", project.Name)

	return nil
}

func resourceNcloudSourceDeployProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	LogCommonRequest("DeleteSourceDeployProject", d.Id())
	resp, err := config.Client.Vsourcedeploy.V1Api.DeleteProject(ctx, ncloud.String(d.Id()))
	if err != nil {
		LogErrorResponse("DeleteSourceDeployProject", err, d.Id())
		return diag.FromErr(err)
	}

	LogResponse("DeleteSourceDeployProject", resp)
	d.SetId("")
	return nil
}

func GetSourceDeployProjectById(ctx context.Context, config *conn.ProviderConfig, id string) (*vsourcedeploy.GetIdNameResponse, error) {
	projectList, err := getSourceDeployProjects(ctx, config)
	if err != nil {
		return nil, err
	}
	for _, project := range projectList {
		if *ncloud.IntString(int(ncloud.Int32Value(project.Id))) == id {
			return project, nil
		}
	}
	return nil, nil
}

func GetSourceDeployProjectByName(ctx context.Context, config *conn.ProviderConfig, name string) (*vsourcedeploy.GetIdNameResponse, error) {
	projectList, err := getSourceDeployProjects(ctx, config)
	if err != nil {
		return nil, err
	}
	for _, project := range projectList {
		if *project.Name == name {
			return project, nil
		}
	}
	return nil, nil
}

func getSourceDeployProjects(ctx context.Context, config *conn.ProviderConfig) ([]*vsourcedeploy.GetIdNameResponse, error) {
	reqParams := make(map[string]interface{})

	LogCommonRequest("GetSourceDeployProjects", reqParams)
	resp, err := config.Client.Vsourcedeploy.V1Api.GetProjects(ctx, reqParams)
	if err != nil {
		LogErrorResponse("GetSourceDeployProjects", err, reqParams)
		return nil, err
	}
	LogResponse("GetSourceDeployProjects", resp)

	return resp.ProjectList, nil
}

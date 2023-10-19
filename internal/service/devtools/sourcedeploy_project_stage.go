package devtools

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcedeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudSourceDeployStage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudSourceDeployStageCreate,
		ReadContext:   resourceNcloudSourceDeployStageRead,
		DeleteContext: resourceNcloudSourceDeployStageDelete,
		UpdateContext: resourceNcloudSourceDeployStageUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ":")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected PROJECT_ID:STAGE_ID", d.Id())
				}
				projectId, _ := strconv.ParseInt(idParts[0], 10, 32)
				stageId := idParts[1]

				d.Set("project_id", projectId)
				d.SetId(stageId)
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
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(regexp.MustCompile(`^[^ !@#$%^&*()+\=\[\]{};':"\\|,.<>\/?]+$`), `Cannot contain special characters ( !@#$%^&*()+\=\[\]{};':"\\|,.<>\/?).`),
				)),
			},
			"target_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Server", "AutoScalingGroup", "KubernetesService", "ObjectStorage"}, false)),
			},
			"config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.StringMatch(regexp.MustCompile(`^[0-9]+$`), `Only "numbers" can be entered.`),
										),
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"auto_scaling_group_no": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"auto_scaling_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNcloudSourceDeployStageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_sourcedeploy_project_stage`"))
	}

	reqParams, paramsErr := getStage(d)
	if paramsErr != nil {
		return diag.FromErr(paramsErr)
	}
	projectId := ncloud.IntString(d.Get("project_id").(int))
	LogCommonRequest("createSourceDeployStage", reqParams)
	resp, err := config.Client.Vsourcedeploy.V1Api.CreateStage(ctx, reqParams, projectId)
	if err != nil {
		LogErrorResponse("createSourceDeployStage", err, reqParams)
		return diag.FromErr(err)
	}
	LogResponse("createSourceDeployStage", resp.Id)

	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(resp.Id))))
	d.Set("project_id", Int32PtrOrNil(d.GetOk("project_id")))

	return resourceNcloudSourceDeployStageRead(ctx, d, meta)
}

func resourceNcloudSourceDeployStageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_sourcedeploy_project_stage`"))
	}
	projectId := ncloud.IntString(d.Get("project_id").(int))
	stage, err := GetSourceDeployStageById(ctx, config, projectId, ncloud.String(d.Id()))

	if err != nil {
		return diag.FromErr(err)
	}

	if stage == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(stage.Id))))
	d.Set("name", stage.Name)
	d.Set("target_type", stage.Type_)
	d.Set("config", makeStageConfig(stage.Config))

	return nil
}

func resourceNcloudSourceDeployStageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	err := changeDeployStage(ctx, d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNcloudSourceDeployStageRead(ctx, d, meta)
}

func resourceNcloudSourceDeployStageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_sourcedeploy_project_stage`"))
	}

	projectId := ncloud.IntString(d.Get("project_id").(int))
	LogCommonRequest("deleteSourceDeployStage", d.Id())
	resp, err := config.Client.Vsourcedeploy.V1Api.DeleteStage(ctx, projectId, ncloud.String(d.Id()))
	if err != nil {
		LogErrorResponse("deleteSourceDeployStage", err, d.Id())
		return diag.FromErr(err)
	}

	LogResponse("deleteSourceDeployStage", resp)
	d.SetId("")
	return nil
}

func getStage(d *schema.ResourceData) (*vsourcedeploy.CreateStage, error) {
	deployTarget, deployTargetErr := getDeployTarget(d)
	if deployTargetErr != nil {
		return nil, deployTargetErr
	}
	reqParams := &vsourcedeploy.CreateStage{
		Name:   StringPtrOrNil(d.GetOk("name")),
		Type_:  StringPtrOrNil(d.GetOk("target_type")),
		Config: deployTarget,
	}
	return reqParams, nil
}

func getDeployTarget(d *schema.ResourceData) (*vsourcedeploy.StageConfig, error) {
	deployTarget := vsourcedeploy.StageConfig{}

	deployTargetType := ncloud.StringValue(StringPtrOrNil(d.GetOk("target_type")))

	switch deployTargetType {
	case "Server":
		if server, ok := d.GetOk("config.0.server"); ok {
			servers, err := expandServerParams(server.([]interface{}))
			if err != nil {
				return nil, err
			}
			deployTarget.ServerNo = servers
		}
		if deployTarget.ServerNo == nil {
			return nil, fmt.Errorf("config(server.id) is required")
		}
	case "AutoScalingGroup":
		deployTarget.AutoScalingGroupNo = Int32PtrOrNil(d.GetOk("config.0.auto_scaling_group_no"))
		if deployTarget.AutoScalingGroupNo == nil {
			return nil, fmt.Errorf("config(auto_scaling_group_no) is required")
		}
	case "KubernetesService":
		deployTarget.ClusterUuid = StringPtrOrNil(d.GetOk("config.0.cluster_uuid"))
		if deployTarget.ClusterUuid == nil {
			return nil, fmt.Errorf("config(cluster_uuid) is required")
		}
	case "ObjectStorage":
		deployTarget.BucketName = StringPtrOrNil(d.GetOk("config.0.bucket_name"))
		if deployTarget.BucketName == nil {
			return nil, fmt.Errorf("config(bucket_name) is required")
		}
	}
	return &deployTarget, nil
}

func GetSourceDeployStageById(ctx context.Context, config *conn.ProviderConfig, projectId *string, id *string) (*vsourcedeploy.GetStageDetailResponse, error) {
	LogCommonRequest("getSourceDeployStage", id)
	resp, err := config.Client.Vsourcedeploy.V1Api.GetStage(ctx, projectId, id)
	if err != nil {
		LogErrorResponse("getSourceDeployStage", err, *id)
		return nil, err
	}
	LogResponse("getSourceDeployStage", resp)

	return resp, nil
}

func makeStageConfig(config *vsourcedeploy.GetStageDetailResponseConfig) []interface{} {
	if config == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["server"] = flattenServer(config)
	values["auto_scaling_group_no"] = ncloud.Int32Value(config.AutoScalingGroupNo)
	values["auto_scaling_group_name"] = ncloud.StringValue(config.AutoScalingGroupName)
	values["cluster_uuid"] = ncloud.StringValue(config.ClusterUuid)
	values["cluster_name"] = ncloud.StringValue(config.ClusterName)
	values["bucket_name"] = ncloud.StringValue(config.BucketName)

	return []interface{}{values}
}

func changeDeployStage(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig) error {

	reqParams, paramErr := getStage(d)

	if paramErr != nil {
		return paramErr
	}
	projectId := ncloud.IntString(d.Get("project_id").(int))
	id := ncloud.String(d.Id())

	LogCommonRequest("changeSourceDeployStage", reqParams)
	resp, err := config.Client.Vsourcedeploy.V1Api.ChangeStage(ctx, reqParams, projectId, id)
	if err != nil {
		LogErrorResponse("changeSourceDeployStage", err, reqParams)
		return err
	}
	LogResponse("changeSourceDeployStage", resp)

	return nil
}

func expandServerParams(servers []interface{}) ([]*int32, error) {
	var list []*int32
	for _, v := range servers {
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "id":
				id, err := strconv.Atoi(value.(string))
				if err == nil {
					list = append(list, ncloud.Int32(int32(id)))
				}
			}
		}
	}

	return list, nil
}

func flattenServer(config *vsourcedeploy.GetStageDetailResponseConfig) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(config.ServerNo))

	for i := 0; i < len(config.ServerNo); i++ {
		values := map[string]interface{}{}
		values["id"] = ncloud.StringValue(ncloud.Int32String(ncloud.Int32Value((config.ServerNo[i]))))
		values["name"] = ncloud.StringValue(config.ServerName[i])

		list = append(list, values)
	}

	return list
}

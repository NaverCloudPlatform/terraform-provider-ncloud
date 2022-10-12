package ncloud

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
)

func init() {
	RegisterResource("ncloud_sourcedeploy_project_stage_scenario", resourceNcloudSourceDeployScenario())
}

func resourceNcloudSourceDeployScenario() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudSourceDeployScenarioCreate,
		ReadContext:   resourceNcloudSourceDeployScenarioRead,
		DeleteContext: resourceNcloudSourceDeployScenarioDelete,
		UpdateContext: resourceNcloudSourceDeployScenarioUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ":")
				if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected PROJECT_ID:STAGE_ID:SCENARIO_ID", d.Id())
				}
				projectId, _ := strconv.ParseInt(idParts[0], 10, 32)
				stageId, _ := strconv.ParseInt(idParts[1], 10, 32)
				scenarioId := idParts[2]

				d.Set("project_id", projectId)
				d.Set("stage_id", stageId)
				d.SetId(scenarioId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultTimeout),
			Read:   schema.DefaultTimeout(DefaultTimeout),
			Update: schema.DefaultTimeout(DefaultTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"stage_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(regexp.MustCompile(`^[^ !@#$%^&*()+\=\[\]{};':"\\|,.<>\/?]+$`), `Cannot contain special characters ( !@#$%^&*()+\=\[\]{};':"\\|,.<>\/?).`),
				)),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(0, 500),
				)),
			},
			"config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"normal", "blueGreen", "rolling", "canary"}, false)),
						},
						"file": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"SourceBuild", "ObjectStorage", "later"}, false)),
									},
									"object_storage": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"object": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"source_build": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"rollback": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"deploy_command": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pre_deploy": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateDiagFunc: ToDiagFunc(validation.All(
														validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
													)),
												},
												"command": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateDiagFunc: ToDiagFunc(validation.All(
														validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
													)),
												},
											},
										},
									},
									"path": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_path": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateDiagFunc: ToDiagFunc(validation.All(
														validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
													)),
												},
												"deploy_path": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateDiagFunc: ToDiagFunc(validation.All(
														validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
													)),
												},
											},
										},
									},
									"post_deploy": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateDiagFunc: ToDiagFunc(validation.All(
														validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
													)),
												},
												"command": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateDiagFunc: ToDiagFunc(validation.All(
														validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
													)),
												},
											},
										},
									},
								},
							},
						},
						"load_balancer": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_target_group_no": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"load_balancer_target_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"delete_server": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"manifest": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"SourceCommit"}, false)),
									},
									"repository_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"branch": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"canary_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"canary_count": {
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 10)),
									},
									"analysis_type": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"manual", "auto"}, false)),
									},
									"timeout": {
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 360)),
									},
									"prometheus": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"env": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"baseline": {
													Type:             schema.TypeString,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 64)),
												},
												"canary": {
													Type:             schema.TypeString,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 64)),
												},
											},
										},
									},
									"metrics": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"success_criteria": {
													Type:             schema.TypeString,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"base", "canary"}, false)),
												},
												"query_type": {
													Type:             schema.TypeString,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"default", "promQL"}, false)),
												},
												"weight": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"metric": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"filter": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"query": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"analysis_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"duration": {
													Type:             schema.TypeInt,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.IntBetween(10, 360)),
												},
												"delay": {
													Type:             schema.TypeInt,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 60)),
												},
												"interval": {
													Type:             schema.TypeInt,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 360)),
												},
												"step": {
													Type:             schema.TypeInt,
													Optional:         true,
													ValidateDiagFunc: ToDiagFunc(validation.IntBetween(10, 360)),
												},
											},
										},
									},
									"pass_score": {
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 100)),
									},
								},
							},
						},
						"path": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_path": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: ToDiagFunc(validation.All(
											validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
										)),
									},
									"deploy_path": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: ToDiagFunc(validation.All(
											validation.StringMatch(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), `Cannot Spaces and empty strings`),
										)),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceNcloudSourceDeployScenarioCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_sourcedeploy_project_stage_scenario`"))
	}
	projectId := ncloud.IntString(d.Get("project_id").(int))
	stageId := ncloud.IntString(d.Get("stage_id").(int))

	stageResp, stageErr := getSourceDeployStageById(ctx, config, projectId, stageId)
	if stageErr != nil {
		return diag.FromErr(stageErr)
	}

	reqParams, paramsErr := getScenario(ncloud.StringValue(stageResp.Type_), d)
	if paramsErr != nil {
		return diag.FromErr(paramsErr)
	}

	logCommonRequest("createSourceDeployScenario", reqParams)
	scenarioCreateResp, scenarioCreateRespErr := config.Client.vsourcedeploy.V1Api.CreateScenario(ctx, reqParams, projectId, stageId)
	if scenarioCreateRespErr != nil {
		logErrorResponse("createSourceDeployScenario", scenarioCreateRespErr, reqParams)
		return diag.FromErr(scenarioCreateRespErr)
	}
	logResponse("createSourceDeployScenario", scenarioCreateResp.Id)

	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(scenarioCreateResp.Id))))

	return resourceNcloudSourceDeployScenarioRead(ctx, d, meta)
}

func resourceNcloudSourceDeployScenarioRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_sourcedeploy_project_stage_scenario`"))
	}
	projectId := ncloud.IntString(d.Get("project_id").(int))
	stageId := ncloud.IntString(d.Get("stage_id").(int))
	scenario, err := getSourceDeployScenarioById(ctx, config, projectId, stageId, ncloud.String(d.Id()))

	if err != nil {
		return diag.FromErr(err)
	}

	if scenario == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(scenario.Id))))
	d.Set("name", scenario.Name)
	d.Set("description", scenario.Description)
	d.Set("config", makeScenarioConfig(scenario.Config))

	return nil
}

func getScenario(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.CreateScenario, error) {
	commonScenarioConfig, commonScenarioConfigErr := commonScenario(deployTargetType, d)
	if commonScenarioConfigErr != nil {
		return nil, commonScenarioConfigErr
	}
	reqParams := &vsourcedeploy.CreateScenario{
		Name:        StringPtrOrNil(d.GetOk("name")),
		Description: commonScenarioConfig.Description,
		Config:      commonScenarioConfig.Config,
	}

	if reqParams.Name == nil {
		return nil, fmt.Errorf("name is required")
	}
	return reqParams, nil
}

func commonScenario(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.ChangeScenario, error) {
	scenarioConfig, scenarioConfigErr := getScenarioConfig(deployTargetType, d)
	if scenarioConfigErr != nil {
		return nil, scenarioConfigErr
	}
	reqParams := &vsourcedeploy.ChangeScenario{
		Description: StringPtrOrNil(d.GetOk("description")),
		Config:      scenarioConfig,
	}
	return reqParams, nil
}

func getScenarioConfig(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.ScenarioConfig, error) {
	file, fileErr := getFile(deployTargetType, d)
	if fileErr != nil {
		return nil, fileErr
	}
	cmd, cmdErr := getCmd(deployTargetType, d)
	if cmdErr != nil {
		return nil, cmdErr
	}
	lb, lbErr := getLoadBalnacer(deployTargetType, d)
	if lbErr != nil {
		return nil, lbErr
	}
	manifest, manifestErr := getManifest(deployTargetType, d)
	if manifestErr != nil {
		return nil, manifestErr
	}
	canaryConfig, canaryConfigErr := getCanaryConfig(deployTargetType, d)
	if canaryConfigErr != nil {
		return nil, canaryConfigErr
	}
	path, pathErr := expandDeployPathParams(d.Get("config.0.path").([]interface{}))
	if pathErr != nil {
		return nil, pathErr
	}

	scenarioConfig := vsourcedeploy.ScenarioConfig{
		Strategy:     StringPtrOrNil(d.GetOk("config.0.strategy")),
		File:         file,
		Rollback:     ncloud.Bool(d.Get("config.0.rollback").(bool)),
		Cmd:          cmd,
		LoadBalancer: lb,
		Manifest:     manifest,
		CanaryConfig: canaryConfig,
		Path:         path,
	}

	if (deployTargetType == "Server" || deployTargetType == "AutoScalingGroup") && (scenarioConfig.Rollback == nil || scenarioConfig.Strategy == nil) {
		return nil, fmt.Errorf("config parameters (strategy, rollback) is required")
	}
	return &scenarioConfig, nil
}

func getFile(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigFile, error) {
	if deployTargetType == "Server" || deployTargetType == "AutoScalingGroup" || deployTargetType == "ObjectStorage" {
		objectstorage, objectstorageErr := getFileObjectStorage(d)
		if objectstorageErr != nil {
			return nil, objectstorageErr
		}
		sourcebuild, sourcebuildErr := getFileSourceBuild(d)
		if sourcebuildErr != nil {
			return nil, sourcebuildErr
		}
		reqParams := vsourcedeploy.ScenarioConfigFile{
			Type_:         StringPtrOrNil(d.GetOk("config.0.file.0.type")),
			ObjectStorage: objectstorage,
			SourceBuild:   sourcebuild,
		}

		if reqParams.Type_ == nil {
			return nil, fmt.Errorf("config.file parameters (type) is required")
		}
		return &reqParams, nil
	}
	return nil, nil
}

func getFileObjectStorage(d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigFileObjectStorage, error) {
	fileType := ncloud.StringValue(StringPtrOrNil(d.GetOk("config.0.file.0.type")))
	if fileType == "ObjectStorage" {
		reqParams := vsourcedeploy.ScenarioConfigFileObjectStorage{
			Bucket: StringPtrOrNil(d.GetOk("config.0.file.0.object_storage.0.bucket")),
			Object: StringPtrOrNil(d.GetOk("config.0.file.0.object_storage.0.object")),
		}

		if reqParams.Bucket == nil || reqParams.Object == nil {
			return nil, fmt.Errorf("config.file.object_storage parameters (bucket, object) is required")
		}
		return &reqParams, nil
	}
	return nil, nil
}
func getFileSourceBuild(d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigFileSourceBuild, error) {
	fileType := ncloud.StringValue(StringPtrOrNil(d.GetOk("config.0.file.0.type")))
	if fileType == "SourceBuild" {
		reqParams := vsourcedeploy.ScenarioConfigFileSourceBuild{
			Id: Int32PtrOrNil(d.GetOk("config.0.file.0.source_build.0.id")),
		}

		if reqParams.Id == nil {
			return nil, fmt.Errorf("config.file.source_build parameters (id) is required")
		}

		return &reqParams, nil
	}
	return nil, nil
}

func expandCmdPrePostParams(cmdPrePosts []interface{}) ([]*vsourcedeploy.ScenarioConfigCmdPrePost, error) {
	list := make([]*vsourcedeploy.ScenarioConfigCmdPrePost, 0, len(cmdPrePosts))

	for _, v := range cmdPrePosts {
		cmdPrePost := new(vsourcedeploy.ScenarioConfigCmdPrePost)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "user":
				cmdPrePost.User = ncloud.String(value.(string))
			case "command":
				cmdPrePost.Cmd = ncloud.String(value.(string))
			}
		}

		if ncloud.StringValue(cmdPrePost.User) == "" {
			return nil, fmt.Errorf("please write config.deploy_command.pre_deploy.user or config.deploy_command.post_deploy.user")
		} else if ncloud.StringValue(cmdPrePost.Cmd) == "" {
			return nil, fmt.Errorf("please write config.deploy_command.pre_deploy.command or config.deploy_command.post_deploy.command")
		}

		list = append(list, cmdPrePost)
	}

	return list, nil
}

func expandDeployPathParams(deployPaths []interface{}) ([]*vsourcedeploy.ScenarioConfigCmdDeploy, error) {
	list := make([]*vsourcedeploy.ScenarioConfigCmdDeploy, 0, len(deployPaths))

	for _, v := range deployPaths {
		deployPath := new(vsourcedeploy.ScenarioConfigCmdDeploy)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "source_path":
				deployPath.SourcePath = ncloud.String(value.(string))
			case "deploy_path":
				deployPath.DeployPath = ncloud.String(value.(string))
			}
		}

		if ncloud.StringValue(deployPath.SourcePath) == "" {
			return nil, fmt.Errorf("please write config.deploy_command.path.source_path or config.deploy_command.path.source_path")
		} else if ncloud.StringValue(deployPath.DeployPath) == "" {
			return nil, fmt.Errorf("please write config.deploy_command.path.deploy_path or config.deploy_command.path.deploy_path")
		}

		list = append(list, deployPath)
	}

	return list, nil
}

func getCmd(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigCmd, error) {
	if deployTargetType == "Server" || deployTargetType == "AutoScalingGroup" || deployTargetType == "ObjecStorage" {
		pre, preErr := expandCmdPrePostParams(d.Get("config.0.deploy_command.0.pre_deploy").([]interface{}))
		if preErr != nil {
			return nil, preErr
		}
		post, postErr := expandCmdPrePostParams(d.Get("config.0.deploy_command.0.post_deploy").([]interface{}))
		if postErr != nil {
			return nil, postErr
		}
		deployPath, deployPathErr := expandDeployPathParams(d.Get("config.0.deploy_command.0.path").([]interface{}))
		if deployPathErr != nil {
			return nil, deployPathErr
		}

		reqParams := vsourcedeploy.ScenarioConfigCmd{
			Pre:    pre,
			Deploy: deployPath,
			Post:   post,
		}
		return &reqParams, nil
	}
	return nil, nil
}

func getLoadBalnacer(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigLoadBalancer, error) {
	strategy := ncloud.StringValue(StringPtrOrNil(d.GetOk("config.0.strategy")))
	if deployTargetType == "AutoScalingGroup" && strategy == "blueGreen" {
		reqParams := vsourcedeploy.ScenarioConfigLoadBalancer{
			LoadBalancerTargetGroupNo: Int32PtrOrNil(d.GetOk("config.0.load_balancer.0.load_balancer_target_group_no")),
			DeleteServer:              ncloud.Bool(d.Get("config.0.load_balancer.0.delete_server").(bool)),
		}

		if reqParams.LoadBalancerTargetGroupNo == nil || reqParams.DeleteServer == nil {
			return nil, fmt.Errorf("config.loadBalancer parameters(load_balancer_target_group_no, delete_server) is required")
		}
		return &reqParams, nil
	}
	return nil, nil
}

func getManifest(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigManifest, error) {
	if deployTargetType == "KubernetesService" {
		reqParams := vsourcedeploy.ScenarioConfigManifest{
			Type_:      StringPtrOrNil(d.GetOk("config.0.manifest.0.type")),
			Repository: StringPtrOrNil(d.GetOk("config.0.manifest.0.repository_name")),
			Branch:     StringPtrOrNil(d.GetOk("config.0.manifest.0.branch")),
		}

		if param, ok := d.GetOk("config.0.manifest.0.path"); ok {
			reqParams.Path = expandStringInterfaceList(param.([]interface{}))
		}

		if reqParams.Type_ == nil || reqParams.Repository == nil || reqParams.Branch == nil || reqParams.Path == nil {
			return nil, fmt.Errorf("config.manifest parameters(type, repository_name, branch, path) is required")
		}
		return &reqParams, nil
	}
	return nil, nil
}

func getCanaryConfig(deployTargetType string, d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigCanaryConfig, error) {
	strategy := ncloud.StringValue(StringPtrOrNil(d.GetOk("config.0.strategy")))
	if deployTargetType == "KubernetesService" && strategy == "canary" {
		env, envErr := getCanaryConfigEnv(d)
		if envErr != nil {
			return nil, envErr
		}
		metrics, metricsErr := expandCanaryConfigMetricsParams(d.Get("config.0.canary_config.0.metrics").([]interface{}))
		if metricsErr != nil {
			return nil, metricsErr
		}
		analysisConfig, analysisConfigErr := getCanaryConfigAnalysisConfig(d)
		if analysisConfigErr != nil {
			return nil, analysisConfigErr
		}

		reqParams := vsourcedeploy.ScenarioConfigCanaryConfig{
			AnalysisType:   StringPtrOrNil(d.GetOk("config.0.canary_config.0.analysis_type")),
			CanaryCount:    Int32PtrOrNil(d.GetOk("config.0.canary_config.0.canary_count")),
			Timeout:        Int32PtrOrNil(d.GetOk("config.0.canary_config.0.timeout")),
			Prometheus:     StringPtrOrNil(d.GetOk("config.0.canary_config.0.prometheus")),
			Env:            env,
			Metrics:        metrics,
			AnalysisConfig: analysisConfig,
			PassScore:      Int32PtrOrNil(d.GetOk("config.0.canary_config.0.pass_score")),
		}

		if reqParams.AnalysisType == nil {
			return nil, fmt.Errorf("config.canary_config parameters (analysis_type) is required")
		} else if ncloud.StringValue(reqParams.AnalysisType) == "manual" && (reqParams.CanaryCount == nil || reqParams.Timeout == nil) {
			return nil, fmt.Errorf("config.canary_config parameters (analysis_type, canary_count, timeout) is required")
		} else if ncloud.StringValue(reqParams.AnalysisType) == "auto" && (reqParams.CanaryCount == nil || reqParams.Prometheus == nil || reqParams.PassScore == nil) {
			return nil, fmt.Errorf("config.canary_config parameters (canary_count, prometheus, pass_score) is required")
		}

		return &reqParams, nil
	}
	return nil, nil
}

func getCanaryConfigEnv(d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigCanaryConfigEnv, error) {
	analysisType := ncloud.StringValue(StringPtrOrNil(d.GetOk("config.0.canary_config.0.analysis_type")))
	if analysisType == "auto" {
		reqParams := vsourcedeploy.ScenarioConfigCanaryConfigEnv{
			Baseline: StringPtrOrNil(d.GetOk("config.0.canary_config.0.env.0.baseline")),
			Canary:   StringPtrOrNil(d.GetOk("config.0.canary_config.0.env.0.canary")),
		}

		if reqParams.Baseline == nil || reqParams.Canary == nil {
			return nil, fmt.Errorf("config.canary_config.env parameters (baseline, canary) is required")
		}

		return &reqParams, nil
	}
	return nil, nil
}

func expandCanaryConfigMetricsParams(metrics []interface{}) ([]*vsourcedeploy.ScenarioConfigCanaryConfigMetrics, error) {
	var list []*vsourcedeploy.ScenarioConfigCanaryConfigMetrics

	for _, v := range metrics {
		m := v.(map[string]interface{})

		if len(m["name"].(string)) == 0 || len(m["success_criteria"].(string)) == 0 || m["weight"] == 0 || len(m["query_type"].(string)) == 0 {
			return nil, fmt.Errorf("config.canary_config.metrics parameters (name, success_criteria, weight, query_type) is required")
		}

		metric := &vsourcedeploy.ScenarioConfigCanaryConfigMetrics{
			Name:            ncloud.String(m["name"].(string)),
			SuccessCriteria: ncloud.String(m["success_criteria"].(string)),
			Weight:          ncloud.Int32(int32(m["weight"].(int))),
			QueryType:       ncloud.String(m["query_type"].(string)),
		}

		switch m["query_type"] {
		case "promQL":
			if len(m["query"].(string)) == 0 {
				return nil, fmt.Errorf("config.canary_config.metrics parameters (query) is required if config.canary_config.metrics.query_type is promQL")
			}
			metric.Query = ncloud.String(m["query"].(string))
		case "default":
			if len(m["metric"].(string)) == 0 || len(m["filter"].(string)) == 0 {
				return nil, fmt.Errorf("config.canary_config.metrics parameters (metric, filter) is required if config.anary_config.metrics.query_type is default")
			}
			metric.Metric = ncloud.String(m["metric"].(string))
			metric.Filter = ncloud.String(m["filter"].(string))
		}

		list = append(list, metric)
	}

	return list, nil
}

func getCanaryConfigAnalysisConfig(d *schema.ResourceData) (*vsourcedeploy.ScenarioConfigCanaryConfigAnalysisConfig, error) {
	analysisType := ncloud.StringValue(StringPtrOrNil(d.GetOk("config.0.canary_config.0.analysis_type")))
	if analysisType == "auto" {
		reqParams := vsourcedeploy.ScenarioConfigCanaryConfigAnalysisConfig{
			Duration: Int32PtrOrNil(d.GetOk("config.0.canary_config.0.analysis_config.0.duration")),
			Delay:    Int32PtrOrNil(d.GetOk("config.0.canary_config.0.analysis_config.0.delay")),
			Interval: Int32PtrOrNil(d.GetOk("config.0.canary_config.0.analysis_config.0.interval")),
			Step:     Int32PtrOrNil(d.GetOk("config.0.canary_config.0.analysis_config.0.step")),
		}

		if reqParams.Duration == nil || reqParams.Delay == nil || reqParams.Interval == nil || reqParams.Step == nil {
			return nil, fmt.Errorf("config.canary_config.analysis_config parameters(duration, delay, interval, step) is required")
		}
		return &reqParams, nil
	}
	return nil, nil
}

func getSourceDeployScenarioById(ctx context.Context, config *ProviderConfig, projectId *string, stageId *string, id *string) (*vsourcedeploy.GetScenarioDetailResponse, error) {
	logCommonRequest("getSourceDeployScenario", id)
	resp, err := config.Client.vsourcedeploy.V1Api.GetScenario(ctx, projectId, stageId, id)
	if err != nil {
		logErrorResponse("getSourceDeployScenario", err, *id)
		return nil, err
	}
	logResponse("getSourceDeployScenario", resp)

	return resp, nil
}

func resourceNcloudSourceDeployScenarioDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_sourcedeploy_project_stage_scenario`"))
	}

	projectId := ncloud.IntString(d.Get("project_id").(int))
	stageId := ncloud.IntString(d.Get("stage_id").(int))
	logCommonRequest("deleteSourceDeployScenario", d.Id())
	resp, err := config.Client.vsourcedeploy.V1Api.DeleteScenario(ctx, projectId, stageId, ncloud.String(d.Id()))
	if err != nil {
		logErrorResponse("deleteSourceDeployScenario", err, d.Id())
		return diag.FromErr(err)
	}
	logResponse("deleteSourceDeployScenario", resp)
	d.SetId("")
	return nil
}

func resourceNcloudSourceDeployScenarioUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	err := changeDeployScenario(ctx, d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNcloudSourceDeployScenarioRead(ctx, d, meta)
}

func changeDeployScenario(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	projectId := ncloud.IntString(d.Get("project_id").(int))
	stageId := ncloud.IntString(d.Get("stage_id").(int))
	stage, stageErr := getSourceDeployStageById(ctx, config, projectId, stageId)
	if stageErr != nil {
		return stageErr
	}

	reqParams, paramsErr := commonScenario(ncloud.StringValue(stage.Type_), d)
	if paramsErr != nil {
		return paramsErr
	}

	logCommonRequest("changeSourceDeployScenario", reqParams)
	resp, err := config.Client.vsourcedeploy.V1Api.ChangeScenario(ctx, reqParams, projectId, stageId, ncloud.String(d.Id()))
	if err != nil {
		logErrorResponse("changeSourceDeployScenario", err, reqParams)
		return err
	}
	logResponse("changeSourceDeployScenario", resp)

	return nil
}

func makeScenarioConfig(config *vsourcedeploy.GetScenarioConfig) []interface{} {
	if config == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["strategy"] = ncloud.StringValue(config.Strategy)
	values["file"] = makeConfigFile(config.File)
	values["rollback"] = ncloud.BoolValue(config.Rollback)
	values["deploy_command"] = makeConfigCmd(config.Cmd)
	values["load_balancer"] = makeConfigLoadBalancer(config.LoadBalancer)
	values["manifest"] = makeConfigManifest(config.Manifest)
	values["canary_config"] = makeConfigCanaryConfig(config.CanaryConfig)
	values["path"] = flattenCmdDeploy(config.Path)

	return []interface{}{values}
}

func makeConfigFile(file *vsourcedeploy.GetScenarioConfigFile) []interface{} {
	if file == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["type"] = ncloud.StringValue(file.Type_)
	values["object_storage"] = makeFileObjectStorage(file.ObjectStorage)
	values["source_build"] = makeFileSourceBuild(file.SourceBuild)

	return []interface{}{values}
}

func makeFileObjectStorage(objectStorage *vsourcedeploy.ScenarioConfigFileObjectStorage) []interface{} {
	if objectStorage == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["bucket"] = ncloud.StringValue(objectStorage.Bucket)
	values["object"] = ncloud.StringValue(objectStorage.Object)

	return []interface{}{values}
}

func makeFileSourceBuild(sourcebuild *vsourcedeploy.GetIdNameResponse) []interface{} {
	if sourcebuild == nil {
		return nil
	}
	values := map[string]interface{}{}
	values["id"] = ncloud.Int32Value(sourcebuild.Id)
	values["name"] = ncloud.StringValue(sourcebuild.Name)

	return []interface{}{values}
}

func makeConfigCmd(cmd *vsourcedeploy.ScenarioConfigCmd) []interface{} {
	if cmd == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["pre_deploy"] = flattenCmdPrePost(cmd.Pre)
	values["path"] = flattenCmdDeploy(cmd.Deploy)
	values["post_deploy"] = flattenCmdPrePost(cmd.Post)

	return []interface{}{values}
}

func flattenCmdPrePost(prePosts []*vsourcedeploy.ScenarioConfigCmdPrePost) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(prePosts))

	for _, prePost := range prePosts {
		values := map[string]interface{}{}
		values["user"] = ncloud.StringValue(prePost.User)
		values["command"] = ncloud.StringValue(prePost.Cmd)

		list = append(list, values)
	}

	return list
}

func flattenCmdDeploy(deploys []*vsourcedeploy.ScenarioConfigCmdDeploy) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(deploys))

	for _, deploy := range deploys {
		values := map[string]interface{}{}
		values["source_path"] = ncloud.StringValue(deploy.SourcePath)
		values["deploy_path"] = ncloud.StringValue(deploy.DeployPath)

		list = append(list, values)
	}

	return list
}

func makeConfigLoadBalancer(lb *vsourcedeploy.GetScenarioConfigLoadBalancer) []interface{} {
	if lb == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["load_balancer_target_group_no"] = ncloud.Int32Value(lb.LoadBalancerTargetGroupNo)
	values["load_balancer_target_group_name"] = ncloud.StringValue(lb.LoadBalancerTargetGroupName)
	values["delete_server"] = ncloud.BoolValue(lb.DeleteServer)

	return []interface{}{values}
}

func makeConfigManifest(manifest *vsourcedeploy.ScenarioConfigManifest) []interface{} {
	if manifest == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["type"] = ncloud.StringValue(manifest.Type_)
	values["repository_name"] = ncloud.StringValue(manifest.Repository)
	values["branch"] = ncloud.StringValue(manifest.Branch)
	values["path"] = ncloud.StringListValue(manifest.Path)

	return []interface{}{values}
}

func makeConfigCanaryConfig(canaryConfig *vsourcedeploy.ScenarioConfigCanaryConfig) []interface{} {
	if canaryConfig == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["canary_count"] = ncloud.Int32Value(canaryConfig.CanaryCount)
	values["analysis_type"] = ncloud.StringValue(canaryConfig.AnalysisType)
	values["timeout"] = ncloud.Int32Value(canaryConfig.Timeout)
	values["prometheus"] = ncloud.StringValue(canaryConfig.Prometheus)
	values["env"] = makeCanaryConfigEnv(canaryConfig.Env)
	values["metrics"] = flattenCanaryConfigMetrics(canaryConfig.Metrics)
	values["analysis_config"] = makeCanaryConfigAnalysisConfig(canaryConfig.AnalysisConfig)
	values["pass_score"] = ncloud.Int32Value(canaryConfig.PassScore)

	return []interface{}{values}
}

func makeCanaryConfigEnv(env *vsourcedeploy.ScenarioConfigCanaryConfigEnv) []interface{} {
	if env == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["baseline"] = ncloud.StringValue(env.Baseline)
	values["canary"] = ncloud.StringValue(env.Canary)

	return []interface{}{values}
}

func flattenCanaryConfigMetrics(metrics []*vsourcedeploy.ScenarioConfigCanaryConfigMetrics) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(metrics))

	for _, metric := range metrics {
		values := map[string]interface{}{}
		values["name"] = ncloud.StringValue(metric.Name)
		values["success_criteria"] = ncloud.StringValue(metric.SuccessCriteria)
		values["query_type"] = ncloud.StringValue(metric.QueryType)
		values["weight"] = ncloud.Int32Value(metric.Weight)
		values["metric"] = ncloud.StringValue(metric.Metric)
		values["filter"] = ncloud.StringValue(metric.Filter)
		values["query"] = ncloud.StringValue(metric.Query)

		list = append(list, values)
	}

	return list
}

func makeCanaryConfigAnalysisConfig(analysisConfig *vsourcedeploy.ScenarioConfigCanaryConfigAnalysisConfig) []interface{} {
	if analysisConfig == nil {
		return nil
	}
	values := map[string]interface{}{}

	values["duration"] = ncloud.Int32Value(analysisConfig.Duration)
	values["delay"] = ncloud.Int32Value(analysisConfig.Delay)
	values["interval"] = ncloud.Int32Value(analysisConfig.Interval)
	values["step"] = ncloud.Int32Value(analysisConfig.Step)

	return []interface{}{values}
}

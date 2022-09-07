package ncloud

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcebuild"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_sourcebuild_project", resourceNcloudSourceBuildProject())
}

func resourceNcloudSourceBuildProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudSourceBuildProjectCreate,
		ReadContext:   resourceNcloudSourceBuildProjectRead,
		DeleteContext: resourceNcloudSourceBuildProjectDelete,
		UpdateContext: resourceNcloudSourceBuildProjectUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultTimeout),
			Read:   schema.DefaultTimeout(DefaultTimeout),
			Update: schema.DefaultTimeout(DefaultTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"project_no": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 80),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]+$`), "Composed of alphabets, numbers, hyphen (-) and underbar (_)"),
				)),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(0, 500),
				)),
			},
			"source": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repository_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"branch": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"env": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"cpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"mem": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"platform": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"SourceBuild", "ContainerRegistry", "PublicRegistry"}, false)),
									},
									"config": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"os": {
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
															"version": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"archi": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"runtime": {
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
															"version": {
																Type:     schema.TypeList,
																Required: true,
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
												"registry": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"image": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"tag": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateDiagFunc: ToDiagFunc(validation.All(
														validation.StringLenBetween(1, 80),
														validation.StringMatch(regexp.MustCompile(`^([\w#][\w#.-]*)$`),
															"Composed of alphabets, numbers, hash (#), dot (.), hyphen (-) and underbar (_)"),
													)),
												},
											},
										},
									},
								},
							},
						},
						"docker_engine": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use": {
										Type:     schema.TypeBool,
										Optional: true,
									},
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
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ValidateDiagFunc: ToDiagFunc(validation.All(
								validation.IntBetween(5, 540),
							)),
						},
						"env_var": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
										ValidateDiagFunc: ToDiagFunc(validation.All(
											validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_]+$`), "Composed of alphabets, numbers and underbar (_)"),
										)),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"build_command": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pre_build": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"in_build": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"post_build": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"docker_image_build": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"dockerfile": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"registry": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"image": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: ToDiagFunc(validation.All(
											validation.StringLenBetween(1, 200),
											validation.StringMatch(regexp.MustCompile(`^[a-z0-9]+(([._]|__|[-]*)[a-z0-9]+)*$`),
												"Composed of alphabets(lowercase), numbers, dot (.), hyphen (-) and underbar (_)\n"+
													"The dot (.) cannot be used consecutively, and underbar (_) is allowed only twice in a row\n"+
													"Also, the dot (.), hyphen (-) and underbar (_) cannot be used for the last character"),
										)),
									},
									"tag": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: ToDiagFunc(validation.All(
											validation.StringLenBetween(1, 80),
											validation.StringMatch(regexp.MustCompile(`^([\w#][\w#.-]*)$`),
												"Composed of alphabets, numbers, hash (#), dot (.), hyphen (-) and underbar (_)"),
										)),
									},
									"latest": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"artifact": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"object_storage_to_upload": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"filename": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"backup": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"build_image_upload": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"container_registry_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateDiagFunc: ToDiagFunc(validation.All(
								validation.StringLenBetween(1, 200),
							)),
						},
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateDiagFunc: ToDiagFunc(validation.All(
								validation.StringLenBetween(1, 80),
								validation.StringMatch(regexp.MustCompile(`^([\w#][\w#.-]*)$`),
									"Composed of alphabets, numbers, hash (#), dot (.), hyphen (-) and underbar (_)"),
							)),
						},
						"latest": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"linked": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_log_analytics": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"file_safer": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"last_build": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"created": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNcloudSourceBuildProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	id, err := SourceBuildProjectCreate(d, config)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(id))))
	log.Printf("[INFO] Project ID: %s", d.Id())

	return resourceNcloudSourceBuildProjectRead(ctx, d, meta)
}

func SourceBuildProjectCreate(d *schema.ResourceData, config *ProviderConfig) (*int32, error) {
	commonParams, paramErr := getCommonProjectParams(d)
	if paramErr != nil {
		return nil, paramErr
	}

	reqParams := &sourcebuild.CreateProject{
		Name:        StringPtrOrNil(d.GetOk("name")),
		Description: commonParams.Description,
		Source:      commonParams.Source,
		Env:         commonParams.Env,
		Cmd:         commonParams.Cmd,
		Artifact:    commonParams.Artifact,
		Cache:       commonParams.Cache,
		Linked:      commonParams.Linked,
	}

	var resp *sourcebuild.CreateProjectResponse
	logCommonRequest("createSourceBuildProject", reqParams)
	resp, err := config.Client.sourcebuild.V1Api.CreateProject(context.Background(), reqParams)
	if err != nil {
		logErrorResponse("createSourceBuildProject", err, reqParams)
		return nil, err
	}
	logResponse("createSourceBuildProject", resp)

	return resp.Id, nil
}

func resourceNcloudSourceBuildProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	id := ncloud.String(d.Id())

	logCommonRequest("deleteSourceBuildProject", id)
	err := config.Client.sourcebuild.V1Api.DeleteProject(ctx, id)
	if err != nil {
		logErrorResponse("deleteSourceBuildProject", err, id)
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceNcloudSourceBuildProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	err := changeBuildProject(ctx, d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNcloudSourceBuildProjectRead(ctx, d, meta)
}

func changeBuildProject(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	reqParams, paramErr := getCommonProjectParams(d)
	if paramErr != nil {
		return paramErr
	}

	id := ncloud.String(d.Id())

	var resp *sourcebuild.CreateProjectResponse
	logCommonRequest("updateSourceBuildProject", reqParams)
	resp, err := config.Client.sourcebuild.V1Api.ChangeProject(ctx, reqParams, id)
	if err != nil {
		logErrorResponse("updateSourceBuildProject", err, id)
		return err
	}
	logResponse("updateSourceBuildProject", resp)

	return nil
}

func getEnvPlatformConfig(d *schema.ResourceData, platformType string) (*sourcebuild.EnvPlatformConfigRequest, error) {
	envPlatformConfigOs := sourcebuild.EnvPlatformConfigRequestOs{
		Id: Int32PtrOrNil(d.GetOk("env.0.platform.0.config.0.os.0.id")),
	}

	envPlatformConfigRuntimeVersion := sourcebuild.EnvPlatformConfigRequestRuntimeVersion{
		Id: Int32PtrOrNil(d.GetOk("env.0.platform.0.config.0.runtime.0.version.0.id")),
	}

	envPlatformConfigRuntime := sourcebuild.EnvPlatformConfigRequestRuntime{
		Id:      Int32PtrOrNil(d.GetOk("env.0.platform.0.config.0.runtime.0.id")),
		Version: &envPlatformConfigRuntimeVersion,
	}

	envPlatformConfig := sourcebuild.EnvPlatformConfigRequest{
		Os:       &envPlatformConfigOs,
		Runtime:  &envPlatformConfigRuntime,
		Registry: StringPtrOrNil(d.GetOk("env.0.platform.0.config.0.registry")),
		Image:    StringPtrOrNil(d.GetOk("env.0.platform.0.config.0.image")),
		Tag:      StringPtrOrNil(d.GetOk("env.0.platform.0.config.0.tag")),
	}

	switch platformType {
	case "SourceBuild":
		if envPlatformConfigOs.Id == nil || envPlatformConfigRuntime.Id == nil || envPlatformConfigRuntimeVersion.Id == nil {
			return nil, fmt.Errorf("env.platform.config(os.id, runtime.id, runtime.version.id) is required")
		}
		if envPlatformConfig.Registry != nil || envPlatformConfig.Image != nil || envPlatformConfig.Tag != nil {
			return nil, fmt.Errorf("env.platform.config requires only os.id, runtime.id, runtime.version.id")
		}
	case "ContainerRegistry":
		if envPlatformConfig.Registry == nil || envPlatformConfig.Image == nil || envPlatformConfig.Tag == nil {
			return nil, fmt.Errorf("env.platform.config(registry, image, tag) is required")
		}
		if envPlatformConfigOs.Id != nil || envPlatformConfigRuntime.Id != nil || envPlatformConfigRuntimeVersion.Id != nil {
			return nil, fmt.Errorf("env.platform.config requires only registry, image, tag")
		}
	case "PublicRegistry":
		if envPlatformConfig.Image == nil || envPlatformConfig.Tag == nil {
			return nil, fmt.Errorf("env.platform.config(image, tag) is required")
		}
		if envPlatformConfigOs.Id != nil || envPlatformConfigRuntime.Id != nil || envPlatformConfigRuntimeVersion.Id != nil || envPlatformConfig.Registry != nil {
			return nil, fmt.Errorf("env.platform.config requires only image, tag")
		}
	}

	return &envPlatformConfig, nil
}

func getCommonProjectParams(d *schema.ResourceData) (*sourcebuild.ChangeProject, error) {
	sourceConfig := sourcebuild.ProjectSourceConfig{
		Repository: StringPtrOrNil(d.GetOk("source.0.config.0.repository_name")),
		Branch:     StringPtrOrNil(d.GetOk("source.0.config.0.branch")),
	}

	source := sourcebuild.ProjectSource{
		Type_:  StringPtrOrNil(d.GetOk("source.0.type")),
		Config: &sourceConfig,
	}

	envCompute := sourcebuild.ProjectEnvCompute{
		Id: Int32PtrOrNil(d.GetOk("env.0.compute.0.id")),
	}

	envPlatform := sourcebuild.ProjectEnvPlatform{
		Type_: StringPtrOrNil(d.GetOk("env.0.platform.0.type")),
	}

	envPlatformConfig, envPlatformConfigErr := getEnvPlatformConfig(d, *envPlatform.Type_)
	if envPlatformConfigErr != nil {
		return nil, envPlatformConfigErr
	}
	envPlatform.Config = envPlatformConfig

	envDocker := sourcebuild.ProjectEnvDocker{
		Use: ncloud.Bool(d.Get("env.0.docker_engine.0.use").(bool)),
		Id:  Int32PtrOrNil(d.GetOk("env.0.docker_engine.0.id")),
	}
	if *envDocker.Use && envDocker.Id == nil {
		return nil, fmt.Errorf("env.docker_engine.id is required if env.docker_engine.use is true")
	}
	if !*envDocker.Use && envDocker.Id != nil {
		return nil, fmt.Errorf("env.docker_engine.id must not be defined if env.docker_engine.use is false")
	}

	envVars, envErr := expandSourceBuildEnvVarsParams(d.Get("env.0.env_var").([]interface{}))
	if envErr != nil {
		return nil, envErr
	}

	env := sourcebuild.ProjectEnv{
		Compute:  &envCompute,
		Platform: &envPlatform,
		Docker:   &envDocker,
		Timeout:  Int32PtrOrNil(d.GetOk("env.0.timeout")),
		EnvVars:  envVars,
	}

	cmdDockerbuild := sourcebuild.ProjectCmdDockerbuild{
		Use:        ncloud.Bool(d.Get("build_command.0.docker_image_build.0.use").(bool)),
		Dockerfile: StringPtrOrNil(d.GetOk("build_command.0.docker_image_build.0.dockerfile")),
		Registry:   StringPtrOrNil(d.GetOk("build_command.0.docker_image_build.0.registry")),
		Image:      StringPtrOrNil(d.GetOk("build_command.0.docker_image_build.0.image")),
		Tag:        StringPtrOrNil(d.GetOk("build_command.0.docker_image_build.0.tag")),
		Latest:     ncloud.Bool(d.Get("build_command.0.docker_image_build.0.latest").(bool)),
	}

	if *cmdDockerbuild.Use && (cmdDockerbuild.Dockerfile == nil || cmdDockerbuild.Registry == nil || cmdDockerbuild.Image == nil || cmdDockerbuild.Tag == nil) {
		return nil, fmt.Errorf("build_command.docker_image_build parameters(dockerfile, registry, image, tag) are required if build_command.docker_image_build.use is true")
	}
	if !*cmdDockerbuild.Use && (cmdDockerbuild.Dockerfile != nil || cmdDockerbuild.Registry != nil || cmdDockerbuild.Image != nil || cmdDockerbuild.Tag != nil || *cmdDockerbuild.Latest) {
		return nil, fmt.Errorf("build_command.docker_image_build parameters(dockerfile, registry, image, tag, latest) must not be defined if build_command.docker_image_build.use is false")
	}

	cmd := sourcebuild.ProjectCmd{
		Dockerbuild: &cmdDockerbuild,
	}

	if param, ok := d.GetOk("build_command.0.pre_build"); ok {
		cmd.Pre = expandStringInterfaceList(param.([]interface{}))
	}

	if param, ok := d.GetOk("build_command.0.in_build"); ok {
		cmd.Build = expandStringInterfaceList(param.([]interface{}))
	}

	if param, ok := d.GetOk("build_command.0.post_build"); ok {
		cmd.Post = expandStringInterfaceList(param.([]interface{}))
	}

	artifactStorage := sourcebuild.ProjectArtifactStorage{
		Bucket:   StringPtrOrNil(d.GetOk("artifact.0.object_storage_to_upload.0.bucket")),
		Path:     StringPtrOrNil(d.GetOk("artifact.0.object_storage_to_upload.0.path")),
		Filename: StringPtrOrNil(d.GetOk("artifact.0.object_storage_to_upload.0.filename")),
	}

	artifact := sourcebuild.ProjectArtifact{
		Use:     ncloud.Bool(d.Get("artifact.0.use").(bool)),
		Storage: &artifactStorage,
		Backup:  ncloud.Bool(d.Get("artifact.0.backup").(bool)),
	}

	if param, ok := d.GetOk("artifact.0.path"); ok {
		artifact.Path = expandStringInterfaceList(param.([]interface{}))
	}

	if *artifact.Use && (len(artifact.Path) == 0 || artifactStorage.Bucket == nil || artifactStorage.Path == nil || artifactStorage.Filename == nil) {
		return nil, fmt.Errorf("artifact.path and artifact.object_storage_to_upload parameters(bucket, path, filename) are required if artifact.use is true")
	}
	if !*artifact.Use && (artifact.Path != nil || artifactStorage.Bucket != nil || artifactStorage.Path != nil || artifactStorage.Filename != nil || *artifact.Backup) {
		return nil, fmt.Errorf("artifact.path, artifact.backup and artifact.object_storage_to_upload parameters(bucket, path, filename) must not be defined if artifact.use is false")
	}

	cache := sourcebuild.ProjectCache{
		Use:      ncloud.Bool(d.Get("build_image_upload.0.use").(bool)),
		Registry: StringPtrOrNil(d.GetOk("build_image_upload.0.container_registry_name")),
		Image:    StringPtrOrNil(d.GetOk("build_image_upload.0.image_name")),
		Tag:      StringPtrOrNil(d.GetOk("build_image_upload.0.tag")),
		Latest:   ncloud.Bool(d.Get("build_image_upload.0.latest").(bool)),
	}

	if *cache.Use && (cache.Registry == nil || cache.Image == nil || cache.Tag == nil) {
		return nil, fmt.Errorf("build_image_upload parameters(container_registry_name, image_name, tag) are required if build_image_upload.use is true")
	}
	if !*cache.Use && (cache.Registry != nil || cache.Image != nil || cache.Tag != nil || *cache.Latest) {
		return nil, fmt.Errorf("build_image_upload parameters(container_registry_name, image_name, tag, latest) must not be defined if build_image_upload.use is false")
	}

	linked := sourcebuild.ProjectLinked{
		CloudLogAnalytics: ncloud.Bool(d.Get("linked.0.cloud_log_analytics").(bool)),
		FileSafer:         ncloud.Bool(d.Get("linked.0.file_safer").(bool)),
	}

	return &sourcebuild.ChangeProject{
		Description: StringPtrOrNil(d.GetOk("description")),
		Source:      &source,
		Env:         &env,
		Cmd:         &cmd,
		Artifact:    &artifact,
		Cache:       &cache,
		Linked:      &linked,
	}, nil
}

func resourceNcloudSourceBuildProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)

	project, err := getBuildProject(ctx, config, ncloud.String(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}
	if project == nil {
		d.SetId("")
		return nil
	}

	setProjectData(d, project)

	return nil
}

func getBuildProject(ctx context.Context, config *ProviderConfig, id *string) (*sourcebuild.GetProjectDetailResponse, error) {
	logCommonRequest("getSourceBuildProjectDetail", id)
	resp, err := config.Client.sourcebuild.V1Api.GetProject(ctx, id)

	if err != nil {
		logErrorResponse("getSourceBuildProjectDetail", err, id)
		return nil, err
	}

	logResponse("getSourceBuildProjectDetail", resp)

	return convertBuildProject(resp), nil
}

func convertBuildProject(r *sourcebuild.GetProjectDetailResponse) *sourcebuild.GetProjectDetailResponse {
	if r == nil {
		return nil
	}

	return &sourcebuild.GetProjectDetailResponse{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Source:      r.Source,
		Env:         r.Env,
		Cmd:         r.Cmd,
		Artifact:    r.Artifact,
		Cache:       r.Cache,
		Linked:      r.Linked,
		LastBuild:   r.LastBuild,
		Created:     r.Created,
	}
}

func setProjectData(d *schema.ResourceData, project *sourcebuild.GetProjectDetailResponse) {
	d.SetId(*ncloud.IntString(int(ncloud.Int32Value(project.Id))))
	d.Set("name", ncloud.StringValue(project.Name))
	d.Set("description", ncloud.StringValue(project.Description))
	d.Set("source", makeSource(project.Source))
	d.Set("artifact", makeArtifact(project.Artifact))
	d.Set("build_image_upload", makeBuildImageUpload(project.Cache))
	d.Set("build_command", makeCmd(project.Cmd))
	d.Set("env", makeEnv(project.Env))
	d.Set("linked", makeLinked(project.Linked))
	d.Set("last_build", makeLastBuild(project.LastBuild))
	d.Set("created", makeCreated(project.Created))
	d.Set("project_no", ncloud.Int32Value(project.Id))
}

/* for deep copy */
func makeSource(source *sourcebuild.GetProjectDetailResponseSource) []interface{} {
	if source == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["type"] = ncloud.StringValue(source.Type_)
	values["config"] = makeSourceConfig(source.Config)

	return []interface{}{values}
}

func makeSourceConfig(config *sourcebuild.GetProjectDetailResponseSourceConfig) []interface{} {
	if config == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["repository_name"] = ncloud.StringValue(config.Repository)
	values["branch"] = ncloud.StringValue(config.Branch)

	return []interface{}{values}
}

func makeArtifact(artifact *sourcebuild.GetProjectDetailResponseArtifact) []interface{} {
	if artifact == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["use"] = ncloud.BoolValue(artifact.Use)
	values["path"] = ncloud.StringListValue(artifact.Path)
	values["object_storage_to_upload"] = makeArtifactStorage(artifact.Storage)
	values["backup"] = ncloud.BoolValue(artifact.Backup)

	return []interface{}{values}
}

func makeArtifactStorage(storage *sourcebuild.GetProjectDetailResponseArtifactStorage) []interface{} {
	if storage == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["bucket"] = ncloud.StringValue(storage.Bucket)
	values["filename"] = ncloud.StringValue(storage.Filename)
	values["path"] = ncloud.StringValue(storage.Path)

	return []interface{}{values}
}

func makeBuildImageUpload(cache *sourcebuild.ProjectCache) []interface{} {
	if cache == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["use"] = ncloud.BoolValue(cache.Use)
	values["container_registry_name"] = ncloud.StringValue(cache.Registry)
	values["image_name"] = ncloud.StringValue(cache.Image)
	values["tag"] = ncloud.StringValue(cache.Tag)
	values["latest"] = ncloud.BoolValue(cache.Latest)

	return []interface{}{values}
}

func makeCmd(cmd *sourcebuild.ProjectCmd) []interface{} {
	if cmd == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["pre_build"] = ncloud.StringListValue(cmd.Pre)
	values["in_build"] = ncloud.StringListValue(cmd.Build)
	values["post_build"] = ncloud.StringListValue(cmd.Post)
	values["docker_image_build"] = makeCmdDockerbuild(cmd.Dockerbuild)

	return []interface{}{values}
}

func makeCmdDockerbuild(dockerbuild *sourcebuild.ProjectCmdDockerbuild) []interface{} {
	if dockerbuild == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["use"] = ncloud.BoolValue(dockerbuild.Use)
	values["dockerfile"] = ncloud.StringValue(dockerbuild.Dockerfile)
	values["registry"] = ncloud.StringValue(dockerbuild.Registry)
	values["image"] = ncloud.StringValue(dockerbuild.Image)
	values["tag"] = ncloud.StringValue(dockerbuild.Tag)
	values["latest"] = ncloud.BoolValue(dockerbuild.Latest)

	return []interface{}{values}
}

func makeEnv(env *sourcebuild.GetProjectDetailResponseEnv) []interface{} {
	if env == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["compute"] = makeEnvCompute(env.Compute)
	values["platform"] = makeEnvPlatform(env.Platform)
	values["docker_engine"] = makeEnvDocker(env.Docker)
	values["timeout"] = ncloud.Int32Value(env.Timeout)
	values["env_var"] = makeEnvEnvVars(env.EnvVars)

	return []interface{}{values}
}

func makeEnvCompute(compute *sourcebuild.GetProjectDetailResponseEnvCompute) []interface{} {
	if compute == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["id"] = ncloud.Int32Value(compute.Id)
	values["cpu"] = ncloud.Int32Value(compute.Cpu)
	values["mem"] = ncloud.Int32Value(compute.Mem)

	return []interface{}{values}
}

func makeEnvPlatform(platform *sourcebuild.GetProjectDetailResponseEnvPlatform) []interface{} {
	if platform == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["type"] = ncloud.StringValue(platform.Type_)
	values["config"] = makeEnvPlatformConfig(platform.Config)

	return []interface{}{values}
}

func makeEnvPlatformConfig(config *sourcebuild.EnvPlatformConfigResponse) []interface{} {
	if config == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["registry"] = ncloud.StringValue(config.Registry)
	values["image"] = ncloud.StringValue(config.Image)
	values["tag"] = ncloud.StringValue(config.Tag)
	values["runtime"] = makeEnvPlatformConfigRuntime(config.Runtime)
	values["os"] = makeEnvPlatformConfigOs(config.Os)

	return []interface{}{values}
}

func makeEnvPlatformConfigRuntime(runtime *sourcebuild.EnvPlatformConfigResponseRuntime) []interface{} {
	if runtime == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["id"] = ncloud.Int32Value(runtime.Id)
	values["name"] = ncloud.StringValue(runtime.Name)
	values["version"] = makeEnvPlatformConfigRuntimeVersion(runtime.Version)

	return []interface{}{values}
}

func makeEnvPlatformConfigRuntimeVersion(version *sourcebuild.EnvPlatformConfigResponseRuntimeVersion) []interface{} {
	if version == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["id"] = ncloud.Int32Value(version.Id)
	values["name"] = ncloud.StringValue(version.Name)

	return []interface{}{values}
}

func makeEnvPlatformConfigOs(os *sourcebuild.EnvPlatformConfigResponseOs) []interface{} {
	if os == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["id"] = ncloud.Int32Value(os.Id)
	values["name"] = ncloud.StringValue(os.Name)
	values["version"] = ncloud.StringValue(os.Version)
	values["archi"] = ncloud.StringValue(os.Archi)

	return []interface{}{values}
}

func makeEnvDocker(docker *sourcebuild.GetProjectDetailResponseEnvDocker) []interface{} {
	if docker == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["use"] = ncloud.BoolValue(docker.Use)
	values["id"] = ncloud.Int32Value(docker.Id)
	values["name"] = ncloud.StringValue(docker.Name)

	return []interface{}{values}
}

func makeEnvEnvVars(envVars []*sourcebuild.ProjectEnvEnvVars) []interface{} {
	if envVars == nil {
		return nil
	}

	values := []interface{}{}
	for _, v := range envVars {
		values = append(values, makeEnvVar(v))
	}

	return values
}

func makeEnvVar(envVar *sourcebuild.ProjectEnvEnvVars) map[string]interface{} {
	if envVar == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["key"] = ncloud.StringValue(envVar.Key)
	values["value"] = ncloud.StringValue(envVar.Value)

	return values
}

func makeLinked(linked *sourcebuild.ProjectLinked) []interface{} {
	if linked == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["cloud_log_analytics"] = ncloud.BoolValue(linked.CloudLogAnalytics)
	values["file_safer"] = ncloud.BoolValue(linked.FileSafer)

	return []interface{}{values}
}

func makeLastBuild(linked *sourcebuild.GetProjectDetailResponseLastBuild) []interface{} {
	if linked == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["id"] = ncloud.StringValue(linked.Id)
	values["status"] = ncloud.StringValue(linked.Status)
	values["timestamp"] = ncloud.Int64Value(linked.Timestamp)

	return []interface{}{values}
}

func makeCreated(linked *sourcebuild.GetProjectDetailResponseCreated) []interface{} {
	if linked == nil {
		return nil
	}

	values := map[string]interface{}{}
	values["user"] = ncloud.StringValue(linked.User)
	values["timestamp"] = ncloud.Int64Value(linked.Timestamp)

	return []interface{}{values}
}

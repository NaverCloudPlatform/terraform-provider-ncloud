package hadoop

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &hadoopResource{}
	_ resource.ResourceWithConfigure   = &hadoopResource{}
	_ resource.ResourceWithImportState = &hadoopResource{}
)

func NewHadoopResource() resource.Resource {
	return &hadoopResource{}
}

type hadoopResource struct {
	config *conn.ProviderConfig
}

func (h *hadoopResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (h *hadoopResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	h.config = config
}

func (h *hadoopResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop"
}

func (h *hadoopResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_product_code": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "default: latest version",
			},
			"master_node_product_code": schema.StringAttribute{
				Optional:    true,
				Description: "default: minimum spec",
			},
			"edge_node_product_code": schema.StringAttribute{
				Optional:    true,
				Description: "default: minimum spec",
			},
			"worker_node_product_code": schema.StringAttribute{
				Optional:    true,
				Description: "default: minimum spec",
			},
			"cluster_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(3, 15),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[a-z0-9]+[ㄱ-ㅣ가-힣A-Za-z0-9-]+[a-z0-9]$`),
							"Composed of alphabets, korean, numbers, hyphen (-). Must start with an alphabetic character and number, must end with an alphabetic character or number",
						),
					),
				},
			},
			"cluster_type_code": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Computed: true,
			},
			"add_on_code_list": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("PRESTO", "HBASE", "IMPALA", "KUDU"),
					),
				},
				Description: "this attribute can used over 1.5 version",
			},
			"admin_user_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(3, 15),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[a-z0-9]+[a-z0-9-]+[a-z0-9]$`),
							"Composed of lowercase alphabets, numbers, hyphen (-). Must start with an lowercase alphabetic character or number, must end with an lowercase alphabetic character or number",
						),
					),
				},
			},
			"admin_user_password": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(8, 20),
						stringvalidator.RegexMatches(regexp.MustCompile(`[A-Z]+`), "Must have at least one uppercase alphabet"),
						stringvalidator.RegexMatches(regexp.MustCompile(`\d+`), "Must have at least one number"),
						stringvalidator.RegexMatches(regexp.MustCompile(`[~!@#$%^*()\-_=\[\]\{\};:,.<>?]+`), "Must have at least one special character"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[^&+\\"'/\s`+"`"+`]*$`), "Must not have ` & + \\ \" ' / and white space."),
					),
				},
				Sensitive: true,
			},
			"login_key_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"edge_node_subnet_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"master_node_subnet_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"worker_node_subnet_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"master_node_data_storage_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("SSD", "HDD"),
				},
				Description: "default: SSD",
			},
			"worker_node_data_storage_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("SSD", "HDD"),
				},
				Description: "default: SSD",
			},
			"master_node_data_storage_size": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(100, 2000),
						int64validator.OneOf(4000, 6000),
					),
				},
			},
			"worker_node_data_storage_size": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(100, 2000),
						int64validator.OneOf(4000, 6000),
					),
				},
			},
			"worker_node_count": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(2, 8),
				},
				Description: "default: 2",
			},
			"use_kdc": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: false",
			},
			"kdc_realm": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("use_kdc"),
					}...),
				},
			},
			"kdc_password": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("use_kdc"),
					}...),
				},
			},
			"use_bootstrap_script": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: false",
			},
			"bootstrap_script": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("use_bootstrap_script"),
					}...),
					stringvalidator.LengthAtMost(1024),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]$`),
						"Composed of alphabets.",
					),
				},
			},
			"use_data_catalog": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "this attribute can used over 2.0 version",
			},
			"ambari_server_host": schema.StringAttribute{
				Computed: true,
			},
			"cluster_direct_access_account": schema.StringAttribute{
				Computed: true,
			},
			"is_ha": schema.BoolAttribute{
				Computed: true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"hadoop_server_instance_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hadoop_server_name": schema.StringAttribute{
							Computed: true,
						},
						"hadoop_server_role": schema.StringAttribute{
							Computed: true,
						},
						"hadoop_product_code": schema.StringAttribute{
							Computed: true,
						},
						"region_code": schema.StringAttribute{
							Computed: true,
						},
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"vpc_no": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"is_public_subnet": schema.BoolAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"id": framework.IDAttribute(),
		},
	}
}

func (h *hadoopResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hadoopResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !h.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.Config.Schema.Type().String()),
		)
		return
	}

	reqParams := &vhadoop.CreateCloudHadoopInstanceRequest{
		RegionCode:                    &h.config.RegionCode,
		VpcNo:                         plan.VpcNo.ValueStringPointer(),
		CloudHadoopClusterName:        plan.ClusterName.ValueStringPointer(),
		CloudHadoopClusterTypeCode:    plan.ClusterTypeCode.ValueStringPointer(),
		CloudHadoopAdminUserName:      plan.AdminUserName.ValueStringPointer(),
		CloudHadoopAdminUserPassword:  plan.AdminUserPassword.ValueStringPointer(),
		LoginKeyName:                  plan.LoginKeyName.ValueStringPointer(),
		EdgeNodeSubnetNo:              plan.EdgeNodeSubnetNo.ValueStringPointer(),
		MasterNodeSubnetNo:            plan.MasterNodeSubnetNo.ValueStringPointer(),
		BucketName:                    plan.BucketName.ValueStringPointer(),
		WorkerNodeSubnetNo:            plan.WorkerNodeSubnetNo.ValueStringPointer(),
		MasterNodeDataStorageTypeCode: plan.MasterNodeDataStorageType.ValueStringPointer(),
		WorkerNodeDataStorageTypeCode: plan.WorkerNodeDataStorageType.ValueStringPointer(),
		MasterNodeDataStorageSize:     ncloud.Int32(int32(plan.MasterNodeDataStorageSize.ValueInt64())),
		WorkerNodeDataStorageSize:     ncloud.Int32(int32(plan.WorkerNodeDataStorageSize.ValueInt64())),
	}

	if !plan.ImageProductCode.IsNull() && !plan.ImageProductCode.IsUnknown() {
		reqParams.CloudHadoopImageProductCode = plan.ImageProductCode.ValueStringPointer()
	}

	if !plan.MasterNodeProductCode.IsNull() && !plan.MasterNodeProductCode.IsUnknown() {
		reqParams.MasterNodeProductCode = plan.MasterNodeProductCode.ValueStringPointer()
	}

	if !plan.EdgeNodeProductCode.IsNull() && !plan.EdgeNodeProductCode.IsUnknown() {
		reqParams.EdgeNodeProductCode = plan.EdgeNodeProductCode.ValueStringPointer()
	}

	if !plan.WorkerNodeProductCode.IsNull() && !plan.WorkerNodeProductCode.IsUnknown() {
		reqParams.WorkerNodeProductCode = plan.EdgeNodeProductCode.ValueStringPointer()
	}

	if !plan.AddOnCodeList.IsNull() && !plan.AddOnCodeList.IsUnknown() {
		addOnList := plan.AddOnCodeList.Elements()
		var addOnListReq []*string
		for _, addon := range addOnList {
			addOnListReq = append(addOnListReq, ncloud.String(addon.String()))
		}
		reqParams.CloudHadoopAddOnCodeList = addOnListReq
	}

	if !plan.WorkerNodeCount.IsNull() && !plan.WorkerNodeCount.IsUnknown() {
		reqParams.WorkerNodeCount = ncloud.Int32(int32(plan.WorkerNodeCount.ValueInt64()))
	}

	if !plan.UseKdc.IsNull() {
		reqParams.UseKdc = plan.UseKdc.ValueBoolPointer()
	}

	if !plan.UseKdc.IsNull() && plan.UseKdc.ValueBool() {
		if !plan.KdcRealm.IsNull() && !plan.KdcRealm.IsUnknown() {
			reqParams.KdcRealm = plan.KdcRealm.ValueStringPointer()
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("`use_kdc` = %t, `kdc_realm` is nil", plan.UseKdc.ValueBool()),
				errors.New("when `use_kdc` is true, `kdc_realm` must be inputted`").Error(),
			)
			return
		}

		if !plan.KdcPassword.IsNull() {
			reqParams.KdcPassword = plan.KdcPassword.ValueStringPointer()
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("`use_kdc` = %t, `Kdc_password` is nil", plan.UseKdc.ValueBool()),
				errors.New("when `use_kdc` is true, `kdc_realm` must be entered`").Error(),
			)
			return
		}
	}

	if !plan.UseBootstrapScript.IsNull() {
		reqParams.UseBootstrapScript = plan.UseBootstrapScript.ValueBoolPointer()
	}

	if !plan.UseBootstrapScript.IsNull() && plan.UseBootstrapScript.ValueBool() {
		if !plan.BootstrapScript.IsNull() {
			reqParams.BootstrapScript = plan.BootstrapScript.ValueStringPointer()
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("`use_bootstrap_script` = %t, `bootstrap_script` is nil", plan.UseBootstrapScript.ValueBool()),
				errors.New("when `use_bootstrap_script` is true, `bootstrap_script` must be entered`").Error(),
			)
			return
		}
	}

	tflog.Info(ctx, "CreateHadoop", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := h.config.Client.Vhadoop.V2Api.CreateCloudHadoopInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Hadoop Instance, err params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "CreateHadoop response", map[string]any{
		"createHadoopResponse": common.MarshalUncheckedString(response),
	})

	hadoopInstance := response.CloudHadoopInstanceList[0]
	plan.ID = types.StringPointerValue(hadoopInstance.CloudHadoopInstanceNo)
	tflog.Info(ctx, "Hadoop ID", map[string]any{
		"HadoopNo": plan.ID,
	})
	output, err := waitHadoopForCreation(ctx, h.config, *hadoopInstance.CloudHadoopInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("waiting for Hadoop creation", err.Error())
		return
	}

	if err := plan.refreshFromOutput(ctx, output); err.ErrorsCount() > 0 {
		resp.Diagnostics.Append(err...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (h *hadoopResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state hadoopResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetHadoopInstance(ctx, h.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetHadoop", err.Error())
		return
	}
	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if errors := state.refreshFromOutput(ctx, output); errors.ErrorsCount() > 0 {
		resp.Diagnostics.Append(errors...)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (h *hadoopResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state hadoopResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.WorkerNodeCount.Equal(state.WorkerNodeCount) {

		reqParams := &vhadoop.ChangeCloudHadoopNodeCountRequest{
			RegionCode:            &h.config.RegionCode,
			CloudHadoopInstanceNo: plan.ID.ValueStringPointer(),
			WorkerNodeCount:       ncloud.Int32(int32(plan.WorkerNodeCount.ValueInt64())),
		}

		tflog.Info(ctx, "ChangeHadoopWorkerNodeCount", map[string]any{
			"reqParams": common.MarshalUncheckedString(reqParams),
		})
		response, err := h.config.Client.Vhadoop.V2Api.ChangeCloudHadoopNodeCount(reqParams)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ChangeHadoopWorkerNodeCount params=%v", *reqParams),
				err.Error(),
			)
			return
		}

		tflog.Info(ctx, "ChangeHadoopWorkerNodeCount", map[string]any{
			"changeHadoopWorkerNodeResponse": common.MarshalUncheckedString(response),
		})

		output, err := waitForHadoopUpdate(ctx, h.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"fail to wait for updating hadoop worker node count",
				err.Error(),
			)
		}

		if err := state.refreshFromOutput(ctx, output); err.HasError() {
			resp.Diagnostics.Append(err...)
		}
		state.WorkerNodeCount = plan.WorkerNodeCount
	}

	if !plan.MasterNodeProductCode.Equal(state.MasterNodeProductCode) ||
		!plan.EdgeNodeProductCode.Equal(state.EdgeNodeProductCode) ||
		!plan.WorkerNodeProductCode.Equal(state.WorkerNodeProductCode) {

		reqParams := &vhadoop.ChangeCloudHadoopNodeSpecRequest{
			RegionCode:            &h.config.RegionCode,
			CloudHadoopInstanceNo: plan.ID.ValueStringPointer(),

			MasterNodeProductCode: state.MasterNodeProductCode.ValueStringPointer(),
			EdgeNodeProductCode:   state.EdgeNodeProductCode.ValueStringPointer(),
			WorkerNodeProductCode: state.WorkerNodeProductCode.ValueStringPointer(),
		}

		if !plan.MasterNodeProductCode.Equal(state.MasterNodeProductCode) {
			reqParams.MasterNodeProductCode = plan.MasterNodeProductCode.ValueStringPointer()
			state.MasterNodeProductCode = types.StringPointerValue(plan.MasterNodeProductCode.ValueStringPointer())
		}

		if !plan.EdgeNodeProductCode.Equal(state.EdgeNodeProductCode) {
			reqParams.EdgeNodeProductCode = plan.EdgeNodeProductCode.ValueStringPointer()
			state.EdgeNodeProductCode = types.StringPointerValue(plan.EdgeNodeProductCode.ValueStringPointer())
		}

		if !plan.WorkerNodeProductCode.Equal(state.WorkerNodeProductCode) {
			reqParams.WorkerNodeProductCode = plan.WorkerNodeProductCode.ValueStringPointer()
			state.WorkerNodeProductCode = types.StringPointerValue(plan.WorkerNodeProductCode.ValueStringPointer())
		}

		tflog.Info(ctx, "ChangeHadoopSpec", map[string]any{
			"reqParams": common.MarshalUncheckedString(reqParams),
		})

		response, err := h.config.Client.Vhadoop.V2Api.ChangeCloudHadoopNodeSpec(reqParams)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ChangeHadoopSpec params=%v", *reqParams),
				err.Error(),
			)
			return
		}

		tflog.Info(ctx, "ChangeHadoopSpec", map[string]any{
			"changeHadoopSpecResponse": common.MarshalUncheckedString(response),
		})

		output, err := waitForHadoopUpdate(ctx, h.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"fail to wait for updating hadoop spec",
				err.Error(),
			)
		}

		if err := state.refreshFromOutput(ctx, output); err.HasError() {
			resp.Diagnostics.Append(err...)
		}

	}

	resp.State.Set(ctx, state)
}

func (h *hadoopResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state hadoopResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqParams := &vhadoop.DeleteCloudHadoopInstanceRequest{
		RegionCode:            &h.config.RegionCode,
		CloudHadoopInstanceNo: state.ID.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteHadoop", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	response, err := h.config.Client.Vhadoop.V2Api.DeleteCloudHadoopInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete Hadoop Instance params=%v", *reqParams),
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "DeleteCloudHadoop response", map[string]any{
		"deleteHadoopResponse": common.MarshalUncheckedString(response),
	})

	if err := waitForHadoopDeletion(ctx, h.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail to wait for hadoop deletion",
			err.Error(),
		)
	}
}

func GetHadoopInstance(ctx context.Context, config *conn.ProviderConfig, id string) (*vhadoop.CloudHadoopInstance, error) {
	reqParams := &vhadoop.GetCloudHadoopInstanceDetailRequest{
		RegionCode:            &config.RegionCode,
		CloudHadoopInstanceNo: ncloud.String(id),
	}

	tflog.Info(ctx, "GetHadoop", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vhadoop.V2Api.GetCloudHadoopInstanceDetail(reqParams)
	if err != nil {
		if strings.Contains(err.Error(), "Unable to lookup cluster instance information") {
			return nil, nil
		}
		return nil, err
	}

	tflog.Info(ctx, "GetHadoop response", map[string]any{
		"getHadoopResponse": common.MarshalUncheckedString(resp),
	})

	if len(resp.CloudHadoopInstanceList) > 0 {
		return resp.CloudHadoopInstanceList[0], nil
	}
	return nil, nil
}

func waitHadoopForCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vhadoop.CloudHadoopInstance, error) {
	var hadoopInstance *vhadoop.CloudHadoopInstance
	var err error
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"CREAT"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			if hadoopInstance, err = GetHadoopInstance(ctx, config, id); err != nil {
				return 0, "", err
			}

			status := hadoopInstance.CloudHadoopInstanceStatus.Code
			op := hadoopInstance.CloudHadoopInstanceOperation.Code

			if *status == "INIT" && *op == "CREAT" {
				return hadoopInstance, "CREAT", nil
			}
			if *status == "CREAT" && *op == "SETUP" {
				return hadoopInstance, "CREAT", nil
			}
			if *status == "CREAT" && *op == "NULL" {
				return hadoopInstance, "RUN", nil
			}
			return 0, "", fmt.Errorf("")
		},
		Timeout:    15 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, err
	}

	return hadoopInstance, nil
}

func waitForHadoopUpdate(ctx context.Context, config *conn.ProviderConfig, id string) (*vhadoop.CloudHadoopInstance, error) {
	var hadoopInstance *vhadoop.CloudHadoopInstance
	var err error

	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"SET", "UPGD"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			hadoopInstance, err = GetHadoopInstance(ctx, config, id)
			if err != nil {
				return 0, "", err
			}

			status := *hadoopInstance.CloudHadoopInstanceStatus.Code
			op := *hadoopInstance.CloudHadoopInstanceOperation.Code
			if status == "CREAT" && op == "SETUP" {
				return hadoopInstance, "SET", nil
			}
			if status == "CREAT" && op == "UPGD" {
				return hadoopInstance, "UPGD", nil
			}
			if status == "CREAT" && op == "NULL" {
				return hadoopInstance, "RUN", nil
			}

			return 0, "", fmt.Errorf("")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, err
	}

	return hadoopInstance, nil
}

func waitForHadoopDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"PEND"},
		Target:  []string{"DEL"},
		Refresh: func() (interface{}, string, error) {
			hadoopInstance, err := GetHadoopInstance(ctx, config, id)
			if err != nil && !strings.Contains(err.Error(), `"returnCode": "5001017"`) {
				return 0, "", err
			}
			if hadoopInstance == nil {
				return hadoopInstance, "DEL", nil
			}

			status := *hadoopInstance.CloudHadoopInstanceStatus.Code
			op := *hadoopInstance.CloudHadoopInstanceOperation.Code

			if status == "DEL" && op == "DEL" {
				return hadoopInstance, "PEND", nil
			}

			return 0, "", fmt.Errorf("")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}

type hadoopResourceModel struct {
	VpcNo                      types.String `tfsdk:"vpc_no"`
	ImageProductCode           types.String `tfsdk:"image_product_code"`
	MasterNodeProductCode      types.String `tfsdk:"master_node_product_code"`
	EdgeNodeProductCode        types.String `tfsdk:"edge_node_product_code"`
	WorkerNodeProductCode      types.String `tfsdk:"worker_node_product_code"`
	ClusterName                types.String `tfsdk:"cluster_name"`
	Version                    types.String `tfsdk:"version"`
	ClusterTypeCode            types.String `tfsdk:"cluster_type_code"`
	AddOnCodeList              types.List   `tfsdk:"add_on_code_list"`
	AdminUserName              types.String `tfsdk:"admin_user_name"`
	AdminUserPassword          types.String `tfsdk:"admin_user_password"`
	LoginKeyName               types.String `tfsdk:"login_key_name"`
	EdgeNodeSubnetNo           types.String `tfsdk:"edge_node_subnet_no"`
	MasterNodeSubnetNo         types.String `tfsdk:"master_node_subnet_no"`
	WorkerNodeSubnetNo         types.String `tfsdk:"worker_node_subnet_no"`
	BucketName                 types.String `tfsdk:"bucket_name"`
	MasterNodeDataStorageType  types.String `tfsdk:"master_node_data_storage_type"`
	WorkerNodeDataStorageType  types.String `tfsdk:"worker_node_data_storage_type"`
	MasterNodeDataStorageSize  types.Int64  `tfsdk:"master_node_data_storage_size"`
	WorkerNodeDataStorageSize  types.Int64  `tfsdk:"worker_node_data_storage_size"`
	WorkerNodeCount            types.Int64  `tfsdk:"worker_node_count"`
	UseKdc                     types.Bool   `tfsdk:"use_kdc"`
	KdcRealm                   types.String `tfsdk:"kdc_realm"`
	KdcPassword                types.String `tfsdk:"kdc_password"`
	UseBootstrapScript         types.Bool   `tfsdk:"use_bootstrap_script"`
	BootstrapScript            types.String `tfsdk:"bootstrap_script"`
	UseDataCatalog             types.Bool   `tfsdk:"use_data_catalog"`
	AmbariServerHost           types.String `tfsdk:"ambari_server_host"`
	ClusterDirectAccessAccount types.String `tfsdk:"cluster_direct_access_account"`
	IsHa                       types.Bool   `tfsdk:"is_ha"`
	AccessControlGroupNoList   types.List   `tfsdk:"access_control_group_no_list"`
	HadoopServerInstanceList   types.List   `tfsdk:"hadoop_server_instance_list"`
	ID                         types.String `tfsdk:"id"`
}

func (m *hadoopResourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.CloudHadoopInstance) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ImageProductCode = types.StringPointerValue(output.CloudHadoopImageProductCode)
	m.ClusterName = types.StringPointerValue(output.CloudHadoopClusterName)
	m.ClusterTypeCode = types.StringPointerValue(output.CloudHadoopClusterType.Code)
	m.AddOnCodeList, _ = types.ListValueFrom(ctx, types.StringType, output.CloudHadoopAddOnList)
	m.KdcRealm = types.StringPointerValue(output.KdcRealm)
	m.ID = types.StringPointerValue(output.CloudHadoopInstanceNo)
	m.AmbariServerHost = types.StringPointerValue(output.AmbariServerHost)
	m.ClusterDirectAccessAccount = types.StringPointerValue(output.ClusterDirectAccessAccount)
	m.IsHa = types.BoolPointerValue(output.IsHa)
	m.Version = types.StringPointerValue(output.CloudHadoopVersion.Code)

	if addOnCodeList, err := types.ListValueFrom(ctx, types.StringType, output.CloudHadoopAddOnList); err.HasError() {
		m.AddOnCodeList = addOnCodeList
	} else {
		diagnostics.Append(err...)
	}

	if acgl, err := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList); err.HasError() {
		m.AccessControlGroupNoList = acgl
	} else {
		diagnostics.Append(err...)
	}

	if diagnostics.ErrorsCount() > 0 {
		return diagnostics
	}

	m.AccessControlGroupNoList, _ = types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)

	m.HadoopServerInstanceList, _ = listValueFromHadoopServerInatanceList(ctx, output.CloudHadoopServerInstanceList)
	return nil
}

type hadoopServer struct {
	HadoopServerName  types.String `tfsdk:"hadoop_server_name"`
	HadoopServerRole  types.String `tfsdk:"hadoop_server_role"`
	HadoopProductCode types.String `tfsdk:"hadoop_product_code"`
	RegionCode        types.String `tfsdk:"region_code"`
	ZoneCode          types.String `tfsdk:"zone_code"`
	VpcNo             types.String `tfsdk:"vpc_no"`
	SubnetNo          types.String `tfsdk:"subnet_no"`
	IsPublicSubnet    types.Bool   `tfsdk:"is_public_subnet"`
	DataStorageSize   types.Int64  `tfsdk:"data_storage_size"`
	CpuCount          types.Int64  `tfsdk:"cpu_count"`
	MemorySize        types.Int64  `tfsdk:"memory_size"`
}

func (h hadoopServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hadoop_server_name":  types.StringType,
		"hadoop_server_role":  types.StringType,
		"hadoop_product_code": types.StringType,
		"region_code":         types.StringType,
		"zone_code":           types.StringType,
		"vpc_no":              types.StringType,
		"subnet_no":           types.StringType,
		"is_public_subnet":    types.BoolType,
		"data_storage_size":   types.Int64Type,
		"cpu_count":           types.Int64Type,
		"memory_size":         types.Int64Type,
	}
}

func listValueFromHadoopServerInatanceList(ctx context.Context, serverInatances []*vhadoop.CloudHadoopServerInstance) (basetypes.ListValue, diag.Diagnostics) {
	var hadoopServerList []hadoopServer
	for _, serverInstance := range serverInatances {
		hadoopServerList = append(hadoopServerList, hadoopServer{
			HadoopServerName:  types.StringPointerValue(serverInstance.CloudHadoopServerName),
			HadoopServerRole:  types.StringPointerValue(serverInstance.CloudHadoopServerRole.CodeName),
			HadoopProductCode: types.StringPointerValue(serverInstance.CloudHadoopProductCode),
			RegionCode:        types.StringPointerValue(serverInstance.RegionCode),
			ZoneCode:          types.StringPointerValue(serverInstance.ZoneCode),
			VpcNo:             types.StringPointerValue(serverInstance.VpcNo),
			SubnetNo:          types.StringPointerValue(serverInstance.SubnetNo),
			IsPublicSubnet:    types.BoolPointerValue(serverInstance.IsPublicSubnet),
			DataStorageSize:   types.Int64Value(*serverInstance.DataStorageSize),
			CpuCount:          types.Int64Value(int64(*serverInstance.CpuCount)),
			MemorySize:        types.Int64Value(*serverInstance.MemorySize),
		})
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: hadoopServer{}.attrTypes()}, hadoopServerList)
}

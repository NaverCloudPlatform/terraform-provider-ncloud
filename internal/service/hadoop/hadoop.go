package hadoop

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

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

func (r *hadoopResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop"
}

func (r *hadoopResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
			"login_key": schema.StringAttribute{
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
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Description: "default: latest version",
			},
			"edge_node_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "default: minimum spec",
			},
			"master_node_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "default: minimum spec",
			},
			"worker_node_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "default: minimum spec",
			},
			"add_on_code_list": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Description: "this attribute can used over 1.5 version",
			},
			"worker_node_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(2),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(2),
				},
			},
			"use_kdc": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"kdc_realm": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("use_kdc"),
					}...),
					stringvalidator.All(
						stringvalidator.LengthBetween(1, 15),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[A-Z.]+$`),
							"Only uppercase letters (A-Z) are allowed and up to 15 digits are allowed. Only one dot(.) is allowed (ex. EXAMPLE.COM).",
						),
					),
				},
			},
			"kdc_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("use_kdc"),
					}...),
					stringvalidator.All(
						stringvalidator.LengthBetween(8, 20),
						stringvalidator.RegexMatches(regexp.MustCompile(`[A-Z]+`), "Must have at least one uppercase alphabet"),
						stringvalidator.RegexMatches(regexp.MustCompile(`\d+`), "Must have at least one number"),
						stringvalidator.RegexMatches(regexp.MustCompile(`[~!@#$%^*()\-_=\[\]\{\};:,.<>?]+`), "Must have at least one special character"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[^&+\\"'/\s`+"`"+`]*$`), "Must not have ` & + \\ \" ' / and white space."),
					),
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
						regexp.MustCompile(`^[a-zA-Z]+$`),
						"Composed of alphabets.",
					),
				},
			},
			// Available only `public` site
			"use_data_catalog": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "this attribute can used over 2.0 version",
			},
			"region_code": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ambari_server_host": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_direct_access_account": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_ha": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_control_group_no_list": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"hadoop_server_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"server_instance_no": schema.StringAttribute{
							Computed: true,
						},
						"server_name": schema.StringAttribute{
							Computed: true,
						},
						"server_role": schema.StringAttribute{
							Computed: true,
						},
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"is_public_subnet": schema.BoolAttribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"data_storage_type": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"uptime": schema.StringAttribute{
							Computed: true,
						},
						"create_date": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *hadoopResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.config = config
}

func (r *hadoopResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hadoopResourceModel

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"resource does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.CreateCloudHadoopInstanceRequest{
		RegionCode:                    &r.config.RegionCode,
		VpcNo:                         plan.VpcNo.ValueStringPointer(),
		CloudHadoopClusterName:        plan.ClusterName.ValueStringPointer(),
		CloudHadoopClusterTypeCode:    plan.ClusterTypeCode.ValueStringPointer(),
		CloudHadoopAdminUserName:      plan.AdminUserName.ValueStringPointer(),
		CloudHadoopAdminUserPassword:  plan.AdminUserPassword.ValueStringPointer(),
		LoginKeyName:                  plan.LoginKey.ValueStringPointer(),
		EdgeNodeSubnetNo:              plan.EdgeNodeSubnetNo.ValueStringPointer(),
		MasterNodeSubnetNo:            plan.MasterNodeSubnetNo.ValueStringPointer(),
		WorkerNodeSubnetNo:            plan.WorkerNodeSubnetNo.ValueStringPointer(),
		BucketName:                    plan.BucketName.ValueStringPointer(),
		MasterNodeDataStorageTypeCode: plan.MasterNodeDataStorageType.ValueStringPointer(),
		WorkerNodeDataStorageTypeCode: plan.WorkerNodeDataStorageType.ValueStringPointer(),
		MasterNodeDataStorageSize:     ncloud.Int32(int32(plan.MasterNodeDataStorageSize.ValueInt64())),
		WorkerNodeDataStorageSize:     ncloud.Int32(int32(plan.WorkerNodeDataStorageSize.ValueInt64())),
		WorkerNodeCount:               ncloud.Int32(int32(plan.WorkerNodeCount.ValueInt64())),
		UseKdc:                        plan.UseKdc.ValueBoolPointer(),
		UseBootstrapScript:            plan.UseBootstrapScript.ValueBoolPointer(),
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
		reqParams.WorkerNodeProductCode = plan.WorkerNodeProductCode.ValueStringPointer()
	}

	if !plan.AddOnCodeList.IsNull() && !plan.AddOnCodeList.IsUnknown() {
		addOnList := plan.AddOnCodeList.Elements()
		for _, addon := range addOnList {
			newStr := strings.Replace(addon.String(), "\"", "", -1)
			reqParams.CloudHadoopAddOnCodeList = append(reqParams.CloudHadoopAddOnCodeList, &newStr)
		}
	}

	if !plan.UseKdc.IsNull() && plan.UseKdc.ValueBool() {
		if !plan.KdcRealm.IsNull() && !plan.KdcRealm.IsUnknown() {
			reqParams.KdcRealm = plan.KdcRealm.ValueStringPointer()
		} else {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `use_kdc` is true, `kdc_realm` must be inputted`",
			)
			return
		}

		if !plan.KdcPassword.IsNull() {
			reqParams.KdcPassword = plan.KdcPassword.ValueStringPointer()
		} else {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `use_kdc` is true, `kdc_password` must be entered`",
			)
			return
		}
	} else {
		kdcRealmHasValue := !plan.KdcRealm.IsNull() && !plan.KdcRealm.IsUnknown()
		kdcPasswordHasValue := !plan.KdcPassword.IsNull() && !plan.KdcPassword.IsUnknown()
		if kdcRealmHasValue || kdcPasswordHasValue {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `use_kdc` is false, `kdc_realm` and `kdc_password` must not be entered`",
			)
			return
		}
	}

	if !plan.UseBootstrapScript.IsNull() && plan.UseBootstrapScript.ValueBool() {
		if !plan.BootstrapScript.IsNull() {
			reqParams.BootstrapScript = plan.BootstrapScript.ValueStringPointer()
		} else {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `use_bootstrap_script` is true, `bootstrap_script` must be entered`",
			)
			return
		}
	} else {
		if !plan.BootstrapScript.IsNull() {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `use_bootstrap_script` is false, `bootstrap_script` must not be entered`",
			)
			return
		}
	}

	// Available only `public` site
	if !plan.UseDataCatalog.IsNull() && plan.UseDataCatalog.ValueBool() {
		reqParams.UseDataCatalog = plan.UseDataCatalog.ValueBoolPointer()
	}

	tflog.Info(ctx, "CreateHadoop reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vhadoop.V2Api.CreateCloudHadoopInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateHadoop response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudHadoopInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	hadoopInstance := response.CloudHadoopInstanceList[0]
	plan.ID = types.StringPointerValue(hadoopInstance.CloudHadoopInstanceNo)

	output, err := waitHadoopCreation(ctx, r.config, *hadoopInstance.CloudHadoopInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *hadoopResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state hadoopResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetHadoopInstance(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *hadoopResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state hadoopResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.WorkerNodeCount.Equal(state.WorkerNodeCount) {
		reqParams := &vhadoop.ChangeCloudHadoopNodeCountRequest{
			RegionCode:            &r.config.RegionCode,
			CloudHadoopInstanceNo: state.ID.ValueStringPointer(),
			WorkerNodeCount:       ncloud.Int32(int32(plan.WorkerNodeCount.ValueInt64())),
		}
		tflog.Info(ctx, "ChangeHadoopWorkerNodeCount reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := r.config.Client.Vhadoop.V2Api.ChangeCloudHadoopNodeCount(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeHadoopWorkerNodeCount response="+common.MarshalUncheckedString(response))

		if response == nil || len(response.CloudHadoopInstanceList) < 1 {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		output, err := waitHadoopUpdate(ctx, r.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		state.refreshFromOutput(ctx, output)
	}

	if !plan.MasterNodeProductCode.Equal(state.MasterNodeProductCode) ||
		!plan.EdgeNodeProductCode.Equal(state.EdgeNodeProductCode) ||
		!plan.WorkerNodeProductCode.Equal(state.WorkerNodeProductCode) {
		reqParams := &vhadoop.ChangeCloudHadoopNodeSpecRequest{
			RegionCode:            &r.config.RegionCode,
			CloudHadoopInstanceNo: state.ID.ValueStringPointer(),
		}

		if !plan.MasterNodeProductCode.Equal(state.MasterNodeProductCode) {
			reqParams.MasterNodeProductCode = plan.MasterNodeProductCode.ValueStringPointer()
		}

		if !plan.EdgeNodeProductCode.Equal(state.EdgeNodeProductCode) {
			reqParams.EdgeNodeProductCode = plan.EdgeNodeProductCode.ValueStringPointer()
		}

		if !plan.WorkerNodeProductCode.Equal(state.WorkerNodeProductCode) {
			reqParams.WorkerNodeProductCode = plan.WorkerNodeProductCode.ValueStringPointer()
		}
		tflog.Info(ctx, "ChangeHadoopNodeSpec reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := r.config.Client.Vhadoop.V2Api.ChangeCloudHadoopNodeSpec(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeHadoopNodeSpec response="+common.MarshalUncheckedString(response))

		if response == nil || len(response.CloudHadoopInstanceList) < 1 {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		output, err := waitHadoopUpdate(ctx, r.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		state.refreshFromOutput(ctx, output)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *hadoopResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state hadoopResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.DeleteCloudHadoopInstanceRequest{
		RegionCode:            &r.config.RegionCode,
		CloudHadoopInstanceNo: state.ID.ValueStringPointer(),
	}
	tflog.Info(ctx, "DeleteHadoop reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vhadoop.V2Api.DeleteCloudHadoopInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteHadoop response="+common.MarshalUncheckedString(response))

	if err := waitHadoopDeletion(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (r *hadoopResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GetHadoopInstance(ctx context.Context, config *conn.ProviderConfig, id string) (*vhadoop.CloudHadoopInstance, error) {
	reqParams := &vhadoop.GetCloudHadoopInstanceDetailRequest{
		RegionCode:            &config.RegionCode,
		CloudHadoopInstanceNo: ncloud.String(id),
	}
	tflog.Info(ctx, "GetHadoopDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vhadoop.V2Api.GetCloudHadoopInstanceDetail(reqParams)
	// If the lookup result is 0 or already deleted, it will respond with a 400 error with a 5001017 return code.
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) {
		return nil, err
	}
	tflog.Info(ctx, "GetHadoopDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudHadoopInstanceList) < 1 || len(resp.CloudHadoopInstanceList[0].CloudHadoopServerInstanceList) < 1 {
		return nil, nil
	}

	return resp.CloudHadoopInstanceList[0], nil
}

func waitHadoopCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vhadoop.CloudHadoopInstance, error) {
	var hadoopInstance *vhadoop.CloudHadoopInstance
	var err error
	stateConf := &retry.StateChangeConf{
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
			return 0, "", fmt.Errorf("error occurred while waiting to create")
		},
		Timeout:    90 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, err
	}

	return hadoopInstance, nil
}

func waitHadoopUpdate(ctx context.Context, config *conn.ProviderConfig, id string) (*vhadoop.CloudHadoopInstance, error) {
	var hadoopInstance *vhadoop.CloudHadoopInstance
	var err error

	stateConf := &retry.StateChangeConf{
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
		Timeout:    6 * conn.DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, err
	}

	return hadoopInstance, nil
}

func waitHadoopDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PEND"},
		Target:  []string{"DEL"},
		Refresh: func() (interface{}, string, error) {
			hadoopInstance, err := GetHadoopInstance(ctx, config, id)
			if err != nil {
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

			return 0, "", fmt.Errorf("error occurred while waiting to delete")
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
	ID                         types.String `tfsdk:"id"`
	VpcNo                      types.String `tfsdk:"vpc_no"`
	ClusterName                types.String `tfsdk:"cluster_name"`
	ClusterTypeCode            types.String `tfsdk:"cluster_type_code"`
	AdminUserName              types.String `tfsdk:"admin_user_name"`
	AdminUserPassword          types.String `tfsdk:"admin_user_password"`
	LoginKey                   types.String `tfsdk:"login_key"`
	EdgeNodeSubnetNo           types.String `tfsdk:"edge_node_subnet_no"`
	MasterNodeSubnetNo         types.String `tfsdk:"master_node_subnet_no"`
	WorkerNodeSubnetNo         types.String `tfsdk:"worker_node_subnet_no"`
	BucketName                 types.String `tfsdk:"bucket_name"`
	MasterNodeDataStorageType  types.String `tfsdk:"master_node_data_storage_type"`
	WorkerNodeDataStorageType  types.String `tfsdk:"worker_node_data_storage_type"`
	MasterNodeDataStorageSize  types.Int64  `tfsdk:"master_node_data_storage_size"`
	WorkerNodeDataStorageSize  types.Int64  `tfsdk:"worker_node_data_storage_size"`
	ImageProductCode           types.String `tfsdk:"image_product_code"`
	EdgeNodeProductCode        types.String `tfsdk:"edge_node_product_code"`
	MasterNodeProductCode      types.String `tfsdk:"master_node_product_code"`
	WorkerNodeProductCode      types.String `tfsdk:"worker_node_product_code"`
	AddOnCodeList              types.List   `tfsdk:"add_on_code_list"`
	WorkerNodeCount            types.Int64  `tfsdk:"worker_node_count"`
	UseKdc                     types.Bool   `tfsdk:"use_kdc"`
	KdcRealm                   types.String `tfsdk:"kdc_realm"`
	KdcPassword                types.String `tfsdk:"kdc_password"`
	UseBootstrapScript         types.Bool   `tfsdk:"use_bootstrap_script"`
	BootstrapScript            types.String `tfsdk:"bootstrap_script"`
	UseDataCatalog             types.Bool   `tfsdk:"use_data_catalog"`
	RegionCode                 types.String `tfsdk:"region_code"`
	AmbariServerHost           types.String `tfsdk:"ambari_server_host"`
	ClusterDirectAccessAccount types.String `tfsdk:"cluster_direct_access_account"`
	Version                    types.String `tfsdk:"version"`
	IsHa                       types.Bool   `tfsdk:"is_ha"`
	Domain                     types.String `tfsdk:"domain"`
	AccessControlGroupNoList   types.List   `tfsdk:"access_control_group_no_list"`
	HadoopServerList           types.List   `tfsdk:"hadoop_server_list"`
}

type hadoopServer struct {
	ServerInstanceNo types.String `tfsdk:"server_instance_no"`
	ServerName       types.String `tfsdk:"server_name"`
	ServerRole       types.String `tfsdk:"server_role"`
	ZoneCode         types.String `tfsdk:"zone_code"`
	SubnetNo         types.String `tfsdk:"subnet_no"`
	ProductCode      types.String `tfsdk:"product_code"`
	IsPublicSubnet   types.Bool   `tfsdk:"is_public_subnet"`
	CpuCount         types.Int64  `tfsdk:"cpu_count"`
	MemorySize       types.Int64  `tfsdk:"memory_size"`
	DataStorageType  types.String `tfsdk:"data_storage_type"`
	DataStorageSize  types.Int64  `tfsdk:"data_storage_size"`
	Uptime           types.String `tfsdk:"uptime"`
	CreateDate       types.String `tfsdk:"create_date"`
}

func (h hadoopServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_instance_no": types.StringType,
		"server_name":        types.StringType,
		"server_role":        types.StringType,
		"zone_code":          types.StringType,
		"subnet_no":          types.StringType,
		"product_code":       types.StringType,
		"is_public_subnet":   types.BoolType,
		"cpu_count":          types.Int64Type,
		"memory_size":        types.Int64Type,
		"data_storage_type":  types.StringType,
		"data_storage_size":  types.Int64Type,
		"uptime":             types.StringType,
		"create_date":        types.StringType,
	}
}

func (m *hadoopResourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.CloudHadoopInstance) {
	m.ID = types.StringPointerValue(output.CloudHadoopInstanceNo)
	m.VpcNo = types.StringPointerValue(output.CloudHadoopServerInstanceList[0].VpcNo)
	m.ClusterName = types.StringPointerValue(output.CloudHadoopClusterName)
	m.ClusterTypeCode = types.StringPointerValue(output.CloudHadoopClusterType.Code)
	m.LoginKey = types.StringPointerValue(output.LoginKey)
	m.BucketName = types.StringPointerValue(output.ObjectStorageBucket)
	m.ImageProductCode = types.StringPointerValue(output.CloudHadoopImageProductCode)
	m.KdcRealm = types.StringPointerValue(output.KdcRealm)
	m.RegionCode = types.StringPointerValue(output.CloudHadoopServerInstanceList[0].RegionCode)
	m.AmbariServerHost = types.StringPointerValue(output.AmbariServerHost)
	m.ClusterDirectAccessAccount = types.StringPointerValue(output.ClusterDirectAccessAccount)
	m.Version = types.StringPointerValue(output.CloudHadoopVersion.Code)
	m.IsHa = types.BoolPointerValue(output.IsHa)
	m.Domain = types.StringPointerValue(output.Domain)

	if output.KdcRealm != nil {
		m.UseKdc = types.BoolValue(true)
	} else {
		m.UseKdc = types.BoolValue(false)
	}

	var count int64
	var storageSize int64
	for _, server := range output.CloudHadoopServerInstanceList {
		if server.CloudHadoopServerRole != nil {
			if *server.CloudHadoopServerRole.Code == "E" {
				m.EdgeNodeProductCode = types.StringPointerValue(server.CloudHadoopProductCode)
				m.EdgeNodeSubnetNo = types.StringPointerValue(server.SubnetNo)
			}
			if *server.CloudHadoopServerRole.Code == "M" {
				m.MasterNodeProductCode = types.StringPointerValue(server.CloudHadoopProductCode)
				m.MasterNodeSubnetNo = types.StringPointerValue(server.SubnetNo)
				if server.DataStorageType != nil {
					m.MasterNodeDataStorageType = types.StringPointerValue(server.DataStorageType.Code)
				}
				// Byte to GBi
				storageSize = *server.DataStorageSize / 1024 / 1024 / 1024
				m.MasterNodeDataStorageSize = types.Int64Value(storageSize)
			}
			if *server.CloudHadoopServerRole.Code == "D" {
				m.WorkerNodeProductCode = types.StringPointerValue(server.CloudHadoopProductCode)
				m.WorkerNodeSubnetNo = types.StringPointerValue(server.SubnetNo)
				if server.DataStorageType != nil {
					m.WorkerNodeDataStorageType = types.StringPointerValue(server.DataStorageType.Code)
				}
				// Byte to GBi
				storageSize = *server.DataStorageSize / 1024 / 1024 / 1024
				m.WorkerNodeDataStorageSize = types.Int64Value(storageSize)
				count++
			}
		}
	}
	m.WorkerNodeCount = types.Int64Value(count)

	var addOnList []string
	for _, addOn := range output.CloudHadoopAddOnList {
		addOnList = append(addOnList, *addOn.Code)
	}
	m.AddOnCodeList, _ = types.ListValueFrom(ctx, types.StringType, addOnList)
	m.AccessControlGroupNoList, _ = types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.HadoopServerList, _ = listValueFromHadoopServerInatanceList(ctx, output.CloudHadoopServerInstanceList)
}

func listValueFromHadoopServerInatanceList(ctx context.Context, serverInatances []*vhadoop.CloudHadoopServerInstance) (basetypes.ListValue, diag.Diagnostics) {
	var hadoopServerList []hadoopServer
	for _, serverInstance := range serverInatances {
		hadoopServerInstance := hadoopServer{
			ServerInstanceNo: types.StringPointerValue(serverInstance.CloudHadoopServerInstanceNo),
			ServerName:       types.StringPointerValue(serverInstance.CloudHadoopServerName),
			ZoneCode:         types.StringPointerValue(serverInstance.ZoneCode),
			SubnetNo:         types.StringPointerValue(serverInstance.SubnetNo),
			ProductCode:      types.StringPointerValue(serverInstance.CloudHadoopProductCode),
			IsPublicSubnet:   types.BoolPointerValue(serverInstance.IsPublicSubnet),
			CpuCount:         common.Int64ValueFromInt32(serverInstance.CpuCount),
			MemorySize:       types.Int64PointerValue(serverInstance.MemorySize),
			DataStorageSize:  types.Int64PointerValue(serverInstance.DataStorageSize),
			Uptime:           types.StringPointerValue(serverInstance.Uptime),
			CreateDate:       types.StringPointerValue(serverInstance.CreateDate),
		}

		if serverInstance.CloudHadoopServerRole != nil {
			hadoopServerInstance.ServerRole = types.StringPointerValue(serverInstance.CloudHadoopServerRole.CodeName)
		}
		if serverInstance.DataStorageType != nil {
			hadoopServerInstance.DataStorageType = types.StringPointerValue(serverInstance.DataStorageType.Code)
		}
		hadoopServerList = append(hadoopServerList, hadoopServerInstance)
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: hadoopServer{}.attrTypes()}, hadoopServerList)
}

package mssql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmssql"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
)

var (
	_ resource.Resource                = &mssqlResource{}
	_ resource.ResourceWithConfigure   = &mssqlResource{}
	_ resource.ResourceWithImportState = &mssqlResource{}
)

func NewMssqlResource() resource.Resource {
	return &mssqlResource{}
}

type mssqlResource struct {
	config *conn.ProviderConfig
}

func (m *mssqlResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mssql"
}

func (m *mssqlResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"subnet_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 15),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[ㄱ-ㅣ가-힣A-Za-z0-9-]+$`),
						"Composed of alphabets, numbers, hyphen (-).",
					),
				},
			},
			"is_ha": schema.BoolAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: true",
			},
			"user_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 16),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_,]+$`),
						"Composed of alphabets, numbers, hyphen (-), (\\), (_), (,). Must start with an alphabetic character.",
					),
				},
			},
			"user_password": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(8, 20),
						stringvalidator.RegexMatches(regexp.MustCompile(`[a-zA-Z]+`), "Must have at least one alphabet"),
						stringvalidator.RegexMatches(regexp.MustCompile(`\d+`), "Must have at least one number"),
						stringvalidator.RegexMatches(regexp.MustCompile(`[~!@#$%^*()\-_=\[\]\{\};:,.<>?]+`), "Must have at least one special character"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[^&+\\"'/\s`+"`"+`]*$`), "Must not have ` & + \\ \" ' / and white space."),
					),
				},
				Sensitive: true,
			},
			"config_group_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("0"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_storage_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SSD", "HDD"}...),
				},
				Description: "default: SSD",
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 30),
				},
			},
			"backup_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3])(:?(00|15|30|45))$`), "Must be in the format HHMM and 15 minutes internvals."),
				},
				Description: "ex) 01:15",
			},
			"is_automatic_backup": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1433),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(1433),
					),
				},
			},
			"character_set_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Korean_Wansung_CI_AS"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"engine_version": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_code": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"mssql_server_list": schema.ListNestedAttribute{
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
						"public_domain": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"used_data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
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

func (r *mssqlResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

func (r *mssqlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mssqlResourceModel

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

	subnet, err := vpc.GetSubnetInstance(r.config, plan.SubnetNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"CREATING ERROR",
			err.Error(),
		)
	}

	reqParams := &vmssql.CreateCloudMssqlInstanceRequest{
		RegionCode:                &r.config.RegionCode,
		VpcNo:                     subnet.VpcNo,
		SubnetNo:                  subnet.SubnetNo,
		CloudMssqlServiceName:     plan.ServiceName.ValueStringPointer(),
		CloudMssqlUserName:        plan.UserName.ValueStringPointer(),
		CloudMssqlUserPassword:    plan.UserPassword.ValueStringPointer(),
		IsHa:                      plan.IsHa.ValueBoolPointer(),
		ConfigGroupNo:             plan.ConfigGroupNo.ValueStringPointer(),
		BackupFileRetentionPeriod: ncloud.Int32(int32(plan.BackupFileRetentionPeriod.ValueInt64())),
		CloudMssqlPort:            ncloud.Int32(int32(plan.Port.ValueInt64())),
		CharacterSetName:          plan.CharacterSetName.ValueStringPointer(),
	}
	plan.VpcNo = types.StringPointerValue(subnet.VpcNo)

	if !plan.DataStorageTypeCode.IsNull() && !plan.DataStorageTypeCode.IsUnknown() {
		reqParams.DataStorageTypeCode = plan.DataStorageTypeCode.ValueStringPointer()
	}

	if !plan.ProductCode.IsNull() && !plan.ProductCode.IsUnknown() {
		reqParams.CloudMssqlProductCode = plan.ProductCode.ValueStringPointer()
	}

	if !plan.ImageProductCode.IsNull() && !plan.ImageProductCode.IsUnknown() {
		reqParams.CloudMssqlImageProductCode = plan.ImageProductCode.ValueStringPointer()
	}

	if !plan.IsAutomaticBackup.IsNull() && !plan.IsAutomaticBackup.IsUnknown() {
		reqParams.IsAutomaticBackup = plan.IsAutomaticBackup.ValueBoolPointer()
	}

	if reqParams.IsAutomaticBackup == nil || *reqParams.IsAutomaticBackup {
		if !plan.BackupTime.IsNull() && !plan.BackupTime.IsUnknown() {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_automactic_backup` is true, `backup_time` is not used",
			)
			return
		}
	} else {
		if plan.BackupTime.IsNull() || plan.BackupTime.IsUnknown() {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_automactic_backup` is false, `backup_time` must be entered",
			)
			return
		}
		reqParams.BackupTime = plan.BackupTime.ValueStringPointer()
	}

	response, err := r.config.Client.Vmssql.V2Api.CreateCloudMssqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMssql response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudMssqlInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	mssqlIns := response.CloudMssqlInstanceList[0]
	plan.ID = types.StringPointerValue(mssqlIns.CloudMssqlInstanceNo)

	output, err := waitMssqlCreation(ctx, r.config, *mssqlIns.CloudMssqlInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mssqlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mssqlResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMssqlInstance(ctx, r.config, state.ID.ValueString())
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

func (m *mssqlResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *mssqlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mssqlResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmssql.DeleteCloudMssqlInstanceRequest{
		RegionCode:           &r.config.RegionCode,
		CloudMssqlInstanceNo: state.ID.ValueStringPointer(),
	}
	tflog.Info(ctx, "DeleteMssql reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmssql.V2Api.DeleteCloudMssqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteMssql response="+common.MarshalUncheckedString(response))

	if err := waitMssqlDeletion(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (r *mssqlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GetMssqlInstance(ctx context.Context, config *conn.ProviderConfig, no string) (*vmssql.CloudMssqlInstance, error) {
	reqParams := &vmssql.GetCloudMssqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMssqlInstanceNo: &no,
	}
	tflog.Info(ctx, "GetMssqlDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vmssql.V2Api.GetCloudMssqlInstanceDetail(reqParams)
	// If the lookup result is 0, it will respond with a 400 error with a 5001017 return code.
	// MSSQL deleted, it will respond with a 400 error with a 5001269 return code.
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) && !strings.Contains(err.Error(), `"returnCode": "5001269"`) {
		return nil, err
	}
	tflog.Info(ctx, "GetMssqlDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudMssqlInstanceList) < 1 || len(resp.CloudMssqlInstanceList[0].CloudMssqlServerInstanceList) < 1 {
		return nil, nil
	}

	return resp.CloudMssqlInstanceList[0], nil
}

func waitMssqlCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vmssql.CloudMssqlInstance, error) {
	var mssqlInstance *vmssql.CloudMssqlInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMssqlInstance(ctx, config, id)
			mssqlInstance = instance
			if err != nil {
				return 0, "", err
			}

			status := instance.CloudMssqlInstanceStatus.Code
			op := instance.CloudMssqlInstanceOperation.Code

			if *status == "INIT" && *op == "CREAT" {
				return instance, "creating", nil
			}

			if *status == "CREAT" && *op == "SETUP" {
				return instance, "settingUp", nil
			}

			if *status == "CREAT" && *op == "NULL" {
				return instance, "running", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create mssql")
		},
		Timeout:    90 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("error waiting for MssqlInstance state to be \"CREAT\": %s", err)
	}

	return mssqlInstance, nil
}

func waitMssqlDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMssqlInstance(ctx, config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, "deleted", nil
			}

			status := instance.CloudMssqlInstanceStatus.Code
			op := instance.CloudMssqlInstanceOperation.Code

			if *status == "DEL" && *op == "DEL" {
				return instance, "deleting", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete mssql")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for mssql (%s) to become termintaing: %s", id, err)
	}

	return nil
}

type mssqlResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ServiceName               types.String `tfsdk:"service_name"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	ConfigGroupNo             types.String `tfsdk:"config_group_no"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	IsAutomaticBackup         types.Bool   `tfsdk:"is_automatic_backup"`
	Port                      types.Int64  `tfsdk:"port"`
	CharacterSetName          types.String `tfsdk:"character_set_name"`
	EngineVersion             types.String `tfsdk:"engine_version"`
	RegionCode                types.String `tfsdk:"region_code"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MssqlServerList           types.List   `tfsdk:"mssql_server_list"`
}

type mssqlServer struct {
	ServerInstanceNo    types.String `tfsdk:"server_instance_no"`
	ServerName          types.String `tfsdk:"server_name"`
	ServerRole          types.String `tfsdk:"server_role"`
	ZoneCode            types.String `tfsdk:"zone_code"`
	SubnetNo            types.String `tfsdk:"subnet_no"`
	ProductCode         types.String `tfsdk:"product_code"`
	IsPublicSubnet      types.Bool   `tfsdk:"is_public_subnet"`
	PublicDomain        types.String `tfsdk:"public_domain"`
	PrivateDomain       types.String `tfsdk:"private_domain"`
	DataStorageSize     types.Int64  `tfsdk:"data_storage_size"`
	UsedDataStorageSize types.Int64  `tfsdk:"used_data_storage_size"`
	CpuCount            types.Int64  `tfsdk:"cpu_count"`
	MemorySize          types.Int64  `tfsdk:"memory_size"`
	Uptime              types.String `tfsdk:"uptime"`
	CreateDate          types.String `tfsdk:"create_date"`
}

func (m mssqlServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_instance_no":     types.StringType,
		"server_name":            types.StringType,
		"server_role":            types.StringType,
		"zone_code":              types.StringType,
		"subnet_no":              types.StringType,
		"product_code":           types.StringType,
		"is_public_subnet":       types.BoolType,
		"public_domain":          types.StringType,
		"private_domain":         types.StringType,
		"data_storage_size":      types.Int64Type,
		"used_data_storage_size": types.Int64Type,
		"cpu_count":              types.Int64Type,
		"memory_size":            types.Int64Type,
		"uptime":                 types.StringType,
		"create_date":            types.StringType,
	}
}

func (m *mssqlResourceModel) refreshFromOutput(ctx context.Context, output *vmssql.CloudMssqlInstance) {
	m.ID = types.StringPointerValue(output.CloudMssqlInstanceNo)
	m.ServiceName = types.StringPointerValue(output.CloudMssqlServiceName)
	m.ImageProductCode = types.StringPointerValue(output.CloudMssqlImageProductCode)
	m.DataStorageTypeCode = types.StringPointerValue(output.CloudMssqlServerInstanceList[0].DataStorageType.Code)
	m.IsHa = types.BoolPointerValue(output.IsHa)
	m.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.ConfigGroupNo = types.StringPointerValue(output.ConfigGroupNo)
	m.Port = common.Int64ValueFromInt32(output.CloudMssqlPort)
	m.EngineVersion = types.StringPointerValue(output.EngineVersion)
	m.CharacterSetName = types.StringPointerValue(output.DbCollation)
	m.RegionCode = types.StringPointerValue(output.CloudMssqlServerInstanceList[0].RegionCode)
	m.VpcNo = types.StringPointerValue(output.CloudMssqlServerInstanceList[0].VpcNo)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.AccessControlGroupNoList = acgList

	var serverList []mssqlServer
	for _, server := range output.CloudMssqlServerInstanceList {
		mssqlServerInstance := mssqlServer{
			ServerInstanceNo: types.StringPointerValue(server.CloudMssqlServerInstanceNo),
			ServerName:       types.StringPointerValue(server.CloudMssqlServerName),
			ServerRole:       types.StringPointerValue(server.CloudMssqlServerRole.Code),
			ZoneCode:         types.StringPointerValue(server.ZoneCode),
			SubnetNo:         types.StringPointerValue(server.SubnetNo),
			ProductCode:      types.StringPointerValue(server.CloudMssqlProductCode),
			IsPublicSubnet:   types.BoolPointerValue(server.IsPublicSubnet),
			PrivateDomain:    types.StringPointerValue(server.PrivateDomain),
			DataStorageSize:  types.Int64PointerValue(server.DataStorageSize),
			CpuCount:         common.Int64ValueFromInt32(server.CpuCount),
			MemorySize:       types.Int64PointerValue(server.MemorySize),
			Uptime:           types.StringPointerValue(server.Uptime),
			CreateDate:       types.StringPointerValue(server.CreateDate),
		}

		if server.PublicDomain != nil {
			mssqlServerInstance.PublicDomain = types.StringPointerValue(server.PublicDomain)
		}

		if server.UsedDataStorageSize != nil {
			mssqlServerInstance.UsedDataStorageSize = types.Int64PointerValue(server.UsedDataStorageSize)
		}
		serverList = append(serverList, mssqlServerInstance)
	}

	mssqlServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mssqlServer{}.attrTypes()}, serverList)

	m.MssqlServerList = mssqlServers
}

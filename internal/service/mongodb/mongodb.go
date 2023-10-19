package mongodb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &mongodbResource{}
	_ resource.ResourceWithConfigure   = &mongodbResource{}
	_ resource.ResourceWithImportState = &mongodbResource{}
)

func NewMongoDbResource() resource.Resource {
	return &mongodbResource{}
}

type mongodbResource struct {
	config *conn.ProviderConfig
}

func (s *mongodbResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (m *mongodbResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb"
}

func (m *mongodbResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	m.config = config
}

func (m *mongodbResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(3, 15),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[가-힣A-Za-z0-9-]+$`), "Allows only hangeuls, alphabets, numbers, hyphen (-)."),
					),
				},
				Description: "Service Name of Cloud DB for MongoDb instance.",
			},
			"id": framework.IDAttribute(),
			"instance_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"user_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(4, 16),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-_]+$`), "Allows only alphabets, numbers, hyphen (-) and underbar (_). Must start with an alphabetic character"),
					),
				},
				Description: "Access username, which will be used for DB admin.",
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
				Description: "Access password for user, which will be used for DB admin.",
				Sensitive:   true,
			},
			"vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_type_code": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"STAND_ALONE", "SINGLE_REPLICA_SET", "SHARDED_CLUSTER"}...),
				},
			},
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"member_product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"arbiter_product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mongos_product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config_product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_storage_type_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shard_count": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(2, 3),
				},
			},
			"member_server_count": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(2, 7),
				},
			},
			"arbiter_server_count": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 1),
				},
			},
			"mongos_server_count": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(2, 5),
				},
			},
			"config_server_count": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(3, 3),
				},
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 30),
				},
			},
			"backup_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3])(:?(00|15|30|45))$`), "Must be in the format HHMM and 15 minutes internvals."),
				},
			},
			"arbiter_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(27017),
					),
				},
			},
			"member_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(27017),
					),
				},
			},
			"mongos_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(27017),
					),
				},
			},
			"config_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(27017),
					),
				},
			},
			"compress_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SNPP", "ZLIB", "ZSTD", "NONE"}...),
				},
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (m *mongodbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mongodbResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !m.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.Config.Schema.Type().String()),
		)
		return
	}

	reqParams := &vmongodb.CreateCloudMongoDbInstanceRequest{
		RegionCode:                   &m.config.RegionCode,
		VpcNo:                        plan.VpcNo.ValueStringPointer(),
		CloudMongoDbImageProductCode: plan.CloudMongoDbImageProductCode.ValueStringPointer(),
		MemberProductCode:            plan.MemberProductCode.ValueStringPointer(),
		ArbiterProductCode:           plan.ArbiterProductCode.ValueStringPointer(),
		MongosProductCode:            plan.MongosProductCode.ValueStringPointer(),
		ConfigProductCode:            plan.ConfigProductCode.ValueStringPointer(),
		CloudMongoDbUserName:         plan.CloudMongoDbUserName.ValueStringPointer(),
		CloudMongoDbUserPassword:     plan.CloudMongoDbUserPassword.ValueStringPointer(),
		CloudMongoDbServiceName:      plan.CloudMongoDbServiceName.ValueStringPointer(),
		SubnetNo:                     plan.SubnetNo.ValueStringPointer(),
		ClusterTypeCode:              plan.ClusterTypeCode.ValueStringPointer(),
		CompressCode:                 plan.CompressCode.ValueStringPointer(),
	}

	if !plan.BackupTime.IsNull() && !plan.BackupTime.IsUnknown() {
		reqParams.BackupTime = plan.BackupTime.ValueStringPointer()
	}

	if !plan.BackupFileRetentionPeriod.IsNull() && !plan.BackupFileRetentionPeriod.IsUnknown() {
		reqParams.BackupFileRetentionPeriod = ncloud.Int32(int32(plan.BackupFileRetentionPeriod.ValueInt64()))
	}

	if !plan.MemberPort.IsNull() && !plan.MemberPort.IsUnknown() {
		reqParams.MemberPort = ncloud.Int32(int32(plan.MemberPort.ValueInt64()))
	}

	if !plan.ArbiterPort.IsNull() && !plan.ArbiterPort.IsUnknown() {
		reqParams.ArbiterPort = ncloud.Int32(int32(plan.ArbiterPort.ValueInt64()))
	}

	if !plan.MongosPort.IsNull() && !plan.MongosPort.IsUnknown() {
		reqParams.MongosPort = ncloud.Int32(int32(plan.MongosPort.ValueInt64()))
	}

	if !plan.ConfigPort.IsNull() && !plan.ConfigPort.IsUnknown() {
		reqParams.ConfigPort = ncloud.Int32(int32(plan.ConfigPort.ValueInt64()))
	}

	if !plan.DataStorageTypeCode.IsNull() {
		reqParams.DataStorageTypeCode = plan.DataStorageTypeCode.ValueStringPointer()
	}

	if *reqParams.ClusterTypeCode == "SHARDED_CLUSTER" {
		reqParams.ShardCount = ncloud.Int32(int32(plan.ShardCount.ValueInt64()))
		reqParams.MemberServerCount = ncloud.Int32(int32(plan.MemberServerCount.ValueInt64()))
		reqParams.ArbiterServerCount = ncloud.Int32(int32(plan.ArbiterServerCount.ValueInt64()))
		reqParams.MongosServerCount = ncloud.Int32(int32(plan.MongosServerCount.ValueInt64()))
		reqParams.ConfigServerCount = ncloud.Int32(int32(plan.ConfigServerCount.ValueInt64()))
	}

	tflog.Info(ctx, "CreateMongoDb", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	response, err := m.config.Client.Vmongodb.V2Api.CreateCloudMongoDbInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create MongoDb Instance, err params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "CreateMongoDb response", map[string]any{
		"createMongoDbResponse": common.MarshalUncheckedString(response),
	})

	mongodbInstance := response.CloudMongoDbInstanceList[0]
	plan.ID = types.StringPointerValue(mongodbInstance.CloudMongoDbInstanceNo)
	tflog.Info(ctx, "MongoDb ID", map[string]any{"CloudMongoDbInstanceNo": *mongodbInstance.CloudMongoDbInstanceNo})

	output, err := waitForNcloudMongoDbActive(ctx, m.config, *mongodbInstance.CloudMongoDbInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("fail to wait for mongodb active", err.Error())
		return
	}

	if err := plan.refreshFromOutput(ctx, output); err != nil {
		resp.Diagnostics.AddError("refreshing mongodb details", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (m *mongodbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mongodbResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetCloudMongoDbInstance(ctx, m.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetCloudMongoDb", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := state.refreshFromOutput(ctx, output); err != nil {
		resp.Diagnostics.AddError("refreshing mongodb details", err.Error())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (m *mongodbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (m *mongodbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mongodbResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmongodb.DeleteCloudMongoDbInstanceRequest{
		RegionCode:             &m.config.RegionCode,
		CloudMongoDbInstanceNo: state.CloudMongoDbInstanceNo.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteMongoDb", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	response, err := m.config.Client.Vmongodb.V2Api.DeleteCloudMongoDbInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete Mongodb Instance params=%v", *reqParams),
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "DeleteMongoDb response", map[string]any{
		"deleteMongoDbResponse": common.MarshalUncheckedString(response),
	})

	if err := WaitForNcloudMongoDbDeletion(ctx, m.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail to wait for mongodb deletion",
			err.Error(),
		)
	}
}

func waitForNcloudMongoDbActive(ctx context.Context, config *conn.ProviderConfig, id string) (*vmongodb.CloudMongoDbInstance, error) {
	var mongodbInstance *vmongodb.CloudMongoDbInstance
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetCloudMongoDbInstance(ctx, config, id)
			mongodbInstance = instance
			if err != nil && !strings.Contains(err.Error(), `"returnCode": "5001017"`) {
				return 0, "", err
			}

			status := instance.CloudMongoDbInstanceStatus.Code
			op := instance.CloudMongoDbInstanceOperation.Code

			if *status == "INIT" && *op == "CREAT" {
				return instance, "creating", nil
			}

			if *status == "CREAT" && *op == "SETUP" {
				return instance, "settingUp", nil
			}

			if *status == "CREAT" && *op == "NULL" {
				return instance, "running", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create mongodb")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("error waiting for MysqlInstance state to be \"CREAT\": %s", err)
	}

	return mongodbInstance, nil
}

func WaitForNcloudMongoDbDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetCloudMongoDbInstance(ctx, config, id)

			if err != nil && !strings.Contains(err.Error(), `"returnCode": "5001017"`) {
				return 0, "", err
			}

			if instance == nil {
				return instance, "deleted", nil
			}

			status := instance.CloudMongoDbInstanceStatus.Code
			op := instance.CloudMongoDbInstanceOperation.Code

			if *status == "DEL" && *op == "DEL" {
				return instance, "deleting", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete mongodb")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for mongodb (%s) to become terminating: %s", id, err)
	}

	return nil
}

func GetCloudMongoDbInstance(ctx context.Context, config *conn.ProviderConfig, id string) (*vmongodb.CloudMongoDbInstance, error) {
	reqParams := &vmongodb.GetCloudMongoDbInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		CloudMongoDbInstanceNo: ncloud.String(id),
	}

	tflog.Info(ctx, "GetMongoDb", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vmongodb.V2Api.GetCloudMongoDbInstanceDetail(reqParams)
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) {
		return nil, err
	}

	tflog.Info(ctx, "GetMongoDb response", map[string]any{
		"getMongoDbResponse": common.MarshalUncheckedString(resp),
	})

	if len(resp.CloudMongoDbInstanceList) > 0 {
		mongodb := resp.CloudMongoDbInstanceList[0]
		return mongodb, nil
	}

	return nil, nil
}

type mongodbResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	VpcNo                        types.String `tfsdk:"vpc_no"`
	SubnetNo                     types.String `tfsdk:"subnet_no"`
	CloudMongoDbInstanceNo       types.String `tfsdk:"instance_no"`
	CloudMongoDbServiceName      types.String `tfsdk:"service_name"`
	CloudMongoDbUserName         types.String `tfsdk:"user_name"`
	CloudMongoDbUserPassword     types.String `tfsdk:"user_password"`
	ClusterTypeCode              types.String `tfsdk:"cluster_type_code"`
	CloudMongoDbImageProductCode types.String `tfsdk:"image_product_code"`
	MemberProductCode            types.String `tfsdk:"member_product_code"`
	ArbiterProductCode           types.String `tfsdk:"arbiter_product_code"`
	MongosProductCode            types.String `tfsdk:"mongos_product_code"`
	ConfigProductCode            types.String `tfsdk:"config_product_code"`
	ShardCount                   types.Int64  `tfsdk:"shard_count"`
	MemberServerCount            types.Int64  `tfsdk:"member_server_count"`
	ArbiterServerCount           types.Int64  `tfsdk:"arbiter_server_count"`
	MongosServerCount            types.Int64  `tfsdk:"mongos_server_count"`
	ConfigServerCount            types.Int64  `tfsdk:"config_server_count"`
	BackupFileRetentionPeriod    types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                   types.String `tfsdk:"backup_time"`
	DataStorageTypeCode          types.String `tfsdk:"data_storage_type_code"`
	ArbiterPort                  types.Int64  `tfsdk:"arbiter_port"`
	MemberPort                   types.Int64  `tfsdk:"member_port"`
	MongosPort                   types.Int64  `tfsdk:"mongos_port"`
	ConfigPort                   types.Int64  `tfsdk:"config_port"`
	CompressCode                 types.String `tfsdk:"compress_code"`
	AccessControlGroupNoList     types.List   `tfsdk:"access_control_group_no_list"`
}

func (m *mongodbResourceModel) refreshFromOutput(ctx context.Context, output *vmongodb.CloudMongoDbInstance) error {
	m.ID = types.StringPointerValue(output.CloudMongoDbInstanceNo)
	m.CloudMongoDbInstanceNo = types.StringPointerValue(output.CloudMongoDbInstanceNo)
	m.CloudMongoDbImageProductCode = types.StringPointerValue(output.CloudMongoDbImageProductCode)
	m.ShardCount = int32PointerValue(output.ShardCount)
	m.BackupFileRetentionPeriod = int32PointerValue(output.BackupFileRetentionPeriod)
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.CloudMongoDbServiceName = types.StringPointerValue(output.CloudMongoDbServiceName)
	m.ArbiterPort = int32PointerValue(output.ArbiterPort)
	m.MemberPort = int32PointerValue(output.MemberPort)
	m.MongosPort = int32PointerValue(output.MongosPort)
	m.ConfigPort = int32PointerValue(output.ConfigPort)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.AccessControlGroupNoList = acgList

	return nil
}

func int32PointerValue(value *int32) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Null()
	}
	newVal := int64(*value)

	return basetypes.NewInt64Value(newVal)
}

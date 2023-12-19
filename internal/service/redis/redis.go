package redis

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vredis"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ resource.Resource                = &redisResource{}
	_ resource.ResourceWithConfigure   = &redisResource{}
	_ resource.ResourceWithImportState = &redisResource{}
)

func NewRedisResource() resource.Resource {
	return &redisResource{}
}

type redisResource struct {
	config *conn.ProviderConfig
}

func (r *redisResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redis"
}

func (r *redisResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 15),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[가-힣A-Za-z0-9-]+$`), "Allows only hangeuls, alphabets, numbers, hyphen (-)."),
				},
			},
			"server_name_prefix": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 15),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`),
						"Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number.",
					),
				},
			},
			// Available only `gov` site
			"user_name": schema.StringAttribute{
				Optional: true,
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
			// Available only `gov` site
			"user_password": schema.StringAttribute{
				Optional: true,
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
			"id": framework.IDAttribute(),
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
			"config_group_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mode": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"CLUSTER", "SIMPLE"}...),
				},
			},
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shard_count": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(3, 10),
				},
				Description: "default: 3",
			},
			"shard_copy_count": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 4),
				},
				Description: "default: 0",
			},
			"is_ha": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: true",
			},
			"is_backup": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: false",
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 7),
				},
				Description: "default: 1",
			},
			"backup_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(5),
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
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(6379),
					),
				},
				Description: "default: 6379",
			},
			"backup_schedule": schema.StringAttribute{
				Computed: true,
			},
			"region_code": schema.StringAttribute{
				Computed: true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"redis_server_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"redis_server_instance_no": schema.StringAttribute{
							Computed: true,
						},
						"redis_server_name": schema.StringAttribute{
							Computed: true,
						},
						"redis_server_role": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"os_memory_size": schema.Int64Attribute{
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
				Computed: true,
			},
		},
	}
}

func (r *redisResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *redisResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan redisResourceModel

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

	reqParams := &vredis.CreateCloudRedisInstanceRequest{
		RegionCode:                 &r.config.RegionCode,
		CloudRedisServiceName:      plan.ServiceName.ValueStringPointer(),
		CloudRedisServerNamePrefix: plan.ServerNamePrefix.ValueStringPointer(),
		VpcNo:                      plan.VpcNo.ValueStringPointer(),
		SubnetNo:                   plan.SubnetNo.ValueStringPointer(),
		ConfigGroupNo:              plan.ConfigGroupNo.ValueStringPointer(),
		CloudRedisModeCode:         plan.Mode.ValueStringPointer(),
	}

	if !plan.ImageProductCode.IsNull() && !plan.ImageProductCode.IsUnknown() {
		reqParams.CloudRedisImageProductCode = plan.ImageProductCode.ValueStringPointer()
	}

	if !plan.ProductCode.IsNull() && !plan.ProductCode.IsUnknown() {
		reqParams.CloudRedisProductCode = plan.ProductCode.ValueStringPointer()
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		reqParams.CloudRedisPort = ncloud.Int32(int32(plan.Port.ValueInt64()))
	}

	if !plan.ShardCount.IsNull() && !plan.ShardCount.IsUnknown() {
		if *reqParams.CloudRedisModeCode == "SIMPLE" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`shard_count` invalid. Necessary only if the mode is CLUSTER.",
			)
			return
		}
		reqParams.ShardCount = ncloud.Int32(int32(plan.ShardCount.ValueInt64()))
	}

	if !plan.ShardCopyCount.IsNull() && !plan.ShardCopyCount.IsUnknown() {
		if *reqParams.CloudRedisModeCode == "SIMPLE" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`shard_copy_count` invalid. Necessary only if the mode is CLUSTER.",
			)
			return
		}
		reqParams.ShardCopyCount = ncloud.Int32(int32(plan.ShardCopyCount.ValueInt64()))
	}

	if !plan.IsHa.IsNull() && !plan.IsHa.IsUnknown() {
		if plan.IsHa.ValueBool() {
			if !plan.IsBackup.IsNull() && !plan.IsBackup.IsUnknown() && !plan.IsBackup.ValueBool() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is true, `is_backup` must be true or not be set",
				)
				return
			}

			if *reqParams.CloudRedisModeCode == "CLUSTER" && (plan.ShardCopyCount.IsNull() || plan.ShardCopyCount.IsUnknown() || *reqParams.ShardCopyCount == 0) {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					errors.New("when `is_ha` is true, `shard_copy_count` must be set to 1~4.").Error(),
				)
				return
			}
		} else {
			if !plan.ShardCopyCount.IsNull() && !plan.ShardCopyCount.IsUnknown() && *reqParams.ShardCopyCount > 0 {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					errors.New("when `is_ha` is false, `shard_copy_count` must be set to 0.").Error(),
				)
				return
			}
		}
		reqParams.IsHa = plan.IsHa.ValueBoolPointer()
	}

	if !plan.BackupFileRetentionPeriod.IsNull() && !plan.BackupFileRetentionPeriod.IsUnknown() {
		reqParams.BackupFileRetentionPeriod = ncloud.Int32(int32(plan.BackupFileRetentionPeriod.ValueInt64()))
	}

	if !plan.BackupTime.IsNull() && !plan.BackupTime.IsUnknown() {
		if plan.IsAutomaticBackup.IsNull() || plan.IsAutomaticBackup.IsUnknown() {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_automatic_backup` not set, `backup_time` must not be set",
			)
			return
		}
		reqParams.BackupTime = plan.BackupTime.ValueStringPointer()
	}

	if !plan.IsAutomaticBackup.IsNull() && !plan.IsAutomaticBackup.IsUnknown() {
		if plan.IsAutomaticBackup.ValueBool() && reqParams.BackupTime != nil {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_automatic_backup` is true, `backup_time` must not be set",
			)
			return
		}
		if !plan.IsAutomaticBackup.ValueBool() && reqParams.BackupTime == nil {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_automatic_backup` is false, `backup_time` must be set",
			)
			return
		}
		reqParams.IsAutomaticBackup = plan.IsAutomaticBackup.ValueBoolPointer()
	}

	if !plan.IsBackup.IsNull() && !plan.IsBackup.IsUnknown() {
		if !plan.IsBackup.ValueBool() {
			if reqParams.BackupFileRetentionPeriod != nil || reqParams.IsAutomaticBackup != nil || reqParams.BackupTime != nil {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_backup` is false, `backup_file_retention_period`, `is_automatic_backup`, `backup_time`  must not be set",
				)
				return
			}
		}
		reqParams.IsBackup = plan.IsBackup.ValueBoolPointer()
	}

	// Available only `gov` site
	if !plan.UserName.IsNull() {
		reqParams.CloudRedisUserName = plan.UserName.ValueStringPointer()
	}

	// Available only `gov` site
	if !plan.UserPassword.IsNull() {
		reqParams.CloudRedisUserPassword = plan.UserPassword.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateCloudRedisInstance reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vredis.V2Api.CreateCloudRedisInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateCloudRedisInstance response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudRedisInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	redisInstance := response.CloudRedisInstanceList[0]
	plan.ID = types.StringPointerValue(redisInstance.CloudRedisInstanceNo)

	output, err := waitRedisCreated(ctx, r.config, *redisInstance.CloudRedisInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *redisResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state redisResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetRedisDetail(ctx, r.config, state.ID.ValueString())
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

func (r *redisResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *redisResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state redisResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vredis.DeleteCloudRedisInstanceRequest{
		RegionCode:           &r.config.RegionCode,
		CloudRedisInstanceNo: state.ID.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteCloudRedis reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vredis.V2Api.DeleteCloudRedisInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "DeleteCloudRedis response="+common.MarshalUncheckedString(response))

	if err := waitRedisDeleted(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (r *redisResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GetRedisDetail(ctx context.Context, config *conn.ProviderConfig, no string) (*vredis.CloudRedisInstance, error) {
	reqParams := &vredis.GetCloudRedisInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudRedisInstanceNo: &no,
	}
	tflog.Info(ctx, "GetRedisDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vredis.V2Api.GetCloudRedisInstanceDetail(reqParams)
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "GetRedisDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudRedisInstanceList) < 1 {
		return nil, nil
	}

	return resp.CloudRedisInstanceList[0], nil
}

func waitRedisDeleted(ctx context.Context, config *conn.ProviderConfig, no string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetRedisDetail(ctx, config, no)
			if err != nil && !strings.Contains(err.Error(), `"returnCode": "5001017"`) {
				return 0, "", err
			}

			if resp == nil {
				return resp, "deleted", nil
			}

			return resp, "deleting", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for Redis (%s) to become termintaing: %s", no, err)
	}

	return nil
}

func waitRedisCreated(ctx context.Context, config *conn.ProviderConfig, no string) (*vredis.CloudRedisInstance, error) {
	var redisInstance *vredis.CloudRedisInstance
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetRedisDetail(ctx, config, no)
			redisInstance = resp
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "", fmt.Errorf("GetRedisDetail is nil")
			}

			if *resp.CloudRedisInstanceStatusName == "running" {
				return resp, "running", nil
			} else if *resp.CloudRedisInstanceStatusName == "creating" {
				return resp, "creating", nil
			} else if *resp.CloudRedisInstanceStatusName == "settingUp" {
				return resp, "settingUp", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("Error waiting for Redis (%s) to become available: %s", no, err)
	}

	return redisInstance, nil
}

type redisResourceModel struct {
	ServiceName               types.String `tfsdk:"service_name"`
	ServerNamePrefix          types.String `tfsdk:"server_name_prefix"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	ID                        types.String `tfsdk:"id"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ConfigGroupNo             types.String `tfsdk:"config_group_no"`
	Mode                      types.String `tfsdk:"mode"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	ShardCount                types.Int64  `tfsdk:"shard_count"`
	ShardCopyCount            types.Int64  `tfsdk:"shard_copy_count"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	IsAutomaticBackup         types.Bool   `tfsdk:"is_automatic_backup"`
	Port                      types.Int64  `tfsdk:"port"`
	BackupSchedule            types.String `tfsdk:"backup_schedule"`
	RegionCode                types.String `tfsdk:"region_code"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	RedisServerList           types.List   `tfsdk:"redis_server_list"`
}

type redisServer struct {
	RedisServerNo   types.String `tfsdk:"redis_server_instance_no"`
	RedisServerName types.String `tfsdk:"redis_server_name"`
	RedisServerRole types.String `tfsdk:"redis_server_role"`
	PrivateDomain   types.String `tfsdk:"private_domain"`
	MemorySize      types.Int64  `tfsdk:"memory_size"`
	OsMemorySize    types.Int64  `tfsdk:"os_memory_size"`
	Uptime          types.String `tfsdk:"uptime"`
	CreateDate      types.String `tfsdk:"create_date"`
}

func (r redisServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"redis_server_instance_no": types.StringType,
		"redis_server_name":        types.StringType,
		"redis_server_role":        types.StringType,
		"private_domain":           types.StringType,
		"memory_size":              types.Int64Type,
		"os_memory_size":           types.Int64Type,
		"uptime":                   types.StringType,
		"create_date":              types.StringType,
	}
}

func (r *redisResourceModel) refreshFromOutput(ctx context.Context, output *vredis.CloudRedisInstance) {
	r.ID = types.StringPointerValue(output.CloudRedisInstanceNo)
	r.ServiceName = types.StringPointerValue(output.CloudRedisServiceName)
	r.ServerNamePrefix = types.StringPointerValue(output.CloudRedisServerPrefix)
	r.VpcNo = types.StringPointerValue(output.CloudRedisServerInstanceList[0].VpcNo)
	r.SubnetNo = types.StringPointerValue(output.CloudRedisServerInstanceList[0].SubnetNo)
	r.ConfigGroupNo = types.StringPointerValue(output.ConfigGroupNo)
	r.Mode = types.StringPointerValue(output.Role.Code)
	r.ImageProductCode = types.StringPointerValue(output.CloudRedisImageProductCode)
	r.ProductCode = types.StringPointerValue(output.CloudRedisServerInstanceList[0].CloudRedisProductCode)
	r.IsHa = types.BoolPointerValue(output.IsHa)
	r.IsBackup = types.BoolPointerValue(output.IsBackup)
	r.BackupTime = types.StringPointerValue(output.BackupTime)
	r.BackupSchedule = types.StringPointerValue(output.BackupSchedule)
	r.RegionCode = types.StringPointerValue(output.CloudRedisServerInstanceList[0].RegionCode)

	if output.BackupFileRetentionPeriod != nil {
		r.BackupFileRetentionPeriod = types.Int64Value(int64(*output.BackupFileRetentionPeriod))
	}

	if output.CloudRedisPort != nil {
		r.Port = types.Int64Value(int64(*output.CloudRedisPort))
	}

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	r.AccessControlGroupNoList = acgList

	var serverList []redisServer
	for _, server := range output.CloudRedisServerInstanceList {
		redisServerInstance := redisServer{
			RedisServerNo:   types.StringPointerValue(server.CloudRedisServerInstanceNo),
			RedisServerName: types.StringPointerValue(server.CloudRedisServerName),
			RedisServerRole: types.StringPointerValue(server.CloudRedisServerRole.CodeName),
			PrivateDomain:   types.StringPointerValue(server.PrivateDomain),
			MemorySize:      types.Int64Value(*server.MemorySize),
			OsMemorySize:    types.Int64Value(*server.OsMemorySize),
			Uptime:          types.StringPointerValue(server.Uptime),
			CreateDate:      types.StringPointerValue(server.CreateDate),
		}
		serverList = append(serverList, redisServerInstance)
	}

	redisServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: redisServer{}.attrTypes()}, serverList)

	r.RedisServerList = redisServers
	return
}

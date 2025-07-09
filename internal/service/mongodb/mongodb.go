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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

func (m *mongodbResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb"
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
			"id": framework.IDAttribute(),
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
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"member_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"arbiter_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"mongos_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"config_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"shard_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(2, 5),
				},
			},
			"member_server_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(2, 7),
				},
			},
			"arbiter_server_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 1),
				},
			},
			"mongos_server_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(2, 5),
				},
			},
			"config_server_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(3, 7),
				},
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Optional: true,
				Computed: true,
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
			},
			"data_storage_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SSD", "HDD", "CB1", "CB2"}...),
				},
			},
			"member_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 65535),
					),
				},
			},
			"arbiter_port": schema.Int64Attribute{
				Computed: true,
			},
			"mongos_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 65535),
					),
				},
			},
			"config_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 65535),
					),
				},
			},
			"compress_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SNPP", "ZLIB", "ZSTD", "NONE"}...),
				},
			},
			"engine_version_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_code": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_code": schema.StringAttribute{
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
			"mongodb_server_list": schema.ListNestedAttribute{
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
						"cluster_role": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
							Computed: true,
						},
						"public_domain": schema.StringAttribute{
							Computed: true,
						},
						"replica_set_name": schema.StringAttribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
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

func (m *mongodbResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	m.config = config
}

func (m *mongodbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mongodbResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmongodb.CreateCloudMongoDbInstanceRequest{
		RegionCode:                   &m.config.RegionCode,
		CloudMongoDbServiceName:      plan.ServiceName.ValueStringPointer(),
		CloudMongoDbServerNamePrefix: plan.ServerNamePrefix.ValueStringPointer(),
		CloudMongoDbUserName:         plan.UserName.ValueStringPointer(),
		CloudMongoDbUserPassword:     plan.UserPassword.ValueStringPointer(),
		VpcNo:                        plan.VpcNo.ValueStringPointer(),
		SubnetNo:                     plan.SubnetNo.ValueStringPointer(),
		ClusterTypeCode:              plan.ClusterTypeCode.ValueStringPointer(),
		CompressCode:                 plan.CompressCode.ValueStringPointer(),
	}

	if !plan.ImageProductCode.IsNull() && !plan.ImageProductCode.IsUnknown() {
		reqParams.CloudMongoDbImageProductCode = plan.ImageProductCode.ValueStringPointer()
	}

	if !plan.MemberProductCode.IsNull() && !plan.MemberProductCode.IsUnknown() {
		reqParams.MemberProductCode = plan.MemberProductCode.ValueStringPointer()
	}

	if !plan.EngineVersionCode.IsNull() {
		reqParams.EngineVersionCode = plan.EngineVersionCode.ValueStringPointer()
	}

	if !plan.ArbiterProductCode.IsNull() && !plan.ArbiterProductCode.IsUnknown() {
		if *reqParams.ClusterTypeCode == "STAND_ALONE" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`arbiter_product_code` invalid. Necessary only if the cluster_type_code is SINGLE_REPLICA_SET or SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.ArbiterProductCode = plan.ArbiterProductCode.ValueStringPointer()
	}

	if !plan.MongosProductCode.IsNull() && !plan.MongosProductCode.IsUnknown() {
		if *reqParams.ClusterTypeCode != "SHARDED_CLUSTER" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`mongos_product_code` invalid. Necessary only if the cluster_type_code is SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.MongosProductCode = plan.MongosProductCode.ValueStringPointer()
	}

	if !plan.ConfigProductCode.IsNull() && !plan.ConfigProductCode.IsUnknown() {
		if *reqParams.ClusterTypeCode != "SHARDED_CLUSTER" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`config_product_code` invalid. Necessary only if the cluster_type_code is SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.ConfigProductCode = plan.ConfigProductCode.ValueStringPointer()
	}

	if !plan.ShardCount.IsNull() && !plan.ShardCount.IsUnknown() {
		if *reqParams.ClusterTypeCode != "SHARDED_CLUSTER" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`shard_count` invalid. Necessary only if the cluster_type_code is SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.ShardCount = ncloud.Int32(int32(plan.ShardCount.ValueInt64()))
	}

	if !plan.MemberServerCount.IsNull() && !plan.MemberServerCount.IsUnknown() {
		if *reqParams.ClusterTypeCode == "STAND_ALONE" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`member_server_count` invalid. Necessary only if the cluster_type_code is SINGLE_REPLICA_SET or SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.MemberServerCount = ncloud.Int32(int32(plan.MemberServerCount.ValueInt64()))
	}

	if !plan.ArbiterServerCount.IsNull() && !plan.ArbiterServerCount.IsUnknown() {
		if *reqParams.ClusterTypeCode == "STAND_ALONE" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`arbiter_server_count` invalid. Necessary only if the cluster_type_code is SINGLE_REPLICA_SET or SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.ArbiterServerCount = ncloud.Int32(int32(plan.ArbiterServerCount.ValueInt64()))
	}

	if !plan.MongosServerCount.IsNull() && !plan.MongosServerCount.IsUnknown() {
		if *reqParams.ClusterTypeCode != "SHARDED_CLUSTER" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`mongos_server_count` invalid. Necessary only if the cluster_type_code is SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.MongosServerCount = ncloud.Int32(int32(plan.MongosServerCount.ValueInt64()))
	}

	if !plan.ConfigServerCount.IsNull() && !plan.ConfigServerCount.IsUnknown() {
		if *reqParams.ClusterTypeCode != "SHARDED_CLUSTER" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`config_server_count` invalid. Necessary only if the cluster_type_code is SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.ConfigServerCount = ncloud.Int32(int32(plan.ConfigServerCount.ValueInt64()))
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

	if !plan.MongosPort.IsNull() && !plan.MongosPort.IsUnknown() {
		if *reqParams.ClusterTypeCode != "SHARDED_CLUSTER" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`mongos_port` invalid. Necessary only if the cluster_type_code is SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.MongosPort = ncloud.Int32(int32(plan.MongosPort.ValueInt64()))
	}

	if !plan.ConfigPort.IsNull() && !plan.ConfigPort.IsUnknown() {
		if *reqParams.ClusterTypeCode != "SHARDED_CLUSTER" {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`config_port` invalid. Necessary only if the cluster_type_code is SHARDED_CLUSTER.",
			)
			return
		}
		reqParams.ConfigPort = ncloud.Int32(int32(plan.ConfigPort.ValueInt64()))
	}

	if !plan.DataStorageType.IsNull() && !plan.DataStorageType.IsUnknown() {
		reqParams.DataStorageTypeCode = plan.DataStorageType.ValueStringPointer()
	}

	response, err := m.config.Client.Vmongodb.V2Api.CreateCloudMongoDbInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMongoDb response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudMongoDbInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	mongodbInstance := response.CloudMongoDbInstanceList[0]
	plan.ID = types.StringPointerValue(mongodbInstance.CloudMongoDbInstanceNo)

	output, err := waitMongoDbCreated(ctx, m.config, *mongodbInstance.CloudMongoDbInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (m *mongodbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mongodbResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetCloudMongoDbInstance(ctx, m.config, state.ID.ValueString())
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

func (m *mongodbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state mongodbResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.ConfigServerCount.Equal(state.ConfigServerCount) {
		reqParams := &vmongodb.ChangeCloudMongoDbConfigCountRequest{
			RegionCode:             &m.config.RegionCode,
			CloudMongoDbInstanceNo: state.ID.ValueStringPointer(),
			ConfigServerCount:      ncloud.Int32(int32(plan.ConfigServerCount.ValueInt64())),
		}
		tflog.Info(ctx, "ChangeCloudMongoDbConfigCount reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := m.config.Client.Vmongodb.V2Api.ChangeCloudMongoDbConfigCount(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeCloudMongoDbConfigCount response="+common.MarshalUncheckedString(response))

		if response == nil || len(response.CloudMongoDbInstanceList) < 1 {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		mongodbInstance := response.CloudMongoDbInstanceList[0]

		output, err := waitMongoDbUpdate(ctx, m.config, *mongodbInstance.CloudMongoDbInstanceNo)
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		state.refreshFromOutput(ctx, output)
	}

	if !plan.MongosServerCount.Equal(state.MongosServerCount) {
		reqParams := &vmongodb.ChangeCloudMongoDbMongosCountRequest{
			RegionCode:             &m.config.RegionCode,
			CloudMongoDbInstanceNo: state.ID.ValueStringPointer(),
			MongosServerCount:      ncloud.Int32(int32(plan.MongosServerCount.ValueInt64())),
		}
		tflog.Info(ctx, "ChangeCloudMongoDbMongosCount reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := m.config.Client.Vmongodb.V2Api.ChangeCloudMongoDbMongosCount(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeCloudMongoDbMongosCount response="+common.MarshalUncheckedString(response))

		if response == nil || len(response.CloudMongoDbInstanceList) < 1 {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		mongodbInstance := response.CloudMongoDbInstanceList[0]

		output, err := waitMongoDbUpdate(ctx, m.config, *mongodbInstance.CloudMongoDbInstanceNo)
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		state.refreshFromOutput(ctx, output)
	}

	if !plan.MemberServerCount.Equal(state.MemberServerCount) ||
		!plan.ArbiterServerCount.Equal(state.ArbiterServerCount) {
		reqParams := &vmongodb.ChangeCloudMongoDbSecondaryCountRequest{
			RegionCode:             &m.config.RegionCode,
			CloudMongoDbInstanceNo: state.ID.ValueStringPointer(),
			MemberServerCount:      ncloud.Int32(int32(plan.MemberServerCount.ValueInt64())),
			ArbiterServerCount:     ncloud.Int32(int32(plan.ArbiterServerCount.ValueInt64())),
		}
		tflog.Info(ctx, "ChangeCloudMongoDbSecondaryCount reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := m.config.Client.Vmongodb.V2Api.ChangeCloudMongoDbSecondaryCount(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeCloudMongoDbSecondaryCount response="+common.MarshalUncheckedString(response))

		if response == nil || len(response.CloudMongoDbInstanceList) < 1 {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		mongodbInstance := response.CloudMongoDbInstanceList[0]

		output, err := waitMongoDbUpdate(ctx, m.config, *mongodbInstance.CloudMongoDbInstanceNo)
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		state.refreshFromOutput(ctx, output)
	}

	if !plan.ShardCount.Equal(state.ShardCount) {
		reqParams := &vmongodb.ChangeCloudMongoDbShardCountRequest{
			RegionCode:             &m.config.RegionCode,
			CloudMongoDbInstanceNo: state.ID.ValueStringPointer(),
			ShardCount:             ncloud.Int32(int32(plan.ShardCount.ValueInt64())),
		}
		tflog.Info(ctx, "ChangeCloudMongoDbShardCount reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := m.config.Client.Vmongodb.V2Api.ChangeCloudMongoDbShardCount(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeCloudMongoDbShardCount response="+common.MarshalUncheckedString(response))

		if response == nil || len(response.CloudMongoDbInstanceList) < 1 {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		mongodbInstance := response.CloudMongoDbInstanceList[0]

		output, err := waitMongoDbUpdate(ctx, m.config, *mongodbInstance.CloudMongoDbInstanceNo)
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		state.refreshFromOutput(ctx, output)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (m *mongodbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mongodbResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmongodb.DeleteCloudMongoDbInstanceRequest{
		RegionCode:             &m.config.RegionCode,
		CloudMongoDbInstanceNo: state.ID.ValueStringPointer(),
	}
	tflog.Info(ctx, "DeleteMongoDb reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := m.config.Client.Vmongodb.V2Api.DeleteCloudMongoDbInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteMongoDb response="+common.MarshalUncheckedString(response))

	if err := waitMongoDbDeleted(ctx, m.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (s *mongodbResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GetCloudMongoDbInstance(ctx context.Context, config *conn.ProviderConfig, no string) (*vmongodb.CloudMongoDbInstance, error) {
	reqParams := &vmongodb.GetCloudMongoDbInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		CloudMongoDbInstanceNo: ncloud.String(no),
	}
	tflog.Info(ctx, "GetMongoDbDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vmongodb.V2Api.GetCloudMongoDbInstanceDetail(reqParams)
	// If the lookup result is 0 or already deleted, it will respond with a 400 error with a 5001017 return code.
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) {
		return nil, err
	}
	tflog.Info(ctx, "GetMongoDbDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudMongoDbInstanceList) < 1 || len(resp.CloudMongoDbInstanceList[0].CloudMongoDbServerInstanceList) < 1 {
		return nil, nil
	}

	return resp.CloudMongoDbInstanceList[0], nil
}

func waitMongoDbCreated(ctx context.Context, config *conn.ProviderConfig, id string) (*vmongodb.CloudMongoDbInstance, error) {
	var mongodbInstance *vmongodb.CloudMongoDbInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetCloudMongoDbInstance(ctx, config, id)
			mongodbInstance = instance
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("CloudMongoDbInstance is nil")
			}

			if *instance.CloudMongoDbInstanceStatusName == "creating" {
				return instance, "creating", nil
			} else if *instance.CloudMongoDbInstanceStatusName == "settingUp" {
				return instance, "settingUp", nil
			} else if *instance.CloudMongoDbInstanceStatusName == "running" {
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
		return nil, fmt.Errorf("error waiting for MongoDbInstance state to be \"CREAT\": %s", err)
	}

	return mongodbInstance, nil
}

func waitMongoDbUpdate(ctx context.Context, config *conn.ProviderConfig, id string) (*vmongodb.CloudMongoDbInstance, error) {
	var mongodbInstance *vmongodb.CloudMongoDbInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetCloudMongoDbInstance(ctx, config, id)
			mongodbInstance = instance
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("CloudMongoDbInstance is nil")
			}

			if *instance.CloudMongoDbInstanceStatusName == "creating" {
				return instance, "creating", nil
			} else if *instance.CloudMongoDbInstanceStatusName == "settingUp" {
				return instance, "settingUp", nil
			}

			for _, server := range instance.CloudMongoDbServerInstanceList {
				if *server.CloudMongoDbServerInstanceStatusName == "running" {
					continue
				} else if *server.CloudMongoDbServerInstanceStatusName == "creating" {
					return instance, "creating", nil
				} else if *server.CloudMongoDbServerInstanceStatusName == "settingUp" {
					return instance, "settingUp", nil
				} else {
					return 0, "", fmt.Errorf("error occurred while waiting to update mongodb")
				}
			}

			return instance, "running", nil
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("error waiting for MongoDbInstance state to be \"CREAT\": %s", err)
	}

	return mongodbInstance, nil
}

func waitMongoDbDeleted(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetCloudMongoDbInstance(ctx, config, id)
			if err != nil {
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
		Timeout:    2 * conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for mongodb (%s) to become terminating: %s", id, err)
	}

	return nil
}

type mongodbResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ServiceName               types.String `tfsdk:"service_name"`
	ServerNamePrefix          types.String `tfsdk:"server_name_prefix"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	ClusterTypeCode           types.String `tfsdk:"cluster_type_code"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	MemberProductCode         types.String `tfsdk:"member_product_code"`
	ArbiterProductCode        types.String `tfsdk:"arbiter_product_code"`
	MongosProductCode         types.String `tfsdk:"mongos_product_code"`
	ConfigProductCode         types.String `tfsdk:"config_product_code"`
	ShardCount                types.Int64  `tfsdk:"shard_count"`
	MemberServerCount         types.Int64  `tfsdk:"member_server_count"`
	ArbiterServerCount        types.Int64  `tfsdk:"arbiter_server_count"`
	MongosServerCount         types.Int64  `tfsdk:"mongos_server_count"`
	ConfigServerCount         types.Int64  `tfsdk:"config_server_count"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	DataStorageType           types.String `tfsdk:"data_storage_type"`
	MemberPort                types.Int64  `tfsdk:"member_port"`
	ArbiterPort               types.Int64  `tfsdk:"arbiter_port"`
	MongosPort                types.Int64  `tfsdk:"mongos_port"`
	ConfigPort                types.Int64  `tfsdk:"config_port"`
	CompressCode              types.String `tfsdk:"compress_code"`
	EngineVersionCode         types.String `tfsdk:"engine_version_code"`
	RegionCode                types.String `tfsdk:"region_code"`
	ZoneCode                  types.String `tfsdk:"zone_code"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MongoDbServerList         types.List   `tfsdk:"mongodb_server_list"`
}

type mongoServer struct {
	ServerNo        types.String `tfsdk:"server_instance_no"`
	ServerName      types.String `tfsdk:"server_name"`
	ServerRole      types.String `tfsdk:"server_role"`
	ClusterRole     types.String `tfsdk:"cluster_role"`
	ProductCode     types.String `tfsdk:"product_code"`
	PrivateDomain   types.String `tfsdk:"private_domain"`
	PublicDomain    types.String `tfsdk:"public_domain"`
	ReplicaSetName  types.String `tfsdk:"replica_set_name"`
	MemorySize      types.Int64  `tfsdk:"memory_size"`
	CpuCount        types.Int64  `tfsdk:"cpu_count"`
	DataStorageSize types.Int64  `tfsdk:"data_storage_size"`
	Uptime          types.String `tfsdk:"uptime"`
	CreateDate      types.String `tfsdk:"create_date"`
}

func (m mongoServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_instance_no": types.StringType,
		"server_name":        types.StringType,
		"server_role":        types.StringType,
		"cluster_role":       types.StringType,
		"product_code":       types.StringType,
		"private_domain":     types.StringType,
		"public_domain":      types.StringType,
		"replica_set_name":   types.StringType,
		"memory_size":        types.Int64Type,
		"cpu_count":          types.Int64Type,
		"data_storage_size":  types.Int64Type,
		"uptime":             types.StringType,
		"create_date":        types.StringType,
	}
}

func (m *mongodbResourceModel) refreshFromOutput(ctx context.Context, output *vmongodb.CloudMongoDbInstance) {
	m.ID = types.StringPointerValue(output.CloudMongoDbInstanceNo)
	m.ServiceName = types.StringPointerValue(output.CloudMongoDbServiceName)
	m.VpcNo = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].VpcNo)
	m.SubnetNo = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].SubnetNo)
	m.ImageProductCode = types.StringPointerValue(output.CloudMongoDbImageProductCode)
	m.EngineVersionCode = types.StringValue(common.ExtractEngineVersion(*output.EngineVersion))
	m.ShardCount = common.Int64ValueFromInt32(output.ShardCount)
	m.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.ArbiterPort = common.Int64FromInt32OrDefault(output.ArbiterPort)
	m.MemberPort = common.Int64FromInt32OrDefault(output.MemberPort)
	m.MongosPort = common.Int64FromInt32OrDefault(output.MongosPort)
	m.ConfigPort = common.Int64FromInt32OrDefault(output.ConfigPort)
	m.RegionCode = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].RegionCode)
	m.ZoneCode = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].ZoneCode)

	if output.CloudMongoDbServerInstanceList[0].DataStorageType != nil {
		m.DataStorageType = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].DataStorageType.Code)
	}
	if output.Compress != nil {
		m.CompressCode = types.StringPointerValue(output.Compress.Code)
	}

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.AccessControlGroupNoList = acgList

	var memberCount, arbiterCount, mongosCount, configCount int64
	var serverList []mongoServer
	for _, server := range output.CloudMongoDbServerInstanceList {
		mongoServerInstance := mongoServer{
			ServerNo:        types.StringPointerValue(server.CloudMongoDbServerInstanceNo),
			ServerName:      types.StringPointerValue(server.CloudMongoDbServerName),
			ProductCode:     types.StringPointerValue(server.CloudMongoDbProductCode),
			PrivateDomain:   types.StringPointerValue(server.PrivateDomain),
			PublicDomain:    types.StringPointerValue(server.PublicDomain),
			ReplicaSetName:  types.StringPointerValue(server.ReplicaSetName),
			MemorySize:      types.Int64PointerValue(server.MemorySize),
			CpuCount:        types.Int64PointerValue(server.CpuCount),
			DataStorageSize: types.Int64PointerValue(server.DataStorageSize),
			Uptime:          types.StringPointerValue(server.Uptime),
			CreateDate:      types.StringPointerValue(server.CreateDate),
		}

		if server.CloudMongoDbServerRole != nil {
			mongoServerInstance.ServerRole = types.StringPointerValue(server.CloudMongoDbServerRole.CodeName)
			if *server.CloudMongoDbServerRole.Code == "A" || *server.CloudMongoDbServerRole.Code == "MB" {
				m.MemberProductCode = types.StringPointerValue(server.CloudMongoDbProductCode)
				memberCount++
			} else if *server.CloudMongoDbServerRole.Code == "AB" {
				m.ArbiterProductCode = types.StringPointerValue(server.CloudMongoDbProductCode)
				arbiterCount++
			} else if *server.CloudMongoDbServerRole.Code == "RT" {
				m.MongosProductCode = types.StringPointerValue(server.CloudMongoDbProductCode)
				mongosCount++
			} else if *server.CloudMongoDbServerRole.Code == "C" {
				m.ConfigProductCode = types.StringPointerValue(server.CloudMongoDbProductCode)
				configCount++
			}
		}
		if server.ClusterRole != nil {
			mongoServerInstance.ClusterRole = types.StringPointerValue(server.ClusterRole.Code)
		}
		serverList = append(serverList, mongoServerInstance)
	}

	if output.ClusterType != nil {
		m.ClusterTypeCode = types.StringPointerValue(output.ClusterType.Code)
		if *output.ClusterType.Code == "SHARDED_CLUSTER" {
			memberCount = memberCount / int64(*output.ShardCount)
			if arbiterCount > 0 {
				arbiterCount = arbiterCount / int64(*output.ShardCount)
			}
		}
	}
	m.MemberServerCount = types.Int64Value(memberCount)
	m.ArbiterServerCount = types.Int64Value(arbiterCount)
	m.MongosServerCount = types.Int64Value(mongosCount)
	m.ConfigServerCount = types.Int64Value(configCount)
	m.ArbiterProductCode = common.StringFrameworkOrDefault(m.ArbiterProductCode)
	m.MongosProductCode = common.StringFrameworkOrDefault(m.MongosProductCode)
	m.ConfigProductCode = common.StringFrameworkOrDefault(m.ConfigProductCode)

	mongoServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mongoServer{}.attrTypes()}, serverList)

	m.MongoDbServerList = mongoServers

}

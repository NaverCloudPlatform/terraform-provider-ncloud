package redis

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vredis"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &redisDataSource{}
	_ datasource.DataSourceWithConfigure = &redisDataSource{}
)

func NewRedisDataSource() datasource.DataSource {
	return &redisDataSource{}
}

type redisDataSource struct {
	config *conn.ProviderConfig
}

func (r *redisDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redis"
}

func (r *redisDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("service_name"),
					),
				},
			},
			"service_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("id"),
					),
				},
			},
			"server_name_prefix": schema.StringAttribute{
				Computed: true,
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"config_group_no": schema.StringAttribute{
				Computed: true,
			},
			"mode": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"product_code": schema.StringAttribute{
				Computed: true,
			},
			"is_ha": schema.BoolAttribute{
				Computed: true,
			},
			"is_backup": schema.BoolAttribute{
				Computed: true,
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Computed: true,
			},
			"backup_time": schema.StringAttribute{
				Computed: true,
			},
			"backup_schedule": schema.StringAttribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
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

func (r *redisDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *redisDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data redisDataSourceModel
	var redisId string

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		redisId = data.ID.ValueString()
	}

	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() {
		reqParams := &vredis.GetCloudRedisInstanceListRequest{
			RegionCode:            &r.config.RegionCode,
			CloudRedisServiceName: data.ServiceName.ValueStringPointer(),
		}
		tflog.Info(ctx, "GetRedisList reqParams="+common.MarshalUncheckedString(reqParams))

		listResp, err := r.config.Client.Vredis.V2Api.GetCloudRedisInstanceList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "GetRedisList response="+common.MarshalUncheckedString(listResp))

		if listResp == nil || len(listResp.CloudRedisInstanceList) < 1 {
			resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
			return
		}
		redisId = *listResp.CloudRedisInstanceList[0].CloudRedisInstanceNo
	}

	output, err := GetRedisDetail(ctx, r.config, redisId)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	data.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type redisDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ServiceName               types.String `tfsdk:"service_name"`
	ServerNamePrefix          types.String `tfsdk:"server_name_prefix"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ConfigGroupNo             types.String `tfsdk:"config_group_no"`
	Mode                      types.String `tfsdk:"mode"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	BackupSchedule            types.String `tfsdk:"backup_schedule"`
	Port                      types.Int64  `tfsdk:"port"`
	RegionCode                types.String `tfsdk:"region_code"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	RedisServerList           types.List   `tfsdk:"redis_server_list"`
}

type redisServerDataSourceModel struct {
	RedisServerNo   types.String `tfsdk:"redis_server_instance_no"`
	RedisServerName types.String `tfsdk:"redis_server_name"`
	RedisServerRole types.String `tfsdk:"redis_server_role"`
	PrivateDomain   types.String `tfsdk:"private_domain"`
	MemorySize      types.Int64  `tfsdk:"memory_size"`
	OsMemorySize    types.Int64  `tfsdk:"os_memory_size"`
	Uptime          types.String `tfsdk:"uptime"`
	CreateDate      types.String `tfsdk:"create_date"`
}

func (r redisServerDataSourceModel) attrTypes() map[string]attr.Type {
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

func (r *redisDataSourceModel) refreshFromOutput(ctx context.Context, output *vredis.CloudRedisInstance) {
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
		r.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	}

	if output.CloudRedisPort != nil {
		r.Port = common.Int64ValueFromInt32(output.CloudRedisPort)
	}

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	r.AccessControlGroupNoList = acgList

	var serverList []redisServerDataSourceModel
	for _, server := range output.CloudRedisServerInstanceList {
		redisServerInstance := redisServerDataSourceModel{
			RedisServerNo:   types.StringPointerValue(server.CloudRedisServerInstanceNo),
			RedisServerName: types.StringPointerValue(server.CloudRedisServerName),
			RedisServerRole: types.StringPointerValue(server.CloudRedisServerRole.CodeName),
			PrivateDomain:   types.StringPointerValue(server.PrivateDomain),
			MemorySize:      types.Int64PointerValue(server.MemorySize),
			OsMemorySize:    types.Int64PointerValue(server.OsMemorySize),
			Uptime:          types.StringPointerValue(server.Uptime),
			CreateDate:      types.StringPointerValue(server.CreateDate),
		}
		serverList = append(serverList, redisServerInstance)
	}
	listValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: redisServerDataSourceModel{}.attrTypes()}, serverList)

	r.RedisServerList = listValue
}

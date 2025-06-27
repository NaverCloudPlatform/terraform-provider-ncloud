package mongodb

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mongodbDataSource{}
	_ datasource.DataSourceWithConfigure = &mongodbDataSource{}
)

func NewMongoDbDataSource() datasource.DataSource {
	return &mongodbDataSource{}
}

type mongodbDataSource struct {
	config *conn.ProviderConfig
}

func (m *mongodbDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb"
}

func (m *mongodbDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"cluster_type_code": schema.StringAttribute{
				Computed: true,
			},
			"engine_version": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Computed: true,
			},
			"backup_time": schema.StringAttribute{
				Computed: true,
			},
			"shard_count": schema.Int64Attribute{
				Computed: true,
			},
			"data_storage_type": schema.StringAttribute{
				Computed: true,
			},
			"member_port": schema.Int64Attribute{
				Computed: true,
			},
			"arbiter_port": schema.Int64Attribute{
				Computed: true,
			},
			"mongos_port": schema.Int64Attribute{
				Computed: true,
			},
			"config_port": schema.Int64Attribute{
				Computed: true,
			},
			"compress_code": schema.StringAttribute{
				Computed: true,
			},
			"region_code": schema.StringAttribute{
				Computed: true,
			},
			"zone_code": schema.StringAttribute{
				Computed: true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
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

func (m *mongodbDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mongodbDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mongodbDataSourceModel
	var mongodbId string

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		mongodbId = data.ID.ValueString()
	}

	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() {
		reqParams := &vmongodb.GetCloudMongoDbInstanceListRequest{
			RegionCode:              &m.config.RegionCode,
			CloudMongoDbServiceName: data.ServiceName.ValueStringPointer(),
		}
		tflog.Info(ctx, "GetMongoDbList reqParams="+common.MarshalUncheckedString(reqParams))

		listResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbInstanceList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "GetMongoDbList response="+common.MarshalUncheckedString(listResp))

		if listResp == nil || len(listResp.CloudMongoDbInstanceList) < 1 {
			resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
			return
		}
		mongodbId = *listResp.CloudMongoDbInstanceList[0].CloudMongoDbInstanceNo
	}

	output, err := GetCloudMongoDbInstance(ctx, m.config, mongodbId)
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

type mongodbDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ServiceName               types.String `tfsdk:"service_name"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ClusterTypeCode           types.String `tfsdk:"cluster_type_code"`
	EngineVersion             types.String `tfsdk:"engine_version"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	ShardCount                types.Int64  `tfsdk:"shard_count"`
	DataStorageType           types.String `tfsdk:"data_storage_type"`
	MemberPort                types.Int64  `tfsdk:"member_port"`
	ArbiterPort               types.Int64  `tfsdk:"arbiter_port"`
	MongosPort                types.Int64  `tfsdk:"mongos_port"`
	ConfigPort                types.Int64  `tfsdk:"config_port"`
	CompressCode              types.String `tfsdk:"compress_code"`
	RegionCode                types.String `tfsdk:"region_code"`
	ZoneCode                  types.String `tfsdk:"zone_code"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MongoDbServerList         types.List   `tfsdk:"mongodb_server_list"`
}

type mongodbServerDataSourceModel struct {
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

func (m mongodbServerDataSourceModel) attrTypes() map[string]attr.Type {
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

func (m *mongodbDataSourceModel) refreshFromOutput(ctx context.Context, output *vmongodb.CloudMongoDbInstance) {
	m.ID = types.StringPointerValue(output.CloudMongoDbInstanceNo)
	m.ServiceName = types.StringPointerValue(output.CloudMongoDbServiceName)
	m.VpcNo = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].VpcNo)
	m.SubnetNo = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].SubnetNo)
	m.ClusterTypeCode = types.StringPointerValue(output.ClusterType.Code)
	m.EngineVersion = types.StringPointerValue(output.EngineVersion)
	m.ImageProductCode = types.StringPointerValue(output.CloudMongoDbImageProductCode)
	m.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.ShardCount = common.Int64ValueFromInt32(output.ShardCount)
	m.DataStorageType = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].DataStorageType.Code)
	m.MemberPort = common.Int64ValueFromInt32(output.MemberPort)
	m.ArbiterPort = common.Int64ValueFromInt32(output.ArbiterPort)
	m.MongosPort = common.Int64ValueFromInt32(output.MongosPort)
	m.ConfigPort = common.Int64ValueFromInt32(output.ConfigPort)
	m.CompressCode = types.StringPointerValue(output.Compress.Code)
	m.RegionCode = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].RegionCode)
	m.ZoneCode = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].ZoneCode)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.AccessControlGroupNoList = acgList

	var serverList []mongodbServerDataSourceModel
	for _, server := range output.CloudMongoDbServerInstanceList {
		mongoServerInstance := mongodbServerDataSourceModel{
			ServerNo:        types.StringPointerValue(server.CloudMongoDbServerInstanceNo),
			ServerName:      types.StringPointerValue(server.CloudMongoDbServerName),
			ServerRole:      types.StringPointerValue(server.CloudMongoDbServerRole.CodeName),
			ClusterRole:     types.StringPointerValue(server.ClusterRole.Code),
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
		serverList = append(serverList, mongoServerInstance)
	}

	mongoServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mongodbServerDataSourceModel{}.attrTypes()}, serverList)

	m.MongoDbServerList = mongoServers

}

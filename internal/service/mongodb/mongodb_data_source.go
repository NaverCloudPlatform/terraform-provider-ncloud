package mongodb

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
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

func (m *mongodbDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"service_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"cluster_type_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"engine_version": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Optional: true,
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
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"server_instance_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"server_instance_no": schema.StringAttribute{
							Computed: true,
						},
						"server_name": schema.StringAttribute{
							Computed: true,
						},
						"cluster_role": schema.StringAttribute{
							Computed: true,
						},
						"server_role": schema.StringAttribute{
							Computed: true,
						},
						"region_code": schema.StringAttribute{
							Computed: true,
						},
						"vpc_no": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"uptime": schema.StringAttribute{
							Computed: true,
						},
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
							Computed: true,
						},
						"public_domain": schema.StringAttribute{
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
						"used_data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"replica_set_name": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mongodbDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !m.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not supported Classic",
			"mongodb data source does not supported in classic",
		)
		return
	}

	var data mongodbDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmongodb.GetCloudMongoDbInstanceListRequest{
		RegionCode: &m.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.CloudMongoDbInstanceNoList = []*string{data.ID.ValueStringPointer()}
	}
	if !data.CloudMongoDbServiceName.IsNull() && !data.CloudMongoDbServiceName.IsUnknown() {
		reqParams.CloudMongoDbServiceName = data.CloudMongoDbServiceName.ValueStringPointer()
	}

	tflog.Info(ctx, "GetMongoDbList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	mongodbResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbInstanceList(reqParams)

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetCloudMongoDbList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetMongoDbList response", map[string]any{
		"mongodbResponse": common.MarshalUncheckedString(mongodbResp),
	})

	mongodbId := mongodbResp.CloudMongoDbInstanceList[0].CloudMongoDbInstanceNo
	detailReqParams := &vmongodb.GetCloudMongoDbInstanceDetailRequest{
		RegionCode:             &m.config.RegionCode,
		CloudMongoDbInstanceNo: mongodbId,
	}
	mongodbDetailResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbInstanceDetail(detailReqParams)

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetCloudMongoDbDetailList",
			fmt.Sprintf("error: %s, detailReqParams: %s", err.Error(), common.MarshalUncheckedString(detailReqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetMongoDbDetailList response", map[string]any{
		"mongodbDetailResponse": common.MarshalUncheckedString(mongodbDetailResp),
	})

	mongodbList, diags := flattenMongoDbs(ctx, mongodbDetailResp.CloudMongoDbInstanceList, m.config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, mongodbList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMongodbList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenMongoDbs(ctx context.Context, mongodbs []*vmongodb.CloudMongoDbInstance, config *conn.ProviderConfig) ([]*mongodbDataSourceModel, diag.Diagnostics) {
	var outputs []*mongodbDataSourceModel

	for _, m := range mongodbs {
		var output mongodbDataSourceModel

		diags := output.refreshFromOutput(ctx, m, config)
		if diags.HasError() {
			return nil, diags
		}

		outputs = append(outputs, &output)
	}

	return outputs, nil
}

type mongodbDataSourceModel struct {
	ID                             types.String `tfsdk:"id"`
	CloudMongoDbServiceName        types.String `tfsdk:"service_name"`
	VpcNo                          types.String `tfsdk:"vpc_no"`
	SubnetNo                       types.String `tfsdk:"subnet_no"`
	ClusterType                    types.String `tfsdk:"cluster_type_code"`
	EngineVersion                  types.String `tfsdk:"engine_version"`
	CloudMongoDbImageProductCode   types.String `tfsdk:"image_product_code"`
	BackupFileRetentionPeriod      types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                     types.String `tfsdk:"backup_time"`
	ShardCount                     types.Int64  `tfsdk:"shard_count"`
	AccessControlGroupNoList       types.List   `tfsdk:"access_control_group_no_list"`
	CloudMongoDbServerInstanceList types.List   `tfsdk:"server_instance_list"`
	Filters                        types.Set    `tfsdk:"filter"`
}

type mongodbServer struct {
	CloudMongoDbServerInstanceNo types.String `tfsdk:"server_instance_no"`
	CloudMongoDbServerName       types.String `tfsdk:"server_name"`
	ClusterRole                  types.String `tfsdk:"cluster_role"`
	CloudMongoDbServerRole       types.String `tfsdk:"server_role"`
	RegionCode                   types.String `tfsdk:"region_code"`
	VpcNo                        types.String `tfsdk:"vpc_no"`
	SubnetNo                     types.String `tfsdk:"subnet_no"`
	Uptime                       types.String `tfsdk:"uptime"`
	ZoneCode                     types.String `tfsdk:"zone_code"`
	PrivateDomain                types.String `tfsdk:"private_domain"`
	PublicDomain                 types.String `tfsdk:"public_domain"`
	MemorySize                   types.Int64  `tfsdk:"memory_size"`
	CpuCount                     types.Int64  `tfsdk:"cpu_count"`
	DataStorageSize              types.Int64  `tfsdk:"data_storage_size"`
	UsedDataStorageSize          types.Int64  `tfsdk:"used_data_storage_size"`
	ProductCode                  types.String `tfsdk:"product_code"`
	ReplicaSetName               types.String `tfsdk:"replica_set_name"`
	DataStorageType              types.String `tfsdk:"data_storage_type"`
}

func (m mongodbServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_instance_no":     types.StringType,
		"server_name":            types.StringType,
		"cluster_role":           types.StringType,
		"server_role":            types.StringType,
		"region_code":            types.StringType,
		"vpc_no":                 types.StringType,
		"subnet_no":              types.StringType,
		"uptime":                 types.StringType,
		"zone_code":              types.StringType,
		"private_domain":         types.StringType,
		"public_domain":          types.StringType,
		"memory_size":            types.Int64Type,
		"cpu_count":              types.Int64Type,
		"data_storage_size":      types.Int64Type,
		"used_data_storage_size": types.Int64Type,
		"product_code":           types.StringType,
		"replica_set_name":       types.StringType,
		"data_storage_type":      types.StringType,
	}
}

func (d *mongodbDataSourceModel) refreshFromOutput(ctx context.Context, output *vmongodb.CloudMongoDbInstance, config *conn.ProviderConfig) diag.Diagnostics {
	d.ID = types.StringPointerValue(output.CloudMongoDbInstanceNo)
	d.VpcNo = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].VpcNo)
	d.SubnetNo = types.StringPointerValue(output.CloudMongoDbServerInstanceList[0].SubnetNo)
	d.CloudMongoDbServiceName = types.StringPointerValue(output.CloudMongoDbServiceName)
	d.CloudMongoDbImageProductCode = types.StringPointerValue(output.CloudMongoDbImageProductCode)
	d.BackupFileRetentionPeriod = int32PointerValue(output.BackupFileRetentionPeriod)
	d.BackupTime = types.StringPointerValue(output.BackupTime)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	d.AccessControlGroupNoList = acgList

	var serverList []mongodbServer
	for _, server := range output.CloudMongoDbServerInstanceList {
		mongodbServerInstance := mongodbServer{
			CloudMongoDbServerInstanceNo: types.StringPointerValue(server.CloudMongoDbServerInstanceNo),
			CloudMongoDbServerName:       types.StringPointerValue(server.CloudMongoDbServerName),
			ClusterRole:                  types.StringPointerValue(server.ClusterRole.Code),
			CloudMongoDbServerRole:       types.StringPointerValue(server.CloudMongoDbServerRole.Code),
			RegionCode:                   types.StringPointerValue(server.RegionCode),
			VpcNo:                        types.StringPointerValue(server.VpcNo),
			SubnetNo:                     types.StringPointerValue(server.SubnetNo),
			Uptime:                       types.StringPointerValue(server.Uptime),
			ZoneCode:                     types.StringPointerValue(server.ZoneCode),
			PrivateDomain:                types.StringPointerValue(server.PrivateDomain),
			PublicDomain:                 types.StringPointerValue(server.PublicDomain),
			MemorySize:                   types.Int64Value(*server.MemorySize),
			CpuCount:                     types.Int64Value(*server.CpuCount),
			DataStorageSize:              types.Int64Value(*server.DataStorageSize),
			ProductCode:                  types.StringPointerValue(server.CloudMongoDbProductCode),
			ReplicaSetName:               types.StringPointerValue(server.ReplicaSetName),
			DataStorageType:              types.StringPointerValue(server.DataStorageType.Code),
		}

		if server.UsedDataStorageSize != nil {
			mongodbServerInstance.UsedDataStorageSize = types.Int64Value(*server.UsedDataStorageSize)
		}

		serverList = append(serverList, mongodbServerInstance)
	}

	mongodbServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mongodbServer{}.attrTypes()}, serverList)
	d.CloudMongoDbServerInstanceList = mongodbServers

	return nil
}

package mysql

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ datasource.DataSource              = &mysqlDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlDataSource{}
)

func NewMysqlDataSource() datasource.DataSource {
	return &mysqlDataSource{}
}

type mysqlDataSource struct {
	config *conn.ProviderConfig
}

func (m *mysqlDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mysqlDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql"
}

func (m *mysqlDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"name_prefix": schema.StringAttribute{
				Computed: true,
			},
			"user_name": schema.StringAttribute{
				Computed: true,
			},
			"user_password": schema.StringAttribute{
				Computed: true,
			},
			"host_ip": schema.StringAttribute{
				Computed: true,
			},
			"database_name": schema.StringAttribute{
				Computed: true,
			},
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"engine_version_code": schema.StringAttribute{
				Computed: true,
			},
			"data_storage_type_code": schema.StringAttribute{
				Computed: true,
			},
			"is_ha": schema.BoolAttribute{
				Computed: true,
			},
			"is_multi_zone": schema.BoolAttribute{
				Computed: true,
			},
			"is_storage_encryption": schema.BoolAttribute{
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
			"is_automatic_backup": schema.BoolAttribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
				Computed: true,
			},
			"standby_master_subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"instance_no": schema.StringAttribute{
				Computed: true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"mysql_config_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"create_date": schema.StringAttribute{
				Computed: true,
			},
			"mysql_server_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mysql_server_instance_no": schema.StringAttribute{
							Computed: true,
						},
						"mysql_server_name": schema.StringAttribute{
							Computed: true,
						},
						"mysql_server_role": schema.StringAttribute{
							Computed: true,
						},
						"mysql_server_product_code": schema.StringAttribute{
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
						"public_domain": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_type": schema.StringAttribute{
							Computed: true,
						},
						"is_storage_encryption": schema.BoolAttribute{
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
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mysqlDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !m.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not Supported Classic",
			"mysql data source does not supported in classic",
		)
		return
	}

	var data mysqlDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmysql.GetCloudMysqlInstanceListRequest{
		RegionCode: &m.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.CloudMysqlInstanceNoList = []*string{data.ID.ValueStringPointer()}
	}
	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() {
		reqParams.CloudMysqlServiceName = data.ServiceName.ValueStringPointer()
	}

	tflog.Info(ctx, "GetMysqlList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	mysqlResp, err := m.config.Client.Vmysql.V2Api.GetCloudMysqlInstanceList(reqParams)
	mysqlId := mysqlResp.CloudMysqlInstanceList[0].CloudMysqlInstanceNo
	detailReqParams := &vmysql.GetCloudMysqlInstanceDetailRequest{
		RegionCode:           &m.config.RegionCode,
		CloudMysqlInstanceNo: mysqlId,
	}
	mysqlDetailResp, err := m.config.Client.Vmysql.V2Api.GetCloudMysqlInstanceDetail(detailReqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMysqlList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetMysqlList response", map[string]any{
		"mysqlResponse": common.MarshalUncheckedString(mysqlResp),
	})

	mysqlList, diags := flattenMysqls(ctx, mysqlDetailResp.CloudMysqlInstanceList, m.config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	filteredList := common.FilterModels(ctx, data.Filters, mysqlList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMysqlList result vaildation: result more thane one",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenMysqls(ctx context.Context, mysqls []*vmysql.CloudMysqlInstance, config *conn.ProviderConfig) ([]*mysqlDataSourceModel, diag.Diagnostics) {
	var outputs []*mysqlDataSourceModel

	for _, v := range mysqls {
		var output mysqlDataSourceModel

		output.refreshFromOutput(ctx, v, config)
		outputs = append(outputs, &output)
	}

	return outputs, nil
}

func (m *mysqlDataSourceModel) refreshFromOutput(ctx context.Context, output *vmysql.CloudMysqlInstance, config *conn.ProviderConfig) diag.Diagnostics {
	m.ID = types.StringPointerValue(output.CloudMysqlInstanceNo)
	m.ServiceName = types.StringPointerValue(output.CloudMysqlServiceName)
	m.EngineVersionCode = types.StringPointerValue(output.EngineVersion)
	m.IsHa = types.BoolPointerValue(output.IsHa)
	m.IsMultiZone = types.BoolPointerValue(output.IsMultiZone)
	m.IsBackup = types.BoolPointerValue(output.IsBackup)
	m.BackupFileRetentionPeriod = types.Int64Value(int64(*output.BackupFileRetentionPeriod))
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.Port = types.Int64Value(int64(*output.CloudMysqlPort))
	m.ImageProductCode = types.StringPointerValue(output.CloudMysqlImageProductCode)
	m.CreateDate = types.StringPointerValue(output.CreateDate)
	m.InstanceNo = types.StringPointerValue(output.CloudMysqlInstanceNo)
	m.VpcNo = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].VpcNo)
	m.DataStorageTypeCode = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].DataStorageType.Code)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	configList, _ := types.ListValueFrom(ctx, types.StringType, output.CloudMysqlConfigList)
	m.AccessControlGroupNoList = acgList
	m.MysqlConfigList = configList

	var serverList []mysqlServerDataSourceModel
	for _, server := range output.CloudMysqlServerInstanceList {
		mysqlServerInstance := mysqlServerDataSourceModel{
			MysqlServerInstanceNo:  types.StringPointerValue(server.CloudMysqlServerInstanceNo),
			MysqlServerName:        types.StringPointerValue(server.CloudMysqlServerName),
			MysqlServerRole:        types.StringPointerValue(server.CloudMysqlServerRole.Code),
			MysqlServerProductCode: types.StringPointerValue(server.CloudMysqlProductCode),
			RegionCode:             types.StringPointerValue(server.RegionCode),
			ZoneCode:               types.StringPointerValue(server.ZoneCode),
			VpcNo:                  types.StringPointerValue(server.VpcNo),
			SubnetNo:               types.StringPointerValue(server.SubnetNo),
			IsPublicSubnet:         types.BoolPointerValue(server.IsPublicSubnet),
			PrivateDomain:          types.StringPointerValue(server.PrivateDomain),
			DataStorageType:        types.StringPointerValue(server.DataStorageType.Code),
			IsStorageEncryption:    types.BoolPointerValue(server.IsStorageEncryption),
			DataStorageSize:        types.Int64Value(*server.DataStorageSize),
			CpuCount:               types.Int64Value(int64(*server.CpuCount)),
			MemorySize:             types.Int64Value(*server.MemorySize),
			Uptime:                 types.StringPointerValue(server.Uptime),
			CreateDate:             types.StringPointerValue(server.CreateDate),
		}
		if server.PublicDomain != nil {
			mysqlServerInstance.PublicDomain = types.StringPointerValue(server.PublicDomain)
		}
		if server.UsedDataStorageSize != nil {
			mysqlServerInstance.UsedDataStorageSize = types.Int64Value(*server.UsedDataStorageSize)
		}
		serverList = append(serverList, mysqlServerInstance)
	}
	listValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlServer{}.attrTypes()}, serverList)
	m.MysqlServerList = listValue
	return nil

}

type mysqlDataSourceModel struct {
	ServiceName               types.String `tfsdk:"service_name"`
	NamePrefix                types.String `tfsdk:"name_prefix"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	HostIp                    types.String `tfsdk:"host_ip"`
	DatabaseName              types.String `tfsdk:"database_name"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	EngineVersionCode         types.String `tfsdk:"engine_version_code"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type_code"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsMultiZone               types.Bool   `tfsdk:"is_multi_zone"`
	IsStorageEncryption       types.Bool   `tfsdk:"is_storage_encryption"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	IsAutomaticBackup         types.Bool   `tfsdk:"is_automatic_backup"`
	Port                      types.Int64  `tfsdk:"port"`
	StandbyMasterSubnetNo     types.String `tfsdk:"standby_master_subnet_no"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	InstanceNo                types.String `tfsdk:"instance_no"`
	ID                        types.String `tfsdk:"id"`
	CreateDate                types.String `tfsdk:"create_date"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MysqlConfigList           types.List   `tfsdk:"mysql_config_list"`
	MysqlServerList           types.List   `tfsdk:"mysql_server_list"`
	Filters                   types.Set    `tfsdk:"filter"`
}

type mysqlServerDataSourceModel struct {
	MysqlServerInstanceNo  types.String `tfsdk:"mysql_server_instance_no"`
	MysqlServerName        types.String `tfsdk:"mysql_server_name"`
	MysqlServerRole        types.String `tfsdk:"mysql_server_role"`
	MysqlServerProductCode types.String `tfsdk:"mysql_server_product_code"`
	RegionCode             types.String `tfsdk:"region_code"`
	ZoneCode               types.String `tfsdk:"zone_code"`
	VpcNo                  types.String `tfsdk:"vpc_no"`
	SubnetNo               types.String `tfsdk:"subnet_no"`
	IsPublicSubnet         types.Bool   `tfsdk:"is_public_subnet"`
	PublicDomain           types.String `tfsdk:"public_domain"`
	PrivateDomain          types.String `tfsdk:"private_domain"`
	DataStorageType        types.String `tfsdk:"data_storage_type"`
	IsStorageEncryption    types.Bool   `tfsdk:"is_storage_encryption"`
	DataStorageSize        types.Int64  `tfsdk:"data_storage_size"`
	UsedDataStorageSize    types.Int64  `tfsdk:"used_data_storage_size"`
	CpuCount               types.Int64  `tfsdk:"cpu_count"`
	MemorySize             types.Int64  `tfsdk:"memory_size"`
	Uptime                 types.String `tfsdk:"uptime"`
	CreateDate             types.String `tfsdk:"create_date"`
}

func (m mysqlServerDataSourceModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mysql_server_instance_no":  types.StringType,
		"mysql_server_name":         types.StringType,
		"mysql_server_role":         types.StringType,
		"mysql_server_product_code": types.StringType,
		"region_code":               types.StringType,
		"zone_code":                 types.StringType,
		"vpc_no":                    types.StringType,
		"subnet_no":                 types.StringType,
		"is_public_subnet":          types.BoolType,
		"public_domain":             types.StringType,
		"private_domain":            types.StringType,
		"data_storage_type":         types.StringType,
		"is_storage_encryption":     types.BoolType,
		"data_storage_size":         types.Int64Type,
		"used_data_storage_size":    types.Int64Type,
		"cpu_count":                 types.Int64Type,
		"memory_size":               types.Int64Type,
		"uptime":                    types.StringType,
		"create_date":               types.StringType,
	}
}

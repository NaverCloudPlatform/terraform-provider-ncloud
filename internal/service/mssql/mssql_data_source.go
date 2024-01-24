package mssql

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmssql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mssqlDataSource{}
	_ datasource.DataSourceWithConfigure = &mssqlDataSource{}
)

func NewMssqlDataSource() datasource.DataSource {
	return &mssqlDataSource{}
}

type mssqlDataSource struct {
	config *conn.ProviderConfig
}

func (m *mssqlDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mssql"
}

func (m *mssqlDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"region_code": schema.StringAttribute{
				Computed: true,
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"data_storage_type": schema.StringAttribute{
				Computed: true,
			},
			"is_ha": schema.BoolAttribute{
				Computed: true,
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Computed: true,
			},
			"backup_time": schema.StringAttribute{
				Computed: true,
			},
			"config_group_no": schema.StringAttribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
				Computed: true,
			},
			"engine_version": schema.StringAttribute{
				Computed: true,
			},
			"character_set_name": schema.StringAttribute{
				Computed: true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
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

func (m *mssqlDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mssqlDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mssqlDataSourceModel
	var mssqlId string

	if !m.config.SupportVPC {
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
		mssqlId = data.ID.ValueString()
	}

	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() {
		reqParams := &vmssql.GetCloudMssqlInstanceListRequest{
			RegionCode:            &m.config.RegionCode,
			CloudMssqlServiceName: data.ServiceName.ValueStringPointer(),
		}
		tflog.Info(ctx, "GetMssqlList reqParams="+common.MarshalUncheckedString(reqParams))

		listResp, err := m.config.Client.Vmssql.V2Api.GetCloudMssqlInstanceList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "GetMssqlList response="+common.MarshalUncheckedString(listResp))

		if listResp == nil || len(listResp.CloudMssqlInstanceList) < 1 {
			resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
			return
		}
		mssqlId = *listResp.CloudMssqlInstanceList[0].CloudMssqlInstanceNo
	}

	output, err := GetMssqlInstance(ctx, m.config, mssqlId)
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

type mssqlDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ServiceName               types.String `tfsdk:"service_name"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	ConfigGroupNo             types.String `tfsdk:"config_group_no"`
	Port                      types.Int64  `tfsdk:"port"`
	EngineVersion             types.String `tfsdk:"engine_version"`
	CharacterSetName          types.String `tfsdk:"character_set_name"`
	RegionCode                types.String `tfsdk:"region_code"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MssqlServerList           types.List   `tfsdk:"mssql_server_list"`
}

type mssqlServerDataSourceModel struct {
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

func (m mssqlServerDataSourceModel) attrTypes() map[string]attr.Type {
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

func (m *mssqlDataSourceModel) refreshFromOutput(ctx context.Context, output *vmssql.CloudMssqlInstance) {
	m.ID = types.StringPointerValue(output.CloudMssqlInstanceNo)
	m.ServiceName = types.StringPointerValue(output.CloudMssqlServiceName)
	m.ImageProductCode = types.StringPointerValue(output.CloudMssqlImageProductCode)
	m.DataStorageTypeCode = types.StringPointerValue(output.CloudMssqlServerInstanceList[0].DataStorageType.Code)
	m.IsHa = types.BoolPointerValue(output.IsHa)
	m.BackupFileRetentionPeriod = types.Int64Value(int64(*output.BackupFileRetentionPeriod))
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.ConfigGroupNo = types.StringPointerValue(output.ConfigGroupNo)
	m.Port = types.Int64Value(int64(*output.CloudMssqlPort))
	m.EngineVersion = types.StringPointerValue(output.EngineVersion)
	m.CharacterSetName = types.StringPointerValue(output.DbCollation)
	m.RegionCode = types.StringPointerValue(output.CloudMssqlServerInstanceList[0].RegionCode)
	m.VpcNo = types.StringPointerValue(output.CloudMssqlServerInstanceList[0].VpcNo)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.AccessControlGroupNoList = acgList

	var serverList []mssqlServerDataSourceModel
	for _, server := range output.CloudMssqlServerInstanceList {
		mssqlServerInstance := mssqlServerDataSourceModel{
			ServerInstanceNo: types.StringPointerValue(server.CloudMssqlServerInstanceNo),
			ServerName:       types.StringPointerValue(server.CloudMssqlServerName),
			ServerRole:       types.StringPointerValue(server.CloudMssqlServerRole.Code),
			ZoneCode:         types.StringPointerValue(server.ZoneCode),
			SubnetNo:         types.StringPointerValue(server.SubnetNo),
			ProductCode:      types.StringPointerValue(server.CloudMssqlProductCode),
			IsPublicSubnet:   types.BoolPointerValue(server.IsPublicSubnet),
			PrivateDomain:    types.StringPointerValue(server.PrivateDomain),
			DataStorageSize:  types.Int64Value(*server.DataStorageSize),
			CpuCount:         types.Int64Value(int64(*server.CpuCount)),
			MemorySize:       types.Int64Value(*server.MemorySize),
			Uptime:           types.StringPointerValue(server.Uptime),
			CreateDate:       types.StringPointerValue(server.CreateDate),
		}

		if server.PublicDomain != nil {
			mssqlServerInstance.PublicDomain = types.StringPointerValue(server.PublicDomain)
		}

		if server.UsedDataStorageSize != nil {
			mssqlServerInstance.UsedDataStorageSize = types.Int64Value(*server.UsedDataStorageSize)
		}
		serverList = append(serverList, mssqlServerInstance)
	}

	listValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mssqlServerDataSourceModel{}.attrTypes()}, serverList)
	m.MssqlServerList = listValue
}

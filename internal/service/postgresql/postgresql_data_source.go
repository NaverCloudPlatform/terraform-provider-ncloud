package postgresql

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ datasource.DataSource              = &postgresqlDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlDataSource{}
)

func NewPostgresqlDataSource() datasource.DataSource {
	return &postgresqlDataSource{}
}

type postgresqlDataSource struct {
	config *conn.ProviderConfig
}

func (d *postgresqlDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql"
}

func (d *postgresqlDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"product_code": schema.StringAttribute{
				Computed: true,
			},
			"data_storage_type_code": schema.StringAttribute{
				Computed: true,
			},
			"client_cidr": schema.StringAttribute{
				Computed: true,
			},
			"is_multi_zone": schema.BoolAttribute{
				Computed: true,
			},
			"is_ha": schema.BoolAttribute{
				Computed: true,
			},
			"is_storage_encryption": schema.BoolAttribute{
				Computed: true,
			},
			"is_backup": schema.BoolAttribute{
				Computed: true,
			},
			"backup_time": schema.StringAttribute{
				Computed: true,
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Computed: true,
			},
			"backup_file_storage_count": schema.Int64Attribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
				Computed: true,
			},
			"engine_version": schema.StringAttribute{
				Computed: true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"postgresql_config_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"postgresql_server_list": schema.ListNestedAttribute{
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
						"private_ip": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"used_data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
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

func (d *postgresqlDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.config = config
}

func (d *postgresqlDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data postgresqlDataSourceModel
	var postgresqlId string

	if !d.config.SupportVPC {
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
		postgresqlId = data.ID.ValueString()
	}

	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() {
		reqParams := &vpostgresql.GetCloudPostgresqlInstanceListRequest{
			RegionCode:                 &d.config.RegionCode,
			CloudPostgresqlServiceName: data.ServiceName.ValueStringPointer(),
		}
		tflog.Info(ctx, "GetPostgresqlList reqParams="+common.MarshalUncheckedString(reqParams))

		listResp, err := d.config.Client.Vpostgresql.V2Api.GetCloudPostgresqlInstanceList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "GetPostgresqlList response="+common.MarshalUncheckedString(listResp))

		if listResp == nil || len(listResp.CloudPostgresqlInstanceList) < 1 {
			resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
			return
		}
		postgresqlId = *listResp.CloudPostgresqlInstanceList[0].CloudPostgresqlInstanceNo
	}

	output, err := GetPostgresqlInstance(ctx, d.config, postgresqlId)
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

type postgresqlDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ServiceName               types.String `tfsdk:"service_name"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	RegionCode                types.String `tfsdk:"region_code"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	IsMultiZone               types.Bool   `tfsdk:"is_multi_zone"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsStorageEncryption       types.Bool   `tfsdk:"is_storage_encryption"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupTime                types.String `tfsdk:"backup_time"`
	BackupFileStorageCount    types.Int64  `tfsdk:"backup_file_storage_count"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	Port                      types.Int64  `tfsdk:"port"`
	ClientCidr                types.String `tfsdk:"client_cidr"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type_code"`
	EngineVersion             types.String `tfsdk:"engine_version"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	PostgresqlConfigList      types.List   `tfsdk:"postgresql_config_list"`
	PostgresqlServerList      types.List   `tfsdk:"postgresql_server_list"`
}

func (d *postgresqlDataSourceModel) refreshFromOutput(ctx context.Context, output *vpostgresql.CloudPostgresqlInstance) {
	d.ID = types.StringPointerValue(output.CloudPostgresqlInstanceNo)
	d.ServiceName = types.StringPointerValue(output.CloudPostgresqlServiceName)
	d.ImageProductCode = types.StringPointerValue(output.CloudPostgresqlImageProductCode)
	d.VpcNo = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].VpcNo)
	d.SubnetNo = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].SubnetNo)
	d.RegionCode = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].RegionCode)
	d.IsMultiZone = types.BoolPointerValue(output.IsMultiZone)
	d.IsHa = types.BoolPointerValue(output.IsHa)
	d.IsBackup = types.BoolPointerValue(output.IsBackup)
	d.BackupTime = types.StringPointerValue(output.BackupTime)
	d.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	d.Port = common.Int64ValueFromInt32(output.CloudPostgresqlPort)
	d.EngineVersion = types.StringPointerValue(output.EngineVersion)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	d.AccessControlGroupNoList = acgList
	configList, _ := types.ListValueFrom(ctx, types.StringType, output.CloudPostgresqlConfigList)
	d.PostgresqlConfigList = configList

	var serverList []postgresqlServer
	for _, server := range output.CloudPostgresqlServerInstanceList {
		postgresqlServerInstance := postgresqlServer{
			ServerInstanceNo: types.StringPointerValue(server.CloudPostgresqlServerInstanceNo),
			ServerName:       types.StringPointerValue(server.CloudPostgresqlServerName),
			ServerRole:       types.StringPointerValue(server.CloudPostgresqlServerRole.Code),
			IsPublicSubnet:   types.BoolPointerValue(server.IsPublicSubnet),
			PublicDomain:     types.StringPointerValue(server.PublicDomain),
			PrivateIp:        types.StringPointerValue(server.PrivateIp),
			DataStorageSize:  types.Int64PointerValue(server.DataStorageSize),
			ProductCode:      types.StringPointerValue(server.CloudPostgresqlProductCode),
			MemorySize:       types.Int64PointerValue(server.MemorySize),
			CpuCount:         common.Int64ValueFromInt32(server.CpuCount),
			Uptime:           types.StringPointerValue(server.Uptime),
			CreateDate:       types.StringPointerValue(server.CreateDate),
		}

		if server.PublicDomain != nil {
			postgresqlServerInstance.PublicDomain = types.StringPointerValue(server.PublicDomain)
		}

		if server.UsedDataStorageSize != nil {
			postgresqlServerInstance.UsedDataStorageSize = types.Int64PointerValue(server.UsedDataStorageSize)
		}
		serverList = append(serverList, postgresqlServerInstance)
	}

	postgresqlServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: postgresqlServer{}.attrTypes()}, serverList)

	d.PostgresqlServerList = postgresqlServers
}

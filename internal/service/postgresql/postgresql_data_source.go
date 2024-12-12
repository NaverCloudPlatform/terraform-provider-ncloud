package postgresql

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"generation_code": schema.StringAttribute{
				Computed: true,
			},
			"engine_version": schema.StringAttribute{
				Computed: true,
			},
			"ha": schema.BoolAttribute{
				Computed: true,
			},
			"multi_zone": schema.BoolAttribute{
				Computed: true,
			},
			"data_storage_type": schema.StringAttribute{
				Computed: true,
			},
			"storage_encryption": schema.BoolAttribute{
				Computed: true,
			},
			"backup": schema.BoolAttribute{
				Computed: true,
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Computed: true,
			},
			"backup_time": schema.StringAttribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
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
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"public_subnet": schema.BoolAttribute{
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

	if diags := data.refreshFromOutput(ctx, output); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type postgresqlDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ServiceName               types.String `tfsdk:"service_name"`
	RegionCode                types.String `tfsdk:"region_code"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	GenerationCode            types.String `tfsdk:"generation_code"`
	EngineVersion             types.String `tfsdk:"engine_version"`
	Ha                        types.Bool   `tfsdk:"ha"`
	MultiZone                 types.Bool   `tfsdk:"multi_zone"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type"`
	StorageEncryption         types.Bool   `tfsdk:"storage_encryption"`
	Backup                    types.Bool   `tfsdk:"backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	Port                      types.Int64  `tfsdk:"port"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	PostgresqlConfigList      types.List   `tfsdk:"postgresql_config_list"`
	PostgresqlServerList      types.List   `tfsdk:"postgresql_server_list"`
}

func (d *postgresqlDataSourceModel) refreshFromOutput(ctx context.Context, output *vpostgresql.CloudPostgresqlInstance) diag.Diagnostics {
	d.ID = types.StringPointerValue(output.CloudPostgresqlInstanceNo)
	d.ServiceName = types.StringPointerValue(output.CloudPostgresqlServiceName)
	d.RegionCode = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].RegionCode)
	d.VpcNo = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].VpcNo)
	d.ImageProductCode = types.StringPointerValue(output.CloudPostgresqlImageProductCode)
	d.GenerationCode = types.StringPointerValue(output.GenerationCode)
	d.EngineVersion = types.StringPointerValue(output.EngineVersion)
	d.Ha = types.BoolPointerValue(output.IsHa)
	d.MultiZone = types.BoolPointerValue(output.IsMultiZone)
	d.DataStorageTypeCode = types.StringPointerValue(common.GetCodePtrByCommonCode(output.CloudPostgresqlServerInstanceList[0].DataStorageType))
	d.StorageEncryption = types.BoolPointerValue(output.CloudPostgresqlServerInstanceList[0].IsStorageEncryption)
	d.Backup = types.BoolPointerValue(output.IsBackup)
	d.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	d.BackupTime = types.StringPointerValue(output.BackupTime)
	d.Port = common.Int64ValueFromInt32(output.CloudPostgresqlPort)

	acgList, diags := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	if diags.HasError() {
		return diags
	}
	d.AccessControlGroupNoList = acgList
	configList, diags := types.ListValueFrom(ctx, types.StringType, output.CloudPostgresqlConfigList)
	if diags.HasError() {
		return diags
	}
	d.PostgresqlConfigList = configList

	d.PostgresqlServerList, diags = listValueFromPostgresqlServerInatanceList(ctx, output.CloudPostgresqlServerInstanceList)

	return diags
}

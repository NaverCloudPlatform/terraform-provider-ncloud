package mysql

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
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
	_ datasource.DataSource              = &mysqlDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlDataSource{}
)

func NewMysqlDataSource() datasource.DataSource {
	return &mysqlDataSource{}
}

type mysqlDataSource struct {
	config *conn.ProviderConfig
}

func (d *mysqlDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql"
}

func (d *mysqlDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"port": schema.Int64Attribute{
				Computed: true,
			},
			"engine_version_code": schema.StringAttribute{
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
			"mysql_server_list": schema.ListNestedAttribute{
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

func (d *mysqlDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *mysqlDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mysqlDataSourceModel
	var mysqlId string

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
		mysqlId = data.ID.ValueString()
	}

	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() {
		reqParams := &vmysql.GetCloudMysqlInstanceListRequest{
			RegionCode:            &d.config.RegionCode,
			CloudMysqlServiceName: data.ServiceName.ValueStringPointer(),
		}
		tflog.Info(ctx, "GetMysqlList reqParams="+common.MarshalUncheckedString(reqParams))

		listResp, err := d.config.Client.Vmysql.V2Api.GetCloudMysqlInstanceList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "GetMysqlList response="+common.MarshalUncheckedString(listResp))

		if listResp == nil || len(listResp.CloudMysqlInstanceList) < 1 {
			resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
			return
		}
		mysqlId = *listResp.CloudMysqlInstanceList[0].CloudMysqlInstanceNo
	}

	output, err := GetMysqlInstance(ctx, d.config, mysqlId)
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

type mysqlDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ServiceName               types.String `tfsdk:"service_name"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsMultiZone               types.Bool   `tfsdk:"is_multi_zone"`
	IsStorageEncryption       types.Bool   `tfsdk:"is_storage_encryption"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	Port                      types.Int64  `tfsdk:"port"`
	EngineVersionCode         types.String `tfsdk:"engine_version_code"`
	RegionCode                types.String `tfsdk:"region_code"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MysqlConfigList           types.List   `tfsdk:"mysql_config_list"`
	MysqlServerList           types.List   `tfsdk:"mysql_server_list"`
}

func (d *mysqlDataSourceModel) refreshFromOutput(ctx context.Context, output *vmysql.CloudMysqlInstance) diag.Diagnostics {
	d.ID = types.StringPointerValue(output.CloudMysqlInstanceNo)
	d.ServiceName = types.StringPointerValue(output.CloudMysqlServiceName)
	d.ImageProductCode = types.StringPointerValue(output.CloudMysqlImageProductCode)
	d.DataStorageTypeCode = types.StringPointerValue(common.GetCodePtrByCommonCode(output.CloudMysqlServerInstanceList[0].DataStorageType))
	d.IsHa = types.BoolPointerValue(output.IsHa)
	d.IsMultiZone = types.BoolPointerValue(output.IsMultiZone)
	d.IsStorageEncryption = types.BoolPointerValue(output.CloudMysqlServerInstanceList[0].IsStorageEncryption)
	d.IsBackup = types.BoolPointerValue(output.IsBackup)
	d.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	d.BackupTime = types.StringPointerValue(output.BackupTime)
	d.Port = common.Int64ValueFromInt32(output.CloudMysqlPort)
	d.EngineVersionCode = types.StringPointerValue(output.EngineVersion)
	d.RegionCode = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].RegionCode)
	d.VpcNo = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].VpcNo)

	acgList, diags := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	if diags.HasError() {
		return diags
	}
	d.AccessControlGroupNoList = acgList
	configList, diags := types.ListValueFrom(ctx, types.StringType, output.CloudMysqlConfigList)
	if diags.HasError() {
		return diags
	}
	d.MysqlConfigList = configList

	d.MysqlServerList, diags = listValueFromMysqlServerList(ctx, output.CloudMysqlServerInstanceList)

	return diags
}

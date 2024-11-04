package mysql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mysqlDatabasesDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlDatabasesDataSource{}
)

func NewMysqlDatabasesDataSource() datasource.DataSource {
	return &mysqlDatabasesDataSource{}
}

type mysqlDatabasesDataSource struct {
	config *conn.ProviderConfig
}

func (d *mysqlDatabasesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_databases"
}

func (d *mysqlDatabasesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *mysqlDatabasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"mysql_instance_no": schema.StringAttribute{
				Required: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"mysql_database_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
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

func (d *mysqlDatabasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mysqlDatabasesDataSourceModel

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

	output, err := GetMysqlDatabaseAllList(ctx, d.config, data.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	mysqlDbList := flattenMysqlDatabases(output)
	fillteredList := common.FilterModels(ctx, data.Filters, mysqlDbList)
	if diags := data.refreshFromOutput(ctx, fillteredList, data.MysqlInstanceNo.ValueString()); diags != nil {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertDbsToJsonStruct(data.MysqlDatabaseList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetMysqlDatabaseAllList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*vmysql.CloudMysqlDatabase, error) {
	var allDbs []*vmysql.CloudMysqlDatabase
	pageNo := int32(0)
	pageSize := int32(100)
	hasMore := true

	for hasMore {
		reqParams := &vmysql.GetCloudMysqlDatabaseListRequest{
			RegionCode:           &config.RegionCode,
			CloudMysqlInstanceNo: ncloud.String(id),
			PageNo:               ncloud.Int32(pageNo),
			PageSize:             ncloud.Int32(pageSize),
		}
		tflog.Info(ctx, "GetMysqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

		resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlDatabaseList(reqParams)
		if err != nil {
			return nil, err
		}

		if resp == nil {
			break
		}

		allDbs = append(allDbs, resp.CloudMysqlDatabaseList...)

		hasMore = len(resp.CloudMysqlDatabaseList) == int(pageSize)
		pageNo++
	}

	if len(allDbs) == 0 {
		return nil, nil
	}

	tflog.Info(ctx, "GetMysqlDatabaseList response="+common.MarshalUncheckedString(allDbs))

	return allDbs, nil
}

type mysqlDatabasesDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	MysqlInstanceNo   types.String `tfsdk:"mysql_instance_no"`
	MysqlDatabaseList types.List   `tfsdk:"mysql_database_list"`
	OutputFile        types.String `tfsdk:"output_file"`
	Filters           types.Set    `tfsdk:"filter"`
}

type mysqlDb struct {
	DatabaseName types.String `tfsdk:"name"`
}

type mysqlDbToJsonConvert struct {
	DatabaseName string `json:"name"`
}

func (d mysqlDb) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}

func convertDbsToJsonStruct(dbs []attr.Value) ([]mysqlDbToJsonConvert, error) {
	var dbToConvert = []mysqlDbToJsonConvert{}

	for _, db := range dbs {
		dbJson := mysqlDbToJsonConvert{}
		if err := json.Unmarshal([]byte(db.String()), &dbJson); err != nil {
			return nil, err
		}
		dbToConvert = append(dbToConvert, dbJson)
	}

	return dbToConvert, nil
}

func flattenMysqlDatabases(list []*vmysql.CloudMysqlDatabase) []*mysqlDb {
	var outputs []*mysqlDb

	for _, v := range list {
		var output mysqlDb
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (d *mysqlDatabasesDataSourceModel) refreshFromOutput(ctx context.Context, output []*mysqlDb, instance string) diag.Diagnostics {
	d.ID = types.StringValue(instance)
	d.MysqlInstanceNo = types.StringValue(instance)
	dbListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlDb{}.attrTypes()}, output)
	if diags.HasError() {
		return diags
	}
	d.MysqlDatabaseList = dbListValue

	return nil
}

func (d *mysqlDb) refreshFromOutput(output *vmysql.CloudMysqlDatabase) {
	d.DatabaseName = types.StringPointerValue(output.DatabaseName)
}

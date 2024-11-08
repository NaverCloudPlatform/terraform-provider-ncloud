package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
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
	_ datasource.DataSource              = &postgresqlDatabasesDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlDatabasesDataSource{}
)

func NewPostgresqlDatabasesDataSource() datasource.DataSource {
	return &postgresqlDatabasesDataSource{}
}

type postgresqlDatabasesDataSource struct {
	config *conn.ProviderConfig
}

func (d *postgresqlDatabasesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_databases"
}

func (d *postgresqlDatabasesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *postgresqlDatabasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"postgresql_instance_no": schema.StringAttribute{
				Required: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"postgresql_database_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"owner": schema.StringAttribute{
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

func (d *postgresqlDatabasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data postgresqlDatabasesDataSourceModel

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

	output, err := GetPostgresqlDatabaseAllList(ctx, d.config, data.PostgresqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	postgresqlDbList := flattenPostgresqlDatabases(output)
	filteredList := common.FilterModels(ctx, data.Filters, postgresqlDbList)
	if diags := data.refreshFromOutput(ctx, filteredList, data.PostgresqlInstanceNo.ValueString()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertDbsToJsonStruct(data.PostgresqlDatabaseList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func GetPostgresqlDatabaseAllList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*vpostgresql.CloudPostgresqlDatabase, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlDatabaseListRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(id),
	}
	tflog.Info(ctx, "GetPostgresqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vpostgresql.V2Api.GetCloudPostgresqlDatabaseList(reqParams)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "GetPostgresqlDatabaseList response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudPostgresqlDatabaseList) < 2 {
		return nil, nil
	}

	return resp.CloudPostgresqlDatabaseList[1:], nil
}

type postgresqlDatabasesDataSourceModel struct {
	ID                     types.String `tfsdk:"id"`
	PostgresqlInstanceNo   types.String `tfsdk:"postgresql_instance_no"`
	PostgresqlDatabaseList types.List   `tfsdk:"postgresql_database_list"`
	OutputFile             types.String `tfsdk:"output_file"`
	Filters                types.Set    `tfsdk:"filter"`
}

func (r *postgresqlDatabasesDataSourceModel) refreshFromOutput(ctx context.Context, output []*postgresqlDb, instanceNo string) diag.Diagnostics {
	r.ID = types.StringValue(instanceNo)
	r.PostgresqlInstanceNo = types.StringValue(instanceNo)
	dbList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: postgresqlDb{}.attrTypes()}, output)
	if diags.HasError() {
		return diags
	}
	r.PostgresqlDatabaseList = dbList

	return nil
}

func (d *postgresqlDb) refreshFromOutput(output *vpostgresql.CloudPostgresqlDatabase) {
	d.DatabaseName = types.StringPointerValue(output.DatabaseName)
	d.Owner = types.StringPointerValue(output.Owner)
}

func convertDbsToJsonStruct(dbs []attr.Value) ([]postgresqlDbToJsonConvert, error) {
	var dbToConvert = []postgresqlDbToJsonConvert{}

	for _, db := range dbs {
		dbJson := postgresqlDbToJsonConvert{}
		if err := json.Unmarshal([]byte(db.String()), &dbJson); err != nil {
			return nil, err
		}
		dbToConvert = append(dbToConvert, dbJson)
	}

	return dbToConvert, nil
}

func flattenPostgresqlDatabases(list []*vpostgresql.CloudPostgresqlDatabase) []*postgresqlDb {
	var outputs []*postgresqlDb

	for _, v := range list {
		var output postgresqlDb
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

type postgresqlDb struct {
	DatabaseName types.String `tfsdk:"name"`
	Owner        types.String `tfsdk:"owner"`
}

type postgresqlDbToJsonConvert struct {
	DatabaseName string `json:"name"`
	Owner        string `json:"owner"`
}

func (d postgresqlDb) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"owner": types.StringType,
	}
}

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
	_ datasource.DataSource              = &postgresqlUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlUsersDataSource{}
)

func NewPostgresqlUsersDataSource() datasource.DataSource {
	return &postgresqlUsersDataSource{}
}

type postgresqlUsersDataSource struct {
	config *conn.ProviderConfig
}

func (d *postgresqlUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_users"
}

func (d *postgresqlUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *postgresqlUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"postgresql_user_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"client_cidr": schema.StringAttribute{
							Computed: true,
						},
						"replication_role": schema.BoolAttribute{
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

func (d *postgresqlUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data postgresqlUsersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetPostgresqlUserAllList(ctx, d.config, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	postgresqlUserList := flattenPostgresqlUsers(output)
	fillteredList := common.FilterModels(ctx, data.Filters, postgresqlUserList)
	if diags := data.refreshFromOutput(ctx, fillteredList, data.ID.ValueString()); diags.HasError() {
		resp.Diagnostics.AddError("READIG EROROR", "refreshFromOutput error")
		return
	}

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertUsersToJsonStruct(data.PostgresqlUserList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetPostgresqlUserAllList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*vpostgresql.CloudPostgresqlUser, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlUserListRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(id),
	}
	tflog.Info(ctx, "GetPostgresqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vpostgresql.V2Api.GetCloudPostgresqlUserList(reqParams)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "GetPostgresqlUserList response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudPostgresqlUserList) < 1 {
		return nil, nil
	}

	return common.ReverseList(resp.CloudPostgresqlUserList), nil
}

type postgresqlUsersDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	PostgresqlUserList types.List   `tfsdk:"postgresql_user_list"`
	OutputFile         types.String `tfsdk:"output_file"`
	Filters            types.Set    `tfsdk:"filter"`
}

type postgresqlUser struct {
	UserName        types.String `tfsdk:"name"`
	ClientCidr      types.String `tfsdk:"client_cidr"`
	ReplicationRole types.Bool   `tfsdk:"replication_role"`
}

type postgresqlUserToJsonConvert struct {
	UserName        string `json:"name"`
	ClientCidr      string `json:"client_cidr"`
	ReplicationRole bool   `json:"replication_role"`
}

func (r postgresqlUser) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":             types.StringType,
		"client_cidr":      types.StringType,
		"replication_role": types.BoolType,
	}
}

func convertUsersToJsonStruct(users []attr.Value) ([]postgresqlUserToJsonConvert, error) {
	var userToConvert = []postgresqlUserToJsonConvert{}

	for _, user := range users {
		userJson := postgresqlUserToJsonConvert{}
		if err := json.Unmarshal([]byte(user.String()), &userJson); err != nil {
			return nil, err
		}
		userToConvert = append(userToConvert, userJson)
	}

	return userToConvert, nil
}

func flattenPostgresqlUsers(list []*vpostgresql.CloudPostgresqlUser) []*postgresqlUser {
	var outputs []*postgresqlUser

	for _, v := range list {
		var output postgresqlUser
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (d *postgresqlUser) refreshFromOutput(output *vpostgresql.CloudPostgresqlUser) {
	d.UserName = types.StringPointerValue(output.UserName)
	d.ClientCidr = types.StringPointerValue(output.ClientCidr)
	d.ReplicationRole = types.BoolPointerValue(output.IsReplicationRole)
}

func (d *postgresqlUsersDataSourceModel) refreshFromOutput(ctx context.Context, output []*postgresqlUser, instance string) diag.Diagnostics {
	d.ID = types.StringValue(instance)
	userListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: postgresqlUser{}.attrTypes()}, output)
	if diags.HasError() {
		return diags
	}

	d.PostgresqlUserList = userListValue

	return diags
}

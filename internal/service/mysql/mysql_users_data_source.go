package mysql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ datasource.DataSource              = &mysqlUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlUsersDataSource{}
)

func NewMysqlUsersDataSource() datasource.DataSource {
	return &mysqlUsersDataSource{}
}

type mysqlUsersDataSource struct {
	config *conn.ProviderConfig
}

func (d *mysqlUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_users"
}

func (d *mysqlUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *mysqlUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("mysql_instance_no"),
					),
				},
			},
			"mysql_instance_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("id"),
					),
				},
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"mysql_user_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"host_ip": schema.StringAttribute{
							Computed: true,
						},
						"authority": schema.StringAttribute{
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

func (d *mysqlUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mysqlUsersDataSourceModel
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
		data.MysqlInstanceNo = data.ID
	}

	if !data.MysqlInstanceNo.IsNull() && !data.MysqlInstanceNo.IsUnknown() {
		mysqlId = data.MysqlInstanceNo.ValueString()
		data.ID = data.MysqlInstanceNo
	}

	output, err := GetMysqlUserList(ctx, d.config, mysqlId)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	mysqlUserList := flattenMysqlUsers(output)
	fillteredList := common.FilterModels(ctx, data.Filters, mysqlUserList)
	data.refreshFromOutput(ctx, fillteredList, data.MysqlInstanceNo.ValueString())

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertUsersToJsonStruct(data.MysqlUserList.Elements()); err != nil {
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

type mysqlUsersDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	MysqlInstanceNo types.String `tfsdk:"mysql_instance_no"`
	MysqlUserList   types.List   `tfsdk:"mysql_user_list"`
	OutputFile      types.String `tfsdk:"output_file"`
	Filters         types.Set    `tfsdk:"filter"`
}

type mysqlUser struct {
	UserName  types.String `tfsdk:"name"`
	HostIp    types.String `tfsdk:"host_ip"`
	Authority types.String `tfsdk:"authority"`
}

type mysqlUserToJsonConvert struct {
	UserName  string `json:"name"`
	HostIp    string `json:"host_ip"`
	Authority string `json:"authority"`
}

func (d mysqlUser) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":      types.StringType,
		"host_ip":   types.StringType,
		"authority": types.StringType,
	}
}

func convertUsersToJsonStruct(users []attr.Value) ([]mysqlUserToJsonConvert, error) {
	var userToConvert = []mysqlUserToJsonConvert{}

	for _, user := range users {
		userJson := mysqlUserToJsonConvert{}
		if err := json.Unmarshal([]byte(user.String()), &userJson); err != nil {
			return nil, err
		}
		userToConvert = append(userToConvert, userJson)
	}

	return userToConvert, nil
}

func flattenMysqlUsers(list []*vmysql.CloudMysqlUser) []*mysqlUser {
	var outputs []*mysqlUser

	for _, v := range list {
		var output mysqlUser
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (d *mysqlUsersDataSourceModel) refreshFromOutput(ctx context.Context, output []*mysqlUser, instance string) {
	d.ID = types.StringValue(instance)
	d.MysqlInstanceNo = types.StringValue(instance)
	userListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlUser{}.attrTypes()}, output)
	d.MysqlUserList = userListValue
}

func (d *mysqlUser) refreshFromOutput(output *vmysql.CloudMysqlUser) {
	d.UserName = types.StringPointerValue(output.UserName)
	d.HostIp = types.StringPointerValue(output.HostIp)
	d.Authority = types.StringPointerValue(output.Authority)
}

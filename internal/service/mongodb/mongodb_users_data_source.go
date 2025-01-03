package mongodb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
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
	_ datasource.DataSource              = &mongodbUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &mongodbUsersDataSource{}
)

func NewMongoDbUsersDataSource() datasource.DataSource {
	return &mongodbUsersDataSource{}
}

type mongodbUsersDataSource struct {
	config *conn.ProviderConfig
}

func (d *mongodbUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_users"
}

func (d *mongodbUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *mongodbUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"mongodb_user_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"database_name": schema.StringAttribute{
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

func (d *mongodbUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mongodbUsersDataSourceModel

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

	output, err := GetMongoDbUserAllList(ctx, d.config, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	mongodbUserList := flattenMongodbUsers(output)
	fillteredList := common.FilterModels(ctx, data.Filters, mongodbUserList)
	if diags := data.refreshFromOutput(ctx, fillteredList, data.ID.ValueString()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertUsersToJsonStruct(data.MongodbUserList.Elements()); err != nil {
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

func GetMongoDbUserAllList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*vmongodb.CloudMongoDbUser, error) {
	reqParams := &vmongodb.GetCloudMongoDbUserListRequest{
		RegionCode:             &config.RegionCode,
		CloudMongoDbInstanceNo: ncloud.String(id),
	}
	tflog.Info(ctx, "GetMongodbUserList reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vmongodb.V2Api.GetCloudMongoDbUserList(reqParams)
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "GetMongodbUserList response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudMongoDbUserList) < 1 {
		return nil, nil
	}

	return resp.CloudMongoDbUserList, nil
}

type mongodbUsersDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	MongodbUserList types.List   `tfsdk:"mongodb_user_list"`
	OutputFile      types.String `tfsdk:"output_file"`
	Filters         types.Set    `tfsdk:"filter"`
}

type mongodbUser struct {
	UserName     types.String `tfsdk:"name"`
	DatabaseName types.String `tfsdk:"database_name"`
	Authority    types.String `tfsdk:"authority"`
}

type mongodbUserToJsonConvert struct {
	UserName     string `json:"name"`
	DatabaseName string `json:"database_name"`
	Authority    string `json:"authority"`
}

func (d mongodbUser) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"database_name": types.StringType,
		"authority":     types.StringType,
	}
}

func convertUsersToJsonStruct(users []attr.Value) ([]mongodbUserToJsonConvert, error) {
	var userToConvert = []mongodbUserToJsonConvert{}

	for _, user := range users {
		userJson := mongodbUserToJsonConvert{}
		if err := json.Unmarshal([]byte(user.String()), &userJson); err != nil {
			return nil, err
		}
		userToConvert = append(userToConvert, userJson)
	}

	return userToConvert, nil
}

func flattenMongodbUsers(list []*vmongodb.CloudMongoDbUser) []*mongodbUser {
	var outputs []*mongodbUser

	for _, v := range list {
		var output mongodbUser
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (d *mongodbUsersDataSourceModel) refreshFromOutput(ctx context.Context, output []*mongodbUser, instance string) diag.Diagnostics {
	d.ID = types.StringValue(instance)
	userListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mongodbUser{}.attrTypes()}, output)
	if diags.HasError() {
		return diags
	}

	d.MongodbUserList = userListValue

	return diags
}

func (d *mongodbUser) refreshFromOutput(output *vmongodb.CloudMongoDbUser) {
	d.UserName = types.StringPointerValue(output.UserName)
	d.DatabaseName = types.StringPointerValue(output.DatabaseName)
	d.Authority = types.StringPointerValue(output.Authority)
}

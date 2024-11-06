package mysql

import (
	"context"
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &mysqlUsersResource{}
	_ resource.ResourceWithConfigure   = &mysqlUsersResource{}
	_ resource.ResourceWithImportState = &mysqlUsersResource{}
)

func NewMysqlUsersResource() resource.Resource {
	return &mysqlUsersResource{}
}

type mysqlUsersResource struct {
	config *conn.ProviderConfig
}

func (r *mysqlUsersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *mysqlUsersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

func (r *mysqlUsersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_users"
}

func (r *mysqlUsersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"mysql_instance_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mysql_user_list": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.LengthBetween(4, 16),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_,]+$`),
									"Composed of alphabets, numbers, hyphen (-), (\\), (_), (,). Must start with an alphabetic character.",
								),
							},
						},
						"host_ip": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"password": schema.StringAttribute{
							Required:  true,
							Sensitive: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.All(
									stringvalidator.LengthBetween(8, 20),
									stringvalidator.RegexMatches(regexp.MustCompile(`[a-zA-Z]+`), "Must have at least one alphabet"),
									stringvalidator.RegexMatches(regexp.MustCompile(`\d+`), "Must have at least one number"),
									stringvalidator.RegexMatches(regexp.MustCompile(`[~!@#$%^*()\-_=\[\]\{\};:,.<>?]+`), "Must have at least one special character"),
									stringvalidator.RegexMatches(regexp.MustCompile(`^[^&+\\"'/\s`+"`"+`]*$`), "Must not have ` & + \\ \" ' / and white space."),
								),
							},
						},
						"authority": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"READ", "CRUD", "DDL"}...),
							},
						},
						"is_system_table_access": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
							Default: booldefault.StaticBool(true),
						},
					},
				},
			},
		},
	}
}

func (r *mysqlUsersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mysqlUsersResourceModel

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"resource does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmysql.AddCloudMysqlUserListRequest{
		RegionCode:           &r.config.RegionCode,
		CloudMysqlInstanceNo: plan.MysqlInstanceNo.ValueStringPointer(),
		CloudMysqlUserList:   convertToCloudMysqlUserParameter(plan.MysqlUserList),
	}

	plan.ID = plan.MysqlInstanceNo

	tflog.Info(ctx, "CreateMysqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.AddCloudMysqlUserList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMysqlUserList response="+common.MarshalUncheckedString(response))

	if response == nil || *response.ReturnCode != "0" {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
	}

	_, err = waitMysqlCreation(ctx, r.config, plan.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR MYSQL CREATING ERROR", err.Error())
		return
	}

	output, err := GetMysqlUserList(ctx, r.config, plan.MysqlInstanceNo.ValueString(), common.ConvertToStringList(plan.MysqlUserList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := plan.refreshFromOutput(ctx, output, plan); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mysqlUsersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysqlUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMysqlUserList(ctx, r.config, state.ID.ValueString(), common.ConvertToStringList(state.MysqlUserList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := state.refreshFromOutput(ctx, output, state); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *mysqlUsersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state mysqlUsersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.MysqlUserList.Equal(state.MysqlUserList) {
		reqParams := &vmysql.ChangeCloudMysqlUserListRequest{
			RegionCode:           &r.config.RegionCode,
			CloudMysqlInstanceNo: state.ID.ValueStringPointer(),
			CloudMysqlUserList:   convertToCloudMysqlUserParameter(plan.MysqlUserList),
		}
		tflog.Info(ctx, "ChangecloudMysqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := r.config.Client.Vmysql.V2Api.ChangeCloudMysqlUserList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeCloudMysqlUserList response="+common.MarshalUncheckedString(response))

		if response == nil || *response.ReturnCode != "0" {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		_, err = waitMysqlCreation(ctx, r.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		output, err := GetMysqlUserList(ctx, r.config, state.ID.ValueString(), common.ConvertToStringList(plan.MysqlUserList, "name"))
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}

		if output == nil {
			resp.State.RemoveResource(ctx)
			return
		}

		if diags := state.refreshFromOutput(ctx, output, plan); diags.HasError() {
			resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *mysqlUsersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mysqlUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := waitMysqlCreation(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR MYSQL CREATION ERROR", err.Error())
		return
	}

	reqParams := &vmysql.DeleteCloudMysqlUserListRequest{
		RegionCode:           &r.config.RegionCode,
		CloudMysqlInstanceNo: state.MysqlInstanceNo.ValueStringPointer(),
		CloudMysqlUserList:   convertToCloudMysqlUserKeyParameter(state.MysqlUserList),
	}
	tflog.Info(ctx, "DeleteMysqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.DeleteCloudMysqlUserList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteMysqlUserList response="+common.MarshalUncheckedString(response))
}

func GetMysqlUserList(ctx context.Context, config *conn.ProviderConfig, id string, users []string) ([]*vmysql.CloudMysqlUser, error) {
	var allUsers []*vmysql.CloudMysqlUser
	pageNo := int32(0)
	pageSize := int32(100)
	hasMore := true

	for hasMore {
		reqParams := &vmysql.GetCloudMysqlUserListRequest{
			RegionCode:           &config.RegionCode,
			CloudMysqlInstanceNo: ncloud.String(id),
			PageNo:               ncloud.Int32(pageNo),
			PageSize:             ncloud.Int32(pageSize),
		}
		tflog.Info(ctx, "GetMysqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

		resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlUserList(reqParams)
		if err != nil {
			return nil, err
		}

		if resp == nil {
			break
		}

		allUsers = append(allUsers, resp.CloudMysqlUserList...)

		hasMore = len(resp.CloudMysqlUserList) == int(pageSize)
		pageNo++
	}

	userMap := make(map[string]*vmysql.CloudMysqlUser)
	for _, user := range allUsers {
		if user != nil && user.UserName != nil {
			userMap[*user.UserName] = user
		}
	}

	var filteredUsers []*vmysql.CloudMysqlUser
	for _, username := range users {
		if user, exists := userMap[username]; exists {
			filteredUsers = append(filteredUsers, user)
		}
	}

	if len(filteredUsers) == 0 {
		return nil, nil
	}

	tflog.Info(ctx, "GetMysqlUserList response="+common.MarshalUncheckedString(filteredUsers))

	return filteredUsers, nil
}

type mysqlUsersResourceModel struct {
	ID              types.String `tfsdk:"id"`
	MysqlInstanceNo types.String `tfsdk:"mysql_instance_no"`
	MysqlUserList   types.List   `tfsdk:"mysql_user_list"`
}

type MysqlUser struct {
	UserName            types.String `tfsdk:"name"`
	UserPassword        types.String `tfsdk:"password"`
	HostIp              types.String `tfsdk:"host_ip"`
	Authority           types.String `tfsdk:"authority"`
	IsSystemTableAccess types.Bool   `tfsdk:"is_system_table_access"`
}

func (r MysqlUser) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                   types.StringType,
		"password":               types.StringType,
		"host_ip":                types.StringType,
		"authority":              types.StringType,
		"is_system_table_access": types.BoolType,
	}
}

func (r *mysqlUsersResourceModel) refreshFromOutput(ctx context.Context, output []*vmysql.CloudMysqlUser, resourceModel mysqlUsersResourceModel) diag.Diagnostics {
	r.ID = resourceModel.ID
	r.MysqlInstanceNo = resourceModel.MysqlInstanceNo

	var userList []MysqlUser

	for idx, user := range output {
		pswd := resourceModel.MysqlUserList.Elements()[idx].(types.Object).Attributes()
		mysqlUser := MysqlUser{
			UserName:            types.StringPointerValue(user.UserName),
			HostIp:              types.StringPointerValue(user.HostIp),
			Authority:           types.StringPointerValue(user.Authority),
			IsSystemTableAccess: types.BoolPointerValue(user.IsSystemTableAccess),
			UserPassword:        pswd["password"].(types.String),
		}

		userList = append(userList, mysqlUser)
	}

	mysqlUsers, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: MysqlUser{}.AttrTypes()}, userList)
	if diags.HasError() {
		return diags
	}

	r.MysqlUserList = mysqlUsers

	return diags
}

func convertToCloudMysqlUserParameter(values basetypes.ListValue) []*vmysql.CloudMysqlUserParameter {
	result := make([]*vmysql.CloudMysqlUserParameter, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		param := &vmysql.CloudMysqlUserParameter{
			Name:      attrs["name"].(types.String).ValueStringPointer(),
			HostIp:    attrs["host_ip"].(types.String).ValueStringPointer(),
			Password:  attrs["password"].(types.String).ValueStringPointer(),
			Authority: attrs["authority"].(types.String).ValueStringPointer(),
		}

		if !attrs["is_system_table_access"].(types.Bool).IsNull() && !attrs["is_system_table_access"].(types.Bool).IsUnknown() {
			param.IsSystemTableAccess = attrs["is_system_table_access"].(types.Bool).ValueBoolPointer()
		}

		result = append(result, param)
	}

	return result
}

func convertToCloudMysqlUserKeyParameter(values basetypes.ListValue) []*vmysql.CloudMysqlUserKeyParameter {
	result := make([]*vmysql.CloudMysqlUserKeyParameter, 0, len(values.Elements())-1)

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		param := &vmysql.CloudMysqlUserKeyParameter{
			Name:   attrs["name"].(types.String).ValueStringPointer(),
			HostIp: attrs["host_ip"].(types.String).ValueStringPointer(),
		}
		result = append(result, param)
	}

	return result
}

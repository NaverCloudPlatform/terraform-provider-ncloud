package postgresql

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifystring"
)

var (
	_ resource.Resource                = &postgresqlUsersResource{}
	_ resource.ResourceWithConfigure   = &postgresqlUsersResource{}
	_ resource.ResourceWithImportState = &postgresqlUsersResource{}
)

func NewPostgresqlUsersResource() resource.Resource {
	return &postgresqlUsersResource{}
}

type postgresqlUsersResource struct {
	config *conn.ProviderConfig
}

func (r *postgresqlUsersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *postgresqlUsersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *postgresqlUsersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_users"
}

func (r *postgresqlUsersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"postgresql_instance_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"postgresql_user_list": schema.ListNestedAttribute{
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
									regexp.MustCompile(`^[a-z]+[a-z0-9_]+$`),
									"Composed of lowercase alphabets, numbers, underbar (_). Must start with an alphabetic character.",
								),
							},
						},
						"client_cidr": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
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
									verifystring.NotContain(path.MatchRoot("name").String()),
								),
							},
						},
						"is_replication_role": schema.BoolAttribute{
							Required: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *postgresqlUsersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan postgresqlUsersResourceModel

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

	reqParams := &vpostgresql.AddCloudPostgresqlUserListRequest{
		RegionCode:                &r.config.RegionCode,
		CloudPostgresqlInstanceNo: plan.PostgresqlInstanceNo.ValueStringPointer(),
		CloudPostgresqlUserList:   convertToCloudPostgresqlUserParameter(plan.PostgresqlUserList),
	}

	plan.ID = plan.PostgresqlInstanceNo

	tflog.Info(ctx, "CreatePostgresqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.AddCloudPostgresqlUserList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreatePostgresqlUserList response="+common.MarshalUncheckedString(response))

	if response == nil || *response.ReturnCode != "0" {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
	}

	_, err = WaitPostgresqlCreation(ctx, r.config, plan.PostgresqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WATING FOR POSTGRESQL CREATION ERROR", err.Error())
		return
	}

	output, err := GetPostgresqlUserList(ctx, r.config, plan.PostgresqlInstanceNo.ValueString(), convertToCloudPostgresqlUserStringList(plan.PostgresqlUserList))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
	}

	if diags := plan.refreshFromOutput(ctx, output, plan); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *postgresqlUsersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state postgresqlUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetPostgresqlUserList(ctx, r.config, state.PostgresqlInstanceNo.ValueString(), convertToCloudPostgresqlUserStringList(state.PostgresqlUserList))
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *postgresqlUsersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state postgresqlUsersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.PostgresqlUserList.Equal(state.PostgresqlUserList) {
		reqParams := &vpostgresql.ChangeCloudPostgresqlUserListRequest{
			RegionCode:                &r.config.RegionCode,
			CloudPostgresqlInstanceNo: state.ID.ValueStringPointer(),
			CloudPostgresqlUserList:   convertToCloudPostgresqlUserParameter(plan.PostgresqlUserList),
		}
		tflog.Info(ctx, "ChangeCloudPostgresqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := r.config.Client.Vpostgresql.V2Api.ChangeCloudPostgresqlUserList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeCloudPostgresqlUserList response="+common.MarshalUncheckedString(response))

		if response == nil || *response.ReturnCode != "0" {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
		}

		_, err = WaitPostgresqlCreation(ctx, r.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		output, err := GetPostgresqlUserList(ctx, r.config, state.ID.ValueString(), convertToCloudPostgresqlUserStringList(plan.PostgresqlUserList))
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
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *postgresqlUsersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postgresqlUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := WaitPostgresqlCreation(ctx, r.config, state.PostgresqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WATING FOR POSTGRESQL CREATION ERROR", err.Error())
		return
	}

	reqParams := &vpostgresql.DeleteCloudPostgresqlUserListRequest{
		RegionCode:                &r.config.RegionCode,
		CloudPostgresqlInstanceNo: state.PostgresqlInstanceNo.ValueStringPointer(),
		CloudPostgresqlUserList:   convertToCloudPostgresqlUserKeyParameter(state.PostgresqlUserList),
	}
	tflog.Info(ctx, "DeletePostgresqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.DeleteCloudPostgresqlUserList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeletePostgresqlUserList response="+common.MarshalUncheckedString(response))

	if err := waitPostgresqlUsersDeletion(ctx, r.config, state.ID.ValueString(), convertToCloudPostgresqlUserStringList(state.PostgresqlUserList)); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func waitPostgresqlUsersDeletion(ctx context.Context, config *conn.ProviderConfig, id string, users []string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			userList, err := GetPostgresqlUserList(ctx, config, id, users)
			if err != nil {
				return 0, "", err
			}

			if len(userList) > 0 {
				return userList, DELETING, nil
			}

			if userList == nil {
				return userList, DELETED, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete postgresql user")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for postgresql user (%s) to become terminating: %s", id, err)
	}

	return nil
}

func GetPostgresqlUserList(ctx context.Context, config *conn.ProviderConfig, id string, users []string) ([]*vpostgresql.CloudPostgresqlUser, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlUserListRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(id),
	}
	tflog.Info(ctx, "GetPostgresqlUserList reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vpostgresql.V2Api.GetCloudPostgresqlUserList(reqParams)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	userMap := make(map[string]*vpostgresql.CloudPostgresqlUser)
	for _, user := range resp.CloudPostgresqlUserList {
		if user != nil && user.UserName != nil {
			userMap[*user.UserName] = user
		}
	}

	var filteredUsers []*vpostgresql.CloudPostgresqlUser
	for _, username := range users {
		if user, exists := userMap[username]; exists {
			filteredUsers = append(filteredUsers, user)
		}
	}

	if len(filteredUsers) == 0 {
		return nil, nil
	}

	tflog.Info(ctx, "GetPostgresqlUserList response="+common.MarshalUncheckedString(resp))

	return filteredUsers, nil
}

type postgresqlUsersResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	PostgresqlInstanceNo types.String `tfsdk:"postgresql_instance_no"`
	PostgresqlUserList   types.List   `tfsdk:"postgresql_user_list"`
}

type PostgresqlUser struct {
	UserName          types.String `tfsdk:"name"`
	UserPassword      types.String `tfsdk:"password"`
	ClientCidr        types.String `tfsdk:"client_cidr"`
	IsReplicationRole types.Bool   `tfsdk:"is_replication_role"`
}

func (r PostgresqlUser) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                types.StringType,
		"password":            types.StringType,
		"client_cidr":         types.StringType,
		"is_replication_role": types.BoolType,
	}
}

func (r *postgresqlUsersResourceModel) refreshFromOutput(ctx context.Context, output []*vpostgresql.CloudPostgresqlUser, resourceModel postgresqlUsersResourceModel) diag.Diagnostics {
	r.ID = resourceModel.ID
	r.PostgresqlInstanceNo = resourceModel.PostgresqlInstanceNo

	var userList []PostgresqlUser

	for idx, user := range output {
		pswd := resourceModel.PostgresqlUserList.Elements()[idx].(types.Object).Attributes()
		postgresqlUser := PostgresqlUser{
			UserName:          types.StringPointerValue(user.UserName),
			UserPassword:      pswd["password"].(types.String),
			ClientCidr:        types.StringPointerValue(user.ClientCidr),
			IsReplicationRole: types.BoolPointerValue(user.IsReplicationRole),
		}

		userList = append(userList, postgresqlUser)
	}

	postgresqlUsers, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: PostgresqlUser{}.AttrTypes()}, userList)
	if diags.HasError() {
		return diags
	}

	r.PostgresqlUserList = postgresqlUsers
	return nil
}

func convertToCloudPostgresqlUserParameter(values basetypes.ListValue) []*vpostgresql.CloudPostgresqlUserParameter {
	result := make([]*vpostgresql.CloudPostgresqlUserParameter, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)

		attrs := obj.Attributes()

		param := &vpostgresql.CloudPostgresqlUserParameter{
			Name:              attrs["name"].(types.String).ValueStringPointer(),
			Password:          attrs["password"].(types.String).ValueStringPointer(),
			ClientCidr:        attrs["client_cidr"].(types.String).ValueStringPointer(),
			IsReplicationRole: attrs["is_replication_role"].(types.Bool).ValueBoolPointer(),
		}
		result = append(result, param)
	}

	return result
}

func convertToCloudPostgresqlUserKeyParameter(values basetypes.ListValue) []*vpostgresql.CloudPostgresqlUserKeyParameter {
	result := make([]*vpostgresql.CloudPostgresqlUserKeyParameter, 0, len(values.Elements())-1)

	for _, v := range values.Elements() {
		obj := v.(types.Object)

		attrs := obj.Attributes()

		param := &vpostgresql.CloudPostgresqlUserKeyParameter{
			Name: attrs["name"].(types.String).ValueStringPointer(),
		}
		result = append(result, param)
	}

	return result
}

func convertToCloudPostgresqlUserStringList(values basetypes.ListValue) []string {
	result := make([]string, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		name := attrs["name"].(types.String).ValueString()
		result = append(result, name)
	}

	return result
}
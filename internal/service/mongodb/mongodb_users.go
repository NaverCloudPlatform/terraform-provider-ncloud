package mongodb

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
)

const (
	DELETING = "deleting"
	DELETED  = "deleted"
)

var (
	_ resource.Resource                = &mongodbUsersResource{}
	_ resource.ResourceWithConfigure   = &mongodbUsersResource{}
	_ resource.ResourceWithImportState = &mongodbUsersResource{}
)

func NewMongoDbUsersResource() resource.Resource {
	return &mongodbUsersResource{}
}

type mongodbUsersResource struct {
	config *conn.ProviderConfig
}

func (r *mongodbUsersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *mongodbUsersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *mongodbUsersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_users"
}

func (r *mongodbUsersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"mongodb_instance_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mongodb_user_list": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.All(
									stringvalidator.LengthBetween(4, 16),
									stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]+$`), "Allows only alphabets, numbers and underbar (_). Must start with an alphabetic character"),
								),
							},
						},
						"password": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Sensitive: true,
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
						"database_name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.All(
									stringvalidator.LengthBetween(4, 30),
									stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]+$`), "Allows only alphabets, numbers and underbar (_). Must start with an alphabetic character"),
								),
							},
						},
						"authority": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"READ", "READ_WRITE"}...),
							},
						},
					},
				},
			},
		},
	}
}

func (r *mongodbUsersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mongodbUsersResourceModel

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

	reqParams := &vmongodb.AddCloudMongoDbUserListRequest{
		RegionCode:             &r.config.RegionCode,
		CloudMongoDbInstanceNo: plan.MongoDbInstanceNo.ValueStringPointer(),
		CloudMongoDbUserList:   convertToCloudMongodbUserParameter(plan.MongoDbUserList),
	}

	plan.ID = plan.MongoDbInstanceNo

	tflog.Info(ctx, "CreateMongodbUserList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmongodb.V2Api.AddCloudMongoDbUserList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMongodbUserList response="+common.MarshalUncheckedString(response))

	if response == nil || *response.ReturnCode != "0" {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	_, err = waitMongoDbCreated(ctx, r.config, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR MONGODB CREATING ERROR", err.Error())
		return
	}

	output, err := GetMongoDbUserList(ctx, r.config, plan.ID.ValueString(), common.ConvertToStringList(plan.MongoDbUserList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	plan.refreshFromOutput(ctx, output, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mongodbUsersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mongodbUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMongoDbUserList(ctx, r.config, state.ID.ValueString(), common.ConvertToStringList(state.MongoDbUserList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(ctx, output, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *mongodbUsersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state mongodbUsersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.MongoDbUserList.Equal(state.MongoDbUserList) {
		reqParams := &vmongodb.ChangeCloudMongoDbUserListRequest{
			RegionCode:             &r.config.RegionCode,
			CloudMongoDbInstanceNo: state.ID.ValueStringPointer(),
			CloudMongoDbUserList:   convertToCloudMongodbUserParameter(plan.MongoDbUserList),
		}
		tflog.Info(ctx, "ChangeCloudMongoDbUserList reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := r.config.Client.Vmongodb.V2Api.ChangeCloudMongoDbUserList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATE ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "ChangeCloudMongoDbUserList response="+common.MarshalUncheckedString(response))

		if response == nil || *response.ReturnCode != "0" {
			resp.Diagnostics.AddError("UPDATE ERROR", "response invalid")
			return
		}

		_, err = waitMongoDbUpdate(ctx, r.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("WAITING FOR UPDATE ERROR", err.Error())
			return
		}

		output, err := GetMongoDbUserList(ctx, r.config, state.ID.ValueString(), common.ConvertToStringList(plan.MongoDbUserList, "name"))
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}

		state.refreshFromOutput(ctx, output, plan)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *mongodbUsersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mongodbUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := waitMongoDbCreated(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WATING FOR MONGODB CREATION ERROR", err.Error())
		return
	}

	reqParams := &vmongodb.DeleteCloudMongoDbUserListRequest{
		RegionCode:             &r.config.RegionCode,
		CloudMongoDbInstanceNo: state.MongoDbInstanceNo.ValueStringPointer(),
		CloudMongoDbUserList:   convertToCloudMongodbUser(state.MongoDbUserList),
	}
	tflog.Info(ctx, "DeleteMongodbUserList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmongodb.V2Api.DeleteCloudMongoDbUserList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteMongodbUserList response="+common.MarshalUncheckedString(response))

	if err := waitMongodbUsersDeletion(ctx, r.config, state.ID.ValueString(), common.ConvertToStringList(state.MongoDbUserList, "name")); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func waitMongodbUsersDeletion(ctx context.Context, config *conn.ProviderConfig, id string, users []string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			userList, err := GetMongoDbUserList(ctx, config, id, users)
			if err != nil {
				return 0, "", err
			}

			if len(userList) == 1 || userList == nil {
				return userList, DELETED, nil
			}

			for idx, v := range userList {
				if users[idx] != *v.UserName {
					return userList, DELETED, nil
				} else {
					return userList, DELETING, nil
				}
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete mongodb user")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for mongodb user (%s) to become terminating: %s", id, err)
	}

	return nil
}

func GetMongoDbUserList(ctx context.Context, config *conn.ProviderConfig, id string, users []string) ([]*vmongodb.CloudMongoDbUser, error) {
	reqParams := &vmongodb.GetCloudMongoDbUserListRequest{
		RegionCode:             &config.RegionCode,
		CloudMongoDbInstanceNo: ncloud.String(id),
	}
	tflog.Info(ctx, "GetMongodbUserList reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vmongodb.V2Api.GetCloudMongoDbUserList(reqParams)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.CloudMongoDbUserList) < 1 {
		return nil, nil
	}

	userMap := make(map[string]*vmongodb.CloudMongoDbUser)
	for _, user := range resp.CloudMongoDbUserList {
		if user != nil && user.UserName != nil {
			userMap[*user.UserName] = user
		}
	}

	var filteredUsers []*vmongodb.CloudMongoDbUser
	for _, username := range users {
		if user, exists := userMap[username]; exists {
			filteredUsers = append(filteredUsers, user)
		}
	}

	if len(filteredUsers) == 0 {
		return nil, nil
	}

	for _, user := range resp.CloudMongoDbUserList {
		if user != nil && user.UserName != nil {
			if !containsInUsergList(*user.UserName, filteredUsers) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	tflog.Info(ctx, "GetMongodbUserList response="+common.MarshalUncheckedString(resp))

	return filteredUsers, nil
}

type mongodbUsersResourceModel struct {
	ID                types.String `tfsdk:"id"`
	MongoDbInstanceNo types.String `tfsdk:"mongodb_instance_no"`
	MongoDbUserList   types.List   `tfsdk:"mongodb_user_list"`
}

type MongodbUser struct {
	UserName     types.String `tfsdk:"name"`
	Password     types.String `tfsdk:"password"`
	DatabaseName types.String `tfsdk:"database_name"`
	Authority    types.String `tfsdk:"authority"`
}

func (r MongodbUser) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"database_name": types.StringType,
		"authority":     types.StringType,
		"password":      types.StringType,
	}
}

func (r *mongodbUsersResourceModel) refreshFromOutput(ctx context.Context, output []*vmongodb.CloudMongoDbUser, plan mongodbUsersResourceModel) {
	r.ID = plan.ID
	r.MongoDbInstanceNo = plan.ID

	var userList []MongodbUser

	for idx, user := range output[:len(plan.MongoDbUserList.Elements())] {
		pswd := plan.MongoDbUserList.Elements()[idx].(types.Object).Attributes()
		mongodbUser := MongodbUser{
			UserName:     types.StringPointerValue(user.UserName),
			DatabaseName: types.StringPointerValue(user.DatabaseName),
			Authority:    types.StringPointerValue(user.Authority),
			Password:     pswd["password"].(types.String),
		}

		userList = append(userList, mongodbUser)
	}

	mongodbUsers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: MongodbUser{}.AttrTypes()}, userList)

	r.MongoDbUserList = mongodbUsers
}

func convertToCloudMongodbUserParameter(values basetypes.ListValue) []*vmongodb.AddOrChangeCloudMongoDbUserParameter {
	result := make([]*vmongodb.AddOrChangeCloudMongoDbUserParameter, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		param := &vmongodb.AddOrChangeCloudMongoDbUserParameter{
			UserName:     attrs["name"].(types.String).ValueStringPointer(),
			Password:     attrs["password"].(types.String).ValueStringPointer(),
			DatabaseName: attrs["database_name"].(types.String).ValueStringPointer(),
			Authority:    attrs["authority"].(types.String).ValueStringPointer(),
		}
		result = append(result, param)
	}

	return result
}

func convertToCloudMongodbUser(values basetypes.ListValue) []*vmongodb.DeleteCloudMongoDbUserParameter {
	result := make([]*vmongodb.DeleteCloudMongoDbUserParameter, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		param := &vmongodb.DeleteCloudMongoDbUserParameter{
			UserName:     attrs["name"].(types.String).ValueStringPointer(),
			DatabaseName: attrs["database_name"].(types.String).ValueStringPointer(),
		}
		result = append(result, param)
	}

	return result
}

func containsInUsergList(userName string, users []*vmongodb.CloudMongoDbUser) bool {
	for _, v := range users {
		if *v.UserName == userName {
			return true
		}
	}
	return false
}

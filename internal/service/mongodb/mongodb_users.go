package mongodb

import (
	"context"
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
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
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mongodb_user_set": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.All(
									stringvalidator.LengthBetween(4, 16),
									stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]+$`), "Allows only alphabets, numbers and underbar (_). Must start with an alphabetic character"),
								),
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
						"database_name": schema.StringAttribute{
							Required: true,
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
		CloudMongoDbInstanceNo: plan.ID.ValueStringPointer(),
		CloudMongoDbUserList:   convertToAddOrChangeParameters(plan.MongoDbUserSet),
	}

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

	output, err := GetMongoDbUserAllList(ctx, r.config, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := plan.refreshFromOutput(ctx, output); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mongodbUsersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mongodbUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMongoDbUserAllList(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := state.refreshFromOutput(ctx, output); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *mongodbUsersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state mongodbUsersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !state.MongoDbUserSet.Equal(plan.MongoDbUserSet) {
		var planUserList, stateUserList []MongodbUser
		resp.Diagnostics.Append(plan.MongoDbUserSet.ElementsAs(ctx, &planUserList, false)...)
		resp.Diagnostics.Append(state.MongoDbUserSet.ElementsAs(ctx, &stateUserList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		err := addOrChangeUserList(ctx, r.config, state.ID.ValueStringPointer(), planUserList, stateUserList)
		if err != nil {
			resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
			return
		}

		err = deleteUserList(ctx, r.config, state.ID.ValueStringPointer(), planUserList, stateUserList)
		if err != nil {
			resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
			return
		}

		output, err := GetMongoDbUserAllList(ctx, r.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}

		if diags := plan.refreshFromOutput(ctx, output); diags.HasError() {
			resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mongodbUsersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mongodbUsersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := waitMongoDbCreated(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete. Please try again later", err.Error())
		return
	}

	reqParams := &vmongodb.DeleteCloudMongoDbUserListRequest{
		RegionCode:             &r.config.RegionCode,
		CloudMongoDbInstanceNo: state.ID.ValueStringPointer(),
		CloudMongoDbUserList:   convertToDeleteParameters(state.MongoDbUserSet),
	}
	tflog.Info(ctx, "DeleteMongodbUserList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmongodb.V2Api.DeleteCloudMongoDbUserList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteMongodbUserList response="+common.MarshalUncheckedString(response))

	_, err = waitMongoDbCreated(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETION ERROR", err.Error())
		return
	}
}

type mongodbUsersResourceModel struct {
	ID             types.String `tfsdk:"id"`
	MongoDbUserSet types.Set    `tfsdk:"mongodb_user_set"`
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

func (r *mongodbUsersResourceModel) refreshFromOutput(ctx context.Context, output []*vmongodb.CloudMongoDbUser) diag.Diagnostics {
	var refreshUserList, resourceUserList []MongodbUser

	diags := r.MongoDbUserSet.ElementsAs(ctx, &resourceUserList, false)
	if diags.HasError() {
		return diags
	}

	for _, o := range output {
		password := types.StringNull()

		for _, rv := range resourceUserList {
			if rv.UserName.Equal(types.StringPointerValue(o.UserName)) && rv.DatabaseName.Equal(types.StringPointerValue(o.DatabaseName)) {
				password = rv.Password
			}
		}

		mongodbUser := MongodbUser{
			UserName:     types.StringPointerValue(o.UserName),
			DatabaseName: types.StringPointerValue(o.DatabaseName),
			Authority:    types.StringPointerValue(o.Authority),
			Password:     password,
		}

		refreshUserList = append(refreshUserList, mongodbUser)
	}

	mongodbUsers, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: MongodbUser{}.AttrTypes()}, refreshUserList)
	if diags.HasError() {
		return diags
	}

	r.MongoDbUserSet = mongodbUsers

	return diags
}

func convertToAddOrChangeParameters(values basetypes.SetValue) []*vmongodb.AddOrChangeCloudMongoDbUserParameter {
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

func convertToDeleteParameters(values basetypes.SetValue) []*vmongodb.DeleteCloudMongoDbUserParameter {
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

func convertToSingleAddOrChangeParameter(v MongodbUser) *vmongodb.AddOrChangeCloudMongoDbUserParameter {
	param := &vmongodb.AddOrChangeCloudMongoDbUserParameter{
		UserName:     v.UserName.ValueStringPointer(),
		Password:     v.Password.ValueStringPointer(),
		DatabaseName: v.DatabaseName.ValueStringPointer(),
		Authority:    v.Authority.ValueStringPointer(),
	}

	return param
}

func convertToSingleDeleteParameter(v MongodbUser) *vmongodb.DeleteCloudMongoDbUserParameter {
	param := &vmongodb.DeleteCloudMongoDbUserParameter{
		UserName:     v.UserName.ValueStringPointer(),
		DatabaseName: v.DatabaseName.ValueStringPointer(),
	}

	return param
}

func addOrChangeUserList(ctx context.Context, config *conn.ProviderConfig, id *string, planUserList []MongodbUser, stateUserList []MongodbUser) error {
	changeParameters := make([]*vmongodb.AddOrChangeCloudMongoDbUserParameter, 0)
	addParameters := make([]*vmongodb.AddOrChangeCloudMongoDbUserParameter, 0)

	for _, pv := range planUserList {
		found := false
		for _, sv := range stateUserList {
			if pv.UserName.Equal(sv.UserName) && pv.DatabaseName.Equal(sv.DatabaseName) {
				if !pv.Password.Equal(sv.Password) || !pv.Authority.Equal(sv.Authority) {
					changeParameters = append(changeParameters, convertToSingleAddOrChangeParameter(pv))
				}
				found = true
				break
			}
		}

		if !found {
			addParameters = append(addParameters, convertToSingleAddOrChangeParameter(pv))
		}
	}

	if len(changeParameters) > 0 {
		reqParams := &vmongodb.ChangeCloudMongoDbUserListRequest{
			RegionCode:             &config.RegionCode,
			CloudMongoDbInstanceNo: id,
			CloudMongoDbUserList:   changeParameters,
		}

		response, err := config.Client.Vmongodb.V2Api.ChangeCloudMongoDbUserList(reqParams)
		if err != nil {
			return err
		}
		tflog.Info(ctx, "ChangeCloudMongoDbUserList response="+common.MarshalUncheckedString(response))

		if response == nil || *response.ReturnCode != "0" {
			return fmt.Errorf("ChangeCloudMongoDbUserList response invalid")
		}

		_, err = waitMongoDbUpdate(ctx, config, *id)
		if err != nil {
			return err
		}
	}

	if len(addParameters) > 0 {
		reqParams := &vmongodb.AddCloudMongoDbUserListRequest{
			RegionCode:             &config.RegionCode,
			CloudMongoDbInstanceNo: id,
			CloudMongoDbUserList:   addParameters,
		}

		response, err := config.Client.Vmongodb.V2Api.AddCloudMongoDbUserList(reqParams)
		if err != nil {
			return err
		}
		tflog.Info(ctx, "AddCloudMongoDbUserList response="+common.MarshalUncheckedString(response))

		if response == nil || *response.ReturnCode != "0" {
			return fmt.Errorf("AddCloudMongoDbUserList response invalid")
		}

		_, err = waitMongoDbUpdate(ctx, config, *id)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteUserList(ctx context.Context, config *conn.ProviderConfig, id *string, planUserList []MongodbUser, stateUserList []MongodbUser) error {
	deleteParameters := make([]*vmongodb.DeleteCloudMongoDbUserParameter, 0)

	for _, sv := range stateUserList {
		found := false
		for _, pv := range planUserList {
			if sv.UserName.Equal(pv.UserName) && sv.DatabaseName.Equal(pv.DatabaseName) {
				found = true
				break
			}
		}

		if !found {
			deleteParameters = append(deleteParameters, convertToSingleDeleteParameter(sv))
		}
	}

	if len(deleteParameters) > 0 {
		reqParams := &vmongodb.DeleteCloudMongoDbUserListRequest{
			RegionCode:             &config.RegionCode,
			CloudMongoDbInstanceNo: id,
			CloudMongoDbUserList:   deleteParameters,
		}
		tflog.Info(ctx, "DeleteMongodbUserList reqParams="+common.MarshalUncheckedString(reqParams))

		response, err := config.Client.Vmongodb.V2Api.DeleteCloudMongoDbUserList(reqParams)
		if err != nil {
			return err
		}
		tflog.Info(ctx, "DeleteMongodbUserList response="+common.MarshalUncheckedString(response))

		_, err = waitMongoDbUpdate(ctx, config, *id)
		if err != nil {
			return err
		}
	}
	return nil
}

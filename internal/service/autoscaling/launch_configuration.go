package autoscaling

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &launchConfigurationResource{}
	_ resource.ResourceWithConfigure   = &launchConfigurationResource{}
	_ resource.ResourceWithImportState = &launchConfigurationResource{}
)

type launchConfigurationResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	LaunchConfigurationNo       types.String `tfsdk:"launch_configuration_no"`
	LaunchConfigurationName     types.String `tfsdk:"name"`
	ServerImageProductCode      types.String `tfsdk:"server_image_product_code"`
	ServerProductCode           types.String `tfsdk:"server_product_code"`
	MemberServerImageInstanceNo types.String `tfsdk:"member_server_image_no"`
	LoginKeyName                types.String `tfsdk:"login_key_name"`
	InitScriptNo                types.String `tfsdk:"init_script_no"`
	UserData                    types.String `tfsdk:"user_data"`
	AccessControlGroupNoList    types.List   `tfsdk:"access_control_group_no_list"`
	IsEncryptedVolume           types.Bool   `tfsdk:"is_encrypted_volume"`
}

func NewLaunchConfigResource() resource.Resource {
	return &launchConfigurationResource{}
}

type launchConfigurationResource struct {
	config *conn.ProviderConfig
}

func (l *launchConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (l *launchConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_launch_configuration"
}

func (l *launchConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"launch_configuration_no": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(255),
				},
			},
			"server_image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("member_server_image_no"),
					),
				},
			},
			"server_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"member_server_image_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("server_image_product_code"),
					),
				},
			},
			"login_key_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"init_script_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_data": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access_control_group_no_list": schema.ListAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				ElementType: types.StringType,
				Description: "This parameter cannot be duplicated in classic type.",
			},
			"is_encrypted_volume": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"id": framework.IDAttribute(),
		},
	}
}

func (l *launchConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	l.config = config
}

func (l *launchConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan launchConfigurationResourceModel
	var err error
	var launchConfigNo *string

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if l.config.SupportVPC {
		// [createVPCLaunchConfiguration] -> refactor to function
		launchConfigNo, err = createVpcLaunchConfiguration(ctx, l.config, &plan)

	} else {
		// [createClassicLaunchConfiguration]
		launchConfigNo, err = createClassicLaunchConfiguration(ctx, l.config, &plan)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating LaunchCOnfiguration",
			err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(ncloud.StringValue(launchConfigNo))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func createVpcLaunchConfiguration(ctx context.Context, config *conn.ProviderConfig, plan *launchConfigurationResourceModel) (*string, error) {
	reqParams := &vautoscaling.CreateLaunchConfigurationRequest{
		RegionCode:                  &config.RegionCode,
		ServerImageProductCode:      plan.ServerImageProductCode.ValueStringPointer(),
		MemberServerImageInstanceNo: plan.MemberServerImageInstanceNo.ValueStringPointer(),
		ServerProductCode:           plan.ServerProductCode.ValueStringPointer(),
		IsEncryptedVolume:           plan.IsEncryptedVolume.ValueBoolPointer(),
		InitScriptNo:                plan.InitScriptNo.ValueStringPointer(),
		LaunchConfigurationName:     plan.LaunchConfigurationName.ValueStringPointer(),
		LoginKeyName:                plan.LoginKeyName.ValueStringPointer(),
	}
	tflog.Info(ctx, "CreateVpcLaunchConfiguration", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vautoscaling.V2Api.CreateLaunchConfiguration(reqParams)
	tflog.Info(ctx, "CreateVpcLaunchConfiguration response", map[string]any{
		"createVpcLaunchConfiguration": common.MarshalUncheckedString(resp),
	})

	return resp.LaunchConfigurationList[0].LaunchConfigurationNo, err
}

func createClassicLaunchConfiguration(ctx context.Context, config *conn.ProviderConfig, plan *launchConfigurationResourceModel) (*string, error) {
	reqParams := &autoscaling.CreateLaunchConfigurationRequest{
		LaunchConfigurationName: plan.LaunchConfigurationName.ValueStringPointer(),
		ServerImageProductCode:  plan.ServerImageProductCode.ValueStringPointer(),
		ServerProductCode:       plan.ServerProductCode.ValueStringPointer(),
		MemberServerImageNo:     plan.MemberServerImageInstanceNo.ValueStringPointer(),
		LoginKeyName:            plan.LoginKeyName.ValueStringPointer(),
		UserData:                plan.UserData.ValueStringPointer(),
		RegionNo:                &config.RegionNo,
	}

	if !plan.AccessControlGroupNoList.IsNull() {
		reqParams.AccessControlGroupConfigurationNoList = expandStringList(plan.AccessControlGroupNoList)
	}

	tflog.Info(ctx, "CreateClassicLaunchConfiguration", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Autoscaling.V2Api.CreateLaunchConfiguration(reqParams)
	tflog.Info(ctx, "CreateClassicLaunchConfiguration response", map[string]any{
		"createClassicLaunchConfiguration": common.MarshalUncheckedString(resp),
	})

	return resp.LaunchConfigurationList[0].LaunchConfigurationNo, err
}

func expandStringList(attr types.List) []*string {
	if attr.IsNull() || attr.IsUnknown() {
		return nil
	}

	var result []*string
	for _, v := range attr.Elements() {
		str := v.(types.String).ValueStringPointer()
		result = append(result, str)
	}

	return result
}

func (l *launchConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state launchConfigurationResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	launchConfig, err := GetLaunchConfiguration(ctx, l.config, state.LaunchConfigurationNo.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error Reading LaunchConfiguration", err.Error())
		return
	}

	if launchConfig == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(ctx, launchConfig)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func GetLaunchConfiguration(ctx context.Context, config *conn.ProviderConfig, id string) (*LaunchConfiguration, error) {
	if config.SupportVPC {
		return GetVpcLaunchConfiguration(ctx, config, id)

	} else {
		return GetClassicLaunchConfiguration(ctx, config, id)
	}
}

func GetVpcLaunchConfiguration(ctx context.Context, config *conn.ProviderConfig, id string) (*LaunchConfiguration, error) {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LaunchConfigurationNoList = []*string{ncloud.String(id)}
	}
	tflog.Info(ctx, "getVpcLaunchConfiguration reqParams", map[string]any{
		"getVpcLaunchConfiguration": common.MarshalUncheckedString(reqParams),
	})
	resp, err := config.Client.Vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)

	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "getVpcLaunchConfiguration response", map[string]any{
		"getVpcLaunchConfiguration": common.MarshalUncheckedString(resp),
	})

	if len(resp.LaunchConfigurationList) < 1 {
		return nil, nil
	}

	l := resp.LaunchConfigurationList[0]

	return &LaunchConfiguration{
		LaunchConfigurationName:     l.LaunchConfigurationName,
		ServerImageProductCode:      l.ServerImageProductCode,
		MemberServerImageInstanceNo: l.MemberServerImageInstanceNo,
		ServerProductCode:           l.ServerProductCode,
		LoginKeyName:                l.LoginKeyName,
		InitScriptNo:                l.InitScriptNo,
		IsEncryptedVolume:           l.IsEncryptedVolume,
		LaunchConfigurationNo:       l.LaunchConfigurationNo,
		AccessControlGroupNoList:    []*string{},
	}, nil
}

func GetClassicLaunchConfiguration(ctx context.Context, config *conn.ProviderConfig, id string) (*LaunchConfiguration, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		RegionNo: &config.RegionNo,
	}

	tflog.Info(ctx, "getClassicLaunchConfiguration reqParams", map[string]any{
		"getClassicLaunchConfiguration reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "getClassicLaunchConfiguration response", map[string]any{
		"getClassicLaunchConfiguration response": common.MarshalUncheckedString(resp),
	})

	for _, l := range resp.LaunchConfigurationList {
		if *l.LaunchConfigurationNo == *no {
			return &LaunchConfiguration{
				LaunchConfigurationNo:       l.LaunchConfigurationNo,
				LaunchConfigurationName:     l.LaunchConfigurationName,
				ServerImageProductCode:      l.ServerImageProductCode,
				MemberServerImageInstanceNo: l.MemberServerImageNo,
				ServerProductCode:           l.ServerProductCode,
				LoginKeyName:                l.LoginKeyName,
				UserData:                    l.UserData,
				AccessControlGroupNoList:    flattenAccessControlGroupList(l.AccessControlGroupList),
			}, nil
		}
	}

	return nil, nil
}

func (l *launchConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (l *launchConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state launchConfigurationResourceModel
	var err error

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if l.config.SupportVPC {
		err = deleteVpcLaunchConfiguration(ctx, l.config, state.LaunchConfigurationNo.ValueString())

	} else {
		err = deleteClassicLaunchConfiguration(ctx, l.config, state.LaunchConfigurationNo.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("Error Deleting LaunchConfiguration", err.Error())
		return
	}
}

func deleteVpcLaunchConfiguration(ctx context.Context, config *conn.ProviderConfig, id string) error {
	reqParams := &vautoscaling.DeleteLaunchConfigurationRequest{
		LaunchConfigurationNo: ncloud.String(id),
	}
	tflog.Info(ctx, "deleteVpcLaunchConfiguration reqParams", map[string]any{
		"deleteVpcLaunchConfiguration reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vautoscaling.V2Api.DeleteLaunchConfiguration(reqParams)
	if err != nil {
		return err
	}
	tflog.Info(ctx, "deleteVpcLaunchConfiguration response", map[string]any{
		"deleteVpcLaunchConfiguration response": common.MarshalUncheckedString(resp),
	})

	return nil
}

func deleteClassicLaunchConfiguration(ctx context.Context, config *conn.ProviderConfig, id string) error {
	launchConfig, err := GetClassicLaunchConfiguration(ctx, config, id)
	if err != nil {
		return err
	}

	if launchConfig == nil {
		return nil
	}

	reqParams := &autoscaling.DeleteAutoScalingLaunchConfigurationRequest{
		LaunchConfigurationName: launchConfig.LaunchConfigurationName,
	}
	tflog.Info(ctx, "deleteClassicLaunchConfiguration reqParams", map[string]any{
		"deleteClassicLaunchConfiguration reqParams": common.MarshalUncheckedString(reqParams),
	})
	resp, err := config.Client.Autoscaling.V2Api.DeleteAutoScalingLaunchConfiguration(reqParams)
	if err != nil {
		return err
	}
	tflog.Info(ctx, "deleteClassicLaunchConfiguration response", map[string]any{
		"deleteClassicLaunchConfiguration response": common.MarshalUncheckedString(resp),
	})
	return nil
}

func (l *launchConfigurationResourceModel) refreshFromOutput(ctx context.Context, output *LaunchConfiguration) {
	l.ID = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationName = types.StringPointerValue(output.LaunchConfigurationName)
	l.ServerImageProductCode = types.StringPointerValue(output.ServerImageProductCode)
	l.ServerProductCode = types.StringPointerValue(output.ServerProductCode)
	l.MemberServerImageInstanceNo = types.StringPointerValue(output.MemberServerImageInstanceNo)
	l.LoginKeyName = types.StringPointerValue(output.LoginKeyName)
	l.InitScriptNo = types.StringPointerValue(output.InitScriptNo)
	l.UserData = types.StringPointerValue(output.UserData)
	accessControlGroupNoList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	l.AccessControlGroupNoList = accessControlGroupNoList
	l.IsEncryptedVolume = types.BoolPointerValue(output.IsEncryptedVolume)
}

type LaunchConfiguration struct {
	LaunchConfigurationNo       *string   `json:"launch_configuration_no,omitempty"`
	LaunchConfigurationName     *string   `json:"name,omitempty"`
	ServerImageProductCode      *string   `json:"server_image_product_code,omitempty"`
	MemberServerImageInstanceNo *string   `json:"member_server_image_no,omitempty"`
	ServerProductCode           *string   `json:"server_product_code,omitempty"`
	LoginKeyName                *string   `json:"login_key_name,omitempty"`
	InitScriptNo                *string   `json:"init_script_no,omitempty"`
	IsEncryptedVolume           *bool     `json:"is_encrypted_volume,omitempty"`
	UserData                    *string   `json:"user_data,omitempty"`
	AccessControlGroupNoList    []*string `json:"access_control_group_no_list"`
}

func GetClassicLaunchConfigurationByNo(no *string, config *conn.ProviderConfig) (*LaunchConfiguration, error) {
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		RegionNo: &config.RegionNo,
	}
	resp, err := config.Client.Autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		return nil, err
	}

	for _, l := range resp.LaunchConfigurationList {
		if *l.LaunchConfigurationNo == *no {
			return &LaunchConfiguration{
				LaunchConfigurationNo:       l.LaunchConfigurationNo,
				LaunchConfigurationName:     l.LaunchConfigurationName,
				ServerImageProductCode:      l.ServerImageProductCode,
				MemberServerImageInstanceNo: l.MemberServerImageNo,
				ServerProductCode:           l.ServerProductCode,
				LoginKeyName:                l.LoginKeyName,
				UserData:                    l.UserData,
			}, nil
		}
	}
	return nil, fmt.Errorf("Not found LaunchConfiguration(%s)", ncloud.StringValue(no))
}

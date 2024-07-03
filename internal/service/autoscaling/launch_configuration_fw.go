package autoscaling

import (
	"context"
	"fmt"
	"time"
	// "strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	// sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	// "github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &launchConfigurationResource{}
	_ resource.ResourceWithConfigure   = &launchConfigurationResource{}
	_ resource.ResourceWithImportState = &launchConfigurationResource{}
)

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
	var id *string

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if l.config.SupportVPC {
		id, err = createVpcLaunchConfiguration(ctx, l.config, &plan)
	} else {
		id, err = createClassicLaunchConfiguration(ctx, l.config, &plan)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating LaunchConfiguration",
			err.Error(),
		)
		return
	}

	output, err := waitForNcloudLaunchConfiguration(l.config, *id)
	if err != nil {
		resp.Diagnostics.AddError(
			"wating for LaunchConfiguration creation",
			err.Error(),
		)
		return
	}

	plan.refreshFromOutput(output)
	plan.id = types.StringValue(ncloud.StringValue(id)) // [TODO] check it

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

}

func (l *launchConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state launchConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !l.config.SupportVPC {
		// [TODO] GetLaunchConfiguration
	} else {
		// [TODO] getVpcLaunchConfiguration
	}
}

// [TODO] Framework should implement Udpate function even if resource does not support Update function
func (l *launchConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (l *launchConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state launchConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !l.config.SupportVPC {
		// [TODO] deleteClassicLaunchConfiguration
	} else {
		// [TODO] deleteVpcLaunchConfiguration
	}
}

type launchConfigurationResourceModel struct {
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

func (l *launchConfigurationResourceModel) refreshFromOutput(output *LaunchConfiguration){
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationName = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)

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

	tflog.Info(ctx, "DeleteVpcLaunchConfiguration", map[string]any{
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

	// [TODO]
	// if param, ok := d.GetOk("access_control_group_no_list"); ok {
	// 	reqParams.AccessControlGroupConfigurationNoList = ExpandStringInterfaceList(param.([]interface{}))
	// } -> ?
	// if !data.AccessControlGroupNoList.IsNull() && !data.AccessControlGroupNoList(){
	// 	reqParams.AccessControlGroupConfigurationNoList = ExpandStringInterfaceList(param.([]interaface{}))
	// }

	tflog.Info(ctx, "DeleteClassicLaunchConfiguration", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Autoscaling.V2Api.CreateLaunchConfiguration(reqParams)
	tflog.Info(ctx, "CreateClassicLaunchConfiguration response", map[string]any{
		"createClassicLaunchConfiguration": common.MarshalUncheckedString(resp),
	})

	return resp.LaunchConfigurationList[0].LaunchConfigurationNo, err
}

func waitForNcloudLaunchConfiguration(config *conn.ProviderConfig, id string) (*vpc.LaunchConfiguration, error) {
	var launchConfiguration *LaunchConfiguration
	stateConf := &retry.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetLaunchConfiguration(ctx, config, id)
			launchConfiguration = resp
			if err != nil {
				return 0, "", err
			}
			return resp, "", nil
			// [TODO] what should return ? 
			// return VpcCommonStateRefreshFunc(launchConfiguration, err, "LaunchConfigurationStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("Error wating for LaunchConfiguration (%s) to become available: %s", id, err)
	}
	return launchConfiguration, nil
}


// func GetLaunchConfiguration(config *conn.ProviderConfig, )
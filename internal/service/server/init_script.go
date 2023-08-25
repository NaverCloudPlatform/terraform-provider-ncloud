package server

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ resource.Resource                = &initScriptResource{}
	_ resource.ResourceWithConfigure   = &initScriptResource{}
	_ resource.ResourceWithImportState = &initScriptResource{}
)

func NewInitScriptResource() resource.Resource {
	return &initScriptResource{}
}

type initScriptResource struct {
	config *conn.ProviderConfig
}

func (i *initScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (i *initScriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_init_script"
}

func (i *initScriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: verify.InstanceNameValidator(),
			},
			"id": framework.IDAttribute(),
			"content": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1000),
				},
			},
			"os_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"LNX", "WND"}...),
				},
			},
			"init_script_no": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (i *initScriptResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	i.config = config
}

func (i *initScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan initScriptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !i.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.Config.Schema.Type().String()),
		)
		return
	}

	reqParams := &vserver.CreateInitScriptRequest{
		RegionCode:        &i.config.RegionCode,
		InitScriptContent: plan.Content.ValueStringPointer(),
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		reqParams.InitScriptName = plan.Name.ValueStringPointer()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqParams.InitScriptDescription = plan.Description.ValueStringPointer()
	}
	if !plan.OsType.IsNull() && !plan.OsType.IsUnknown() {
		reqParams.OsTypeCode = plan.OsType.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateVpcInitScript", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := i.config.Client.Vserver.V2Api.CreateInitScript(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Vpc Init Script, err params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "CreateVpcInitScript response", map[string]any{
		"createVpcInitScriptResponse": common.MarshalUncheckedString(response),
	})

	initScriptInstance := response.InitScriptList[0]
	plan.ID = types.StringPointerValue(initScriptInstance.InitScriptNo)
	tflog.Info(ctx, "InitScript ID", map[string]any{"initScriptNo": *initScriptInstance.InitScriptNo})

	if err := plan.refreshFromOutput(initScriptInstance); err != nil {
		resp.Diagnostics.AddError("refreshing init script details", err.Error())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (i *initScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state initScriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !i.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.ProviderMeta.Schema.Type().String()),
		)
		return
	}
	output, err := GetInitScript(i.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetInitScript", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := state.refreshFromOutput(output); err != nil {
		resp.Diagnostics.AddError("refreshing init script details", err.Error())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (i *initScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (i *initScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state initScriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !i.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.ProviderMeta.Schema.Type().String()),
		)
		return
	}

	if err := DeleteInitScript(i.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail for init script deletion",
			err.Error())
	}

}

func GetInitScript(config *conn.ProviderConfig, id string) (*vserver.InitScript, error) {
	reqParams := &vserver.GetInitScriptDetailRequest{
		RegionCode:   &config.RegionCode,
		InitScriptNo: ncloud.String(id),
	}

	common.LogCommonRequest("GetInitScriptDetail", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetInitScriptDetail(reqParams)
	if err != nil {
		common.LogErrorResponse("GetInitScriptDetail", err, reqParams)
		return nil, err
	}
	common.LogResponse("GetInitScriptDetail", resp)

	if len(resp.InitScriptList) > 0 {
		return resp.InitScriptList[0], nil
	}

	return nil, nil
}

func DeleteInitScript(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.DeleteInitScriptsRequest{
		RegionCode:       &config.RegionCode,
		InitScriptNoList: []*string{ncloud.String(id)},
	}

	common.LogCommonRequest("deleteVpcInitScript", reqParams)
	resp, err := config.Client.Vserver.V2Api.DeleteInitScripts(reqParams)
	if err != nil {
		common.LogErrorResponse("deleteVpcInitScript", err, reqParams)
		return err
	}
	common.LogResponse("deleteVpcInitScript", resp)

	return nil
}

type initScriptResourceModel struct {
	InitScriptNo types.String `tfsdk:"init_script_no"`
	OsType       types.String `tfsdk:"os_type"`
	ID           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	Name         types.String `tfsdk:"name"`
	Content      types.String `tfsdk:"content"`
}

func (m *initScriptResourceModel) refreshFromOutput(output *vserver.InitScript) error {
	m.ID = types.StringPointerValue(output.InitScriptNo)
	m.Name = types.StringPointerValue(output.InitScriptName)
	m.Description = framework.EmptyStringToNull(types.StringPointerValue(output.InitScriptDescription))
	m.OsType = types.StringPointerValue(output.InitScriptContent)
	m.Content = types.StringPointerValue(output.InitScriptContent)
	m.InitScriptNo = types.StringPointerValue(output.InitScriptNo)

	return nil
}

package subaccount

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	subaccountsdk "github.com/terraform-providers/terraform-provider-ncloud/internal/sdk/subaccount"
)

var (
	_ resource.Resource                = &subAccountResource{}
	_ resource.ResourceWithConfigure   = &subAccountResource{}
	_ resource.ResourceWithImportState = &subAccountResource{}
)

func NewSubAccountResource() resource.Resource {
	return &subAccountResource{}
}

type subAccountResource struct {
	config *conn.ProviderConfig
}

type subAccountResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	LoginId             types.String `tfsdk:"login_id"`
	Name                types.String `tfsdk:"name"`
	Email               types.String `tfsdk:"email"`
	Memo                types.String `tfsdk:"memo"`
	CanConsoleAccess    types.Bool   `tfsdk:"can_console_access"`
	CanAPIGatewayAccess types.Bool   `tfsdk:"can_api_gateway_access"`
	IsMfaMandatory      types.Bool   `tfsdk:"is_mfa_mandatory"`
	ConsolePermitIps    types.List   `tfsdk:"console_permit_ips"`
	ApiAllowSources     types.List   `tfsdk:"api_allow_sources"`
	GeneratedPassword   types.String `tfsdk:"generated_password"`
	SubAccountNo        types.Int64  `tfsdk:"sub_account_no"`
	Nrn                 types.String `tfsdk:"nrn"`
	Active              types.Bool   `tfsdk:"active"`
	CreateTime          types.String `tfsdk:"create_time"`
}

type apiAllowSourceModel struct {
	Type   types.String `tfsdk:"type"`
	Source types.String `tfsdk:"source"`
}

var apiAllowSourceAttrTypes = map[string]attr.Type{
	"type":   types.StringType,
	"source": types.StringType,
}

func (r *subAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount"
}

func (r *subAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"login_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Login ID of the sub account. Changing this creates a new sub account.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the sub account.",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: "Email address of the sub account.",
			},
			"memo": schema.StringAttribute{
				Optional:    true,
				Description: "Memo for the sub account.",
			},
			"can_console_access": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Whether the sub account can access the console. When enabled at creation, " +
					"an initial password is generated and exposed once via `generated_password`.",
			},
			"can_api_gateway_access": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the sub account can access APIs through API Gateway. Required to issue access keys.",
			},
			"is_mfa_mandatory": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether two-factor authentication is mandatory for console login.",
			},
			"console_permit_ips": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				Description: "List of IP addresses allowed to access the console. Omit to allow all.",
			},
			"generated_password": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				Description: "Initial console password generated when `can_console_access` is enabled at creation. " +
					"Available only at creation time.",
			},
			"sub_account_no": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of the sub account.",
			},
			"nrn": schema.StringAttribute{
				Computed:    true,
				Description: "NCP Resource Name of the sub account.",
			},
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the sub account is active.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sub account.",
			},
			"api_allow_sources": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("IP", "VPC", "VPC_SERVER"),
							},
							Description: "Type of the allowed API source (IP, VPC or VPC_SERVER).",
						},
						"source": schema.StringAttribute{
							Required:    true,
							Description: "Source value: an IP address (CIDR allowed) for IP, or an instance number for VPC/VPC_SERVER.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				Description: "Sources allowed to call APIs with this sub account's access keys. Omit to allow all.",
			},
		},
	}
}

func (r *subAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *subAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *subAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	consolePermitIps, apiAllowSources, diags := expandAccessSources(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &subaccountsdk.CreateSubAccountRequest{
		LoginId:             plan.LoginId.ValueString(),
		Name:                plan.Name.ValueString(),
		Email:               plan.Email.ValueString(),
		Memo:                plan.Memo.ValueString(),
		CanConsoleAccess:    plan.CanConsoleAccess.ValueBool(),
		CanAPIGatewayAccess: plan.CanAPIGatewayAccess.ValueBool(),
		IsMfaMandatory:      plan.IsMfaMandatory.ValueBool(),
		UseConsolePermitIp:  len(consolePermitIps) > 0,
		ConsolePermitIps:    consolePermitIps,
		UseApiAllowSource:   len(apiAllowSources) > 0,
		ApiAllowSources:     apiAllowSources,
	}

	// Password flows are intentionally kept out of Terraform configuration.
	// For console-enabled accounts the API generates the initial password,
	// which is surfaced once through the generated_password attribute.
	if createReq.CanConsoleAccess {
		createReq.NeedPasswordGenerate = true
		createReq.NeedPasswordReset = true
	}

	createResp, err := r.config.Client.SubAccount.CreateSubAccount(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating SubAccount", err.Error())
		return
	}

	plan.GeneratedPassword = types.StringNull()
	if createResp.GeneratedPassword != "" {
		plan.GeneratedPassword = types.StringValue(createResp.GeneratedPassword)
	}

	detail, err := waitForSubAccount(ctx, r.config, createResp.Id)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating SubAccount", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.refreshFromOutput(ctx, detail)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *subAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	detail, err := r.config.Client.SubAccount.GetSubAccount(ctx, state.ID.ValueString())
	if err != nil {
		if subaccountsdk.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading SubAccount", err.Error())
		return
	}

	resp.Diagnostics.Append(state.refreshFromOutput(ctx, detail)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *subAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state subAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	consolePermitIps, apiAllowSources, diags := expandAccessSources(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &subaccountsdk.UpdateSubAccountRequest{
		Name:                plan.Name.ValueString(),
		Email:               plan.Email.ValueString(),
		Memo:                plan.Memo.ValueString(),
		CanConsoleAccess:    plan.CanConsoleAccess.ValueBool(),
		CanAPIGatewayAccess: plan.CanAPIGatewayAccess.ValueBool(),
		UseConsolePermitIp:  len(consolePermitIps) > 0,
		ConsolePermitIps:    consolePermitIps,
		UseApiAllowSource:   len(apiAllowSources) > 0,
		ApiAllowSources:     apiAllowSources,
	}
	if !plan.IsMfaMandatory.IsNull() {
		updateReq.IsMfaMandatory = plan.IsMfaMandatory.ValueBoolPointer()
	}
	if updateReq.ConsolePermitIps == nil {
		updateReq.ConsolePermitIps = []string{}
	}
	if updateReq.ApiAllowSources == nil {
		updateReq.ApiAllowSources = []subaccountsdk.ApiAllowSource{}
	}

	id := state.ID.ValueString()
	if err := r.config.Client.SubAccount.UpdateSubAccount(ctx, id, updateReq); err != nil {
		resp.Diagnostics.AddError("Error Updating SubAccount", err.Error())
		return
	}

	detail, err := r.config.Client.SubAccount.GetSubAccount(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating SubAccount", err.Error())
		return
	}

	plan.GeneratedPassword = state.GeneratedPassword

	resp.Diagnostics.Append(plan.refreshFromOutput(ctx, detail)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *subAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.config.Client.SubAccount.DeleteSubAccount(ctx, state.ID.ValueString()); err != nil && !subaccountsdk.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting SubAccount", err.Error())
	}
}

func expandAccessSources(ctx context.Context, plan subAccountResourceModel) ([]string, []subaccountsdk.ApiAllowSource, diag.Diagnostics) {
	var diags diag.Diagnostics

	var consolePermitIps []string
	if !plan.ConsolePermitIps.IsNull() && !plan.ConsolePermitIps.IsUnknown() {
		diags.Append(plan.ConsolePermitIps.ElementsAs(ctx, &consolePermitIps, false)...)
	}

	var apiAllowSources []subaccountsdk.ApiAllowSource
	if !plan.ApiAllowSources.IsNull() && !plan.ApiAllowSources.IsUnknown() {
		var sourceModels []apiAllowSourceModel
		diags.Append(plan.ApiAllowSources.ElementsAs(ctx, &sourceModels, false)...)
		for _, s := range sourceModels {
			apiAllowSources = append(apiAllowSources, subaccountsdk.ApiAllowSource{
				Type:   s.Type.ValueString(),
				Source: s.Source.ValueString(),
			})
		}
	}

	return consolePermitIps, apiAllowSources, diags
}

func waitForSubAccount(ctx context.Context, config *conn.ProviderConfig, id string) (*subaccountsdk.SubAccountDetail, error) {
	var detail *subaccountsdk.SubAccountDetail

	stateConf := &retry.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (any, string, error) {
			resp, err := config.Client.SubAccount.GetSubAccount(ctx, id)
			if err != nil {
				if subaccountsdk.IsNotFound(err) {
					return 0, "", nil
				}
				return 0, "", err
			}

			detail = resp
			return resp, "OK", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return nil, fmt.Errorf("error waiting for SubAccount (%s) to become available: %s", id, err)
	}

	return detail, nil
}

// refreshFromOutput does not touch is_mfa_mandatory and generated_password:
// neither is returned by the Get Sub Account API, so their values can only
// come from configuration and the creation response respectively.
func (m *subAccountResourceModel) refreshFromOutput(ctx context.Context, detail *subaccountsdk.SubAccountDetail) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(detail.SubAccountId)
	m.LoginId = types.StringValue(detail.LoginId)
	m.Name = types.StringValue(detail.Name)
	m.CanConsoleAccess = types.BoolValue(detail.CanConsoleAccess)
	m.CanAPIGatewayAccess = types.BoolValue(detail.CanAPIGatewayAccess)
	m.SubAccountNo = types.Int64Value(detail.SubAccountNo)
	m.Nrn = types.StringValue(detail.Nrn)
	m.Active = types.BoolValue(detail.Active)
	m.CreateTime = types.StringValue(detail.CreateTime)

	m.Email = types.StringNull()
	if detail.Email != "" {
		m.Email = types.StringValue(detail.Email)
	}

	m.Memo = types.StringNull()
	if detail.Memo != "" {
		m.Memo = types.StringValue(detail.Memo)
	}

	m.ConsolePermitIps = types.ListNull(types.StringType)
	if detail.UseConsolePermitIp && len(detail.ConsolePermitIps) > 0 {
		ips, d := types.ListValueFrom(ctx, types.StringType, detail.ConsolePermitIps)
		diags.Append(d...)
		m.ConsolePermitIps = ips
	}

	m.ApiAllowSources = types.ListNull(types.ObjectType{AttrTypes: apiAllowSourceAttrTypes})
	if detail.UseApiAllowSource && len(detail.ApiAllowSources) > 0 {
		sources := make([]attr.Value, 0, len(detail.ApiAllowSources))
		for _, s := range detail.ApiAllowSources {
			obj, d := types.ObjectValue(apiAllowSourceAttrTypes, map[string]attr.Value{
				"type":   types.StringValue(s.Type),
				"source": types.StringValue(s.Source),
			})
			diags.Append(d...)
			sources = append(sources, obj)
		}
		list, d := types.ListValue(types.ObjectType{AttrTypes: apiAllowSourceAttrTypes}, sources)
		diags.Append(d...)
		m.ApiAllowSources = list
	}

	return diags
}

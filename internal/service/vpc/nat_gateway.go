package vpc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ resource.Resource                = &natGatewayResource{}
	_ resource.ResourceWithConfigure   = &natGatewayResource{}
	_ resource.ResourceWithImportState = &natGatewayResource{}
)

func NewNatGatewayResource() resource.Resource {
	return &natGatewayResource{}
}

type natGatewayResource struct {
	config *conn.ProviderConfig
}

func (n *natGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (n *natGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nat_gateway"
}

func (n *natGatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1000),
				},
			},
			"vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_ip": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_ip_no": schema.StringAttribute{
				Computed: true,
			},
			"nat_gateway_no": schema.StringAttribute{
				Computed: true,
			},
			"public_ip": schema.StringAttribute{
				Computed: true,
			},
			"subnet_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
func (n *natGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	n.config = config
}

func (n *natGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan natGatewayResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !n.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.Config.Schema.Type().String()),
		)
		return
	}

	reqParams := &vpc.CreateNatGatewayInstanceRequest{
		RegionCode: &n.config.RegionCode,
		VpcNo:      plan.VpcNo.ValueStringPointer(),
		ZoneCode:   plan.Zone.ValueStringPointer(),
		SubnetNo:   plan.SubnetNo.ValueStringPointer(),
	}

	if !plan.Name.IsNull() {
		reqParams.NatGatewayName = plan.Name.ValueStringPointer()
	}

	if !plan.Description.IsNull() {
		reqParams.NatGatewayDescription = plan.Description.ValueStringPointer()
	}

	if !plan.PrivateIp.IsNull() {
		reqParams.PrivateIp = plan.PrivateIp.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateNatGateway", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	response, err := n.config.Client.Vpc.V2Api.CreateNatGatewayInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create NatGateway Instance, err params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "CreateNatGateway response", map[string]any{
		"createNatGatewayResponse": common.MarshalUncheckedString(response),
	})

	natGatewayInstance := response.NatGatewayInstanceList[0]
	plan.ID = types.StringPointerValue(natGatewayInstance.NatGatewayInstanceNo)
	tflog.Info(ctx, "NAT GATEWAY ID", map[string]any{"natGatewayNo": *natGatewayInstance.NatGatewayInstanceNo})

	output, err := waitForNcloudNatGatewayCreation(n.config, *natGatewayInstance.NatGatewayInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("waiting for Nat Gateway creation", err.Error())
		return
	}

	if err := plan.refreshFromOutput(output); err != nil {
		resp.Diagnostics.AddError("refreshing nat gateway details", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (n *natGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state natGatewayResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetNatGatewayInstance(n.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetNatGateway", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := state.refreshFromOutput(output); err != nil {
		resp.Diagnostics.AddError("refreshing nat gateway details", err.Error())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *natGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state natGatewayResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Description.Equal(state.Description) {

		reqParams := &vpc.SetNatGatewayDescriptionRequest{
			RegionCode:            &n.config.RegionCode,
			NatGatewayInstanceNo:  state.NatGatewayNo.ValueStringPointer(),
			NatGatewayDescription: plan.Description.ValueStringPointer(),
		}

		tflog.Info(ctx, "SetNatGateway", map[string]any{
			"reqParams": common.MarshalUncheckedString(reqParams),
		})

		response, err := n.config.Client.Vpc.V2Api.SetNatGatewayDescription(reqParams)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("SetNatGateway params=%v", *reqParams),
				err.Error(),
			)
			return
		}

		tflog.Info(ctx, "SetNatGateway", map[string]any{
			"updateNatGatewayResponse": common.MarshalUncheckedString(response),
		})

		output, err := GetNatGatewayInstance(n.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("GetNatGateway", err.Error())
			return
		}

		if err := state.refreshFromOutput(output); err != nil {
			resp.Diagnostics.AddError("refreshing nat gateway details", err.Error())
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (n *natGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state natGatewayResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.DeleteNatGatewayInstanceRequest{
		RegionCode:           &n.config.RegionCode,
		NatGatewayInstanceNo: state.NatGatewayNo.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteNatGateway", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := n.config.Client.Vpc.V2Api.DeleteNatGatewayInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("DeleteNatGateway NatGateway Instance params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "DeleteNatGateway response", map[string]any{
		"deleteNatGatewayResponse": common.MarshalUncheckedString(response),
	})

	if err := WaitForNcloudNatGatewayDeletion(n.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail to wait for nat gateway deletion",
			err.Error(),
		)
	}
}

func waitForNcloudNatGatewayCreation(config *conn.ProviderConfig, id string) (*vpc.NatGatewayInstance, error) {
	var natGatewayInstance *vpc.NatGatewayInstance
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNatGatewayInstance(config, id)
			natGatewayInstance = instance
			return VpcCommonStateRefreshFunc(instance, err, "NatGatewayInstanceStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("Error waiting for NAT GATEWAY (%s) to become available: %s", id, err)
	}

	return natGatewayInstance, nil
}

func WaitForNcloudNatGatewayDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNatGatewayInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NatGatewayInstanceStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetNatGatewayInstance(config *conn.ProviderConfig, id string) (*vpc.NatGatewayInstance, error) {
	reqParams := &vpc.GetNatGatewayInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		NatGatewayInstanceNo: ncloud.String(id),
	}

	common.LogCommonRequest("GetNatGatewayInstanceDetail", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetNatGatewayInstanceDetail(reqParams)
	if err != nil {
		common.LogErrorResponse("GetNatGatewayInstanceDetail", err, reqParams)
		return nil, err
	}
	common.LogResponse("GetNatGatewayInstanceDetail", resp)

	if len(resp.NatGatewayInstanceList) > 0 {
		instance := resp.NatGatewayInstanceList[0]
		return instance, nil
	}

	return nil, nil
}

type natGatewayResourceModel struct {
	Description  types.String `tfsdk:"description"`
	VpcNo        types.String `tfsdk:"vpc_no"`
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Zone         types.String `tfsdk:"zone"`
	SubnetNo     types.String `tfsdk:"subnet_no"`
	PrivateIp    types.String `tfsdk:"private_ip"`
	PublicIpNo   types.String `tfsdk:"public_ip_no"`
	NatGatewayNo types.String `tfsdk:"nat_gateway_no"`
	PublicIp     types.String `tfsdk:"public_ip"`
	SubnetName   types.String `tfsdk:"subnet_name"`
}

func (m *natGatewayResourceModel) refreshFromOutput(output *vpc.NatGatewayInstance) error {
	m.ID = types.StringPointerValue(output.NatGatewayInstanceNo)
	m.NatGatewayNo = types.StringPointerValue(output.NatGatewayInstanceNo)
	m.Name = types.StringPointerValue(output.NatGatewayName)
	m.Description = framework.EmptyStringToNull(types.StringPointerValue(output.NatGatewayDescription))
	m.VpcNo = types.StringPointerValue(output.VpcNo)
	m.Zone = types.StringPointerValue(output.ZoneCode)
	m.SubnetNo = types.StringPointerValue(output.SubnetNo)
	m.PrivateIp = types.StringPointerValue(output.PrivateIp)
	m.PublicIpNo = types.StringPointerValue(output.PublicIpInstanceNo)
	m.PublicIp = types.StringPointerValue(output.PublicIp)
	m.SubnetName = types.StringPointerValue(output.SubnetName)

	return nil
}

package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ resource.Resource                = &vpcPeeringResource{}
	_ resource.ResourceWithConfigure   = &vpcPeeringResource{}
	_ resource.ResourceWithImportState = &vpcPeeringResource{}
)

func NewVpcPeeringResource() resource.Resource {
	return &vpcPeeringResource{}
}

type vpcPeeringResource struct {
	config *conn.ProviderConfig
}

func (v *vpcPeeringResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

func (v *vpcPeeringResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_peering"
}

func (v *vpcPeeringResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(0),
					stringvalidator.LengthAtMost(1000),
				},
			},
			"source_vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_vpc_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_vpc_login_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_peering_no": schema.StringAttribute{
				Computed: true,
			},
			"has_reverse_vpc_peering": schema.BoolAttribute{
				Computed: true,
			},
			"is_between_accounts": schema.BoolAttribute{
				Computed: true,
			},
			"id": framework.IDAttribute(),
		},
	}
}

func (v *vpcPeeringResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	v.config = config
}

func (v *vpcPeeringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcPeeringResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !v.config.SupportVPC {
		resp.Diagnostics.AddError(
			"not support classic",
			fmt.Sprintf("resource %s does not support classic", req.Config.Schema.Type().String()),
		)
		return
	}

	reqParams := &vpc.CreateVpcPeeringInstanceRequest{
		RegionCode:  &v.config.RegionCode,
		SourceVpcNo: plan.SourceVpcNo.ValueStringPointer(),
		TargetVpcNo: plan.TargetVpcNo.ValueStringPointer(),
	}

	if !plan.Name.IsNull() {
		reqParams.VpcPeeringName = plan.Name.ValueStringPointer()
	}

	if !plan.Description.IsNull() {
		reqParams.VpcPeeringDescription = plan.Description.ValueStringPointer()
	}

	if !plan.TargetVpcName.IsNull() {
		reqParams.TargetVpcName = plan.TargetVpcName.ValueStringPointer()
	}

	if !plan.TargetVpcLoginId.IsNull() {
		reqParams.TargetVpcLoginId = plan.TargetVpcLoginId.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateVpcPeering", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := v.config.Client.Vpc.V2Api.CreateVpcPeeringInstance(reqParams)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("create vpc peering instance, err params=%v", *reqParams),
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "CreateVpcPeering response", map[string]any{
		"createVpcPeeringResponse": common.MarshalUncheckedString(resp),
	})

	instance := response.VpcPeeringInstanceList[0]
	plan.ID = types.StringPointerValue(instance.VpcPeeringInstanceNo)
	tflog.Info(ctx, "VPC Peering ID: %s", map[string]any{"vpcPeeringNo": *instance.VpcPeeringInstanceNo})

	output, err := waitForNcloudVpcPeeringCreation(ctx, v.config, *instance.VpcPeeringInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("waiting for Vpc peering creation", err.Error())
		return
	}

	plan.refreshFromOutput(output)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (v *vpcPeeringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcPeeringResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetVpcPeeringInstance(ctx, v.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetVpcPeering", err.Error())
		return
	}
	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(output)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (v *vpcPeeringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state vpcPeeringResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Description.Equal(state.Description) {
		reqParams := &vpc.SetVpcPeeringDescriptionRequest{
			RegionCode:            &v.config.RegionCode,
			VpcPeeringInstanceNo:  state.VpcPeeringNo.ValueStringPointer(),
			VpcPeeringDescription: plan.Description.ValueStringPointer(),
		}

		tflog.Info(ctx, "setVpcPeering", map[string]any{
			"reqParams": common.MarshalUncheckedString(reqParams),
		})

		response, err := v.config.Client.Vpc.V2Api.SetVpcPeeringDescription(reqParams)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("SetVpcPeeringDescription  params=%v", *reqParams),
				err.Error(),
			)
			return
		}

		tflog.Info(ctx, "SetVpcPeeringDescription", map[string]any{
			"updateVpcPeeringResponse": common.MarshalUncheckedString(response),
		})

		output, err := GetVpcPeeringInstance(ctx, v.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("GetVpcPeering", err.Error())
			return
		}

		state.refreshFromOutput(output)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (v *vpcPeeringResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcPeeringResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.DeleteVpcPeeringInstanceRequest{
		RegionCode:           &v.config.RegionCode,
		VpcPeeringInstanceNo: state.VpcPeeringNo.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteVpcPeering", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := v.config.Client.Vpc.V2Api.DeleteVpcPeeringInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("DeleteVpcPeering Instance params=%v", *reqParams),
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "DeleteVpcPeering response", map[string]any{
		"deleteVpcPeeringResponse": common.MarshalUncheckedString(response),
	})

	if err := WaitForNcloudVpcPeeringDeletion(ctx, v.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail to wait for vpc peering deletion",
			err.Error(),
		)
	}

}

func (m *vpcPeeringResourceModel) refreshFromOutput(output *vpc.VpcPeeringInstance) {
	m.ID = types.StringPointerValue(output.VpcPeeringInstanceNo)
	m.Name = types.StringPointerValue(output.VpcPeeringName)
	m.VpcPeeringNo = types.StringPointerValue(output.VpcPeeringInstanceNo)
	m.TargetVpcName = types.StringPointerValue(output.TargetVpcName)
	m.Description = types.StringPointerValue(output.VpcPeeringDescription)
	m.SourceVpcNo = types.StringPointerValue(output.SourceVpcNo)
	m.TargetVpcNo = types.StringPointerValue(output.TargetVpcNo)
	m.TargetVpcLoginId = types.StringPointerValue(output.TargetVpcLoginId)
	m.HasReverseVpcPeering = types.BoolPointerValue(output.HasReverseVpcPeering)
	m.IsBetweenAccounts = types.BoolPointerValue(output.IsBetweenAccounts)
}

func waitForNcloudVpcPeeringCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vpc.VpcPeeringInstance, error) {
	var vpcPeeringInstance *vpc.VpcPeeringInstance
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetVpcPeeringInstance(ctx, config, id)
			vpcPeeringInstance = instance
			return VpcCommonStateRefreshFunc(instance, err, "VpcPeeringInstanceStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("Error waiting for VPC Peering (%s) to become available: %s", id, err)
	}

	return vpcPeeringInstance, nil
}

func WaitForNcloudVpcPeeringDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {

	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetVpcPeeringInstance(ctx, config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcPeeringInstanceStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Peering (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetVpcPeeringInstance(ctx context.Context, config *conn.ProviderConfig, id string) (*vpc.VpcPeeringInstance, error) {
	reqParams := &vpc.GetVpcPeeringInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		VpcPeeringInstanceNo: ncloud.String(id),
	}

	tflog.Info(ctx, "GetVpcPeeringInstanceDetail", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vpc.V2Api.GetVpcPeeringInstanceDetail(reqParams)
	if err != nil {
		tflog.Error(ctx, "GetVpcPeeringInstanceDetail", map[string]any{
			"reqParams": common.MarshalUncheckedString(reqParams),
		})
		return nil, err
	}

	tflog.Info(ctx, "GetVpcPeeringInstanceDetail", map[string]any{
		"respParams": common.MarshalUncheckedString(resp),
	})

	if len(resp.VpcPeeringInstanceList) > 0 {
		instance := resp.VpcPeeringInstanceList[0]
		return instance, nil
	}

	return nil, nil
}

type vpcPeeringResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	SourceVpcNo          types.String `tfsdk:"source_vpc_no"`
	TargetVpcNo          types.String `tfsdk:"target_vpc_no"`
	TargetVpcName        types.String `tfsdk:"target_vpc_name"`
	TargetVpcLoginId     types.String `tfsdk:"target_vpc_login_id"`
	VpcPeeringNo         types.String `tfsdk:"vpc_peering_no"`
	HasReverseVpcPeering types.Bool   `tfsdk:"has_reverse_vpc_peering"`
	IsBetweenAccounts    types.Bool   `tfsdk:"is_between_accounts"`
}

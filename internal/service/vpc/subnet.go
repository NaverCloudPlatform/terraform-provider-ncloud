package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

const (
	SubnetPleaseTryAgainErrorCode = "3000"
)

var (
	_ resource.Resource                = &subnetResource{}
	_ resource.ResourceWithConfigure   = &subnetResource{}
	_ resource.ResourceWithImportState = &subnetResource{}
)

func NewSubnetResource() resource.Resource {
	return &subnetResource{}
}

type subnetResource struct {
	config *conn.ProviderConfig
}

func (s *subnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (s *subnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (s *subnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators:  verify.InstanceNameValidator(),
				Description: "Subnet name to create. default: Assigned by NAVER CLOUD PLATFORM",
			},
			"id": framework.IDAttribute(),
			"vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The id of the VPC that the desired subnet belongs to.",
			},
			"subnet": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: verify.CidrBlockValidator(),
			},
			"zone": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_acl_no": schema.StringAttribute{
				Required: true,
			},
			"subnet_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"PUBLIC", "PRIVATE"}...),
				},
			},
			"usage_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"GEN", "LOADB", "BM", "NATGW"}...),
				},
			},
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (s *subnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	s.config = config
}

func (s *subnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subnetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.CreateSubnetRequest{
		RegionCode:     &s.config.RegionCode,
		Subnet:         plan.Subnet.ValueStringPointer(),
		SubnetTypeCode: plan.SubnetType.ValueStringPointer(),
		UsageTypeCode:  plan.UsageType.ValueStringPointer(),
		NetworkAclNo:   plan.NetworkAclNo.ValueStringPointer(),
		VpcNo:          plan.VpcNo.ValueStringPointer(),
		ZoneCode:       plan.Zone.ValueStringPointer(),
	}

	if !plan.Name.IsNull() {
		reqParams.SubnetName = plan.Name.ValueStringPointer()
	}

	if !plan.UsageType.IsNull() {
		reqParams.UsageTypeCode = plan.UsageType.ValueStringPointer()
	}

	timeout := time.Minute * 20
	var response *vpc.CreateSubnetResponse

	err := sdkresource.RetryContext(ctx, timeout, func() *sdkresource.RetryError {
		var err error
		tflog.Info(ctx, "CreateSubnet", map[string]any{
			"reqParams": common.MarshalUncheckedString(reqParams),
		})
		response, err = s.config.Client.Vpc.V2Api.CreateSubnet(reqParams)

		if err != nil {
			errBody, _ := common.GetCommonErrorBody(err)
			if errBody.ReturnCode == "1001015" || errBody.ReturnCode == SubnetPleaseTryAgainErrorCode {
				common.LogErrorResponse("retry CreateSubnet", err, reqParams)
				time.Sleep(time.Second * 5)
				return sdkresource.RetryableError(err)
			}
			return sdkresource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError("fail to create subnet", err.Error())
		return
	}

	subnetInstance := response.SubnetList[0]
	plan.ID = types.StringPointerValue(subnetInstance.SubnetNo)

	output, err := waitForNcloudSubnetCreation(s.config, *subnetInstance.SubnetNo)
	if err != nil {
		resp.Diagnostics.AddError("waiting for Subnet creation", err.Error())
		return
	}

	if err := plan.refreshFromOutput(output); err != nil {
		resp.Diagnostics.AddError("refreshing subnet details", err.Error())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (s *subnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subnetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetSubnetInstance(s.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetSubnet", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := state.refreshFromOutput(output); err != nil {
		resp.Diagnostics.AddError("refreshing subnet details", err.Error())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *subnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state subnetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.NetworkAclNo.Equal(state.NetworkAclNo) {
		reqParams := &vpc.SetSubnetNetworkAclRequest{
			RegionCode:   &s.config.RegionCode,
			NetworkAclNo: plan.NetworkAclNo.ValueStringPointer(),
			SubnetNo:     state.SubnetNo.ValueStringPointer(),
		}

		tflog.Info(ctx, "SetSubnetNetworkAcl", map[string]any{
			"reqParams": common.MarshalUncheckedString(reqParams),
		})
		response, err := s.config.Client.Vpc.V2Api.SetSubnetNetworkAcl(reqParams)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("SetSubnetNetworkAcl params=%v", *reqParams),
				err.Error(),
			)
			return
		}

		tflog.Info(ctx, "SetSubnetNetworkAcl", map[string]any{
			"updateSubnetResponse": common.MarshalUncheckedString(response),
		})

		if err := waitForNcloudNetworkACLUpdate(s.config, plan.NetworkAclNo.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"fail to wait for subnet update",
				err.Error(),
			)
		}

		output, err := GetSubnetInstance(s.config, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("GetSubnet", err.Error())
			return
		}

		if err := state.refreshFromOutput(output); err != nil {
			resp.Diagnostics.AddError("refreshing subnet details", err.Error())
		}

	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (s *subnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subnetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.DeleteSubnetRequest{
		RegionCode: &s.config.RegionCode,
		SubnetNo:   state.SubnetNo.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteSubnet", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := s.config.Client.Vpc.V2Api.DeleteSubnet(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("DeleteSubnet Subnet Instance params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "DeleteSubnet response", map[string]any{
		"deleteSubnetResponse": common.MarshalUncheckedString(response),
	})

	if err := WaitForNcloudSubnetDeletion(s.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail to wait for subnet deletion",
			err.Error(),
		)
	}
}

func waitForNcloudSubnetCreation(config *conn.ProviderConfig, id string) (*vpc.Subnet, error) {
	var subnetInstance *vpc.Subnet
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetSubnetInstance(config, id)
			subnetInstance = instance
			return VpcCommonStateRefreshFunc(instance, err, "SubnetStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("Error waiting for Subnet (%s) to become available: %s", id, err)
	}

	return subnetInstance, nil
}

func waitForNcloudNetworkACLUpdate(config *conn.ProviderConfig, id string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNetworkACLInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Set network ACL for Subnet (%s) to become running: %s", id, err)
	}

	return nil
}

func WaitForNcloudSubnetDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetSubnetInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "SubnetStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Subnet (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetSubnetInstance(config *conn.ProviderConfig, id string) (*vpc.Subnet, error) {
	reqParams := &vpc.GetSubnetDetailRequest{
		RegionCode: &config.RegionCode,
		SubnetNo:   ncloud.String(id),
	}

	common.LogCommonRequest("GetSubnetDetail", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetSubnetDetail(reqParams)
	if err != nil {
		common.LogErrorResponse("GetSubnetDetail", err, reqParams)
		return nil, err
	}
	common.LogResponse("GetSubnetDetail", resp)

	if len(resp.SubnetList) > 0 {
		instance := resp.SubnetList[0]
		return instance, nil
	}

	return nil, nil
}

type subnetResourceModel struct {
	NetworkAclNo types.String `tfsdk:"network_acl_no"`
	VpcNo        types.String `tfsdk:"vpc_no"`
	ID           types.String `tfsdk:"id"`
	Subnet       types.String `tfsdk:"subnet"`
	Zone         types.String `tfsdk:"zone"`
	SubnetType   types.String `tfsdk:"subnet_type"`
	UsageType    types.String `tfsdk:"usage_type"`
	Name         types.String `tfsdk:"name"`
	SubnetNo     types.String `tfsdk:"subnet_no"`
}

func (m *subnetResourceModel) refreshFromOutput(output *vpc.Subnet) error {
	m.ID = types.StringPointerValue(output.SubnetNo)
	m.SubnetNo = types.StringPointerValue(output.SubnetNo)
	m.VpcNo = types.StringPointerValue(output.VpcNo)
	m.Zone = types.StringPointerValue(output.ZoneCode)
	m.Name = types.StringPointerValue(output.SubnetName)
	m.Subnet = types.StringPointerValue(output.Subnet)
	m.SubnetType = types.StringPointerValue(output.SubnetType.Code)
	m.UsageType = types.StringPointerValue(output.UsageType.Code)
	m.NetworkAclNo = types.StringPointerValue(output.NetworkAclNo)

	return nil
}

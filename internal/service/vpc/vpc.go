package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ resource.Resource              = &vpcResource{}
	_ resource.ResourceWithConfigure = &vpcResource{}
)

func NewVpcResource() resource.Resource {
	return &vpcResource{}
}

type vpcResource struct {
	config *conn.ProviderConfig
}

func (v *vpcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc"
}

func (v *vpcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"ipv4_cidr_block": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				// FIXME: Validators:  validation.IsCIDRNetwork(16, 28),
				Description: "The CIDR block for the vpc",
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"default_network_acl_no": schema.StringAttribute{
				Computed: true,
			},
			"default_access_control_group_no": schema.StringAttribute{
				Computed: true,
			},
			"default_public_route_table_no": schema.StringAttribute{
				Computed: true,
			},
			"default_private_route_table_no": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *vpcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.config = config
}

func (r *vpcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.Config.Schema.Type().String()),
		)
		return
	}

	reqParams := &vpc.CreateVpcRequest{
		RegionCode:    &r.config.RegionCode,
		Ipv4CidrBlock: plan.Ipv4CidrBlock.ValueStringPointer(),
	}

	if !plan.Name.IsNull() {
		reqParams.VpcName = plan.Name.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateVpc", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := r.config.Client.Vpc.V2Api.CreateVpc(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Vpc Instance, err params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "CreateVpc response", map[string]any{
		"createVpcResponse": common.MarshalUncheckedString(response),
	})

	vpcInstance := response.VpcList[0]
	plan.ID = types.StringPointerValue(vpcInstance.VpcNo)
	tflog.Info(ctx, "VPC ID", map[string]any{"vpcNo": *vpcInstance.VpcNo})

	output, err := waitForNcloudVpcCreation(r.config, *vpcInstance.VpcNo)
	if err != nil {
		resp.Diagnostics.AddError("waiting for VPC creation", err.Error())
		return
	}

	if err := plan.refreshFromOutput(output, r.config); err != nil {
		resp.Diagnostics.AddError("refreshing vpc details", err.Error())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *vpcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetVpcInstance(r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetVPC", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := state.refreshFromOutput(output, r.config); err != nil {
		resp.Diagnostics.AddError("refreshing vpc details", err.Error())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *vpcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.DeleteVpcRequest{
		RegionCode: &r.config.RegionCode,
		VpcNo:      state.VpcNo.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteVpc", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := r.config.Client.Vpc.V2Api.DeleteVpc(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("DeleteVpc Vpc Instance params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "DeleteVpc response", map[string]any{
		"deleteVpcResponse": common.MarshalUncheckedString(response),
	})

	if err := WaitForNcloudVpcDeletion(r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail to wait for vpc deletion",
			err.Error(),
		)
	}
}

func getDefaultNetworkACL(config *conn.ProviderConfig, id string) (string, error) {
	reqParams := &vpc.GetNetworkAclListRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(id),
	}

	common.LogCommonRequest("GetNetworkAclList", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetNetworkAclList(reqParams)

	if err != nil {
		common.LogErrorResponse("GetNetworkAclList", err, reqParams)
		return "", err
	}

	common.LogResponse("GetNetworkAclList", resp)

	if resp == nil || len(resp.NetworkAclList) == 0 {
		return "", fmt.Errorf("no matching Network ACL found")
	}

	for _, i := range resp.NetworkAclList {
		if *i.IsDefault {
			return *i.NetworkAclNo, nil
		}
	}

	return "", fmt.Errorf("No matching default network ACL found")
}

func GetDefaultAccessControlGroup(config *conn.ProviderConfig, id string) (string, error) {
	reqParams := &vserver.GetAccessControlGroupListRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(id),
	}

	common.LogCommonRequest("getDefaultAccessControlGroup", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetAccessControlGroupList(reqParams)

	if err != nil {
		common.LogErrorResponse("getDefaultAccessControlGroup", err, reqParams)
		return "", err
	}

	common.LogResponse("getDefaultAccessControlGroup", resp)

	if resp == nil || len(resp.AccessControlGroupList) == 0 {
		return "", fmt.Errorf("no matching Access Control Group found")
	}

	for _, i := range resp.AccessControlGroupList {
		if *i.IsDefault {
			return *i.AccessControlGroupNo, nil
		}
	}

	return "", fmt.Errorf("No matching default Access Control Group found")
}

func getDefaultRouteTable(config *conn.ProviderConfig, id string) (publicRouteTableNo string, privateRouteTableNo string, error error) {
	reqParams := &vpc.GetRouteTableListRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(id),
	}

	common.LogCommonRequest("getDefaultRouteTable", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetRouteTableList(reqParams)

	if err != nil {
		common.LogErrorResponse("getDefaultRouteTable", err, reqParams)
		return "", "", err
	}

	common.LogResponse("getDefaultRouteTable", resp)

	for _, i := range resp.RouteTableList {
		if *i.IsDefault && *i.SupportedSubnetType.Code == "PRIVATE" {
			privateRouteTableNo = *i.RouteTableNo
		} else if *i.IsDefault && *i.SupportedSubnetType.Code == "PUBLIC" {
			publicRouteTableNo = *i.RouteTableNo
		}
	}

	return publicRouteTableNo, privateRouteTableNo, nil
}

func waitForNcloudVpcCreation(config *conn.ProviderConfig, id string) (*vpc.Vpc, error) {
	var vpcInstance *vpc.Vpc
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetVpcInstance(config, id)
			vpcInstance = instance
			return VpcCommonStateRefreshFunc(instance, err, "VpcStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("Error waiting for VPC (%s) to become available: %s", id, err)
	}

	return vpcInstance, nil
}

func WaitForNcloudVpcDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetVpcInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetVpcInstance(config *conn.ProviderConfig, id string) (*vpc.Vpc, error) {
	reqParams := &vpc.GetVpcDetailRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(id),
	}

	resp, err := config.Client.Vpc.V2Api.GetVpcDetail(reqParams)
	if err != nil {
		common.LogErrorResponse("Get Vpc Instance", err, reqParams)
		return nil, err
	}
	common.LogResponse("GetVpcDetail", resp)

	if len(resp.VpcList) > 0 {
		vpc := resp.VpcList[0]
		return vpc, nil
	}

	return nil, nil
}

type vpcResourceModel struct {
	DefaultAccessControlGroupNo types.String `tfsdk:"default_access_control_group_no"`
	DefaultNetworkAclNo         types.String `tfsdk:"default_network_acl_no"`
	DefaultPrivateRouteTableNo  types.String `tfsdk:"default_private_route_table_no"`
	DefaultPublicRouteTableNo   types.String `tfsdk:"default_public_route_table_no"`
	ID                          types.String `tfsdk:"id"`
	Ipv4CidrBlock               types.String `tfsdk:"ipv4_cidr_block"`
	Name                        types.String `tfsdk:"name"`
	VpcNo                       types.String `tfsdk:"vpc_no"`
}

func (m *vpcResourceModel) refreshFromOutput(output *vpc.Vpc, config *conn.ProviderConfig) error {
	m.ID = types.StringPointerValue(output.VpcNo)
	m.VpcNo = types.StringPointerValue(output.VpcNo)
	m.Name = types.StringPointerValue(output.VpcName)
	m.Ipv4CidrBlock = types.StringPointerValue(output.Ipv4CidrBlock)

	if *output.VpcStatus.Code != "TERMTING" {
		defaultNetworkACLNo, err := getDefaultNetworkACL(config, m.ID.ValueString())
		if err != nil {
			return fmt.Errorf("error get default network acl for VPC (%s): %s", m.ID.ValueString(), err)
		}

		m.DefaultNetworkAclNo = types.StringValue(defaultNetworkACLNo)

		defaultAcgNo, err := GetDefaultAccessControlGroup(config, m.ID.ValueString())
		if err != nil {
			return fmt.Errorf("error get default Access Control Group for VPC (%s): %s", m.ID.ValueString(), err)
		}
		m.DefaultAccessControlGroupNo = types.StringValue(defaultAcgNo)

		publicRouteTableNo, privateRouteTableNo, err := getDefaultRouteTable(config, m.ID.ValueString())
		if err != nil {
			return fmt.Errorf("error get default Route Table for VPC (%s): %s", m.ID.ValueString(), err)
		}
		m.DefaultPublicRouteTableNo = types.StringValue(publicRouteTableNo)
		m.DefaultPrivateRouteTableNo = types.StringValue(privateRouteTableNo)
	}

	return nil
}

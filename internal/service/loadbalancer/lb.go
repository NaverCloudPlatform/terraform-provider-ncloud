package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	vpcservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	verify "github.com/terraform-providers/terraform-provider-ncloud/internal/verify/int32"
)

var (
	_ resource.Resource = &lbResource{}
)

const (
	LoadBalancerInstanceOperationChangeCode             = "CHANG"
	LoadBalancerInstanceOperationCreateCode             = "CREAT"
	LoadBalancerInstanceOperationDisUseCode             = "DISUS"
	LoadBalancerInstanceOperationNullCode               = "NULL"
	LoadBalancerInstanceOperationPendingTerminationCode = "PTERM"
	LoadBalancerInstanceOperationTerminateCode          = "TERMT"
	LoadBalancerInstanceOperationUseCode                = "USE"
)

func NewLbResource() resource.Resource {
	return &lbResource{}
}

type lbResource struct {
	config *conn.ProviderConfig
}

func (l *lbResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lb"
}

func (l *lbResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"load_balancer_no": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": framework.IDAttribute(),
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"domain": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("PUBLIC", "PRIVATE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"idle_timeout": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 3600),
					verify.ConflictsWithVaule("type", "NETWORK"),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("APPLICATION", "NETWORK", "NETWORK_PROXY"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"throughput_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("SMALL", "MEDIUM", "LARGE", "DYNAMIC"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"ip_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (l *lbResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	l.config = req.ProviderData.(*conn.ProviderConfig)
}

func (l *lbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan lbResourceModel

	if !l.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"The `ncloud_lb` resource is not supported in Classic environment",
		)
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, conn.DefaultCreateTimeout)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// Validate throughput_type for NETWORK load balancer type
	throughputType := plan.ThroughputType.ValueString()
	if plan.Type.ValueString() == "NETWORK" && throughputType != "" && throughputType != "DYNAMIC" {
		resp.Diagnostics.AddError(
			"Invalid Throughput Type",
			"Network Load Balancer throughput_type can only be set to empty or DYNAMIC",
		)
		return
	}

	reqParams := &vloadbalancer.CreateLoadBalancerInstanceRequest{
		RegionCode: &l.config.RegionCode,
		// Optional
		LoadBalancerDescription:     plan.Description.ValueStringPointer(),
		LoadBalancerNetworkTypeCode: plan.NetworkType.ValueStringPointer(),
		LoadBalancerName:            plan.Name.ValueStringPointer(),

		// Required
		LoadBalancerTypeCode: ncloud.String(plan.Type.ValueString()),
		SubnetNoList: func() []*string {
			elements := make([]*string, 0, len(plan.SubnetNoList.Elements()))
			plan.SubnetNoList.ElementsAs(ctx, &elements, true)
			return elements
		}(),
	}

	if !plan.ThroughputType.IsNull() && !plan.ThroughputType.IsUnknown() {
		reqParams.ThroughputTypeCode = plan.ThroughputType.ValueStringPointer()
	}

	if !plan.IdleTimeout.IsNull() && !plan.IdleTimeout.IsUnknown() {
		reqParams.IdleTimeout = plan.IdleTimeout.ValueInt32Pointer()
	}

	vpcNoMap := make(map[string]int)
	subnetList := make([]*vpc.Subnet, 0)
	for _, subnetNo := range reqParams.SubnetNoList {
		subnet, err := vpcservice.GetSubnetInstance(l.config, *subnetNo)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving subnet instance",
				err.Error(),
			)
			return
		}
		if subnet == nil {
			resp.Diagnostics.AddError(
				"Subnet not found",
				fmt.Sprintf("Subnet with ID %s was not found", *subnetNo),
			)
			return
		}
		subnetList = append(subnetList, subnet)
		vpcNoMap[*subnet.VpcNo]++
	}

	if len(vpcNoMap) > 1 {
		resp.Diagnostics.AddError(
			"Invalid subnet configuration",
			"All subnets must belong to the same VPC",
		)
		return
	}

	reqParams.VpcNo = subnetList[0].VpcNo

	LogCommonRequest("createLoadBalancerInstance", reqParams)
	createResp, err := l.config.Client.Vloadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		LogErrorResponse("createLoadBalancerInstance", err, reqParams)
		resp.Diagnostics.AddError(
			"Error creating load balancer instance",
			err.Error(),
		)
		return
	}
	LogResponse("createLoadBalancerInstance", createResp)

	if err := waitForLoadBalancerActive(ctx, l.config, ncloud.StringValue(createResp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo)); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for load balancer to become active",
			err.Error(),
		)
		return
	}

	output, err := GetFwVpcLoadBalancer(ctx, l.config, ncloud.StringValue(createResp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving created load balancer instance",
			err.Error(),
		)
		return
	}

	plan.LoadBalancerNo = types.StringValue(ncloud.StringValue(createResp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
	if err := plan.refreshFromOutput(ctx, output); err != nil {
		resp.Diagnostics.AddError(
			"Error while getting output values of load balancer instance",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

}

func (l *lbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lbResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetFwVpcLoadBalancer(ctx, l.config, state.LoadBalancerNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := state.refreshFromOutput(ctx, output); err != nil {
		resp.Diagnostics.AddError(
			"Error while getting output values of load balancer instance",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (l *lbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state lbResourceModel

	if !l.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"The `ncloud_lb` resource is not supported in Classic environment",
		)
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	updateTimeout, diags := state.Timeouts.Update(ctx, conn.DefaultUpdateTimeout)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.IdleTimeout.Equal(state.IdleTimeout) ||
		!plan.ThroughputType.Equal(state.ThroughputType) ||
		!plan.Description.Equal(state.Description) {

		if err := waitForLoadBalancerActive(ctx, l.config, state.LoadBalancerNo.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Failed to wait for load balancer to become active",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}

		var err error
		if !plan.IdleTimeout.Equal(state.IdleTimeout) || !plan.ThroughputType.Equal(state.ThroughputType) {
			_, err = l.config.Client.Vloadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(&vloadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
				RegionCode:             &l.config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(state.LoadBalancerNo.ValueString()),
				IdleTimeout:            plan.IdleTimeout.ValueInt32Pointer(),
				ThroughputTypeCode:     plan.ThroughputType.ValueStringPointer(),
			})
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to update load balancer configuration", err.Error(),
				)
				return
			}

			state.IdleTimeout = plan.IdleTimeout
			state.ThroughputType = plan.ThroughputType
		}

		if err := waitForLoadBalancerActive(ctx, l.config, state.LoadBalancerNo.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Failed to wait for load balancer to become active",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}

		if !plan.Description.Equal(state.Description) {
			_, err = l.config.Client.Vloadbalancer.V2Api.SetLoadBalancerDescription(&vloadbalancer.SetLoadBalancerDescriptionRequest{
				RegionCode:              &l.config.RegionCode,
				LoadBalancerInstanceNo:  ncloud.String(state.LoadBalancerNo.ValueString()),
				LoadBalancerDescription: plan.Description.ValueStringPointer(),
			})
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to update load balancer description", err.Error(),
				)
				return
			}

			state.Description = plan.Description
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (l *lbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lbResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, conn.DefaultTimeout)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	reqParams := &vloadbalancer.DeleteLoadBalancerInstancesRequest{
		RegionCode:                 &l.config.RegionCode,
		LoadBalancerInstanceNoList: []*string{ncloud.String(state.LoadBalancerNo.ValueString())},
	}

	tflog.Info(ctx, "DeleteLoadBalancer reqParams="+MarshalUncheckedString(reqParams))

	if err := waitForLoadBalancerActive(ctx, l.config, state.LoadBalancerNo.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAIT FOR LOADBALANCER ERROR", err.Error())
		return
	}

	response, err := l.config.Client.Vloadbalancer.V2Api.DeleteLoadBalancerInstances(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteLoadBalancer response="+MarshalUncheckedString(response))

	if err := waitForLoadBalancerDeletion(ctx, l.config, state.LoadBalancerNo.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
		return
	}
}

func GetFwVpcLoadBalancer(ctx context.Context, config *conn.ProviderConfig, id string) (*LoadBalancerInstance, error) {
	reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(id),
	}
	LogCommonRequest("getLoadBalancerInstanceDetail", reqParams)

	resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getLoadBalancerInstanceDetail", err, reqParams)
		return nil, err
	}
	LogResponse("getLoadBalancerInstanceDetail", resp)

	if len(resp.LoadBalancerInstanceList) < 1 {
		return nil, nil
	}

	return convertVpcLoadBalancer(resp.LoadBalancerInstanceList[0]), nil
}

func waitForLoadBalancerActive(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{LoadBalancerInstanceOperationCreateCode, LoadBalancerInstanceOperationChangeCode},
		Target:  []string{LoadBalancerInstanceOperationNullCode},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(id),
			}
			resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
			if err != nil {
				return nil, "", err
			}

			if len(resp.LoadBalancerInstanceList) < 1 {
				return nil, "", fmt.Errorf("not found load balancer instance(%s)", id)
			}

			lb := resp.LoadBalancerInstanceList[0]
			return resp, ncloud.StringValue(lb.LoadBalancerInstanceOperation.Code), nil
		},
		Timeout:    6 * conn.DefaultTimeout,
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for Load Balancer instance (%s) to become activating: %s", id, err)
	}
	return nil
}

func waitForLoadBalancerDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{LoadBalancerInstanceOperationTerminateCode},
		Target:  []string{LoadBalancerInstanceOperationNullCode},
		Refresh: func() (interface{}, string, error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(id),
			}
			resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
			if err != nil {
				return nil, "", err
			}

			if len(resp.LoadBalancerInstanceList) < 1 {
				return resp, LoadBalancerInstanceOperationNullCode, nil
			}

			respCode := resp.LoadBalancerInstanceList[0].LoadBalancerInstanceOperation.Code
			return nil, ncloud.StringValue(respCode), nil
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for Load Balancer instance (%s) to be deleted: %s", id, err)
	}

	return nil
}

func GetVpcLoadBalancer(config *conn.ProviderConfig, id string) (*LoadBalancerInstance, error) {
	reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(id),
	}
	LogCommonRequest("getLoadBalancerInstanceDetail", reqParams)

	resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getLoadBalancerInstanceDetail", err, reqParams)
		return nil, err
	}
	LogResponse("getLoadBalancerInstanceDetail", resp)

	if len(resp.LoadBalancerInstanceList) < 1 {
		return nil, nil
	}

	return convertVpcLoadBalancer(resp.LoadBalancerInstanceList[0]), nil
}

func (l *lbResourceModel) refreshFromOutput(ctx context.Context, output *LoadBalancerInstance) error {
	l.ID = types.StringPointerValue(output.LoadBalancerInstanceNo)
	l.LoadBalancerNo = types.StringPointerValue(output.LoadBalancerInstanceNo)
	l.Name = types.StringPointerValue(output.LoadBalancerName)
	l.Domain = types.StringPointerValue(output.LoadBalancerDomain)
	l.NetworkType = types.StringPointerValue(output.LoadBalancerNetworkType)
	l.IdleTimeout = types.Int32PointerValue(output.IdleTimeout)
	l.Type = types.StringPointerValue(output.LoadBalancerType)
	l.ThroughputType = types.StringPointerValue(output.ThroughputType)
	l.VpcNo = types.StringPointerValue(output.VpcNo)

	if output.LoadBalancerDescription != nil && *output.LoadBalancerDescription != "" {
		l.Description = types.StringPointerValue(output.LoadBalancerDescription)
	}

	subnetNoList := make([]string, 0)
	for _, subnet := range output.SubnetNoList {
		subnetNoList = append(subnetNoList, *subnet)
	}
	subnetValueList, err := types.ListValueFrom(ctx, types.StringType, subnetNoList)
	if err != nil {
		return fmt.Errorf("error creating ListValue for SubnetNoList: %s", err)
	}
	l.SubnetNoList = subnetValueList

	ipList := make([]string, 0)
	for _, ip := range output.LoadBalancerIpList {
		ipList = append(ipList, *ip)
	}
	ipValueList, err := types.ListValueFrom(ctx, types.StringType, ipList)
	if err != nil {
		return fmt.Errorf("error creating ListValue for IpList: %s", err)
	}
	l.IpList = ipValueList

	listenerNoList := make([]string, 0)
	for _, listener := range output.LoadBalancerListenerList {
		listenerNoList = append(listenerNoList, *listener)
	}
	listenerValueList, err := types.ListValueFrom(ctx, types.StringType, listenerNoList)
	if err != nil {
		return fmt.Errorf("error creating ListValue for ListenerNoList: %s", err)
	}
	l.ListenerNoList = listenerValueList

	return nil
}

func convertVpcLoadBalancer(instance *vloadbalancer.LoadBalancerInstance) *LoadBalancerInstance {
	return &LoadBalancerInstance{
		LoadBalancerInstanceNo:   instance.LoadBalancerInstanceNo,
		LoadBalancerDescription:  instance.LoadBalancerDescription,
		LoadBalancerName:         instance.LoadBalancerName,
		LoadBalancerDomain:       instance.LoadBalancerDomain,
		LoadBalancerIpList:       instance.LoadBalancerIpList,
		LoadBalancerType:         instance.LoadBalancerType.Code,
		LoadBalancerNetworkType:  instance.LoadBalancerNetworkType.Code,
		ThroughputType:           instance.ThroughputType.Code,
		IdleTimeout:              instance.IdleTimeout,
		VpcNo:                    instance.VpcNo,
		SubnetNoList:             instance.SubnetNoList,
		LoadBalancerListenerList: instance.LoadBalancerListenerNoList,
	}
}

type lbResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	LoadBalancerNo types.String   `tfsdk:"load_balancer_no"`
	Name           types.String   `tfsdk:"name"`
	Description    types.String   `tfsdk:"description"`
	Domain         types.String   `tfsdk:"domain"`
	NetworkType    types.String   `tfsdk:"network_type"`
	IdleTimeout    types.Int32    `tfsdk:"idle_timeout"`
	Type           types.String   `tfsdk:"type"`
	ThroughputType types.String   `tfsdk:"throughput_type"`
	VpcNo          types.String   `tfsdk:"vpc_no"`
	SubnetNoList   types.List     `tfsdk:"subnet_no_list"`
	IpList         types.List     `tfsdk:"ip_list"`
	ListenerNoList types.List     `tfsdk:"listener_no_list"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

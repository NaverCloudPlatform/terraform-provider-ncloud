package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	vpcservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
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

func (l *lbResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"idle_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 3600),
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
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
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
			},
			"listener_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
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

	var idleTimeoutValue int32
	if !plan.IdleTimeout.IsNull() {
		idleTimeoutValue = int32(plan.IdleTimeout.ValueInt64())
	} else {
		idleTimeoutValue = 0
	}

	reqParams := &vloadbalancer.CreateLoadBalancerInstanceRequest{
		RegionCode: &l.config.RegionCode,
		// Optional
		IdleTimeout:                 ncloud.Int32(idleTimeoutValue),
		LoadBalancerDescription:     plan.Description.ValueStringPointer(),
		LoadBalancerNetworkTypeCode: plan.NetworkType.ValueStringPointer(),
		LoadBalancerName:            plan.Name.ValueStringPointer(),
		ThroughputTypeCode:          plan.ThroughputType.ValueStringPointer(),

		// Required
		LoadBalancerTypeCode: ncloud.String(plan.Type.ValueString()),
		// SubnetNoList:         listValueToSubnetNoList(ctx, plan.SubnetNoList),
		SubnetNoList: func() []*string {
			elements := make([]*string, 0, len(plan.SubnetNoList.Elements()))
			plan.SubnetNoList.ElementsAs(ctx, &elements, true)
			return elements
		}(),
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

	if err := waitForFwLoadBalancerActive(ctx, plan, l.config, ncloud.StringValue(createResp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo)); err != nil {
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
	plan.refreshFromOutput(ctx, output)

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

	state.refreshFromOutput(ctx, output)

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

	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.IdleTimeout.Equal(state.IdleTimeout) {
		if err := waitForFwLoadBalancerActive(ctx, plan, l.config, state.LoadBalancerNo.String()); err != nil {
			resp.Diagnostics.AddError(
				"Failed to wait for load balancer to become active",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
		_, err := l.config.Client.Vloadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(&vloadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
			RegionCode:             &l.config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(state.LoadBalancerNo.String()),
			IdleTimeout:            ncloud.Int32(int32(plan.IdleTimeout.ValueInt64())),
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to change idle timeout configuration",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
	}

	if !plan.ThroughputType.Equal(state.ThroughputType) {
		if err := waitForFwLoadBalancerActive(ctx, plan, l.config, state.LoadBalancerNo.String()); err != nil {
			resp.Diagnostics.AddError(
				"Failed to wait for load balancer to become active",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
		_, err := l.config.Client.Vloadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(&vloadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
			RegionCode:             &l.config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(state.LoadBalancerNo.String()),
			ThroughputTypeCode:     plan.ThroughputType.ValueStringPointer(),
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to change throughput type configuration",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
	}

	if !plan.Description.Equal(state.Description) {
		if err := waitForFwLoadBalancerActive(ctx, plan, l.config, state.LoadBalancerNo.String()); err != nil {
			resp.Diagnostics.AddError(
				"Failed to wait for load balancer to become active",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
		_, err := l.config.Client.Vloadbalancer.V2Api.SetLoadBalancerDescription(&vloadbalancer.SetLoadBalancerDescriptionRequest{
			RegionCode:              &l.config.RegionCode,
			LoadBalancerInstanceNo:  ncloud.String(state.LoadBalancerNo.String()),
			LoadBalancerDescription: plan.Description.ValueStringPointer(),
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to set load balancer description",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *lbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lbResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vloadbalancer.DeleteLoadBalancerInstancesRequest{
		RegionCode:                 &r.config.RegionCode,
		LoadBalancerInstanceNoList: []*string{ncloud.String(state.LoadBalancerNo.ValueString())},
	}

	tflog.Info(ctx, "DeleteLoadBalancer reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vloadbalancer.V2Api.DeleteLoadBalancerInstances(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteLoadBalancer response="+common.MarshalUncheckedString(response))

	if err := waitFwForLoadBalancerDeletion(ctx, r.config, state.LoadBalancerNo.ValueString()); err != nil {
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

func listValueToSubnetNoList(ctx context.Context, list basetypes.ListValue) []*string {
	result := make([]*string, 0, len(list.Elements()))
	for _, v := range list.Elements() {
		if str, ok := v.(basetypes.StringValue); ok {
			value := str.ValueString()
			result = append(result, &value)
		}
	}
	return result
}

func waitForFwLoadBalancerActive(ctx context.Context, plan lbResourceModel, config *conn.ProviderConfig, id string) error {
	createTimeout := 20 * time.Minute

	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
			RegionCode:             &config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(id),
		}
		resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		if len(resp.LoadBalancerInstanceList) < 1 {
			return retry.NonRetryableError(fmt.Errorf("not found load balancer instance(%s)", id))
		}

		lb := resp.LoadBalancerInstanceList[0]
		operation := ncloud.StringValue(lb.LoadBalancerInstanceOperation.Code)

		switch operation {
		case LoadBalancerInstanceOperationCreateCode, LoadBalancerInstanceOperationChangeCode:
			return retry.RetryableError(fmt.Errorf("expected instance to be active, was %s", operation))
		case LoadBalancerInstanceOperationNullCode:
			return nil
		default:
			return retry.NonRetryableError(fmt.Errorf("unexpected load balancer instance operation: %s", operation))
		}
	})

	if err != nil {
		return fmt.Errorf("error waiting for Load Balancer instance (%s) to become active: %s", id, err)
	}

	return nil
}

func waitFwForLoadBalancerDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PEND"},
		Target:  []string{"DEL"},
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
				return resp, "DEL", nil
			}

			lb := resp.LoadBalancerInstanceList[0]
			status := ncloud.StringValue(lb.LoadBalancerInstanceStatus.Code)
			op := ncloud.StringValue(lb.LoadBalancerInstanceOperation.Code)

			if status == "TERMINATED" && op == "NULL" {
				return resp, "DEL", nil
			}
			if op == "TERMINATING" {
				return resp, "PEND", nil
			}

			return nil, "", fmt.Errorf("error occurred while waiting to delete load balancer instance")
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

func (l *lbResourceModel) refreshFromOutput(ctx context.Context, output *LoadBalancerInstance) {
	l.LoadBalancerNo = types.StringPointerValue(output.LoadBalancerInstanceNo)
	l.Name = types.StringPointerValue(output.LoadBalancerName)
	l.Description = types.StringPointerValue(output.LoadBalancerDescription)
	l.Domain = types.StringPointerValue(output.LoadBalancerDomain)
	l.NetworkType = types.StringPointerValue(output.LoadBalancerNetworkType)
	l.IdleTimeout = types.Int64Value(int64(ncloud.Int32Value(output.IdleTimeout)))
	l.Type = types.StringPointerValue(output.LoadBalancerType)
	l.ThroughputType = types.StringPointerValue(output.ThroughputType)
	l.VpcNo = types.StringPointerValue(output.VpcNo)

	subnetNoList := make([]string, 0)
	for _, subnet := range output.SubnetNoList {
		subnetNoList = append(subnetNoList, *subnet)
	}
	l.SubnetNoList, _ = types.ListValueFrom(ctx, types.StringType, subnetNoList)

	ipList := make([]string, 0)
	for _, ip := range output.LoadBalancerIpList {
		ipList = append(ipList, *ip)
	}
	l.IpList, _ = types.ListValueFrom(ctx, types.StringType, ipList)

	listenerNoList := make([]string, 0)
	for _, listener := range output.LoadBalancerListenerList {
		listenerNoList = append(listenerNoList, *listener)
	}
	l.ListenerNoList, _ = types.ListValueFrom(ctx, types.StringType, listenerNoList)
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
	LoadBalancerNo types.String `tfsdk:"load_balancer_no"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Domain         types.String `tfsdk:"domain"`
	NetworkType    types.String `tfsdk:"network_type"`
	IdleTimeout    types.Int64  `tfsdk:"idle_timeout"`
	Type           types.String `tfsdk:"type"`
	ThroughputType types.String `tfsdk:"throughput_type"`
	VpcNo          types.String `tfsdk:"vpc_no"`
	SubnetNoList   types.List   `tfsdk:"subnet_no_list"`
	IpList         types.List   `tfsdk:"ip_list"`
	ListenerNoList types.List   `tfsdk:"listener_no_list"`
}

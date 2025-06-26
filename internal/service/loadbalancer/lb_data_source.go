package loadbalancer

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ datasource.DataSource              = &loadBalancerDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerDataSource{}
)

func NewLoadBalancerDataSource() datasource.DataSource {
	return &loadBalancerDataSource{}
}

type loadBalancerDataSource struct {
	config *conn.ProviderConfig
}

func (l *loadBalancerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lb"
}

func (l *loadBalancerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"load_balancer_no": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Computed: true,
			},
			"network_type": schema.StringAttribute{
				Computed: true,
			},
			"idle_timeout": schema.Int32Attribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"throughput_type": schema.StringAttribute{
				Computed: true,
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"subnet_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
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
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (l *loadBalancerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (l *loadBalancerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data loadBalancerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vloadbalancer.GetLoadBalancerInstanceListRequest{
		RegionCode: &l.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.LoadBalancerInstanceNoList = []*string{data.ID.ValueStringPointer()}
	}

	tflog.Info(ctx, "GetLoadBalancerInstanceList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	lbResp, err := l.config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceList(reqParams)

	if err != nil {
		resp.Diagnostics.AddError(
			"GetLoadBalancerInstanceList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		return
	}
	tflog.Info(ctx, "GetLoadBalancerInstanceList response", map[string]any{
		"lbResponse": common.MarshalUncheckedString(lbResp),
	})

	lbList, diags := flattenLoadBalancers(ctx, lbResp.LoadBalancerInstanceList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filter, lbList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		resp.Diagnostics.AddError(
			"GetLoadBalancerInstanceList result validation",
			err.Error(),
		)
		return
	}

	state := filteredList[0]
	state.Filter = data.Filter

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenLoadBalancers(ctx context.Context, list []*vloadbalancer.LoadBalancerInstance) ([]*loadBalancerDataSourceModel, diag.Diagnostics) {
	var lbList []*loadBalancerDataSourceModel
	var diags diag.Diagnostics

	for _, lb := range list {
		subnetNumList, _ := types.ListValueFrom(ctx, types.StringType, ncloud.StringListValue(lb.SubnetNoList))
		ipList, _ := types.ListValueFrom(ctx, types.StringType, ncloud.StringListValue(lb.LoadBalancerIpList))
		listenerNoList, _ := types.ListValueFrom(ctx, types.StringType, ncloud.StringListValue(lb.LoadBalancerListenerNoList))

		item := &loadBalancerDataSourceModel{
			ID:             types.StringValue(ncloud.StringValue(lb.LoadBalancerInstanceNo)),
			LoadBalancerNo: types.StringValue(ncloud.StringValue(lb.LoadBalancerInstanceNo)),
			Name:           types.StringValue(ncloud.StringValue(lb.LoadBalancerName)),
			Description:    types.StringValue(ncloud.StringValue(lb.LoadBalancerDescription)),
			Domain:         types.StringValue(ncloud.StringValue(lb.LoadBalancerDomain)),
			NetworkType:    types.StringPointerValue(lb.LoadBalancerNetworkType.Code),
			IdleTimeout:    types.Int32Value(ncloud.Int32Value(lb.IdleTimeout)),
			Type:           types.StringValue(ncloud.StringValue(lb.LoadBalancerType.Code)),
			ThroughputType: types.StringValue(ncloud.StringValue(lb.ThroughputType.Code)),
			VpcNo:          types.StringValue(ncloud.StringValue(lb.VpcNo)),
			SubnetNoList:   subnetNumList,
			IpList:         ipList,
			ListenerNoList: listenerNoList,
		}
		lbList = append(lbList, item)
	}

	return lbList, diags
}

type loadBalancerDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	LoadBalancerNo types.String `tfsdk:"load_balancer_no"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Domain         types.String `tfsdk:"domain"`
	NetworkType    types.String `tfsdk:"network_type"`
	IdleTimeout    types.Int32  `tfsdk:"idle_timeout"`
	Type           types.String `tfsdk:"type"`
	ThroughputType types.String `tfsdk:"throughput_type"`
	VpcNo          types.String `tfsdk:"vpc_no"`
	SubnetNoList   types.List   `tfsdk:"subnet_no_list"`
	IpList         types.List   `tfsdk:"ip_list"`
	ListenerNoList types.List   `tfsdk:"listener_no_list"`
	Filter         types.Set    `tfsdk:"filter"`
}

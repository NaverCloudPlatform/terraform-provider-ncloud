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
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"filter": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"values": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
						"regex": schema.BoolAttribute{
							Optional: true,
						},
					},
				},
			},
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

	if !l.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not Supported Classic",
			"load balancer data source does not support classic",
		)
		return
	}

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

	lbList, diags := flattenLoadBalancers(ctx, lbResp.LoadBalancerInstanceList, l.config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lbPointerList := make([]*loadBalancerDataSourceModel, len(lbList))
	for i := range lbList {
		lbPointerList[i] = &lbList[i]
	}
	filteredList := common.FilterModels(ctx, data.Filter, lbPointerList)

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

func flattenLoadBalancers(ctx context.Context, list []*vloadbalancer.LoadBalancerInstance, config *conn.ProviderConfig) ([]loadBalancerDataSourceModel, diag.Diagnostics) {
	var lbList []loadBalancerDataSourceModel
	var diags diag.Diagnostics

	for _, lb := range list {
		item := loadBalancerDataSourceModel{
			ID:          types.StringValue(ncloud.StringValue(lb.LoadBalancerInstanceNo)),
			Description: types.StringValue(ncloud.StringValue(lb.LoadBalancerDescription)),
		}
		lbList = append(lbList, item)
	}

	return lbList, diags
}

func getVpcLoadBalancerList(config *conn.ProviderConfig, id string) ([]*LoadBalancerInstance, error) {
	reqParams := &vloadbalancer.GetLoadBalancerInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LoadBalancerInstanceNoList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceList(reqParams)
	if err != nil {
		return nil, err
	}

	lbList := make([]*LoadBalancerInstance, 0)
	for _, lb := range resp.LoadBalancerInstanceList {
		lbList = append(lbList, convertVpcLoadBalancer(lb))
	}

	return lbList, nil
}

type loadBalancerDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Filter      types.Set    `tfsdk:"filter"`
}

type filterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
	Regex  types.Bool   `tfsdk:"regex"`
}

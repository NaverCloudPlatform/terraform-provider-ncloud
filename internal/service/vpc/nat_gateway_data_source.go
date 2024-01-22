package vpc

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
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
	_ datasource.DataSource              = &natGatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &natGatewayDataSource{}
)

func NewNatGatewayDataSource() datasource.DataSource {
	return &natGatewayDataSource{}
}

type natGatewayDataSource struct {
	config *conn.ProviderConfig
}

func (n *natGatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nat_gateway"
}

func (n *natGatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"vpc_name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"zone": schema.StringAttribute{
				Computed: true,
			},
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"private_ip": schema.StringAttribute{
				Computed: true,
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
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (n *natGatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (n *natGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data natGatewayDataSourceModel

	if !n.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"nat gateway data source does not supported in classic",
		)
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.GetNatGatewayInstanceListRequest{
		RegionCode: &n.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.NatGatewayInstanceNoList = []*string{data.ID.ValueStringPointer()}
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		reqParams.NatGatewayName = data.Name.ValueStringPointer()
	}

	if !data.VpcName.IsNull() && !data.VpcName.IsUnknown() {
		reqParams.VpcName = data.VpcName.ValueStringPointer()
	}
	tflog.Info(ctx, "GetNatGatewayList reqParams="+common.MarshalUncheckedString(reqParams))

	natGatewayResp, err := n.config.Client.Vpc.V2Api.GetNatGatewayInstanceList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetNatGatewayList response="+common.MarshalUncheckedString(natGatewayResp))

	natGatewayList, diags := flattenNatGateways(natGatewayResp.NatGatewayInstanceList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, natGatewayList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetNatGatewayList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenNatGateways(natGateways []*vpc.NatGatewayInstance) ([]*natGatewayDataSourceModel, diag.Diagnostics) {
	var outputs []*natGatewayDataSourceModel

	for _, v := range natGateways {
		var output natGatewayDataSourceModel

		output.refreshFromOutput(v)
		outputs = append(outputs, &output)
	}

	return outputs, nil
}

type natGatewayDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	NatGatewayNo types.String `tfsdk:"nat_gateway_no"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	PublicIp     types.String `tfsdk:"public_ip"`
	VpcNo        types.String `tfsdk:"vpc_no"`
	VpcName      types.String `tfsdk:"vpc_name"`
	Zone         types.String `tfsdk:"zone"`
	SubnetNo     types.String `tfsdk:"subnet_no"`
	SubnetName   types.String `tfsdk:"subnet_name"`
	PrivateIp    types.String `tfsdk:"private_ip"`
	PublicIpNo   types.String `tfsdk:"public_ip_no"`
	Filters      types.Set    `tfsdk:"filter"`
}

func (d *natGatewayDataSourceModel) refreshFromOutput(output *vpc.NatGatewayInstance) {
	d.ID = types.StringPointerValue(output.NatGatewayInstanceNo)
	d.NatGatewayNo = types.StringPointerValue(output.NatGatewayInstanceNo)
	d.Name = types.StringPointerValue(output.NatGatewayName)
	d.Description = types.StringPointerValue(output.NatGatewayDescription)
	d.PublicIp = types.StringPointerValue(output.PublicIp)
	d.VpcNo = types.StringPointerValue(output.VpcNo)
	d.VpcName = types.StringPointerValue(output.VpcName)
	d.Zone = types.StringPointerValue(output.ZoneCode)
	d.SubnetNo = types.StringPointerValue(output.SubnetNo)
	d.SubnetName = types.StringPointerValue(output.SubnetName)
	d.PrivateIp = types.StringPointerValue(output.PrivateIp)
	d.PublicIpNo = types.StringPointerValue(output.PublicIpInstanceNo)
}

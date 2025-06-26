package vpc

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
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
	_ datasource.DataSource              = &vpcsDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcsDataSource{}
)

func NewVpcDataSource() datasource.DataSource {
	return &vpcDataSource{}
}

type vpcDataSource struct {
	config *conn.ProviderConfig
}

func (v *vpcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc"
}

// Schema defines the schema for the data source.
func (v *vpcDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"ipv4_cidr_block": schema.StringAttribute{
				Computed: true,
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
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (v *vpcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (v *vpcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vpcDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.GetVpcListRequest{
		RegionCode: &v.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.VpcNoList = []*string{data.ID.ValueStringPointer()}
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		reqParams.VpcName = data.Name.ValueStringPointer()
	}

	tflog.Info(ctx, "GetVpcList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	vpcResp, err := v.config.Client.Vpc.V2Api.GetVpcList(reqParams)

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetVpcList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetVpcList response", map[string]any{
		"vpcResponse": common.MarshalUncheckedString(vpcResp),
	})

	vpcList, diags := flattenVpcs(ctx, vpcResp.VpcList, v.config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, vpcList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetVpcList result vaildation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenVpcs(ctx context.Context, vpcs []*vpc.Vpc, config *conn.ProviderConfig) ([]*vpcDataSourceModel, diag.Diagnostics) {
	var outputs []*vpcDataSourceModel

	for _, v := range vpcs {
		var output vpcDataSourceModel

		diags := output.refreshFromOutput(v, config)
		if diags.HasError() {
			return nil, diags
		}

		outputs = append(outputs, &output)
	}

	return outputs, nil
}

type vpcDataSourceModel struct {
	DefaultAccessControlGroupNo types.String `tfsdk:"default_access_control_group_no"`
	DefaultNetworkAclNo         types.String `tfsdk:"default_network_acl_no"`
	DefaultPrivateRouteTableNo  types.String `tfsdk:"default_private_route_table_no"`
	DefaultPublicRouteTableNo   types.String `tfsdk:"default_public_route_table_no"`
	Filters                     types.Set    `tfsdk:"filter"`
	ID                          types.String `tfsdk:"id"`
	Ipv4CidrBlock               types.String `tfsdk:"ipv4_cidr_block"`
	Name                        types.String `tfsdk:"name"`
	VpcNo                       types.String `tfsdk:"vpc_no"`
}

func (d *vpcDataSourceModel) refreshFromOutput(output *vpc.Vpc, config *conn.ProviderConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	id := ncloud.StringValue(output.VpcNo)

	defaultNetworkACLNo, err := getDefaultNetworkACL(config, id)
	if err != nil {
		diags.AddError(
			"GetDefaultNetworkAcl info",
			fmt.Sprintf("error get default network acl for VPC (%s): %s", id, err),
		)
	}

	defaultAcgNo, err := GetDefaultAccessControlGroup(config, id)
	if err != nil {
		diags.AddError(
			"GetDefaultAccessControlGroup info",
			fmt.Sprintf("error get default Access Control Group for VPC (%s): %s", id, err),
		)
	}

	publicRouteTableNo, privateRouteTableNo, err := getDefaultRouteTable(config, id)
	if err != nil {
		diags.AddError(
			"GetDefaultRouteTable info",
			fmt.Sprintf("error get default Route Table for VPC (%s): %s", id, err),
		)
	}

	if diags.HasError() {
		return diags
	}

	d.DefaultAccessControlGroupNo = types.StringValue(defaultAcgNo)
	d.DefaultNetworkAclNo = types.StringValue(defaultNetworkACLNo)
	d.DefaultPrivateRouteTableNo = types.StringValue(privateRouteTableNo)
	d.DefaultPublicRouteTableNo = types.StringValue(publicRouteTableNo)
	d.ID = types.StringPointerValue(output.VpcNo)
	d.Ipv4CidrBlock = types.StringPointerValue(output.Ipv4CidrBlock)
	d.Name = types.StringPointerValue(output.VpcName)
	d.VpcNo = types.StringPointerValue(output.VpcNo)

	return diags
}

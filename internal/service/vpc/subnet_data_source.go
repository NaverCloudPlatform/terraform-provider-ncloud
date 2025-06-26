package vpc

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &subnetDataSource{}
	_ datasource.DataSourceWithConfigure = &subnetDataSource{}
)

func NewSubnetDataSource() datasource.DataSource {
	return &subnetDataSource{}
}

type subnetDataSource struct {
	config *conn.ProviderConfig
}

func (s *subnetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

// Schema defines the schema for the data source.
func (s *subnetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"vpc_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"subnet": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"zone": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"network_acl_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"subnet_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"PUBLIC", "PRIVATE"}...),
				},
			},
			"usage_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"GEN", "LOADB", "BM", "NATGW"}...),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (s *subnetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	s.config = config
}

// Read refreshes the Terraform state with the latest data.
func (s *subnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subnetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.GetSubnetListRequest{
		RegionCode: &s.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.SubnetNoList = []*string{data.ID.ValueStringPointer()}
	}

	if !data.VpcNo.IsNull() && !data.VpcNo.IsUnknown() {
		reqParams.VpcNo = data.VpcNo.ValueStringPointer()
	}

	if !data.Subnet.IsNull() && !data.Subnet.IsUnknown() {
		reqParams.Subnet = data.Subnet.ValueStringPointer()
	}

	if !data.Zone.IsNull() && !data.Zone.IsUnknown() {
		reqParams.ZoneCode = data.Zone.ValueStringPointer()
	}

	if !data.NetworkAclNo.IsNull() && !data.NetworkAclNo.IsUnknown() {
		reqParams.NetworkAclNo = data.NetworkAclNo.ValueStringPointer()
	}

	if !data.SubnetType.IsNull() && !data.SubnetType.IsUnknown() {
		reqParams.SubnetTypeCode = data.SubnetType.ValueStringPointer()
	}

	if !data.UsageType.IsNull() && !data.UsageType.IsUnknown() {
		reqParams.UsageTypeCode = data.UsageType.ValueStringPointer()
	}

	tflog.Info(ctx, "GetSubnetList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	subnetResp, err := s.config.Client.Vpc.V2Api.GetSubnetList(reqParams)

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetSubnetList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetSubnetList response", map[string]any{
		"subnetResponse": common.MarshalUncheckedString(subnetResp),
	})

	subnetList, diags := flattenSubnets(subnetResp.SubnetList, s.config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, subnetList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetSubnetList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenSubnets(subnets []*vpc.Subnet, config *conn.ProviderConfig) ([]*subnetDataSourceModel, diag.Diagnostics) {
	var outputs []*subnetDataSourceModel

	for _, v := range subnets {
		var output subnetDataSourceModel

		diags := output.refreshFromOutput(v, config)
		if diags.HasError() {
			return nil, diags
		}

		outputs = append(outputs, &output)
	}

	return outputs, nil
}

type subnetDataSourceModel struct {
	NetworkAclNo types.String `tfsdk:"network_acl_no"`
	VpcNo        types.String `tfsdk:"vpc_no"`
	Filters      types.Set    `tfsdk:"filter"`
	ID           types.String `tfsdk:"id"`
	Subnet       types.String `tfsdk:"subnet"`
	Zone         types.String `tfsdk:"zone"`
	SubnetType   types.String `tfsdk:"subnet_type"`
	UsageType    types.String `tfsdk:"usage_type"`
	Name         types.String `tfsdk:"name"`
	SubnetNo     types.String `tfsdk:"subnet_no"`
}

func (d *subnetDataSourceModel) refreshFromOutput(output *vpc.Subnet, config *conn.ProviderConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	d.ID = types.StringPointerValue(output.SubnetNo)
	d.SubnetNo = types.StringPointerValue(output.SubnetNo)
	d.VpcNo = types.StringPointerValue(output.VpcNo)
	d.Zone = types.StringPointerValue(output.ZoneCode)
	d.Name = types.StringPointerValue(output.SubnetName)
	d.Subnet = types.StringPointerValue(output.Subnet)
	d.SubnetType = types.StringPointerValue(output.SubnetType.Code)
	d.UsageType = types.StringPointerValue(output.UsageType.Code)
	d.NetworkAclNo = types.StringPointerValue(output.NetworkAclNo)

	return diags
}

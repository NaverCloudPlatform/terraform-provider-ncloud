package vpc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &subnetsDataSource{}
	_ datasource.DataSourceWithConfigure = &subnetsDataSource{}
)

func NewSubnetsDataSource() datasource.DataSource {
	return &subnetsDataSource{}
}

type subnetsDataSource struct {
	config *conn.ProviderConfig
}

func (s *subnetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnets"
}

// Schema defines the schema for the data source.
func (s *subnetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"subnet_no": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of subnet ID to retrieve",
			},
			"vpc_no": schema.StringAttribute{
				Optional:    true,
				Description: "The VPC ID that you want to filter from",
			},
			"subnet": schema.StringAttribute{
				Optional:    true,
				Description: "The CIDR block for the subnet.",
			},
			"zone": schema.StringAttribute{
				Optional:    true,
				Description: "Available Zone. Get available values using the `data ncloud_zones`.",
			},
			"network_acl_no": schema.StringAttribute{
				Optional:    true,
				Description: "Network ACL No. Get available values using the `default_network_acl_no` from Resource `ncloud_vpc` or Data source `data.ncloud_network_acls`.",
			},
			"subnet_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"PUBLIC", "PRIVATE"}...),
				},
				Description: "Internet Gateway Only. PUBLC(Yes/Public), PRIVATE(No/Private).",
			},
			"usage_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"GEN", "LOADB", "BM", "NATGW"}...),
				},
				Description: "Usage type. GEN(Normal), LOADB(Load Balance), BM(BareMetal), NATGW(NAT Gateway). default : GEN(Normal).",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
			"subnets": schema.SetNestedBlock{
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"id": framework.IDAttribute(),
						"vpc_no": schema.StringAttribute{
							Computed: true,
						},
						"subnet": schema.StringAttribute{
							Computed: true,
						},
						"zone": schema.StringAttribute{
							Computed: true,
						},
						"network_acl_no": schema.StringAttribute{
							Computed: true,
						},
						"subnet_type": schema.StringAttribute{
							Computed: true,
						},
						"usage_type": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.ListAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (s *subnetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (s *subnetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !s.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not Supported Classic",
			"subnets data source does not supported in classic",
		)
		return
	}

	var data subnetsDataSourceModel
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

	subnetList, diags := flattenSubnets(ctx, subnetResp.SubnetList, s.config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, subnetList)

	state := data
	state.ID = types.StringValue(time.Now().UTC().String())
	resp.Diagnostics.Append(state.refreshFromSubnetOutputModel(ctx, filteredList, s.config)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type subnetsDataSourceModel struct {
	Filters      types.Set    `tfsdk:"filter"`
	ID           types.String `tfsdk:"id"`
	SubnetNo     types.String `tfsdk:"subnet_no"`
	VpcNo        types.String `tfsdk:"vpc_no"`
	Subnet       types.String `tfsdk:"subnet"`
	Zone         types.String `tfsdk:"zone"`
	NetworkAclNo types.String `tfsdk:"network_acl_no"`
	SubnetType   types.String `tfsdk:"subnet_type"`
	UsageType    types.String `tfsdk:"usage_type"`
	Subnets      types.Set    `tfsdk:"subnets"`
}

func (d *subnetsDataSourceModel) refreshFromSubnetOutputModel(ctx context.Context, subnetModels []*subnetDataSourceModel, config *conn.ProviderConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: subnetAttrTypes}
	elems := []attr.Value{}

	for _, model := range subnetModels {
		obj := map[string]attr.Value{
			"network_acl_no": model.NetworkAclNo,
			"vpc_no":         model.VpcNo,
			"id":             model.ID,
			"subnet":         model.Subnet,
			"zone":           model.Zone,
			"subnet_type":    model.SubnetType,
			"usage_type":     model.UsageType,
			"name":           model.Name,
			"subnet_no":      model.SubnetNo,
		}
		objVal, di := types.ObjectValue(subnetAttrTypes, obj)
		diags.Append(di...)

		elems = append(elems, objVal)
	}
	setVal, di := types.SetValue(elemType, elems)
	diags.Append(di...)

	if diags.HasError() {
		return diags
	}

	d.Subnets = setVal
	return diags
}

var (
	subnetAttrTypes = map[string]attr.Type{
		"network_acl_no": types.StringType,
		"vpc_no":         types.StringType,
		"id":             types.StringType,
		"subnet":         types.StringType,
		"zone":           types.StringType,
		"subnet_type":    types.StringType,
		"usage_type":     types.StringType,
		"name":           types.StringType,
		"subnet_no":      types.StringType,
	}
)

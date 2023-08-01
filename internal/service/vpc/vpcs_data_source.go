package vpc

import (
	"context"
	"fmt"
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

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ datasource.DataSource              = &vpcsDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcsDataSource{}
)

func NewVpcsDataSource() datasource.DataSource {
	return &vpcsDataSource{}
}

type vpcsDataSource struct {
	config *conn.ProviderConfig
}

func (v *vpcsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpcs"
}

// Schema defines the schema for the data source.
func (v *vpcsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"name": schema.StringAttribute{
				Optional: true,
			},
			"vpc_no": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
			"vpcs": schema.SetNestedBlock{
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"id": framework.IDAttribute(),
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
				},
			},
		},
	}
}

func (v *vpcsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (v *vpcsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !v.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not Supported Classic",
			"vpcs data source does not supported in classic",
		)
		return
	}

	var data vpcsDataSourceModel
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

	state := data
	state.ID = types.StringValue(time.Now().UTC().String())
	resp.Diagnostics.Append(state.refreshFromVpcOutputModel(ctx, filteredList, v.config)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type vpcsDataSourceModel struct {
	Filters types.Set    `tfsdk:"filter"`
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	VpcNo   types.String `tfsdk:"vpc_no"`
	Vpcs    types.Set    `tfsdk:"vpcs"`
}

func (d *vpcsDataSourceModel) refreshFromVpcOutputModel(ctx context.Context, vpcModels []*vpcDataSourceModel, config *conn.ProviderConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: vpcAttrTypes}
	elems := []attr.Value{}

	for _, model := range vpcModels {
		obj := map[string]attr.Value{
			"default_access_control_group_no": model.DefaultAccessControlGroupNo,
			"default_network_acl_no":          model.DefaultNetworkAclNo,
			"default_private_route_table_no":  model.DefaultPrivateRouteTableNo,
			"default_public_route_table_no":   model.DefaultPublicRouteTableNo,
			"id":                              model.ID,
			"ipv4_cidr_block":                 model.Ipv4CidrBlock,
			"name":                            model.Name,
			"vpc_no":                          model.VpcNo,
		}
		objVal, di := types.ObjectValue(vpcAttrTypes, obj)
		diags.Append(di...)

		elems = append(elems, objVal)
	}
	setVal, di := types.SetValue(elemType, elems)
	diags.Append(di...)

	if diags.HasError() {
		return diags
	}

	d.Vpcs = setVal
	return diags
}

var (
	vpcAttrTypes = map[string]attr.Type{
		"default_access_control_group_no": types.StringType,
		"default_network_acl_no":          types.StringType,
		"default_private_route_table_no":  types.StringType,
		"default_public_route_table_no":   types.StringType,
		"id":                              types.StringType,
		"ipv4_cidr_block":                 types.StringType,
		"name":                            types.StringType,
		"vpc_no":                          types.StringType,
	}
)

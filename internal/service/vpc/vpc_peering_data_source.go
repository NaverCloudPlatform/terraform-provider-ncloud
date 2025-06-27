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
	_ datasource.DataSource              = &vpcPeeringDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcPeeringDataSource{}
)

func NewVpcPeeringDataSource() datasource.DataSource {
	return &vpcPeeringDataSource{}
}

type vpcPeeringDataSource struct {
	config *conn.ProviderConfig
}

func (v *vpcPeeringDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (v *vpcPeeringDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_peering"
}

func (v *vpcPeeringDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"description": schema.StringAttribute{
				Computed: true,
			},
			"source_vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"source_vpc_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"target_vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"target_vpc_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"target_vpc_login_id": schema.StringAttribute{
				Computed: true,
			},
			"vpc_peering_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"has_reverse_vpc_peering": schema.BoolAttribute{
				Computed: true,
			},
			"is_between_accounts": schema.BoolAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (v *vpcPeeringDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vpcPeeringDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpc.GetVpcPeeringInstanceListRequest{
		RegionCode: &v.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.VpcPeeringInstanceNoList = []*string{data.ID.ValueStringPointer()}
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		reqParams.VpcPeeringName = data.Name.ValueStringPointer()
	}

	if !data.TargetVpcName.IsNull() && !data.TargetVpcName.IsUnknown() {
		reqParams.TargetVpcName = data.TargetVpcName.ValueStringPointer()
	}

	if !data.SourceVpcName.IsNull() && !data.SourceVpcName.IsUnknown() {
		reqParams.SourceVpcName = data.SourceVpcName.ValueStringPointer()
	}

	tflog.Info(ctx, "GetVpcPeeringList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	response, err := v.config.Client.Vpc.V2Api.GetVpcPeeringInstanceList(reqParams)

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetVpcPeeringList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetVpcPeeringList response", map[string]any{
		"vpcPeeringResponse": common.MarshalUncheckedString(response),
	})

	vpcPeeringList, diags := flattenVpcPeerings(ctx, response.VpcPeeringInstanceList, v.config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, vpcPeeringList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetVpcPeeringList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenVpcPeerings(ctx context.Context, vpcPeerings []*vpc.VpcPeeringInstance, config *conn.ProviderConfig) ([]*vpcPeeringDataSourceModel, diag.Diagnostics) {
	var outputs []*vpcPeeringDataSourceModel

	for _, v := range vpcPeerings {
		var output vpcPeeringDataSourceModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}

	return outputs, nil
}

func (v *vpcPeeringDataSourceModel) refreshFromOutput(output *vpc.VpcPeeringInstance) {
	v.ID = types.StringPointerValue(output.VpcPeeringInstanceNo)
	v.Name = types.StringPointerValue(output.VpcPeeringName)
	v.SourceVpcName = types.StringPointerValue(output.SourceVpcName)
	v.TargetVpcName = types.StringPointerValue(output.TargetVpcName)
	v.Description = types.StringPointerValue(output.VpcPeeringDescription)
	v.SourceVpcNo = types.StringPointerValue(output.SourceVpcNo)
	v.TargetVpcNo = types.StringPointerValue(output.TargetVpcNo)
	v.TargetVpcLoginId = types.StringPointerValue(output.TargetVpcLoginId)
	v.VpcPeeringNo = types.StringPointerValue(output.VpcPeeringInstanceNo)
	v.HasReverseVpcPeering = types.BoolPointerValue(output.HasReverseVpcPeering)
	v.IsBetweenAccounts = types.BoolPointerValue(output.IsBetweenAccounts)
}

type vpcPeeringDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	SourceVpcName        types.String `tfsdk:"source_vpc_name"`
	TargetVpcName        types.String `tfsdk:"target_vpc_name"`
	Description          types.String `tfsdk:"description"`
	SourceVpcNo          types.String `tfsdk:"source_vpc_no"`
	TargetVpcNo          types.String `tfsdk:"target_vpc_no"`
	TargetVpcLoginId     types.String `tfsdk:"target_vpc_login_id"`
	VpcPeeringNo         types.String `tfsdk:"vpc_peering_no"`
	HasReverseVpcPeering types.Bool   `tfsdk:"has_reverse_vpc_peering"`
	IsBetweenAccounts    types.Bool   `tfsdk:"is_between_accounts"`
	Filters              types.Set    `tfsdk:"filter"`
}

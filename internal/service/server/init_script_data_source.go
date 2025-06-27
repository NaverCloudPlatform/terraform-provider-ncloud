package server

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ datasource.DataSource              = &initScriptDataSource{}
	_ datasource.DataSourceWithConfigure = &initScriptDataSource{}
)

func NewInitScriptDataSource() datasource.DataSource {
	return &initScriptDataSource{}
}

type initScriptDataSource struct {
	config *conn.ProviderConfig
}

func (i *initScriptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_init_script"
}

// Schema defines the schema for the data source.
func (i *initScriptDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"init_script_no": schema.StringAttribute{
				Computed: true,
			},
			"os_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"LNX", "WND"}...),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (i *initScriptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	i.config = config
}

// Read refreshes the Terraform state with the latest data.
func (i *initScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data initScriptDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vserver.GetInitScriptListRequest{
		RegionCode: &i.config.RegionCode,
	}
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.InitScriptNoList = []*string{data.ID.ValueStringPointer()}
	}
	if !data.OsType.IsNull() && !data.OsType.IsUnknown() {
		reqParams.OsTypeCode = data.OsType.ValueStringPointer()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		reqParams.InitScriptName = data.Name.ValueStringPointer()
	}

	tflog.Info(ctx, "GetVpcInitScriptList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	initScriptResp, err := i.config.Client.Vserver.V2Api.GetInitScriptList(reqParams)

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetNatGatewayList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetVpcInitScriptList response", map[string]any{
		"initScriptResponse": common.MarshalUncheckedString(initScriptResp),
	})

	initScriptList, diags := flattenNatGateways(initScriptResp.InitScriptList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, initScriptList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetInitScriptList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenNatGateways(natGateways []*vserver.InitScript) ([]*initScriptDataSourceModel, diag.Diagnostics) {
	var outputs []*initScriptDataSourceModel

	for _, v := range natGateways {
		var output initScriptDataSourceModel

		diags := output.refreshFromOutput(v)
		if diags.HasError() {
			return nil, diags
		}

		outputs = append(outputs, &output)
	}

	return outputs, nil
}

type initScriptDataSourceModel struct {
	Description  types.String `tfsdk:"description"`
	Filters      types.Set    `tfsdk:"filter"`
	ID           types.String `tfsdk:"id"`
	OsType       types.String `tfsdk:"os_type"`
	Name         types.String `tfsdk:"name"`
	InitScriptNo types.String `tfsdk:"init_script_no"`
}

func (m *initScriptDataSourceModel) refreshFromOutput(output *vserver.InitScript) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringPointerValue(output.InitScriptNo)
	m.Name = types.StringPointerValue(output.InitScriptName)
	m.Description = framework.EmptyStringToNull(types.StringPointerValue(output.InitScriptDescription))
	m.OsType = types.StringPointerValue(output.OsType.Code)
	m.InitScriptNo = types.StringPointerValue(output.InitScriptNo)

	return diags
}

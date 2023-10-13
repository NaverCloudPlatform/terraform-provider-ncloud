package mysql

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ datasource.DataSource              = &mysqlImageProductDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlImageProductDataSource{}
)

func NewMysqlImageProductDataSource() datasource.DataSource {
	return &mysqlImageProductDataSource{}
}

type mysqlImageProductDataSource struct {
	config *conn.ProviderConfig
}

func (m *mysqlImageProductDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	m.config = config
}

func (m *mysqlImageProductDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_image_product"
}

func (m *mysqlImageProductDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"product_code": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"exclusion_product_code": schema.StringAttribute{
				Optional: true,
			},
			"generation_code": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("G2", "G3"),
				},
			},
			"product_name": schema.StringAttribute{
				Computed: true,
			},
			"product_type": schema.StringAttribute{
				Computed: true,
			},
			"infra_resource_type": schema.StringAttribute{
				Computed: true,
			},
			"product_description": schema.StringAttribute{
				Computed: true,
			},
			"platform_type": schema.StringAttribute{
				Computed: true,
			},
			"os_information": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mysqlImageProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mysqlImageProductDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmysql.GetCloudMysqlImageProductListRequest{
		RegionCode: &m.config.RegionCode,
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetMysqlImageProductList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	mysqlImageProductResp, err := m.config.Client.Vmysql.V2Api.GetCloudMysqlImageProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMysqlImageProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetMysqlImageProductList response", map[string]any{
		"mysqlImageProductResponse": common.MarshalUncheckedString(mysqlImageProductResp),
	})

	mysqlImageProductList := flattenMysqlImageProducts(ctx, mysqlImageProductResp.ProductList)

	fillteredList := common.FilterModels(ctx, data.Filters, mysqlImageProductList)

	if err := verify.ValidateOneResult(len(fillteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMysqlImageProductList result more than one",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := fillteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenMysqlImageProducts(ctx context.Context, list []*vmysql.Product) []*mysqlImageProductDataSourceModel {
	var outputs []*mysqlImageProductDataSourceModel

	for _, v := range list {
		var output mysqlImageProductDataSourceModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

type mysqlImageProductDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ProductCode          types.String `tfsdk:"product_code"`
	ExclusionProductCode types.String `tfsdk:"exclusion_product_code"`
	GenerationCode       types.String `tfsdk:"generation_code"`
	ProductName          types.String `tfsdk:"product_name"`
	ProductType          types.String `tfsdk:"product_type"`
	InfraResourceType    types.String `tfsdk:"infra_resource_type"`
	PlatformType         types.String `tfsdk:"platform_type"`
	OsInformation        types.String `tfsdk:"os_information"`
	ProductDescription   types.String `tfsdk:"product_description"`
	Filters              types.Set    `tfsdk:"filter"`
}

func (m *mysqlImageProductDataSourceModel) refreshFromOutput(output *vmysql.Product) {
	m.ID = types.StringPointerValue(output.ProductCode)
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.GenerationCode = types.StringPointerValue(output.GenerationCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	m.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	m.OsInformation = types.StringPointerValue(output.OsInformation)
	m.ProductDescription = types.StringPointerValue(output.ProductDescription)
}

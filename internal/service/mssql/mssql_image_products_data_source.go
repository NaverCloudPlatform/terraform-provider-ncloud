package mssql

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmssql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mssqlImageProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &mssqlImageProductsDataSource{}
)

func NewMssqlImageProductsDataSource() datasource.DataSource {
	return &mssqlImageProductsDataSource{}
}

type mssqlImageProductsDataSource struct {
	config *conn.ProviderConfig
}

func (m *mssqlImageProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mssql_image_products"
}

func (m *mssqlImageProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"image_product_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"generation_code": schema.StringAttribute{
							Computed: true,
						},
						"product_name": schema.StringAttribute{
							Computed: true,
						},
						"product_type": schema.StringAttribute{
							Computed: true,
						},
						"platform_type": schema.StringAttribute{
							Computed: true,
						},
						"os_information": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mssqlImageProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mssqlImageProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mssqlImageProductsDataSourceModel

	if !m.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmssql.GetCloudMssqlImageProductListRequest{
		RegionCode: &m.config.RegionCode,
	}
	tflog.Info(ctx, "GetMssqlImageProductList reqParams="+common.MarshalUncheckedString(reqParams))

	mssqlImageProductResp, err := m.config.Client.Vmssql.V2Api.GetCloudMssqlImageProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetMssqlImageProductList response="+common.MarshalUncheckedString(mssqlImageProductResp))

	if mssqlImageProductResp == nil || len(mssqlImageProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	mssqlImageProductList := flattenMssqlImageProduct(mssqlImageProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, mssqlImageProductList)
	data.refreshFromOutput(ctx, fillteredList)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if err := common.WriteImageProductToFile(outputPath, data.ImageProductList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenMssqlImageProduct(list []*vmssql.Product) []*mssqlImageProduct {
	var outputs []*mssqlImageProduct

	for _, v := range list {
		var output mssqlImageProduct
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (m *mssqlImageProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*mssqlImageProduct) {
	imageProductListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mssqlImageProduct{}.attrTypes()}, list)
	m.ImageProductList = imageProductListValue
	m.ID = types.StringValue("")
}

type mssqlImageProductsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductList types.List   `tfsdk:"image_product_list"`
	OutputFile       types.String `tfsdk:"output_file"`
	Filters          types.Set    `tfsdk:"filter"`
}

type mssqlImageProduct struct {
	ProductCode    types.String `tfsdk:"product_code"`
	GenerationCode types.String `tfsdk:"generation_code"`
	ProductName    types.String `tfsdk:"product_name"`
	ProductType    types.String `tfsdk:"product_type"`
	PlatformType   types.String `tfsdk:"platform_type"`
	OsInformation  types.String `tfsdk:"os_information"`
}

type mssqlImageProductToJsonConvert struct {
	ProductCode    string `json:"product_code"`
	GenerationCode string `json:"generation_code"`
	ProductName    string `json:"product_name"`
	ProductType    string `json:"product_type"`
	PlatformType   string `json:"platform_type"`
	OsInformation  string `json:"os_information"`
}

func (m mssqlImageProduct) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":    types.StringType,
		"generation_code": types.StringType,
		"product_name":    types.StringType,
		"product_type":    types.StringType,
		"platform_type":   types.StringType,
		"os_information":  types.StringType,
	}
}

func (m *mssqlImageProduct) refreshFromOutput(output *vmssql.Product) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.GenerationCode = types.StringPointerValue(output.GenerationCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	m.OsInformation = types.StringPointerValue(output.OsInformation)
}

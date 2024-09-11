package mysql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mysqlImageProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlImageProductsDataSource{}
)

func NewMysqlImageProductsDataSource() datasource.DataSource {
	return &mysqlImageProductsDataSource{}
}

type mysqlImageProductsDataSource struct {
	config *conn.ProviderConfig
}

func (m *mysqlImageProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_image_products"
}

func (m *mysqlImageProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (m *mysqlImageProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mysqlImageProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mysqlImageProductsDataSourceModel

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

	reqParams := &vmysql.GetCloudMysqlImageProductListRequest{
		RegionCode: &m.config.RegionCode,
	}
	tflog.Info(ctx, "GetMysqlImageProductList reqParams="+common.MarshalUncheckedString(reqParams))

	mysqlImageProductResp, err := m.config.Client.Vmysql.V2Api.GetCloudMysqlImageProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetMysqlImageProductList response="+common.MarshalUncheckedString(mysqlImageProductResp))

	if mysqlImageProductResp == nil || len(mysqlImageProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	mysqlImageProductList := flattenMysqlImageProduct(mysqlImageProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, mysqlImageProductList)
	data.refreshFromOutput(ctx, fillteredList)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertToJsonStruct(data.ImageProductList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertToJsonStruct(images []attr.Value) ([]mysqlImageProductToJsonConvert, error) {
	var mysqlImagesToConvert = []mysqlImageProductToJsonConvert{}

	for _, image := range images {
		imageJasn := mysqlImageProductToJsonConvert{}
		if err := json.Unmarshal([]byte(image.String()), &imageJasn); err != nil {
			return nil, err
		}
		mysqlImagesToConvert = append(mysqlImagesToConvert, imageJasn)
	}

	return mysqlImagesToConvert, nil
}

func flattenMysqlImageProduct(list []*vmysql.CloudDbProduct) []*mysqlImageProduct {
	var outputs []*mysqlImageProduct

	for _, v := range list {
		var output mysqlImageProduct
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (m *mysqlImageProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*mysqlImageProduct) {
	imageProductListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlImageProduct{}.attrTypes()}, list)
	m.ImageProductList = imageProductListValue
	m.ID = types.StringValue("")
}

type mysqlImageProductsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductList types.List   `tfsdk:"image_product_list"`
	OutputFile       types.String `tfsdk:"output_file"`
	Filters          types.Set    `tfsdk:"filter"`
}

type mysqlImageProduct struct {
	ProductCode    types.String `tfsdk:"product_code"`
	GenerationCode types.String `tfsdk:"generation_code"`
	ProductName    types.String `tfsdk:"product_name"`
	ProductType    types.String `tfsdk:"product_type"`
	PlatformType   types.String `tfsdk:"platform_type"`
	OsInformation  types.String `tfsdk:"os_information"`
}

type mysqlImageProductToJsonConvert struct {
	ProductCode    string `json:"product_code"`
	GenerationCode string `json:"generation_code"`
	ProductName    string `json:"product_name"`
	ProductType    string `json:"product_type"`
	PlatformType   string `json:"platform_type"`
	OsInformation  string `json:"os_information"`
}

func (m mysqlImageProduct) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":    types.StringType,
		"generation_code": types.StringType,
		"product_name":    types.StringType,
		"product_type":    types.StringType,
		"platform_type":   types.StringType,
		"os_information":  types.StringType,
	}
}

func (m *mysqlImageProduct) refreshFromOutput(output *vmysql.CloudDbProduct) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.GenerationCode = types.StringPointerValue(output.GenerationCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	m.OsInformation = types.StringPointerValue(output.OsInformation)
}

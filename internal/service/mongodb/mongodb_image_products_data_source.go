package mongodb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mongodbImageProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &mongodbImageProductsDataSource{}
)

func NewMongoDbImageProductsDataSource() datasource.DataSource {
	return &mongodbImageProductsDataSource{}
}

type mongodbImageProductsDataSource struct {
	config *conn.ProviderConfig
}

func (m *mongodbImageProductsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_image_products"
}

func (m *mongodbImageProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"engine_version_code": schema.StringAttribute{
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

func (m *mongodbImageProductsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mongodbImageProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mongodbImageProductsDataSourceModel

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

	reqParams := &vmongodb.GetCloudMongoDbImageProductListRequest{
		RegionCode: &m.config.RegionCode,
	}
	tflog.Info(ctx, "GetMongoDbImageProductList reqParams="+common.MarshalUncheckedString(reqParams))

	mongodbImageProductResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbImageProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetMongoDbImageProductList response="+common.MarshalUncheckedString(mongodbImageProductResp))

	if mongodbImageProductResp == nil || len(mongodbImageProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	mongodbImageProductList := flattenMongoDbImageProduct(ctx, mongodbImageProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, mongodbImageProductList)
	data.refreshFromOutput(ctx, fillteredList)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertToJsonStruct(data.ImageProductList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertToJsonStruct(images []attr.Value) ([]mongodbImageProductToJsonConvert, error) {
	var mongodbImagesToConvert = []mongodbImageProductToJsonConvert{}

	for _, image := range images {
		imageJasn := mongodbImageProductToJsonConvert{}
		if err := json.Unmarshal([]byte(image.String()), &imageJasn); err != nil {
			return nil, err
		}
		mongodbImagesToConvert = append(mongodbImagesToConvert, imageJasn)
	}

	return mongodbImagesToConvert, nil
}

func flattenMongoDbImageProduct(ctx context.Context, list []*vmongodb.Product) []*mongodbImageProduct {
	var outputs []*mongodbImageProduct

	for _, v := range list {
		var output mongodbImageProduct
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (m *mongodbImageProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*mongodbImageProduct) {
	imageProductListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mongodbImageProduct{}.attrTypes()}, list)
	m.ImageProductList = imageProductListValue
	m.ID = types.StringValue("")
}

type mongodbImageProductsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductList types.List   `tfsdk:"image_product_list"`
	OutputFile       types.String `tfsdk:"output_file"`
	Filters          types.Set    `tfsdk:"filter"`
}

type mongodbImageProduct struct {
	ProductCode       types.String `tfsdk:"product_code"`
	GenerationCode    types.String `tfsdk:"generation_code"`
	ProductName       types.String `tfsdk:"product_name"`
	ProductType       types.String `tfsdk:"product_type"`
	PlatformType      types.String `tfsdk:"platform_type"`
	OsInformation     types.String `tfsdk:"os_information"`
	EngineVersionCode types.String `tfsdk:"engine_version_code"`
}

type mongodbImageProductToJsonConvert struct {
	ProductCode       string `json:"product_code"`
	GenerationCode    string `json:"generation_code"`
	ProductName       string `json:"product_name"`
	ProductType       string `json:"product_type"`
	PlatformType      string `json:"platform_type"`
	OsInformation     string `json:"os_information"`
	EngineVersionCode string `json:"engine_version_code"`
}

func (m mongodbImageProduct) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":        types.StringType,
		"generation_code":     types.StringType,
		"product_name":        types.StringType,
		"product_type":        types.StringType,
		"platform_type":       types.StringType,
		"os_information":      types.StringType,
		"engine_version_code": types.StringType,
	}
}

func (m *mongodbImageProduct) refreshFromOutput(output *vmongodb.Product) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.GenerationCode = types.StringPointerValue(output.GenerationCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	m.OsInformation = types.StringPointerValue(output.OsInformation)
	m.EngineVersionCode = types.StringPointerValue(output.EngineVersionCode)
}

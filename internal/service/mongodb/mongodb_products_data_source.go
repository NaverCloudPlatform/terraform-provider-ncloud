package mongodb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mongodbProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &mongodbProductsDataSource{}
)

func NewMongoDbProductsDataSource() datasource.DataSource {
	return &mongodbProductsDataSource{}
}

type mongodbProductsDataSource struct {
	config *conn.ProviderConfig
}

func (m *mongodbProductsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_products"
}

func (m *mongodbProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Required: true,
			},
			"infra_resource_detail_type_code": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MNGOD", "MNGOS", "ARBIT", "CFGSV"}...),
				},
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"product_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"product_name": schema.StringAttribute{
							Computed: true,
						},
						"product_type": schema.StringAttribute{
							Computed: true,
						},
						"product_description": schema.StringAttribute{
							Computed: true,
						},
						"infra_resource_type": schema.StringAttribute{
							Computed: true,
						},
						"infra_resource_detail_type": schema.StringAttribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"disk_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mongodbProductsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mongodbProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mongodbProductsDataSourceModel

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

	reqParams := &vmongodb.GetCloudMongoDbProductListRequest{
		RegionCode:                   &m.config.RegionCode,
		CloudMongoDbImageProductCode: data.CloudMongoDbImageProductCode.ValueStringPointer(),
	}

	if !data.InfraResourceDetailTypeCode.IsNull() && !data.InfraResourceDetailTypeCode.IsUnknown() {
		reqParams.InfraResourceDetailTypeCode = data.InfraResourceDetailTypeCode.ValueStringPointer()
	}
	tflog.Info(ctx, "GetMongoDbProductsList reqParams="+common.MarshalUncheckedString(reqParams))

	mongodbProductResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetMongoDbProductList response="+common.MarshalUncheckedString(mongodbProductResp))

	if mongodbProductResp == nil || len(mongodbProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	mongodbProductList := flattenMongoDbProductLists(ctx, mongodbProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, mongodbProductList)

	data.refreshFromOutput(ctx, fillteredList)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertProductsToJsonStruct(data.ProductList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertProductsToJsonStruct(images []attr.Value) ([]mongodbProductsToJsonConvert, error) {
	var mongodbProductsToConvert = []mongodbProductsToJsonConvert{}

	for _, image := range images {
		imageJasn := mongodbProductsToJsonConvert{}
		if err := json.Unmarshal([]byte(image.String()), &imageJasn); err != nil {
			return nil, err
		}
		mongodbProductsToConvert = append(mongodbProductsToConvert, imageJasn)
	}

	return mongodbProductsToConvert, nil
}

func flattenMongoDbProductLists(ctx context.Context, list []*vmongodb.Product) []*mongodbProductModel {
	var outputs []*mongodbProductModel

	for _, v := range list {
		var output mongodbProductModel
		output.refreshFromOutput(ctx, v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (m *mongodbProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*mongodbProductModel) {
	productListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mongodbProductModel{}.attrTypes()}, list)
	m.ProductList = productListValue
	m.ID = types.StringValue("")
}

type mongodbProductsDataSourceModel struct {
	ID                           types.String `tfsdk:"id"`
	CloudMongoDbImageProductCode types.String `tfsdk:"image_product_code"`
	InfraResourceDetailTypeCode  types.String `tfsdk:"infra_resource_detail_type_code"`
	ProductList                  types.List   `tfsdk:"product_list"`
	OutputFile                   types.String `tfsdk:"output_file"`
	Filters                      types.Set    `tfsdk:"filter"`
}

type mongodbProductModel struct {
	ProductCode             types.String `tfsdk:"product_code"`
	ProductName             types.String `tfsdk:"product_name"`
	ProductType             types.String `tfsdk:"product_type"`
	ProductDescription      types.String `tfsdk:"product_description"`
	InfraResourceType       types.String `tfsdk:"infra_resource_type"`
	InfraResourceDetailType types.String `tfsdk:"infra_resource_detail_type"`
	CpuCount                types.Int64  `tfsdk:"cpu_count"`
	MemorySize              types.Int64  `tfsdk:"memory_size"`
	DiskType                types.String `tfsdk:"disk_type"`
}

type mongodbProductsToJsonConvert struct {
	ProductCode             string `json:"product_code"`
	ProductName             string `json:"product_name"`
	ProductType             string `json:"product_type"`
	ProductDescription      string `json:"product_description"`
	InfraResourceType       string `json:"infra_resource_type"`
	InfraResourceDetailType string `json:"infra_resource_detail_type"`
	CpuCount                int    `json:"cpu_count"`
	MemorySize              int    `json:"memory_size"`
	DiskType                string `json:"disk_type"`
}

func (m mongodbProductModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":               types.StringType,
		"product_name":               types.StringType,
		"product_type":               types.StringType,
		"product_description":        types.StringType,
		"infra_resource_type":        types.StringType,
		"infra_resource_detail_type": types.StringType,
		"cpu_count":                  types.Int64Type,
		"memory_size":                types.Int64Type,
		"disk_type":                  types.StringType,
	}
}

func (m *mongodbProductModel) refreshFromOutput(ctx context.Context, output *vmongodb.Product) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.ProductDescription = types.StringPointerValue(output.ProductDescription)
	m.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	m.InfraResourceDetailType = types.StringPointerValue(output.InfraResourceDetailType.Code)
	m.CpuCount = types.Int64Value(int64(*output.CpuCount))
	m.MemorySize = types.Int64Value(int64(*output.MemorySize))
	m.DiskType = types.StringPointerValue(output.DiskType.Code)
}

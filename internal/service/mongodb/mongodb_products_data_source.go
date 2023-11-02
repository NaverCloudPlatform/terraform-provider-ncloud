package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
			"product_code": schema.StringAttribute{
				Optional: true,
			},
			"infra_resource_detail_type_code": schema.StringAttribute{
				Optional: true,
			},
			"exclusion_product_code": schema.StringAttribute{
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

func (m *mongodbProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mongodbProductsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmongodb.GetCloudMongoDbProductListRequest{
		RegionCode:                   &m.config.RegionCode,
		CloudMongoDbImageProductCode: data.CloudMongoDbImageProductCode.ValueStringPointer(),
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetMongoDbProductsList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	mongodbProductResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMongoDbProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetMongoDbProductList response", map[string]any{
		"mongodbProductResponse": common.MarshalUncheckedString(mongodbProductResp),
	})

	mongodbProductList := flattenMongoDbProductLists(ctx, mongodbProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, mongodbProductList)

	data.refreshFromOutput(ctx, fillteredList)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

type mongodbProductsDataSourceModel struct {
	ID                           types.String `tfsdk:"id"`
	CloudMongoDbImageProductCode types.String `tfsdk:"image_product_code"`
	ProductCode                  types.String `tfsdk:"product_code"`
	InfraResourceDetailTypeCode  types.String `tfsdk:"infra_resource_detail_type_code"`
	ExclusionProductCode         types.String `tfsdk:"exclusion_product_code"`
	ProductList                  types.List   `tfsdk:"product_list"`
	Filters                      types.Set    `tfsdk:"filter"`
}

func (m *mongodbProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*mongodbProductModel) {
	m.ID = types.StringValue(time.Now().UTC().String())
	productListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mongodbProductModel{}.attrTypes()}, list)
	m.ProductList = productListValue
}

type mongodbProductModel struct {
	ProductCode        types.String `tfsdk:"product_code"`
	ProductName        types.String `tfsdk:"product_name"`
	ProductType        types.String `tfsdk:"product_type"`
	ProductDescription types.String `tfsdk:"product_description"`
	InfraResourceType  types.String `tfsdk:"infra_resource_type"`
	CpuCount           types.Int64  `tfsdk:"cpu_count"`
	MemorySize         types.Int64  `tfsdk:"memory_size"`
	DiskType           types.String `tfsdk:"disk_type"`
}

func (m mongodbProductModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":        types.StringType,
		"product_name":        types.StringType,
		"product_type":        types.StringType,
		"product_description": types.StringType,
		"infra_resource_type": types.StringType,
		"cpu_count":           types.Int64Type,
		"memory_size":         types.Int64Type,
		"disk_type":           types.StringType,
	}
}

func (m *mongodbProductModel) refreshFromOutput(ctx context.Context, output *vmongodb.Product) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.ProductDescription = types.StringPointerValue(output.ProductDescription)
	m.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	m.CpuCount = types.Int64Value(int64(*output.CpuCount))
	m.MemorySize = types.Int64Value(int64(*output.MemorySize))
	m.DiskType = types.StringPointerValue(output.DiskType.Code)
}

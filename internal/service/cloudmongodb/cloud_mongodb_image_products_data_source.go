package cloudmongodb

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"time"
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
	resp.TypeName = req.ProviderTypeName + "_mongodb_image_product_list"
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

func (m *mongodbImageProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"product_code": schema.StringAttribute{
				Optional: true,
			},
			"exclusion_product_code": schema.StringAttribute{
				Optional: true,
			},
			"generation_code": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"G2", "G3"}...),
				},
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
				},
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mongodbImageProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mongodbImageProductsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmongodb.GetCloudMongoDbImageProductListRequest{
		RegionCode: &m.config.RegionCode,
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetMongoDbImageProductList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	mongodbImageProductResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbImageProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMongoDbImageProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetMongoDbImageProductList response", map[string]any{
		"mongodbImageProductResponse": common.MarshalUncheckedString(mongodbImageProductResp),
	})

	mongodbImageProductList := flattenMongoDbImageProduct(ctx, mongodbImageProductResp.ProductList)

	fillteredList := common.FilterModels(ctx, data.Filters, mongodbImageProductList)

	data.refreshFromOutput(ctx, fillteredList)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	m.ID = types.StringValue(time.Now().UTC().String())
}

type mongodbImageProductsDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ImageProductList     types.List   `tfsdk:"image_product_list"`
	ProductCode          types.String `tfsdk:"product_code"`
	GenerationCode       types.String `tfsdk:"generation_code"`
	ExclusionProductCode types.String `tfsdk:"exclusion_product_code"`
	Filters              types.Set    `tfsdk:"filter"`
}

type mongodbImageProduct struct {
	ProductCode        types.String `tfsdk:"product_code"`
	GenerationCode     types.String `tfsdk:"generation_code"`
	ProductName        types.String `tfsdk:"product_name"`
	ProductType        types.String `tfsdk:"product_type"`
	InfraResourceType  types.String `tfsdk:"infra_resource_type"`
	PlatformType       types.String `tfsdk:"platform_type"`
	OsInformation      types.String `tfsdk:"os_information"`
	ProductDescription types.String `tfsdk:"product_description"`
}

func (m mongodbImageProduct) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":        types.StringType,
		"generation_code":     types.StringType,
		"product_name":        types.StringType,
		"product_type":        types.StringType,
		"infra_resource_type": types.StringType,
		"platform_type":       types.StringType,
		"os_information":      types.StringType,
		"product_description": types.StringType,
	}
}
func (m *mongodbImageProduct) refreshFromOutput(output *vmongodb.Product) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.GenerationCode = types.StringPointerValue(output.GenerationCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	m.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	m.OsInformation = types.StringPointerValue(output.OsInformation)
	m.ProductDescription = types.StringPointerValue(output.ProductDescription)
}

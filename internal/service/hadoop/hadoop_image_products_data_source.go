package hadoop

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &hadoopImageProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopImageProductsDataSource{}
)

func NewHadoopImageProductsDataSource() datasource.DataSource {
	return &hadoopImageProductsDataSource{}
}

type hadoopImageProductsDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopImageProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_image_products"
}

func (h *hadoopImageProductsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (h *hadoopImageProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	h.config = config
}

func (h *hadoopImageProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopImageProductsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopImageProductListRequest{
		RegionCode: &h.config.RegionCode,
	}
	tflog.Info(ctx, "GetHadoopImageProductList reqParams="+common.MarshalUncheckedString(reqParams))

	imageProductResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopImageProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetHadoopimageProductList response="+common.MarshalUncheckedString(imageProductResp))

	if imageProductResp == nil || len(imageProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	imagesProductList := flattenHadoopImageList(imageProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, imagesProductList)
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

func flattenHadoopImageList(list []*vhadoop.Product) []*hadoopImageProduct {
	var outputs []*hadoopImageProduct

	for _, v := range list {
		var output hadoopImageProduct
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

type hadoopImageProductsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductList types.List   `tfsdk:"image_product_list"`
	OutputFile       types.String `tfsdk:"output_file"`
	Filters          types.Set    `tfsdk:"filter"`
}

type hadoopImageProduct struct {
	ProductCode       types.String `tfsdk:"product_code"`
	GenerationCode    types.String `tfsdk:"generation_code"`
	ProductName       types.String `tfsdk:"product_name"`
	ProductType       types.String `tfsdk:"product_type"`
	PlatformType      types.String `tfsdk:"platform_type"`
	OsInformation     types.String `tfsdk:"os_information"`
	EngineVersionCode types.String `tfsdk:"engine_version_code"`
}

func (h hadoopImageProduct) attrTypes() map[string]attr.Type {
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

func (h *hadoopImageProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*hadoopImageProduct) {
	imageProductListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: hadoopImageProduct{}.attrTypes()}, list)
	h.ImageProductList = imageProductListValue
	h.ID = types.StringValue("")
}

func (h *hadoopImageProduct) refreshFromOutput(output *vhadoop.Product) {
	h.ProductCode = types.StringPointerValue(output.ProductCode)
	h.GenerationCode = types.StringPointerValue(output.GenerationCode)
	h.ProductName = types.StringPointerValue(output.ProductName)
	h.ProductType = types.StringPointerValue(output.ProductType.Code)
	h.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	h.OsInformation = types.StringPointerValue(output.OsInformation)
	h.EngineVersionCode = types.StringPointerValue(output.EngineVersionCode)
}

package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &postgresqlImageProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlImageProductsDataSource{}
)

func NewPostgresqlImageProductsDataSource() datasource.DataSource {
	return &postgresqlImageProductsDataSource{}
}

type postgresqlImageProductsDataSource struct {
	config *conn.ProviderConfig
}

func (d *postgresqlImageProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_image_products"
}

func (d *postgresqlImageProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"generation_code": schema.StringAttribute{
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

func (d *postgresqlImageProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.config = config
}

func (d *postgresqlImageProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data postgresqlImageProductsDataSourceModel

	if !d.config.SupportVPC {
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

	reqParams := &vpostgresql.GetCloudPostgresqlImageProductListRequest{
		RegionCode: &d.config.RegionCode,
	}
	tflog.Info(ctx, "GetPostgresqlImageProductList reqParams="+common.MarshalUncheckedString(reqParams))

	postgresqlImageProductResp, err := d.config.Client.Vpostgresql.V2Api.GetCloudPostgresqlImageProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetPostgresqlImageProductList response="+common.MarshalUncheckedString(postgresqlImageProductResp))

	if postgresqlImageProductResp == nil || len(postgresqlImageProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	postgresqlImageProductList := flattenPostgresqlImageProduct(postgresqlImageProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, postgresqlImageProductList)
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

func convertToJsonStruct(images []attr.Value) ([]postgresqlImageProductToJsonConvert, error) {
	var postgresqlImagesToConvert = []postgresqlImageProductToJsonConvert{}

	for _, image := range images {
		imageJasn := postgresqlImageProductToJsonConvert{}
		if err := json.Unmarshal([]byte(image.String()), &imageJasn); err != nil {
			return nil, err
		}
		postgresqlImagesToConvert = append(postgresqlImagesToConvert, imageJasn)
	}

	return postgresqlImagesToConvert, nil
}

func flattenPostgresqlImageProduct(list []*vpostgresql.Product) []*postgresqlImageProduct {
	var outputs []*postgresqlImageProduct

	for _, v := range list {
		var output postgresqlImageProduct
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (d *postgresqlImageProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*postgresqlImageProduct) {
	imageProductListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: postgresqlImageProduct{}.attrTypes()}, list)
	d.ImageProductList = imageProductListValue
	d.ID = types.StringValue("")
}

type postgresqlImageProductsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductList types.List   `tfsdk:"image_product_list"`
	OutputFile       types.String `tfsdk:"output_file"`
	Filters          types.Set    `tfsdk:"filter"`
}

type postgresqlImageProduct struct {
	ProductCode    types.String `tfsdk:"product_code"`
	ProductName    types.String `tfsdk:"product_name"`
	ProductType    types.String `tfsdk:"product_type"`
	PlatformType   types.String `tfsdk:"platform_type"`
	OsInformation  types.String `tfsdk:"os_information"`
	GenerationCode types.String `tfsdk:"generation_code"`
}

type postgresqlImageProductToJsonConvert struct {
	ProductCode    string `json:"product_code"`
	ProductName    string `json:"product_name"`
	ProductType    string `json:"product_type"`
	PlatformType   string `json:"platform_type"`
	OsInformation  string `json:"os_information"`
	GenerationCode string `json:"generation_code"`
}

func (d postgresqlImageProduct) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":    types.StringType,
		"product_name":    types.StringType,
		"product_type":    types.StringType,
		"platform_type":   types.StringType,
		"os_information":  types.StringType,
		"generation_code": types.StringType,
	}
}

func (d *postgresqlImageProduct) refreshFromOutput(output *vpostgresql.Product) {
	d.ProductCode = types.StringPointerValue(output.ProductCode)
	d.ProductName = types.StringPointerValue(output.ProductName)
	d.ProductType = types.StringPointerValue(output.ProductType.Code)
	d.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	d.OsInformation = types.StringPointerValue(output.OsInformation)
	d.GenerationCode = types.StringPointerValue(output.GenerationCode)
}

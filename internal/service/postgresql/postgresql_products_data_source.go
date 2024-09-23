package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	_ datasource.DataSource              = &postgresqlProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlProductsDataSource{}
)

func NewPostgresqlProductsDataSource() datasource.DataSource {
	return &postgresqlProductsDataSource{}
}

type postgresqlProductsDataSource struct {
	config *conn.ProviderConfig
}

func (d *postgresqlProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_products"
}

func (d *postgresqlProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Required: true,
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

func (d *postgresqlProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *postgresqlProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data postgresqlProductList

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

	reqParams := &vpostgresql.GetCloudPostgresqlProductListRequest{
		RegionCode:                      &d.config.RegionCode,
		CloudPostgresqlImageProductCode: data.CloudPostgresqlImageProductCode.ValueStringPointer(),
	}
	tflog.Info(ctx, "GetPostgresqlProductsList reqParams="+common.MarshalUncheckedString(reqParams))

	postgresqlProductResp, err := d.config.Client.Vpostgresql.V2Api.GetCloudPostgresqlProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetPostgresqlProductsList response="+common.MarshalUncheckedString(postgresqlProductResp))

	if postgresqlProductResp == nil || len(postgresqlProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	postgresqlProductList := flattenPostgresqlProduct(postgresqlProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, postgresqlProductList)

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

func convertProductsToJsonStruct(images []attr.Value) ([]postgresqlProductsToJsonConvert, error) {
	var postgresqlProductsToConvert = []postgresqlProductsToJsonConvert{}

	for _, image := range images {
		imageJasn := postgresqlProductsToJsonConvert{}
		if err := json.Unmarshal([]byte(image.String()), &imageJasn); err != nil {
			return nil, err
		}
		postgresqlProductsToConvert = append(postgresqlProductsToConvert, imageJasn)
	}

	return postgresqlProductsToConvert, nil
}

func flattenPostgresqlProduct(list []*vpostgresql.Product) []*postgresqlProductModel {
	var outputs []*postgresqlProductModel

	for _, v := range list {
		var output postgresqlProductModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (d *postgresqlProductList) refreshFromOutput(ctx context.Context, list []*postgresqlProductModel) {
	productListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: postgresqlProductModel{}.attrTypes()}, list)
	d.ProductList = productListValue
	d.ID = types.StringValue(time.Now().UTC().String())
}

type postgresqlProductList struct {
	ID                              types.String `tfsdk:"id"`
	CloudPostgresqlImageProductCode types.String `tfsdk:"image_product_code"`
	ProductList                     types.List   `tfsdk:"product_list"`
	OutputFile                      types.String `tfsdk:"output_file"`
	Filters                         types.Set    `tfsdk:"filter"`
}

type postgresqlProductModel struct {
	ProductCode        types.String `tfsdk:"product_code"`
	ProductName        types.String `tfsdk:"product_name"`
	ProductType        types.String `tfsdk:"product_type"`
	ProductDescription types.String `tfsdk:"product_description"`
	InfraResourceType  types.String `tfsdk:"infra_resource_type"`
	CpuCount           types.Int64  `tfsdk:"cpu_count"`
	MemorySize         types.Int64  `tfsdk:"memory_size"`
	DiskType           types.String `tfsdk:"disk_type"`
}
type postgresqlProductsToJsonConvert struct {
	ProductCode        string `json:"product_code"`
	ProductName        string `json:"product_name"`
	ProductType        string `json:"product_type"`
	ProductDescription string `json:"product_description"`
	InfraResourceType  string `json:"infra_resource_type"`
	CpuCount           int    `json:"cpu_count"`
	MemorySize         int    `json:"memory_size"`
	DiskType           string `json:"disk_type"`
}

func (d postgresqlProductModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":        types.StringType,
		"product_name":        types.StringType,
		"product_type":        types.StringType,
		"infra_resource_type": types.StringType,
		"cpu_count":           types.Int64Type,
		"memory_size":         types.Int64Type,
		"disk_type":           types.StringType,
		"product_description": types.StringType,
	}
}
func (d *postgresqlProductModel) refreshFromOutput(output *vpostgresql.Product) {
	d.ProductCode = types.StringPointerValue(output.ProductCode)
	d.ProductName = types.StringPointerValue(output.ProductName)
	d.ProductType = types.StringPointerValue(output.ProductType.Code)
	d.ProductDescription = types.StringPointerValue(output.ProductDescription)
	d.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	d.CpuCount = common.Int64ValueFromInt32(output.CpuCount)
	d.MemorySize = types.Int64PointerValue(output.MemorySize)
	d.DiskType = types.StringPointerValue(output.DiskType.Code)
}

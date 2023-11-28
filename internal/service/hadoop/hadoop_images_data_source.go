package hadoop

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"time"
)

var (
	_ datasource.DataSource              = &hadoopImagesDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopImagesDataSource{}
)

func NewHadoopImagesDataSource() datasource.DataSource {
	return &hadoopImagesDataSource{}
}

type hadoopImagesDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (h *hadoopImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_images"
}

func (h *hadoopImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"output_file": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
			"images": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"product_name": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
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
						"base_block_storage_size": schema.Int64Attribute{
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
			},
		},
	}
}

func (h *hadoopImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopImagesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopImageProductListRequest{
		RegionCode: &h.config.RegionCode,
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetHadoopimageProductList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	imagesProductResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopImageProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetHadoopimageProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetHadoopimageProductList response", map[string]any{
		"imagesProductResponse": common.MarshalUncheckedString(imagesProductResp),
	})

	imagesProductList := flattenHadoopImageList(imagesProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, imagesProductList)
	if diag := data.refreshFromOutput(ctx, fillteredList); diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if diag := witeHadoopImagesToFile(outputPath, data.Images); diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
	}

	state := data
	state.ID = types.StringValue(time.Now().UTC().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func witeHadoopImagesToFile(path string, images types.List) diag.Diagnostics {
	var hadoopImages []hadoopImagesJson
	var diag diag.Diagnostics

	for _, image := range images.Elements() {
		hadoopImage := hadoopImagesJson{}
		if err := json.Unmarshal([]byte(image.String()), &hadoopImage); err != nil {
			diag.AddError(
				"Unmarshal",
				fmt.Sprintf("error: %s", err.Error()),
			)
			return diag
		}
		hadoopImages = append(hadoopImages, hadoopImage)
	}

	if err := common.WriteToFile(path, hadoopImages); err != nil {
		diag.AddError(
			"WriteToFile",
			fmt.Sprintf("error: %s", err.Error()),
		)
		return diag
	}
	return nil
}

type hadoopImagesDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ProductCode          types.String `tfsdk:"product_code"`
	ExclusionProductCode types.String `tfsdk:"exclusion_product_code"`
	OutputFile           types.String `tfsdk:"output_file"`
	Images               types.List   `tfsdk:"images"`
	Filters              types.Set    `tfsdk:"filter"`
}

type hadoopImage struct {
	ProductName          types.String `tfsdk:"product_name"`
	ProductCode          types.String `tfsdk:"product_code"`
	ProductType          types.String `tfsdk:"product_type"`
	ProductDescription   types.String `tfsdk:"product_description"`
	InfraResourceType    types.String `tfsdk:"infra_resource_type"`
	BaseBlockStorageSize types.Int64  `tfsdk:"base_block_storage_size"`
	PlatformType         types.String `tfsdk:"platform_type"`
	OsInformation        types.String `tfsdk:"os_information"`
	GenerationCode       types.String `tfsdk:"generation_code"`
}

type hadoopImagesJson struct {
	ProductName          string `json:"product_name"`
	ProductCode          string `json:"product_code"`
	ProductType          string `json:"product_type"`
	ProductDescription   string `json:"product_description"`
	InfraResourceType    string `json:"infra_resource_type"`
	BaseBlockStorageSize int64  `json:"base_block_storage_size"`
	PlatformType         string `json:"platform_type"`
	OsInformation        string `json:"os_information"`
	GenerationCode       string `json:"generation_code"`
}

func (i hadoopImage) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_name":            types.StringType,
		"product_code":            types.StringType,
		"product_type":            types.StringType,
		"product_description":     types.StringType,
		"infra_resource_type":     types.StringType,
		"base_block_storage_size": types.Int64Type,
		"platform_type":           types.StringType,
		"os_information":          types.StringType,
		"generation_code":         types.StringType,
	}
}

func (h *hadoopImagesDataSourceModel) refreshFromOutput(ctx context.Context, output []*hadoopImageDataSourceModel) diag.Diagnostics {
	var images []hadoopImage
	var diag diag.Diagnostics

	for _, image := range output {
		images = append(images, hadoopImage{
			ProductName:          image.ProductName,
			ProductCode:          image.ProductCode,
			ProductType:          image.ProductType,
			ProductDescription:   image.ProductDescription,
			InfraResourceType:    image.InfraResourceType,
			BaseBlockStorageSize: image.BaseBlockStorageSize,
			PlatformType:         image.PlatformType,
			OsInformation:        image.OsInformation,
			GenerationCode:       image.GenerationCode,
		})
	}

	h.Images, diag = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: hadoopImage{}.attrTypes()}, images)
	if diag.HasError() {
		return diag
	}
	return nil
}

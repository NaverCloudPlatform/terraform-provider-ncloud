package hadoop

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &hadoopAddOnDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopAddOnDataSource{}
)

func NewHadoopAddOnDataSource() datasource.DataSource {
	return &hadoopAddOnDataSource{}
}

type hadoopAddOnDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopAddOnDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_add_on"
}

func (h *hadoopAddOnDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Required: true,
			},
			"cluster_type_code": schema.StringAttribute{
				Required: true,
			},
			"add_on_list": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (h *hadoopAddOnDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (h *hadoopAddOnDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopAddOnDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopAddOnListRequest{
		RegionCode:                  &h.config.RegionCode,
		CloudHadoopImageProductCode: data.ImageProductCode.ValueStringPointer(),
		CloudHadoopClusterTypeCode:  data.ClusterTypeCode.ValueStringPointer(),
	}
	tflog.Info(ctx, "GetHadoopAddOnList reqParams="+common.MarshalUncheckedString(reqParams))

	addOnResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopAddOnList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetHadoopAddOnList response="+common.MarshalUncheckedString(addOnResp))

	if addOnResp == nil || len(addOnResp.CloudHadoopAddOnList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	data.refreshFromOutput(ctx, addOnResp)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if err := common.WriteStringListToFile(outputPath, data.AddOnList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type hadoopAddOnDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductCode types.String `tfsdk:"image_product_code"`
	ClusterTypeCode  types.String `tfsdk:"cluster_type_code"`
	AddOnList        types.List   `tfsdk:"add_on_list"`
	OutputFile       types.String `tfsdk:"output_file"`
}

func (h *hadoopAddOnDataSourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.GetCloudHadoopAddOnListResponse) {
	var addOnList []string
	for _, addOn := range output.CloudHadoopAddOnList {
		// getCloudHadoopAddOnList API response : The Code and Codename are reversed.
		addOnList = append(addOnList, *addOn.CodeName)
	}

	h.AddOnList, _ = types.ListValueFrom(ctx, types.StringType, addOnList)
	h.ID = types.StringValue("")
}

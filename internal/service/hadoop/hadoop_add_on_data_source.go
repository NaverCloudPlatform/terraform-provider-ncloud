package hadoop

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"time"
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

func (h *hadoopAddOnDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_add_on"
}

func (h *hadoopAddOnDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
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
func (h *hadoopAddOnDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data addOnDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopAddOnListRequest{
		RegionCode:                  &h.config.RegionCode,
		CloudHadoopImageProductCode: data.ImageProductCode.ValueStringPointer(),
		CloudHadoopClusterTypeCode:  data.ClusterTypeCode.ValueStringPointer(),
	}

	tflog.Info(ctx, "GetHadoopAddOnList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	addOnResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopAddOnList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetCloudHadoopAddOnList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetHadoopAddOnList response", map[string]any{
		"hadoopAddOnResponse": common.MarshalUncheckedString(addOnResp),
	})
	data.ID = types.StringValue(time.Now().UTC().String())
	data.refreshFromOutput(ctx, addOnResp)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()
		if err := writeStringListToFile(outputPath, data.AddOnList); err != nil {
			var diags diag.Diagnostics
			diags.AddError(
				"WriteToFile",
				fmt.Sprintf("error: %s", err.Error()),
			)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type addOnDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductCode types.String `tfsdk:"image_product_code"`
	ClusterTypeCode  types.String `tfsdk:"cluster_type_code"`
	AddOnList        types.List   `tfsdk:"add_on_list"`
	OutputFile       types.String `tfsdk:"output_file"`
}

func (m *addOnDataSourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.GetCloudHadoopAddOnListResponse) {
	var addOnList []string
	for _, addOn := range output.CloudHadoopAddOnList {
		addOnList = append(addOnList, *addOn.Code)
	}

	m.AddOnList, _ = types.ListValueFrom(ctx, types.StringType, addOnList)
}

func writeStringListToFile(path string, list types.List) error {
	var dataList []string

	for _, v := range list.Elements() {
		var data string
		if err := json.Unmarshal([]byte(v.String()), &data); err != nil {
			return err
		}
		dataList = append(dataList, data)
	}

	if err := common.WriteToFile(path, dataList); err != nil {
		return err
	}
	return nil
}

package hadoop

import (
	"context"
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
	_ datasource.DataSource              = &hadoopBucketDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopBucketDataSource{}
)

func NewHadoopBucketDataSource() datasource.DataSource {
	return &hadoopBucketDataSource{}
}

type hadoopBucketDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopBucketDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (h *hadoopBucketDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_bucket"
}

func (h *hadoopBucketDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"bucket_list": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
func (h *hadoopBucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BucketDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopBucketListRequest{
		RegionCode: &h.config.RegionCode,
	}

	tflog.Info(ctx, "GetHadoopBucketList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	BucketResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopBucketList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetCloudHadoopBucketList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetHadoopBucketList response", map[string]any{
		"hadoopBucketResponse": common.MarshalUncheckedString(BucketResp),
	})
	data.ID = types.StringValue(time.Now().UTC().String())
	data.refreshFromOutput(ctx, BucketResp)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()
		if err := writeStringListToFile(outputPath, data.BucketList); err != nil {
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

type BucketDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	BucketList types.List   `tfsdk:"bucket_list"`
	OutputFile types.String `tfsdk:"output_file"`
}

func (m *BucketDataSourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.GetCloudHadoopBucketListResponse) {
	var BucketList []string
	for _, bucket := range output.CloudHadoopBucketList {
		BucketList = append(BucketList, *bucket.BucketName)
	}

	m.BucketList, _ = types.ListValueFrom(ctx, types.StringType, BucketList)
}

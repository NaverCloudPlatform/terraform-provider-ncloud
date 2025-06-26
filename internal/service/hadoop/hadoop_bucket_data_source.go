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
	_ datasource.DataSource              = &hadoopBucketDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopBucketDataSource{}
)

func NewHadoopBucketDataSource() datasource.DataSource {
	return &hadoopBucketDataSource{}
}

type hadoopBucketDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopBucketDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_bucket"
}

func (h *hadoopBucketDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
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

func (h *hadoopBucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopBucketDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopBucketListRequest{
		RegionCode: &h.config.RegionCode,
	}
	tflog.Info(ctx, "GetHadoopBucketList reqParams="+common.MarshalUncheckedString(reqParams))

	BucketResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopBucketList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetHadoopBucketList response="+common.MarshalUncheckedString(BucketResp))

	data.refreshFromOutput(ctx, BucketResp)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if err := common.WriteStringListToFile(outputPath, data.BucketList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type hadoopBucketDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	BucketList types.List   `tfsdk:"bucket_list"`
	OutputFile types.String `tfsdk:"output_file"`
}

func (h *hadoopBucketDataSourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.GetCloudHadoopBucketListResponse) {
	var BucketList []string
	for _, bucket := range output.CloudHadoopBucketList {
		BucketList = append(BucketList, *bucket.BucketName)
	}

	h.BucketList, _ = types.ListValueFrom(ctx, types.StringType, BucketList)
	h.ID = types.StringValue("")
}

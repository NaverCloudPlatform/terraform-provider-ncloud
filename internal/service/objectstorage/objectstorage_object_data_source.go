package objectstorage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &objectDataSource{}
	_ datasource.DataSourceWithConfigure = &objectDataSource{}
)

func NewObjectDataSource() datasource.DataSource {
	return &objectDataSource{}
}

type objectDataSource struct {
	config *conn.ProviderConfig
}

func (o *objectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	o.config = config
}

func (o *objectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstorage_object"
}

func (o *objectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data objectDataResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, bucketName, key := ObjectIDParser(data.ObjectID.String())

	output, err := o.config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
		Bucket: ncloud.String(bucketName),
		Key:    ncloud.String(key),
	})
	if err != nil {
		return
	}
	defer output.Body.Close()

	data.ID = types.StringValue(ObjectIDGenerator(strings.ToLower(o.config.RegionCode), bucketName, key))
	data.ContentLength = types.Int64PointerValue(output.ContentLength)
	data.ContentType = types.StringPointerValue(output.ContentType)
	data.LastModified = types.StringValue(output.LastModified.String())

	bodyBytes, err := io.ReadAll(output.Body)
	if err != nil {
		return
	}

	data.Body = types.StringValue(string(bodyBytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (o *objectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"object_id": schema.StringAttribute{
				Required: true,
			},
			"bucket": schema.StringAttribute{
				Computed: true,
			},
			"key": schema.StringAttribute{
				Computed: true,
			},
			"source": schema.StringAttribute{
				Computed: true,
			},
			"content_length": schema.Int64Attribute{
				Computed: true,
			},
			"content_type": schema.StringAttribute{
				Computed: true,
			},
			"last_modified": schema.StringAttribute{
				Computed: true,
			},
			"body": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type objectDataResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ObjectID      types.String `tfsdk:"object_id"`
	Bucket        types.String `tfsdk:"bucket"`
	Key           types.String `tfsdk:"key"`
	Source        types.String `tfsdk:"source"`
	ContentLength types.Int64  `tfsdk:"content_length"`
	ContentType   types.String `tfsdk:"content_type"`
	LastModified  types.String `tfsdk:"last_modified"`
	Body          types.String `tfsdk:"body"`
}

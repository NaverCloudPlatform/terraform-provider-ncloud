package objectstorage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
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

	bucketName, key := ObjectIDParser(data.ObjectID.String())

	output, err := o.config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
		Bucket: ncloud.String(bucketName),
		Key:    ncloud.String(key),
	})
	if err != nil {
		return
	}
	defer output.Body.Close()

	data.ID = types.StringValue(ObjectIDGenerator(bucketName, key))
	data.ContentLength = types.Int64PointerValue(output.ContentLength)
	data.ContentType = types.StringPointerValue(output.ContentType)

	if !types.StringPointerValue(output.AcceptRanges).IsNull() || !types.StringPointerValue(output.AcceptRanges).IsUnknown() {
		data.AcceptRanges = types.StringPointerValue(output.AcceptRanges)
	}

	if !types.StringPointerValue(output.ContentEncoding).IsNull() || !types.StringPointerValue(output.ContentEncoding).IsUnknown() {
		data.ContentEncoding = types.StringPointerValue(output.ContentEncoding)
	}

	if !types.StringPointerValue(output.ContentLanguage).IsNull() || !types.StringPointerValue(output.ContentLanguage).IsUnknown() {
		data.ContentLanguage = types.StringPointerValue(output.ContentLanguage)
	}

	if !types.Int64PointerValue(output.ContentLength).IsNull() || !types.Int64PointerValue(output.ContentLength).IsUnknown() {
		data.ContentLength = types.Int64PointerValue(output.ContentLength)
	}

	if !types.StringPointerValue(output.ContentType).IsNull() || !types.StringPointerValue(output.ContentType).IsUnknown() {
		data.ContentType = types.StringPointerValue(output.ContentType)
	}

	if !types.StringPointerValue(output.ETag).IsNull() || !types.StringPointerValue(output.ETag).IsUnknown() {
		data.ETag = types.StringPointerValue(output.ETag)
	}

	if !types.StringPointerValue(output.Expiration).IsNull() || !types.StringPointerValue(output.Expiration).IsUnknown() {
		data.Expiration = types.StringPointerValue(output.Expiration)
	}

	if !types.Int32PointerValue(output.PartsCount).IsNull() || !types.Int32PointerValue(output.PartsCount).IsUnknown() {
		data.PartsCount = common.Int64ValueFromInt32(output.PartsCount)
	}

	if !types.StringPointerValue(output.VersionId).IsNull() || !types.StringPointerValue(output.VersionId).IsUnknown() {
		data.VersionId = types.StringPointerValue(output.VersionId)
	}

	if !types.StringPointerValue(output.WebsiteRedirectLocation).IsNull() || !types.StringPointerValue(output.WebsiteRedirectLocation).IsUnknown() {
		data.WebsiteRedirectLocation = types.StringPointerValue(output.WebsiteRedirectLocation)
	}

	if output.LastModified != nil {
		data.LastModified = types.StringValue(output.LastModified.Format(time.RFC3339))
	}

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
			"accept_ranges": schema.StringAttribute{
				Computed: true,
			},
			"content_encoding": schema.StringAttribute{
				Computed: true,
			},
			"content_language": schema.StringAttribute{
				Computed: true,
			},
			"content_length": schema.Int64Attribute{
				Computed: true,
			},
			"content_type": schema.StringAttribute{
				Computed: true,
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
			"expiration": schema.StringAttribute{
				Computed: true,
			},
			"parts_count": schema.Int64Attribute{
				Computed: true,
			},
			"version_id": schema.StringAttribute{
				Computed: true,
			},
			"website_redirect_location": schema.StringAttribute{
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
	ID                      types.String `tfsdk:"id"`
	ObjectID                types.String `tfsdk:"object_id"`
	Bucket                  types.String `tfsdk:"bucket"`
	Key                     types.String `tfsdk:"key"`
	Source                  types.String `tfsdk:"source"`
	ContentLength           types.Int64  `tfsdk:"content_length"`
	AcceptRanges            types.String `tfsdk:"accept_ranges"`
	ContentEncoding         types.String `tfsdk:"content_encoding"`
	ContentLanguage         types.String `tfsdk:"content_language"`
	ContentType             types.String `tfsdk:"content_type"`
	ETag                    types.String `tfsdk:"etag"`
	Expiration              types.String `tfsdk:"expiration"`
	LastModified            types.String `tfsdk:"last_modified"`
	PartsCount              types.Int64  `tfsdk:"parts_count"`
	VersionId               types.String `tfsdk:"version_id"`
	WebsiteRedirectLocation types.String `tfsdk:"website_redirect_location"`
	Body                    types.String `tfsdk:"body"`
}

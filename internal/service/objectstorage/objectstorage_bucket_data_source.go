package objectstorage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &bucketDataSource{}
	_ datasource.DataSourceWithConfigure = &bucketDataSource{}
)

func NewBucketDataSource() datasource.DataSource {
	return &bucketDataSource{}
}

type bucketDataSource struct {
	config *conn.ProviderConfig
}

func (b *bucketDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	b.config = config
}

func (b *bucketDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstorage_bucket"
}

func (b *bucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bucketDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := b.config.Client.ObjectStorage.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return
	}

	for _, bucket := range output.Buckets {
		if *bucket.Name == *data.BucketName.ValueStringPointer() {
			_, err := b.config.Client.ObjectStorage.HeadBucket(ctx, &s3.HeadBucketInput{
				Bucket: data.BucketName.ValueStringPointer(),
			})
			if err != nil {
				return
			}

			data.ID = types.StringValue(*bucket.Name)
			data.BucketName = types.StringValue(*bucket.Name)

			break
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (b *bucketDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"bucket_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(3, 15),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[가-힣A-Za-z0-9-]+$`), "Allows only hangeuls, alphabets, numbers, hyphen (-)."),
					),
				},
				Description: "Bucket Name for Object Storage",
			},
		},
	}
}

type bucketDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	BucketName types.String `tfsdk:"bucket_name"`
}

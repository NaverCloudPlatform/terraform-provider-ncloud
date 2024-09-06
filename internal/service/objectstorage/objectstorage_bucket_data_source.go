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
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	data.OwnerID = types.StringValue(*output.Owner.ID)
	data.OwnerDisplayName = types.StringValue(*output.Owner.DisplayName)

	for _, bucket := range output.Buckets {
		if *bucket.Name == *data.BucketName.ValueStringPointer() {
			_, err := b.config.Client.ObjectStorage.HeadBucket(ctx, &s3.HeadBucketInput{
				Bucket: data.BucketName.ValueStringPointer(),
			})
			if err != nil {
				resp.Diagnostics.AddError("READING ERROR", err.Error())
				return
			}

			data.ID = types.StringValue(*bucket.Name)
			data.BucketName = types.StringValue(*bucket.Name)

			if !data.CreationDate.IsNull() || !data.CreationDate.IsUnknown() {
				data.CreationDate = types.StringValue(bucket.CreationDate.String())
			}

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
						stringvalidator.LengthBetween(3, 63),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[a-z0-9][a-z0-9\.-]{1,61}[a-z0-9]$`),
							"Bucket name must be between 3 and 63 characters long, can contain lowercase letters, numbers, periods, and hyphens. It must start and end with a letter or number, and cannot have consecutive periods.",
						),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$|.+)`),
							"Bucket name cannot be formatted as an IP address.",
						),
					),
				},
				Description: "Bucket Name for Object Storage",
			},
			"owner_id": schema.StringAttribute{
				Computed: true,
			},
			"owner_displayname": schema.StringAttribute{
				Computed: true,
			},
			"creation_date": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}

type bucketDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	BucketName       types.String `tfsdk:"bucket_name"`
	OwnerID          types.String `tfsdk:"owner_id"`
	OwnerDisplayName types.String `tfsdk:"owner_displayname"`
	CreationDate     types.String `tfsdk:"creation_date"`
}

package objectstorage

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &objectCopyResource{}
	_ resource.ResourceWithConfigure   = &objectCopyResource{}
	_ resource.ResourceWithImportState = &objectCopyResource{}
)

func NewObjectCopyResource() resource.Resource {
	return &objectCopyResource{}
}

type objectCopyResource struct {
	config *conn.ProviderConfig
}

func (o *objectCopyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Exprected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	o.config = config
}

func (o *objectCopyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan objectCopyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &s3.CopyObjectInput{
		Bucket:     plan.Bucket.ValueStringPointer(),
		CopySource: plan.Source.ValueStringPointer(),
		Key:        plan.Key.ValueStringPointer(),
	}

	if !plan.BucketKeyEnabled.IsNull() && !plan.BucketKeyEnabled.IsUnknown() {
		reqParams.BucketKeyEnabled = plan.BucketKeyEnabled.ValueBoolPointer()
	}

	if !plan.ContentEncoding.IsNull() && !plan.ContentEncoding.IsUnknown() {
		reqParams.ContentEncoding = plan.ContentEncoding.ValueStringPointer()
	}

	if !plan.ContentLanguage.IsNull() && !plan.ContentLanguage.IsUnknown() {
		reqParams.ContentLanguage = plan.ContentLanguage.ValueStringPointer()
	}

	if !plan.ContentType.IsNull() && !plan.ContentType.IsUnknown() {
		reqParams.ContentType = plan.ContentType.ValueStringPointer()
	}

	if !plan.ServerSideEncryption.IsNull() && !plan.ServerSideEncryption.IsUnknown() {
		reqParams.ServerSideEncryption = awsTypes.ServerSideEncryption(*plan.ServerSideEncryption.ValueStringPointer())
	}

	if !plan.WebsiteRedirectLocation.IsNull() && !plan.WebsiteRedirectLocation.IsUnknown() {
		reqParams.WebsiteRedirectLocation = plan.WebsiteRedirectLocation.ValueStringPointer()
	}

	tflog.Info(ctx, "CopyObject reqParams="+common.MarshalUncheckedString(reqParams))

	output, err := o.config.Client.ObjectStorage.CopyObject(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("COPYING ERROR", err.Error())
	}
	if output == nil {
		resp.Diagnostics.AddError("COPYING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "CopyObject response="+common.MarshalUncheckedString(output))

	if err := waitObjectCopied(ctx, o.config, plan.Bucket.ValueString(), plan.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError("COPYING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, o.config)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (o *objectCopyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan objectCopyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &s3.DeleteObjectInput{
		Bucket: plan.Bucket.ValueStringPointer(),
		Key:    plan.Key.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteObject reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := o.config.Client.ObjectStorage.DeleteObject(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "DeleteObject response="+common.MarshalUncheckedString(response))

	if err := waitObjectCopyDeleted(ctx, o.config, plan.Bucket.String(), plan.Key.String()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (o *objectCopyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (o *objectCopyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstorage_object_copy"
}

func (o *objectCopyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan objectCopyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.refreshFromOutput(ctx, o.config)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (o *objectCopyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"bucket": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Description: "Bucket name for object",
			},
			"key": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "(Required) Name of the object once it is in the bucket",
			},
			"source": schema.StringAttribute{
				Required:    true,
				Description: "(Required) Path of the object",
			},
			"accept_ranges": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"bucket_key_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"cache_control": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"checksum_crc32": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"checksum_crc32c": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"checksum_sha1": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"checksum_sha256": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"content_encoding": schema.StringAttribute{
				Optional: true,
			},
			"content_language": schema.StringAttribute{
				Optional: true,
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
				Optional: true,
			},
			"parts_count": schema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
			"sse_customer_key_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"server_side_encryption": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"version_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"website_redirect_location": schema.StringAttribute{
				Optional: true,
			},
			"last_modified": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (o *objectCopyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func waitObjectCopied(ctx context.Context, config *conn.ProviderConfig, bucketName string, key string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{CREATING},
		Target:  []string{CREATED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: &bucketName,
				Key:    &key,
			})
			if output != nil {
				return output, CREATED, nil
			}

			if err != nil {
				return output, CREATING, nil
			}

			return output, CREATING, nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for object (%s) to be upload: %s", key, err)
	}
	return nil
}

func waitObjectCopyDeleted(ctx context.Context, config *conn.ProviderConfig, bucketName, key string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: &bucketName,
				Key:    &key,
			})
			if output != nil {
				return output, DELETING, nil
			}

			if err != nil {
				return output, DELETED, nil
			}

			return output, DELETED, nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for object (%s) to be upload: %s", key, err)
	}
	return nil
}

type objectCopyResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Bucket                  types.String `tfsdk:"bucket"`
	Key                     types.String `tfsdk:"key"`
	Source                  types.String `tfsdk:"source"`
	AcceptRanges            types.String `tfsdk:"accept_ranges"`
	BucketKeyEnabled        types.Bool   `tfsdk:"bucket_key_enabled"`
	CacheControl            types.String `tfsdk:"cache_control"`
	ChecksumCRC32           types.String `tfsdk:"checksum_crc32"`
	ChecksumCRC32C          types.String `tfsdk:"checksum_crc32c"`
	ChecksumSHA1            types.String `tfsdk:"checksum_sha1"`
	ChecksumSHA256          types.String `tfsdk:"checksum_sha256"`
	ContentEncoding         types.String `tfsdk:"content_encoding"`
	ContentLanguage         types.String `tfsdk:"content_language"`
	ContentLength           types.Int64  `tfsdk:"content_length"`
	ContentType             types.String `tfsdk:"content_type"`
	ETag                    types.String `tfsdk:"etag"`
	Expiration              types.String `tfsdk:"expiration"`
	LastModified            types.String `tfsdk:"last_modified"`
	PartsCount              types.Int64  `tfsdk:"parts_count"`
	SSECustomerKeyID        types.String `tfsdk:"sse_customer_key_id"`
	ServerSideEncryption    types.String `tfsdk:"server_side_encryption"`
	VersionId               types.String `tfsdk:"version_id"`
	WebsiteRedirectLocation types.String `tfsdk:"website_redirect_location"`
}

func (o *objectCopyResourceModel) refreshFromOutput(ctx context.Context, config *conn.ProviderConfig) {
	output, err := config.Client.ObjectStorage.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: o.Bucket.ValueStringPointer(),
		Key:    o.Key.ValueStringPointer(),
	})
	if err != nil {
		return
	}

	bucketName, key := TrimForParsing(o.Bucket.String()), TrimForParsing(o.Key.String())

	o.ID = types.StringValue(ObjectIDGenerator(bucketName, key))
	if !types.StringPointerValue(output.AcceptRanges).IsNull() || !types.StringPointerValue(output.AcceptRanges).IsUnknown() {
		o.AcceptRanges = types.StringPointerValue(output.AcceptRanges)
	}

	if !types.BoolPointerValue(output.BucketKeyEnabled).IsNull() || !types.BoolPointerValue(output.BucketKeyEnabled).IsUnknown() {
		o.BucketKeyEnabled = types.BoolPointerValue(output.BucketKeyEnabled)
	}

	if !types.StringPointerValue(output.CacheControl).IsNull() || !types.StringPointerValue(output.CacheControl).IsUnknown() {
		o.CacheControl = types.StringPointerValue(output.CacheControl)
	}

	if !types.StringPointerValue(output.ChecksumCRC32).IsNull() || !types.StringPointerValue(output.ChecksumCRC32).IsUnknown() {
		o.ChecksumCRC32 = types.StringPointerValue(output.ChecksumCRC32)
	}

	if !types.StringPointerValue(output.ChecksumCRC32C).IsNull() || !types.StringPointerValue(output.ChecksumCRC32C).IsUnknown() {
		o.ChecksumCRC32C = types.StringPointerValue(output.ChecksumCRC32C)
	}

	if !types.StringPointerValue(output.ChecksumSHA1).IsNull() || !types.StringPointerValue(output.ChecksumSHA1).IsUnknown() {
		o.ChecksumSHA1 = types.StringPointerValue(output.ChecksumSHA1)
	}

	if !types.StringPointerValue(output.ChecksumSHA256).IsNull() || !types.StringPointerValue(output.ChecksumSHA256).IsUnknown() {
		o.ChecksumSHA256 = types.StringPointerValue(output.ChecksumSHA256)
	}

	if !types.StringPointerValue(output.ContentEncoding).IsNull() || !types.StringPointerValue(output.ContentEncoding).IsUnknown() {
		o.ContentEncoding = types.StringPointerValue(output.ContentEncoding)
	}

	if !types.StringPointerValue(output.ContentLanguage).IsNull() || !types.StringPointerValue(output.ContentLanguage).IsUnknown() {
		o.ContentLanguage = types.StringPointerValue(output.ContentLanguage)
	}

	if !types.Int64PointerValue(output.ContentLength).IsNull() || !types.Int64PointerValue(output.ContentLength).IsUnknown() {
		o.ContentLength = types.Int64PointerValue(output.ContentLength)
	}

	if !types.StringPointerValue(output.ContentType).IsNull() || !types.StringPointerValue(output.ContentType).IsUnknown() {
		o.ContentType = types.StringPointerValue(output.ContentType)
	}

	if !types.StringPointerValue(output.ETag).IsNull() || !types.StringPointerValue(output.ETag).IsUnknown() {
		o.ETag = types.StringPointerValue(output.ETag)
	}

	if !types.StringPointerValue(output.Expiration).IsNull() || !types.StringPointerValue(output.Expiration).IsUnknown() {
		o.Expiration = types.StringPointerValue(output.Expiration)
	}

	if !types.Int32PointerValue(output.PartsCount).IsNull() || !types.Int32PointerValue(output.PartsCount).IsUnknown() {
		o.PartsCount = common.Int64ValueFromInt32(output.PartsCount)
	}

	if !types.StringPointerValue(output.SSEKMSKeyId).IsNull() || !types.StringPointerValue(output.SSEKMSKeyId).IsUnknown() {
		o.SSECustomerKeyID = types.StringPointerValue(output.SSEKMSKeyId)
	}

	if !types.StringPointerValue((*string)(&output.ServerSideEncryption)).IsNull() || !types.StringPointerValue((*string)(&output.ServerSideEncryption)).IsUnknown() {
		o.ServerSideEncryption = types.StringPointerValue((*string)(&output.ServerSideEncryption))
	}

	if !types.StringPointerValue(output.VersionId).IsNull() || !types.StringPointerValue(output.VersionId).IsUnknown() {
		o.VersionId = types.StringPointerValue(output.VersionId)
	}

	if !types.StringPointerValue(output.WebsiteRedirectLocation).IsNull() || !types.StringPointerValue(output.WebsiteRedirectLocation).IsUnknown() {
		o.WebsiteRedirectLocation = types.StringPointerValue(output.WebsiteRedirectLocation)
	}

	if output.LastModified != nil {
		o.LastModified = types.StringValue(output.LastModified.Format(time.RFC3339))
	}
}

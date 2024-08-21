package objectstorage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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

const (
	CREATING = "creating"
	CREATED  = "created"
	DELETING = "deleting"
	DELETED  = "deleted"
)

var (
	_ resource.Resource                = &bucketResource{}
	_ resource.ResourceWithConfigure   = &bucketResource{}
	_ resource.ResourceWithImportState = &bucketResource{}
)

func NewBucketResource() resource.Resource {
	return &bucketResource{}
}

type bucketResource struct {
	config *conn.ProviderConfig
}

func (o *bucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bucketResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &s3.CreateBucketInput{
		Bucket: plan.BucketName.ValueStringPointer(),
	}

	tflog.Info(ctx, "CreateObjectStorage reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := o.config.Client.ObjectStorage.CreateBucket(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "CreateObjectStorage response="+common.MarshalUncheckedString(response))

	output, err := o.config.Client.ObjectStorage.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: plan.BucketName.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", "HeadBucket call failed")
	}

	err = waitBucketCreated(ctx, o.config, plan.BucketName.String())
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(output, o.config, plan.BucketName.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (o *bucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan bucketResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &s3.DeleteBucketInput{
		Bucket: plan.BucketName.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteBucket reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := o.config.Client.ObjectStorage.DeleteBucket(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "DeleteBucket response="+common.MarshalUncheckedString(response))

	if err := waitBucketDeleted(ctx, o.config, plan.BucketName.String()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (o *bucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstorage_bucket"
}

func (o *bucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan bucketResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := o.config.Client.ObjectStorage.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return
	}

	for _, bucket := range output.Buckets {
		if *bucket.Name == *plan.BucketName.ValueStringPointer() {
			_, err := o.config.Client.ObjectStorage.HeadBucket(ctx, &s3.HeadBucketInput{
				Bucket: plan.BucketName.ValueStringPointer(),
			})
			if err != nil {
				return
			}

			plan = bucketResourceModel{
				ID:         types.StringValue(*bucket.Name),
				BucketName: types.StringValue(*bucket.Name),
			}

			break
		}
	}
}

func (o *bucketResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"bucket_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

func (o *bucketResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (o *bucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (o *bucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitBucketCreated(ctx context.Context, config *conn.ProviderConfig, bucketName string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{CREATING},
		Target:  []string{CREATED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.ListBuckets(ctx, &s3.ListBucketsInput{})
			if err != nil {
				return 0, "", fmt.Errorf("ListBuckets is nil")
			}

			for _, bucket := range output.Buckets {
				if TrimForParsing(*bucket.Name) == bucketName {
					return bucket, CREATED, nil
				}
			}

			return output.Buckets, CREATING, nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for object storage (%s) to become terminating: %s", bucketName, err)
	}
	return nil
}

func waitBucketDeleted(ctx context.Context, config *conn.ProviderConfig, bucketName string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.ListBuckets(ctx, &s3.ListBucketsInput{})
			if err != nil {
				return 0, "", fmt.Errorf("ListBuckets is nil")
			}

			for _, bucket := range output.Buckets {
				if "\""+*bucket.Name+"\"" == bucketName {
					return bucket, DELETING, nil
				}
			}

			return output.Buckets, DELETED, nil
		},
		Timeout:    2 * conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for object storage (%s) to become terminating: %s", bucketName, err)
	}

	return nil
}

type bucketResourceModel struct {
	ID         types.String `tfsdk:"id"`
	BucketName types.String `tfsdk:"bucket_name"`
}

func (o *bucketResourceModel) refreshFromOutput(output *s3.HeadBucketOutput, config *conn.ProviderConfig, bucketName string) {
	if output == nil {
		return
	}

	o.ID = types.StringValue(fmt.Sprintf("https://%s.object.ncloudstorage.com/%s", strings.ToLower(config.RegionCode), bucketName))
	o.BucketName = types.StringValue(bucketName)

}

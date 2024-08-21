package objectstorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
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

var (
	_ resource.Resource                = &objectResource{}
	_ resource.ResourceWithConfigure   = &objectResource{}
	_ resource.ResourceWithImportState = &objectResource{}
)

func NewObjectResource() resource.Resource {
	return &objectResource{}
}

type objectResource struct {
	config *conn.ProviderConfig
}

func (o *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan objectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	file, err := os.Open(plan.Source.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", "source path invalid")
		return
	}

	reqParams := &s3.PutObjectInput{
		Bucket: plan.Bucket.ValueStringPointer(),
		Key:    plan.Key.ValueStringPointer(),
		Body:   file,
	}

	tflog.Info(ctx, "PutObject reqParams="+common.MarshalUncheckedString(reqParams))

	output, err := o.config.Client.ObjectStorage.PutObject(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
	}
	if output == nil {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "PutObject response="+common.MarshalUncheckedString(output))

	if err := waitObjectUploaded(ctx, o.config, plan.Bucket.ValueString(), plan.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, o.config)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (o *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan objectResourceModel

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

	if err := waitObjectDeleted(ctx, o.config, plan.Bucket.String(), plan.Key.String()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (o *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstorage_object"
}

func (o *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan objectResourceModel

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

func (o *objectResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
						stringvalidator.LengthBetween(3, 15),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[가-힣A-Za-z0-9-]+$`), "Allows only hangeuls, alphabets, numbers, hyphen (-)."),
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

func (o *objectResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (o *objectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (o *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func waitObjectUploaded(ctx context.Context, config *conn.ProviderConfig, bucketName, key string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{CREATING},
		Target:  []string{CREATED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
				Bucket: ncloud.String(bucketName),
				Key:    ncloud.String(key),
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

func waitObjectDeleted(ctx context.Context, config *conn.ProviderConfig, bucketName, key string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
				Bucket: ncloud.String(bucketName),
				Key:    ncloud.String(key),
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

type objectResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Bucket        types.String `tfsdk:"bucket"`
	Key           types.String `tfsdk:"key"`
	Source        types.String `tfsdk:"source"`
	ContentLength types.Int64  `tfsdk:"content_length"`
	ContentType   types.String `tfsdk:"content_type"`
	LastModified  types.String `tfsdk:"last_modified"`
	Body          types.String `tfsdk:"body"`
}

func (o *objectResourceModel) refreshFromOutput(ctx context.Context, config *conn.ProviderConfig) {
	output, err := config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
		Bucket: o.Bucket.ValueStringPointer(),
		Key:    o.Key.ValueStringPointer(),
	})
	if err != nil {
		return
	}
	defer output.Body.Close()

	bucketName, key := TrimForParsing(o.Bucket.String()), TrimForParsing(o.Key.String())

	o.ID = types.StringValue(fmt.Sprintf("https://%s.object.ncloudstorage.com/%s/%s", strings.ToLower(config.RegionCode), bucketName, key))
	o.ContentLength = types.Int64PointerValue(output.ContentLength)
	o.ContentType = types.StringPointerValue(output.ContentType)
	o.LastModified = types.StringValue(output.LastModified.String())

	bodyBytes, err := io.ReadAll(output.Body)
	if err != nil {
		return
	}

	o.Body = types.StringValue(string(bodyBytes))
}

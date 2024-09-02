package objectstorage

import (
	"context"
	"fmt"
	"io"
	"regexp"
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
				Description: "(Required) Name of the object once it is in the bucket.",
			},
			"source": schema.StringAttribute{
				Required:    true,
				Description: "(Required) Specifies the source object for the copy operation. You specify the value in one of two formats. For objects not accessed through an access point, specify the name of the source bucket and the key of the source object, separated by a slash (/)",
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

func (o *objectCopyResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}

func waitObjectCopied(ctx context.Context, config *conn.ProviderConfig, bucketName string, key string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{CREATING},
		Target:  []string{CREATED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
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
			output, err := config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
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
	ID            types.String `tfsdk:"id"`
	Bucket        types.String `tfsdk:"bucket"`
	Key           types.String `tfsdk:"key"`
	Source        types.String `tfsdk:"source"`
	ContentLength types.Int64  `tfsdk:"content_length"`
	ContentType   types.String `tfsdk:"content_type"`
	LastModified  types.String `tfsdk:"last_modified"`
	Body          types.String `tfsdk:"body"`
}

func (o *objectCopyResourceModel) refreshFromOutput(ctx context.Context, config *conn.ProviderConfig) {
	output, err := config.Client.ObjectStorage.GetObject(ctx, &s3.GetObjectInput{
		Bucket: o.Bucket.ValueStringPointer(),
		Key:    o.Key.ValueStringPointer(),
	})
	if err != nil {
		return
	}

	defer output.Body.Close()

	bucketName, key := TrimForParsing(o.Bucket.String()), TrimForParsing(o.Key.String())

	o.ID = types.StringValue(ObjectIDGenerator(bucketName, key))
	o.ContentLength = types.Int64PointerValue(output.ContentLength)
	o.ContentType = types.StringPointerValue(output.ContentType)
	o.LastModified = types.StringValue(output.LastModified.String())

	bodyBytes, err := io.ReadAll(output.Body)
	if err != nil {
		return
	}

	o.Body = types.StringValue(string(bodyBytes))
}
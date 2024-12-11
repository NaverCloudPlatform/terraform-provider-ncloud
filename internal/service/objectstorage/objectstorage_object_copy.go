package objectstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

	if !plan.ContentEncoding.IsNull() && !plan.ContentEncoding.IsUnknown() {
		reqParams.ContentEncoding = plan.ContentEncoding.ValueStringPointer()
	}

	if !plan.ContentLanguage.IsNull() && !plan.ContentLanguage.IsUnknown() {
		reqParams.ContentLanguage = plan.ContentLanguage.ValueStringPointer()
	}

	if !plan.ContentType.IsNull() && !plan.ContentType.IsUnknown() {
		reqParams.ContentType = plan.ContentType.ValueStringPointer()
	}

	if !plan.WebsiteRedirectLocation.IsNull() && !plan.WebsiteRedirectLocation.IsUnknown() {
		reqParams.WebsiteRedirectLocation = plan.WebsiteRedirectLocation.ValueStringPointer()
	}

	tflog.Info(ctx, "CopyObject reqParams="+common.MarshalUncheckedString(reqParams))

	output, err := o.config.Client.ObjectStorage.CopyObject(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("COPYING ERROR", err.Error())
		return
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

	plan.refreshFromOutput(ctx, o.config, &resp.Diagnostics)
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

	if err := waitObjectCopyDeleted(ctx, o.config, plan.Bucket.ValueString(), plan.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (o *objectCopyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	bucketName, key := ObjectIDParser(req.ID)

	if bucketName == "" || key == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: bucket-name/key Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket"), bucketName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), key)...)
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

	plan.refreshFromOutput(ctx, o.config, &resp.Diagnostics)

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
				Validators:  BucketNameValidator(),
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
				Sensitive:   true,
				Required:    true,
				Description: "(Required) Path of the object",
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
				Optional: true,
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
		},
	}
}

func (o *objectCopyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectCopyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Source.Equal(state.Source) {
		reqParams := &s3.CopyObjectInput{
			Bucket:     state.Bucket.ValueStringPointer(),
			CopySource: plan.Source.ValueStringPointer(),
			Key:        state.Key.ValueStringPointer(),
		}

		// attributes that has dependancies with source
		if !plan.ContentEncoding.IsNull() && !plan.ContentEncoding.IsUnknown() {
			reqParams.ContentEncoding = plan.ContentEncoding.ValueStringPointer()
		}

		if !plan.ContentLanguage.IsNull() && !plan.ContentLanguage.IsUnknown() {
			reqParams.ContentLanguage = plan.ContentLanguage.ValueStringPointer()
		}

		if !plan.ContentType.IsNull() && !plan.ContentType.IsUnknown() {
			reqParams.ContentType = plan.ContentType.ValueStringPointer()
		}

		if !plan.WebsiteRedirectLocation.IsNull() && !plan.WebsiteRedirectLocation.IsUnknown() {
			reqParams.WebsiteRedirectLocation = plan.WebsiteRedirectLocation.ValueStringPointer()
		}

		tflog.Info(ctx, "CopyObject at update operation reqParams="+common.MarshalUncheckedString(reqParams))

		output, err := o.config.Client.ObjectStorage.CopyObject(ctx, reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
			return
		}
		if output == nil {
			resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
			return
		}

		tflog.Info(ctx, "CopyObject at update operation response="+common.MarshalUncheckedString(output))

		if err := waitObjectCopied(ctx, o.config, plan.Bucket.ValueString(), plan.Key.ValueString()); err != nil {
			resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
			return
		}

		plan.refreshFromOutput(ctx, o.config, &resp.Diagnostics)
	}

	if !plan.ContentType.Equal(state.ContentType) {

		if o.config.RegionCode == "JPN" {
			resp.Diagnostics.AddError("UPDATING ERROR", "updating object Content-Type is unavailable in this region")
			return
		}

		getReqParams := &s3.GetObjectInput{
			Bucket: state.Bucket.ValueStringPointer(),
			Key:    state.Key.ValueStringPointer(),
		}

		tflog.Info(ctx, "GetObject at update operation reqParams="+common.MarshalUncheckedString(getReqParams))

		getOutput, err := o.config.Client.ObjectStorage.GetObject(ctx, getReqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
			return
		}
		if getOutput == nil {
			resp.Diagnostics.AddError("UPDATING ERROR", "response invalid at get object")
			return
		}

		tflog.Info(ctx, "GetObject at update operation response="+common.MarshalUncheckedString(getOutput))

		reqParams := &s3.PutObjectInput{
			Bucket: plan.Bucket.ValueStringPointer(),
			Key:    plan.Key.ValueStringPointer(),
			Body:   getOutput.Body,
			// this option is unavailable in JPN region
			ContentType: plan.ContentType.ValueStringPointer(),
		}

		tflog.Info(ctx, "PutObject at update operation reqParams="+common.MarshalUncheckedString(reqParams))

		output, err := o.config.Client.ObjectStorage.PutObject(ctx, reqParams)
		if err != nil {
			resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
			return
		}
		if output == nil {
			resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
			return
		}

		tflog.Info(ctx, "PutObject at update operation response="+common.MarshalUncheckedString(output))

		if err := waitObjectUploaded(ctx, o.config, plan.Bucket.ValueString(), plan.Key.ValueString()); err != nil {
			resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
			return
		}

		plan.refreshFromOutput(ctx, o.config, &resp.Diagnostics)
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
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
	ContentEncoding         types.String `tfsdk:"content_encoding"`
	ContentLanguage         types.String `tfsdk:"content_language"`
	ContentLength           types.Int64  `tfsdk:"content_length"`
	ContentType             types.String `tfsdk:"content_type"`
	ETag                    types.String `tfsdk:"etag"`
	Expiration              types.String `tfsdk:"expiration"`
	LastModified            types.String `tfsdk:"last_modified"`
	PartsCount              types.Int64  `tfsdk:"parts_count"`
	VersionId               types.String `tfsdk:"version_id"`
	WebsiteRedirectLocation types.String `tfsdk:"website_redirect_location"`
}

func (o *objectCopyResourceModel) refreshFromOutput(ctx context.Context, config *conn.ProviderConfig, diag *diag.Diagnostics) {
	output, err := config.Client.ObjectStorage.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: o.Bucket.ValueStringPointer(),
		Key:    o.Key.ValueStringPointer(),
	})
	if err != nil {
		diag.AddError("HeadObject ERROR", err.Error())
		return
	}
	if output == nil {
		diag.AddError("HeadObject ERROR", "invalid output")
		return
	}

	o.ID = types.StringValue(ObjectIDGenerator(o.Bucket.ValueString(), o.Key.ValueString()))
	if !types.StringPointerValue(output.AcceptRanges).IsNull() || !types.StringPointerValue(output.AcceptRanges).IsUnknown() {
		o.AcceptRanges = types.StringPointerValue(output.AcceptRanges)
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

package objectstorage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

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

	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	_ resource.Resource                = &objectACLResource{}
	_ resource.ResourceWithConfigure   = &objectACLResource{}
	_ resource.ResourceWithImportState = &objectACLResource{}
)

func NewObjectACLResource() resource.Resource {
	return &objectACLResource{}
}

type objectACLResource struct {
	config *conn.ProviderConfig
}

func (o *objectACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan objectACLResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName, key := ObjectIDParser(plan.ObjectID.String())

	reqParams := &s3.PutObjectAclInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		ACL:    plan.Rule,
	}

	tflog.Info(ctx, "PutObjectACL reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := o.config.Client.ObjectStorage.PutObjectAcl(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
	}

	tflog.Info(ctx, "PutObjectACL response="+common.MarshalUncheckedString(response))

	if err := waitObjectACLApplied(ctx, o.config, bucketName, key); err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	output, err := o.config.Client.ObjectStorage.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(output)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *objectACLResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (o *objectACLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstorage_object_acl"
}

func (o *objectACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan objectACLResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName, key := ObjectIDParser(plan.ObjectID.String())

	output, err := o.config.Client.ObjectStorage.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (o *objectACLResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"object_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.RegexMatches(regexp.MustCompile(`^https:\/\/.*\.object\.ncloudstorage\.com\/[^\/]+\/[^\/]+\.*$`), "Requires pattern with link of target object"),
					),
				},
				Description: "Target object id",
			},
			"rule": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(awsTypes.ObjectCannedACLPrivate),
						string(awsTypes.ObjectCannedACLPublicRead),
						string(awsTypes.ObjectCannedACLPublicReadWrite),
						string(awsTypes.ObjectCannedACLAuthenticatedRead),
						string(awsTypes.ObjectCannedACLAwsExecRead),
						string(awsTypes.ObjectCannedACLBucketOwnerRead),
						string(awsTypes.ObjectCannedACLBucketOwnerFullControl),
					),
				},
			},
			"grants": schema.StringAttribute{
				Computed: true,
			},
			"owner": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (o *objectACLResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (o *objectACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (o *objectACLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func waitObjectACLApplied(ctx context.Context, config *conn.ProviderConfig, bucketName, key string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"applying"},
		Target:  []string{"applied"},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.GetObjectAcl(ctx, &s3.GetObjectAclInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
			})

			if output != nil {
				return output, "applied", nil
			}

			if err != nil {
				return output, "applying", nil
			}

			return output, "applying", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for object acl (%s) to be applied: %s", key, err)
	}
	return nil
}

// TODO: type casting between awsTypes <=> basetypes
type objectACLResourceModel struct {
	ID       types.String             `tfsdk:"id"` // computed
	ObjectID types.String             `tfsdk:"object_id"`
	Rule     awsTypes.ObjectCannedACL `tfsdk:"rule"`
	Grants   types.String             `tfsdk:"grants"` // computed
	Owner    types.String             `tfsdk:"owner"`  // computed
}

func (o *objectACLResourceModel) refreshFromOutput(output *s3.GetObjectAclOutput) {
	if output == nil {
		return
	}

	if len(output.Grants) != 0 {
		o.Grants = types.StringPointerValue(output.Grants[0].Grantee.DisplayName)
	} else {
		o.Grants = types.StringValue("")
	}
	o.ID = types.StringValue(fmt.Sprintf("bucket_acl_%s", o.ObjectID))
	o.Owner = types.StringValue(*output.Owner.ID)
}

func ObjectIDParser(id string) (bucket string, key string) {
	if id == "" {
		return "", ""
	}

	id = strings.TrimPrefix(id, "\"")
	id = strings.TrimSuffix(id, "\"")

	parts := strings.Split(id, "/")
	if len(parts) < 5 {
		return "", ""
	}

	return parts[3], parts[4]
}

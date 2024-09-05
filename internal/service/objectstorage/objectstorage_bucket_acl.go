package objectstorage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

const (
	APPLYING = "applying"
	APPLIED  = "applied"
)

var (
	_ resource.Resource                = &bucketACLResource{}
	_ resource.ResourceWithConfigure   = &bucketACLResource{}
	_ resource.ResourceWithImportState = &bucketACLResource{}
)

func NewBucketACLResource() resource.Resource {
	return &bucketACLResource{}
}

type bucketACLResource struct {
	config *conn.ProviderConfig
}

func (b *bucketACLResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
						stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9][a-z0-9\.-]{1,61}[a-z0-9]$`), "Requires pattern with link of target bucket"),
					),
				},
				Description: "Target bucket name",
			},
			"rule": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(awsTypes.BucketCannedACLPrivate),
						string(awsTypes.BucketCannedACLPublicRead),
						string(awsTypes.BucketCannedACLPublicReadWrite),
						string(awsTypes.BucketCannedACLAuthenticatedRead),
					),
				},
			},
			"grants": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"grantee": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
								},
								"display_name": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"email_address": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"id": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"uri": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
							},
						},
						"permission": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"owner": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (b *bucketACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bucketACLResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := TrimForParsing(plan.BucketName.String())

	reqParams := &s3.PutBucketAclInput{
		Bucket: ncloud.String(bucketName),
		ACL:    plan.Rule,
	}

	tflog.Info(ctx, "PutBucketACL reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := b.config.Client.ObjectStorage.PutBucketAcl(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
	}

	tflog.Info(ctx, "PutObjectACL response="+common.MarshalUncheckedString(response))

	if err := waitBucketACLApplied(ctx, b.config, bucketName); err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	output, err := b.config.Client.ObjectStorage.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: ncloud.String(bucketName),
	})
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (b *bucketACLResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

func (b *bucketACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan bucketACLResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := TrimForParsing(plan.BucketName.String())

	output, err := b.config.Client.ObjectStorage.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: ncloud.String(bucketName),
	})
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (b *bucketACLResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (b *bucketACLResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstorage_bucket_acl"
}

func (b *bucketACLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	b.config = config
}

func (b *bucketACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitBucketACLApplied(ctx context.Context, config *conn.ProviderConfig, bucketName string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{APPLYING},
		Target:  []string{APPLIED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.GetBucketAcl(ctx, &s3.GetBucketAclInput{
				Bucket: ncloud.String(bucketName),
			})

			if output != nil {
				return output, APPLIED, nil
			}

			if err != nil {
				return output, APPLYING, nil
			}

			return output, APPLYING, nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for bucket acl (%s) to be applied: %s", bucketName, err)
	}
	return nil
}

type bucketACLResourceModel struct {
	ID         types.String             `tfsdk:"id"`
	BucketName types.String             `tfsdk:"bucket_name"`
	Rule       awsTypes.BucketCannedACL `tfsdk:"rule"`
	Grants     types.List               `tfsdk:"grants"`
	Owner      types.String             `tfsdk:"owner"`
}

func (b *bucketACLResourceModel) refreshFromOutput(ctx context.Context, output *s3.GetBucketAclOutput) {
	if output == nil {
		return
	}
	var grantList []awsTypes.Grant
	for _, grant := range output.Grants {
		var indivGrant awsTypes.Grant

		indivGrant.Grantee = &awsTypes.Grantee{}
		indivGrant.Grantee.Type = grant.Grantee.Type
		indivGrant.Permission = grant.Permission

		if !types.StringPointerValue(grant.Grantee.ID).IsNull() {
			indivGrant.Grantee.ID = grant.Grantee.ID
		}

		if !types.StringPointerValue(grant.Grantee.DisplayName).IsNull() {
			indivGrant.Grantee.DisplayName = grant.Grantee.DisplayName
		}

		if !types.StringPointerValue(grant.Grantee.EmailAddress).IsNull() {
			indivGrant.Grantee.EmailAddress = grant.Grantee.EmailAddress
		}

		if !types.StringPointerValue(grant.Grantee.URI).IsNull() {
			indivGrant.Grantee.URI = grant.Grantee.URI
		}

		grantList = append(grantList, indivGrant)
	}

	listValueFromGrants, diag := convertGrantsToListValueAtBucket(ctx, grantList)

	if diag.HasError() {
		fmt.Printf("Error Occured with parsing Grants to ListValue")
		return
	}

	b.Grants = listValueFromGrants
	b.ID = types.StringValue(fmt.Sprintf("bucket_acl_%s", b.BucketName))
	b.Owner = types.StringValue(*output.Owner.ID)
}

func convertGrantsToListValueAtBucket(ctx context.Context, grants []awsTypes.Grant) (basetypes.ListValue, diag.Diagnostics) {
	var grantValues []attr.Value

	for _, grant := range grants {
		granteeMap := map[string]attr.Value{
			"type":          types.StringValue(string(grant.Grantee.Type)),
			"display_name":  types.StringPointerValue(grant.Grantee.DisplayName),
			"email_address": types.StringPointerValue(grant.Grantee.EmailAddress),
			"id":            types.StringPointerValue(grant.Grantee.ID),
			"uri":           types.StringPointerValue(grant.Grantee.URI),
		}

		granteeObj, diags := types.ObjectValue(map[string]attr.Type{
			"type":          types.StringType,
			"display_name":  types.StringType,
			"email_address": types.StringType,
			"id":            types.StringType,
			"uri":           types.StringType,
		}, granteeMap)
		if diags.HasError() {
			return basetypes.ListValue{}, diags
		}

		grantMap := map[string]attr.Value{
			"grantee":    granteeObj,
			"permission": types.StringValue(string(grant.Permission)),
		}

		grantObj, diags := types.ObjectValue(map[string]attr.Type{
			"grantee":    granteeObj.Type(ctx),
			"permission": types.StringType,
		}, grantMap)
		if diags.HasError() {
			return basetypes.ListValue{}, diags
		}

		grantValues = append(grantValues, grantObj)
	}

	return types.ListValue(types.ObjectType{AttrTypes: map[string]attr.Type{
		"grantee": types.ObjectType{AttrTypes: map[string]attr.Type{
			"type":          types.StringType,
			"display_name":  types.StringType,
			"email_address": types.StringType,
			"id":            types.StringType,
			"uri":           types.StringType,
		}},
		"permission": types.StringType,
	}}, grantValues)
}

func TrimForParsing(s string) string {
	s = strings.TrimSuffix(s, "\"")
	s = strings.TrimPrefix(s, "\"")

	return s
}

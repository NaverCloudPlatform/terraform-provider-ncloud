package objectstorage

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
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

	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
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

	_, _, bucketName, key := ObjectIDParser(plan.ObjectID.String())

	reqParams := &s3.PutObjectAclInput{
		Bucket: ncloud.String(bucketName),
		Key:    ncloud.String(key),
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
		Bucket: ncloud.String(bucketName),
		Key:    ncloud.String(key),
	})
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)
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

	_, _, bucketName, key := ObjectIDParser(plan.ObjectID.String())

	output, err := o.config.Client.ObjectStorage.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: ncloud.String(bucketName),
		Key:    ncloud.String(key),
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
						stringvalidator.RegexMatches(regexp.MustCompile(`^objectstorage_object::[a-z]{2}::[a-z0-9-_]+::[a-zA-Z0-9_.-]+$`), "Requires pattern with link of target object"),
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
						string(awsTypes.ObjectCannedACLBucketOwnerRead),
						string(awsTypes.ObjectCannedACLBucketOwnerFullControl),
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
			"owner_id": schema.StringAttribute{
				Computed: true,
			},
			"owner_displayname": schema.StringAttribute{
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
		Pending: []string{APPLYING},
		Target:  []string{APPLIED},
		Refresh: func() (interface{}, string, error) {
			output, err := config.Client.ObjectStorage.GetObjectAcl(ctx, &s3.GetObjectAclInput{
				Bucket: ncloud.String(bucketName),
				Key:    ncloud.String(key),
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
		return fmt.Errorf("error waiting for object acl (%s) to be applied: %s", key, err)
	}
	return nil
}

type objectACLResourceModel struct {
	ID               types.String             `tfsdk:"id"`
	ObjectID         types.String             `tfsdk:"object_id"`
	Rule             awsTypes.ObjectCannedACL `tfsdk:"rule"`
	Grants           types.List               `tfsdk:"grants"`
	OwnerID          types.String             `tfsdk:"owner_id"`
	OwnerDisplayName types.String             `tfsdk:"owner_displayname"`
}

func (o *objectACLResourceModel) refreshFromOutput(ctx context.Context, output *s3.GetObjectAclOutput) {
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

	listValueWithGrants, diag := convertGrantsToListValueAtObject(ctx, grantList)
	if diag.HasError() {
		fmt.Printf("Error Occured with parsing Grants to ListValue")
		return
	}

	o.Grants = listValueWithGrants
	o.ID = types.StringValue(fmt.Sprintf("bucket_acl_%s", o.ObjectID))
	o.OwnerID = types.StringValue(*output.Owner.ID)
	o.OwnerDisplayName = types.StringValue(*output.Owner.DisplayName)
}

func convertGrantsToListValueAtObject(ctx context.Context, grants []awsTypes.Grant) (basetypes.ListValue, diag.Diagnostics) {
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

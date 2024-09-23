package postgresql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifybool"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifyint64"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifystring"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	CREATING = "creating"
	SETTING  = "settingUp"
	RUNNING  = "running"
	DELETING = "deleting"
	DELETED  = "deleted"
)

var (
	_ resource.Resource                = &postgresqlResource{}
	_ resource.ResourceWithConfigure   = &postgresqlResource{}
	_ resource.ResourceWithImportState = &postgresqlResource{}
)

func NewPostgresqlResource() resource.Resource {
	return &postgresqlResource{}
}

type postgresqlResource struct {
	config *conn.ProviderConfig
}

func (r *postgresqlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *postgresqlResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

func (r *postgresqlResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql"
}

func (r *postgresqlResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 30),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[가-힣A-Za-z0-9-]+$`),
						"Composed of a alphabets, hangeuls, numbers, hyphen (-).",
					),
				},
			},
			"server_name_prefix": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`),
						"Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number.",
					),
				},
			},
			"id": framework.IDAttribute(),
			"postgresql_instance_no": schema.StringAttribute{
				Computed: true,
			},
			"database_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 30),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9_]+$`),
						"Composed of lowercase alphabets, numbers, underbar (_). Must start with an alphabetic character.",
					),
				},
			},
			"vpc_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secondary_subnet_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						verifystring.RequiresIfTrue(path.Expressions{
							path.MatchRoot("is_ha"),
						}...),
					),
				},
			},
			"user_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 16),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9_]+$`),
						"Composed of lowercase alphabets, numbers, underbar (_). Must start with an alphabetic character.",
					),
				},
			},
			"user_password": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(8, 20),
						stringvalidator.RegexMatches(regexp.MustCompile(`[a-zA-Z]+`), "Must have at least one alphabet"),
						stringvalidator.RegexMatches(regexp.MustCompile(`\d+`), "Must have at least one number"),
						stringvalidator.RegexMatches(regexp.MustCompile(`[~!@#$%^*()\-_=\[\]\{\};:,.<>?]+`), "Must have at least one special character"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[^&+\\"'/\s`+"`"+`]*$`), "Must not have ` & + \\ \" ' / and white space."),
						verifystring.NotContain(path.MatchRoot("user_name").String()),
					),
				},
				Sensitive: true,
			},
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_ha": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(true),
			},
			"is_multi_zone": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					verifybool.RequiresIfTrue(path.Expressions{
						path.MatchRoot("is_ha"),
					}...),
				},
				Description: "default: false",
			},
			"is_storage_encryption": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "default: false",
			},
			"is_backup": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(true),
			},
			"backup_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					verifystring.RequiresIfTrue(path.Expressions{
						path.MatchRoot("is_backup"),
					}...),
				},
				Description: "ex) 01:15",
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.All(
						int64validator.Between(1, 30),
						verifyint64.RequiresIfTrue(path.Expressions{
							path.MatchRoot("is_backup"),
						}...),
					),
				},
			},
			"backup_file_storage_count": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.All(
						int64validator.Between(1, 30),
						verifyint64.RequiresIfTrue(path.Expressions{
							path.MatchRoot("is_backup"),
						}...),
					),
				},
			},
			"is_backup_file_compression": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					verifybool.RequiresIfTrue(path.Expressions{
						path.MatchRoot("is_backup"),
					}...),
				},
				Description: "default: true",
			},
			"is_automatic_backup": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					verifybool.RequiresIfTrue(path.Expressions{
						path.MatchRoot("is_backup"),
					}...),
				},
				Description: "default: true",
			},
			"port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(5432),
					),
				},
				Description: "default: 5432",
			},
			"client_cidr": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: verify.CidrBlockValidator(),
			},
			"data_storage_type_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SSD", "HDD"}...),
				},
				Description: "default: SSD",
			},
			"engine_version": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"postgresql_config_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"postgresql_server_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"server_instance_no": schema.StringAttribute{
							Computed: true,
						},
						"server_name": schema.StringAttribute{
							Computed: true,
						},
						"server_role": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"is_public_subnet": schema.BoolAttribute{
							Computed: true,
						},
						"public_domain": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
							Computed: true,
						},
						"private_ip": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"used_data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"uptime": schema.StringAttribute{
							Computed: true,
						},
						"create_date": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *postgresqlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan postgresqlResourceModel

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"resource does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, err := vpc.GetSubnetInstance(r.config, plan.SubnetNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"CREATING ERROR",
			err.Error(),
		)
	}

	reqParams := &vpostgresql.CreateCloudPostgresqlInstanceRequest{
		RegionCode:                      &r.config.RegionCode,
		CloudPostgresqlServiceName:      plan.ServiceName.ValueStringPointer(),
		CloudPostgresqlServerNamePrefix: plan.ServerNamePrefix.ValueStringPointer(),
		CloudPostgresqlUserName:         plan.UserName.ValueStringPointer(),
		CloudPostgresqlUserPassword:     plan.UserPassword.ValueStringPointer(),
		CloudPostgresqlDatabaseName:     plan.DatabaseName.ValueStringPointer(),
		ClientCidr:                      plan.ClientCidr.ValueStringPointer(),
		VpcNo:                           subnet.VpcNo,
		SubnetNo:                        subnet.SubnetNo,
	}
	plan.VpcNo = types.StringPointerValue(subnet.VpcNo)

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		reqParams.CloudPostgresqlPort = ncloud.Int32(int32(plan.Port.ValueInt64()))
	}

	if !plan.ProductCode.IsNull() {
		reqParams.CloudPostgresqlProductCode = plan.ProductCode.ValueStringPointer()
	}

	if !plan.ImageProductCode.IsNull() {
		reqParams.CloudPostgresqlImageProductCode = plan.ImageProductCode.ValueStringPointer()
	}

	if !plan.DataStorageTypeCode.IsNull() {
		reqParams.DataStorageTypeCode = plan.DataStorageTypeCode.ValueStringPointer()
	}

	if !plan.IsStorageEncryption.IsNull() && !plan.IsStorageEncryption.IsUnknown() {
		reqParams.IsStorageEncryption = plan.IsStorageEncryption.ValueBoolPointer()
	}

	if !plan.IsHa.IsNull() && !plan.IsHa.IsUnknown() {
		reqParams.IsHa = plan.IsHa.ValueBoolPointer()
		if plan.IsHa.ValueBool() {
			if !plan.IsMultiZone.IsNull() && !plan.IsMultiZone.IsUnknown() {
				reqParams.IsMultiZone = plan.IsMultiZone.ValueBoolPointer()
			}
			if !plan.IsBackup.IsNull() && !plan.IsBackup.IsUnknown() && !plan.IsBackup.ValueBool() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is true, `is_backup` must be true or not be input",
				)
				return
			}

		} else {
			if !plan.IsMultiZone.IsNull() && !plan.IsMultiZone.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is false, `is_multi_zone` parameter is not used",
				)
				return
			}
			if !plan.IsBackup.IsNull() && !plan.IsBackup.IsUnknown() {
				reqParams.IsBackup = plan.IsBackup.ValueBoolPointer()
			}
			if !plan.SecondarySubnetNo.IsNull() && !plan.SecondarySubnetNo.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is false, `secondary_subnet_no` is not used",
				)
				return
			}
		}
	}

	if plan.IsMultiZone.ValueBool() {
		if plan.SecondarySubnetNo.IsNull() || plan.SecondarySubnetNo.IsUnknown() {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_multi_zone` is true, `secondary_subnet_no` must be entered",
			)
			return
		}
		reqParams.SecondarySubnetNo = plan.SecondarySubnetNo.ValueStringPointer()
	} else if !plan.SecondarySubnetNo.IsNull() && !plan.SecondarySubnetNo.IsUnknown() {
		resp.Diagnostics.AddError(
			"CREATING ERROR",
			"when `is_multi_zone` is false, `secondary_subnet_no` is not used",
		)
		return
	}

	if !plan.BackupFileRetentionPeriod.IsNull() && !plan.BackupFileRetentionPeriod.IsUnknown() {
		reqParams.BackupFileRetentionPeriod = ncloud.Int32(int32(plan.BackupFileRetentionPeriod.ValueInt64()))
	}

	if !plan.IsAutomaticBackup.IsNull() && !plan.IsAutomaticBackup.IsUnknown() {
		reqParams.IsAutomaticBackup = plan.IsAutomaticBackup.ValueBoolPointer()
	}

	if !plan.IsBackupFileCompression.IsNull() && !plan.IsBackupFileCompression.IsUnknown() {
		reqParams.IsBackupFileCompression = plan.IsBackupFileCompression.ValueBoolPointer()
	}

	if !plan.BackupFileStorageCount.IsNull() && !plan.BackupFileStorageCount.IsUnknown() {
		reqParams.BackupFileStorageCount = ncloud.Int32(int32(plan.BackupFileStorageCount.ValueInt64()))
	}

	if reqParams.IsBackup == nil || *reqParams.IsBackup {
		if reqParams.IsAutomaticBackup == nil || *reqParams.IsAutomaticBackup {
			if !plan.BackupTime.IsNull() && !plan.BackupTime.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_backup` is true and `is_automactic_backup` is true, `backup_time` is not used",
				)
				return
			}
		} else {
			if plan.BackupTime.IsNull() || plan.BackupTime.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_backup` is true and `is_automactic_backup` is false, `backup_time` must be entered",
				)
				return
			}
			reqParams.BackupTime = plan.BackupTime.ValueStringPointer()
		}
	} else {
		backupTimeHasValue := !plan.BackupTime.IsNull() && !plan.BackupTime.IsUnknown()
		if reqParams.IsAutomaticBackup != nil || backupTimeHasValue {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"`is_automatic_backup` or `backup_time` should not be specified when `is_backup` has enabled",
			)
			return
		}
	}

	tflog.Info(ctx, "CreatePostgresql reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.CreateCloudPostgresqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreatePostgresql response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudPostgresqlInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	postgresqlIns := response.CloudPostgresqlInstanceList[0]
	plan.ID = types.StringPointerValue(postgresqlIns.CloudPostgresqlInstanceNo)

	output, err := WaitPostgresqlCreation(ctx, r.config, *postgresqlIns.CloudPostgresqlInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *postgresqlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state postgresqlResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetPostgresqlInstance(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *postgresqlResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *postgresqlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postgresqlResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpostgresql.DeleteCloudPostgresqlInstanceRequest{
		RegionCode:                &r.config.RegionCode,
		CloudPostgresqlInstanceNo: state.ID.ValueStringPointer(),
	}
	tflog.Info(ctx, "DeletePostgresql reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.DeleteCloudPostgresqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeletePostgresql response="+common.MarshalUncheckedString(response))

	if err := waitPostgresqlDeletion(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func GetPostgresqlInstance(ctx context.Context, config *conn.ProviderConfig, no string) (*vpostgresql.CloudPostgresqlInstance, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlInstanceDetailRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(no),
	}
	tflog.Info(ctx, "GetPostgresqlDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vpostgresql.V2Api.GetCloudPostgresqlInstanceDetail(reqParams)
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) {
		return nil, err
	}
	tflog.Info(ctx, "GetPostgresqlDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudPostgresqlInstanceList) < 1 || len(resp.CloudPostgresqlInstanceList[0].CloudPostgresqlServerInstanceList) < 1 {
		return nil, nil
	}

	return resp.CloudPostgresqlInstanceList[0], nil
}

func WaitPostgresqlCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vpostgresql.CloudPostgresqlInstance, error) {
	var postgresqlInstance *vpostgresql.CloudPostgresqlInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{CREATING, SETTING},
		Target:  []string{RUNNING},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetPostgresqlInstance(ctx, config, id)
			postgresqlInstance = instance
			if err != nil {
				return 0, "", err
			}

			status := instance.CloudPostgresqlInstanceStatus.Code
			op := instance.CloudPostgresqlInstanceOperation.Code

			if *status == "INIT" && *op == "CREAT" {
				return instance, CREATING, nil
			}

			if *status == "CREAT" && *op == "SETUP" {
				return instance, SETTING, nil
			}

			if *status == "CREAT" && *op == "NULL" {
				return instance, RUNNING, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create postgresql")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for PostgresqlInstance state to be \"CREAT\": %s", err)
	}

	return postgresqlInstance, nil
}

func waitPostgresqlDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetPostgresqlInstance(ctx, config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, DELETED, nil
			}

			status := instance.CloudPostgresqlInstanceStatusName
			if *status == DELETING {
				return instance, DELETING, nil
			}

			if *status == DELETED {
				return instance, DELETING, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete postgresql")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for postgresql (%s) to become termintaing: %s", id, err)
	}

	return nil
}

type postgresqlResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	PostgresqlInstanceNo      types.String `tfsdk:"postgresql_instance_no"`
	ServiceName               types.String `tfsdk:"service_name"`
	ServerNamePrefix          types.String `tfsdk:"server_name_prefix"`
	DatabaseName              types.String `tfsdk:"database_name"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	SecondarySubnetNo         types.String `tfsdk:"secondary_subnet_no"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	IsMultiZone               types.Bool   `tfsdk:"is_multi_zone"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsStorageEncryption       types.Bool   `tfsdk:"is_storage_encryption"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupTime                types.String `tfsdk:"backup_time"`
	BackupFileStorageCount    types.Int64  `tfsdk:"backup_file_storage_count"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	IsBackupFileCompression   types.Bool   `tfsdk:"is_backup_file_compression"`
	IsAutomaticBackup         types.Bool   `tfsdk:"is_automatic_backup"`
	Port                      types.Int64  `tfsdk:"port"`
	ClientCidr                types.String `tfsdk:"client_cidr"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type_code"`
	EngineVersion             types.String `tfsdk:"engine_version"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	PostgresqlConfigList      types.List   `tfsdk:"postgresql_config_list"`
	PostgresqlServerList      types.List   `tfsdk:"postgresql_server_list"`
}

type postgresqlServer struct {
	ServerInstanceNo    types.String `tfsdk:"server_instance_no"`
	ServerName          types.String `tfsdk:"server_name"`
	ServerRole          types.String `tfsdk:"server_role"`
	IsPublicSubnet      types.Bool   `tfsdk:"is_public_subnet"`
	PublicDomain        types.String `tfsdk:"public_domain"`
	PrivateDomain       types.String `tfsdk:"private_domain"`
	PrivateIp           types.String `tfsdk:"private_ip"`
	ProductCode         types.String `tfsdk:"product_code"`
	DataStorageSize     types.Int64  `tfsdk:"data_storage_size"`
	UsedDataStorageSize types.Int64  `tfsdk:"used_data_storage_size"`
	MemorySize          types.Int64  `tfsdk:"memory_size"`
	CpuCount            types.Int64  `tfsdk:"cpu_count"`
	Uptime              types.String `tfsdk:"uptime"`
	CreateDate          types.String `tfsdk:"create_date"`
}

func (r postgresqlServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_instance_no":     types.StringType,
		"server_name":            types.StringType,
		"server_role":            types.StringType,
		"is_public_subnet":       types.BoolType,
		"public_domain":          types.StringType,
		"private_domain":         types.StringType,
		"private_ip":             types.StringType,
		"product_code":           types.StringType,
		"data_storage_size":      types.Int64Type,
		"used_data_storage_size": types.Int64Type,
		"memory_size":            types.Int64Type,
		"cpu_count":              types.Int64Type,
		"uptime":                 types.StringType,
		"create_date":            types.StringType,
	}
}

func (r *postgresqlResourceModel) refreshFromOutput(ctx context.Context, output *vpostgresql.CloudPostgresqlInstance) {
	r.ID = types.StringPointerValue(output.CloudPostgresqlInstanceNo)
	r.PostgresqlInstanceNo = types.StringPointerValue(output.CloudPostgresqlInstanceNo)
	r.ServiceName = types.StringPointerValue(output.CloudPostgresqlServiceName)
	r.ImageProductCode = types.StringPointerValue(output.CloudPostgresqlImageProductCode)
	r.VpcNo = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].VpcNo)
	r.SubnetNo = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].SubnetNo)
	r.IsMultiZone = types.BoolPointerValue(output.IsMultiZone)
	r.IsHa = types.BoolPointerValue(output.IsHa)
	r.IsBackup = types.BoolPointerValue(output.IsBackup)
	r.Port = common.Int64ValueFromInt32(output.CloudPostgresqlPort)
	r.BackupTime = types.StringPointerValue(output.BackupTime)
	r.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	r.IsStorageEncryption = types.BoolPointerValue(output.CloudPostgresqlServerInstanceList[0].IsStorageEncryption)
	r.EngineVersion = types.StringPointerValue(output.EngineVersion)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	r.AccessControlGroupNoList = acgList
	configList, _ := types.ListValueFrom(ctx, types.StringType, output.CloudPostgresqlConfigList)
	r.PostgresqlConfigList = configList

	var serverList []postgresqlServer
	for _, server := range output.CloudPostgresqlServerInstanceList {
		postgresqlServerInstance := postgresqlServer{
			ServerInstanceNo: types.StringPointerValue(server.CloudPostgresqlServerInstanceNo),
			ServerName:       types.StringPointerValue(server.CloudPostgresqlServerName),
			ServerRole:       types.StringPointerValue(server.CloudPostgresqlServerRole.Code),
			IsPublicSubnet:   types.BoolPointerValue(server.IsPublicSubnet),
			PublicDomain:     types.StringPointerValue(server.PublicDomain),
			PrivateIp:        types.StringPointerValue(server.PrivateIp),
			DataStorageSize:  types.Int64PointerValue(server.DataStorageSize),
			ProductCode:      types.StringPointerValue(server.CloudPostgresqlProductCode),
			MemorySize:       types.Int64PointerValue(server.MemorySize),
			CpuCount:         common.Int64ValueFromInt32(server.CpuCount),
			Uptime:           types.StringPointerValue(server.Uptime),
			CreateDate:       types.StringPointerValue(server.CreateDate),
		}

		if server.PublicDomain != nil {
			postgresqlServerInstance.PublicDomain = types.StringPointerValue(server.PublicDomain)
		}

		if server.UsedDataStorageSize != nil {
			postgresqlServerInstance.UsedDataStorageSize = types.Int64PointerValue(server.UsedDataStorageSize)
		}
		serverList = append(serverList, postgresqlServerInstance)
	}

	postgresqlServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: postgresqlServer{}.attrTypes()}, serverList)

	r.PostgresqlServerList = postgresqlServers
}

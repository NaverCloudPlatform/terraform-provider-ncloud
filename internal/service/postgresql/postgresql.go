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
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifybool"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifyint64"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify/verifystring"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
			"id": framework.IDAttribute(),
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
			"client_cidr": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: verify.CidrBlockValidator(),
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
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"engine_version_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"data_storage_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SSD", "HDD"}...),
				},
				Description: "default: SSD",
			},
			// Available only `pub` and `fin` site. But GOV response message have both values.
			"storage_encryption": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "default: false",
			},
			"ha": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(true),
			},
			// Available only `pub` and `fin` site. But GOV response message have both values.
			"multi_zone": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Bool{
					verifybool.RequiresIfTrue(path.Expressions{
						path.MatchRoot("ha"),
					}...),
				},
				Description: "default: false",
			},
			// Available only `pub` and `fin` site. All sites response message have neither values.
			"secondary_subnet_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						verifystring.RequiresIfTrue(path.Expressions{
							path.MatchRoot("ha"),
						}...),
						verifystring.RequiresIfTrue(path.Expressions{
							path.MatchRoot("multi_zone"),
						}...),
					),
				},
			},
			"backup": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(true),
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
							path.MatchRoot("backup"),
						}...),
					),
				},
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
						path.MatchRoot("backup"),
					}...),
				},
				Description: "ex) 01:15",
			},
			"backup_file_storage_count": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.All(
						int64validator.Between(1, 30),
						verifyint64.RequiresIfTrue(path.Expressions{
							path.MatchRoot("backup"),
						}...),
					),
				},
			},
			"backup_file_compression": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					verifybool.RequiresIfTrue(path.Expressions{
						path.MatchRoot("backup"),
					}...),
				},
				Description: "default: true",
			},
			"automatic_backup": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					verifybool.RequiresIfTrue(path.Expressions{
						path.MatchRoot("backup"),
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
			"region_code": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"generation_code": schema.StringAttribute{
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
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"public_subnet": schema.BoolAttribute{
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
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
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

	reqParams := &vpostgresql.CreateCloudPostgresqlInstanceRequest{
		RegionCode:                      &r.config.RegionCode,
		CloudPostgresqlServiceName:      plan.ServiceName.ValueStringPointer(),
		CloudPostgresqlServerNamePrefix: plan.ServerNamePrefix.ValueStringPointer(),
		CloudPostgresqlUserName:         plan.UserName.ValueStringPointer(),
		CloudPostgresqlUserPassword:     plan.UserPassword.ValueStringPointer(),
		VpcNo:                           plan.VpcNo.ValueStringPointer(),
		SubnetNo:                        plan.SubnetNo.ValueStringPointer(),
		ClientCidr:                      plan.ClientCidr.ValueStringPointer(),
		CloudPostgresqlDatabaseName:     plan.DatabaseName.ValueStringPointer(),
	}

	if !plan.ImageProductCode.IsNull() && !plan.ImageProductCode.IsUnknown() {
		reqParams.CloudPostgresqlImageProductCode = plan.ImageProductCode.ValueStringPointer()
	}

	if !plan.ProductCode.IsNull() {
		reqParams.CloudPostgresqlProductCode = plan.ProductCode.ValueStringPointer()
	}

	if !plan.EngineVersionCode.IsNull() && !plan.EngineVersionCode.IsUnknown() {
		reqParams.EngineVersionCode = plan.EngineVersionCode.ValueStringPointer()
	}

	if !plan.DataStorageType.IsNull() && !plan.DataStorageType.IsUnknown() {
		reqParams.DataStorageTypeCode = plan.DataStorageType.ValueStringPointer()
	}

	if !plan.StorageEncryption.IsNull() && !plan.StorageEncryption.IsUnknown() {
		if r.config.Site == "gov" {
			resp.Diagnostics.AddError(
				"NOT SUPPORT GOV SITE",
				"`storage_encryption` does not support gov site",
			)
			return
		}
		reqParams.IsStorageEncryption = plan.StorageEncryption.ValueBoolPointer()
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		reqParams.CloudPostgresqlPort = ncloud.Int32(int32(plan.Port.ValueInt64()))
	}

	if !plan.Ha.IsNull() && !plan.Ha.IsUnknown() {
		reqParams.IsHa = plan.Ha.ValueBoolPointer()
	}

	if !plan.MultiZone.IsNull() && !plan.MultiZone.IsUnknown() {
		if r.config.Site == "gov" {
			resp.Diagnostics.AddError(
				"NOT SUPPORT GOV SITE",
				"`multi_zone` does not support gov site",
			)
			return
		}
		reqParams.IsMultiZone = plan.MultiZone.ValueBoolPointer()
	}

	if !plan.Backup.IsNull() && !plan.Backup.IsUnknown() {
		reqParams.IsBackup = plan.Backup.ValueBoolPointer()
	}

	if plan.Ha.ValueBool() {
		if !plan.Backup.IsNull() && !plan.Backup.IsUnknown() && !plan.Backup.ValueBool() {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `ha` is true, `backup` must be true or not be input",
			)
			return
		}

		if plan.MultiZone.ValueBool() {
			if plan.SecondarySubnetNo.IsNull() || plan.SecondarySubnetNo.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `multi_zone` is true, `secondary_subnet_no` must be entered",
				)
				return
			}
			reqParams.SecondarySubnetNo = plan.SecondarySubnetNo.ValueStringPointer()
		}
	}

	if !plan.BackupFileRetentionPeriod.IsNull() && !plan.BackupFileRetentionPeriod.IsUnknown() {
		reqParams.BackupFileRetentionPeriod = ncloud.Int32(int32(plan.BackupFileRetentionPeriod.ValueInt64()))
	}

	if !plan.AutomaticBackup.IsNull() && !plan.AutomaticBackup.IsUnknown() {
		reqParams.IsAutomaticBackup = plan.AutomaticBackup.ValueBoolPointer()
	}

	if !plan.BackupFileCompression.IsNull() && !plan.BackupFileCompression.IsUnknown() {
		reqParams.IsBackupFileCompression = plan.BackupFileCompression.ValueBoolPointer()
	}

	if !plan.BackupFileStorageCount.IsNull() && !plan.BackupFileStorageCount.IsUnknown() {
		reqParams.BackupFileStorageCount = ncloud.Int32(int32(plan.BackupFileStorageCount.ValueInt64()))
	}

	if reqParams.IsBackup == nil || *reqParams.IsBackup {
		if reqParams.IsAutomaticBackup == nil || *reqParams.IsAutomaticBackup {
			if !plan.BackupTime.IsNull() && !plan.BackupTime.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `backup` and `automatic_backup` is true, `backup_time` must not be entered",
				)
				return
			}
		} else {
			if plan.BackupTime.IsNull() || plan.BackupTime.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `backup` is true and `automatic_backup` is false, `backup_time` must be entered",
				)
				return
			}
			reqParams.BackupTime = plan.BackupTime.ValueStringPointer()
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

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := plan.refreshFromOutput(ctx, output); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

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

	if diags := state.refreshFromOutput(ctx, output); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

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
	// If the lookup result is 0 or already deleted, it will respond with a 400 error with a 5001017 return code.
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

			if instance == nil {
				return 0, "", fmt.Errorf("Instance is nil")
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
			if *status == DELETING || *status == DELETED {
				return instance, DELETING, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete postgresql")
		},
		Timeout:    2 * conn.DefaultTimeout,
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
	ServiceName               types.String `tfsdk:"service_name"`
	ServerNamePrefix          types.String `tfsdk:"server_name_prefix"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ClientCidr                types.String `tfsdk:"client_cidr"`
	DatabaseName              types.String `tfsdk:"database_name"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	EngineVersionCode         types.String `tfsdk:"engine_version_code"`
	DataStorageType           types.String `tfsdk:"data_storage_type"`
	StorageEncryption         types.Bool   `tfsdk:"storage_encryption"`
	Ha                        types.Bool   `tfsdk:"ha"`
	MultiZone                 types.Bool   `tfsdk:"multi_zone"`
	SecondarySubnetNo         types.String `tfsdk:"secondary_subnet_no"`
	Backup                    types.Bool   `tfsdk:"backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	BackupFileStorageCount    types.Int64  `tfsdk:"backup_file_storage_count"`
	BackupFileCompression     types.Bool   `tfsdk:"backup_file_compression"`
	AutomaticBackup           types.Bool   `tfsdk:"automatic_backup"`
	Port                      types.Int64  `tfsdk:"port"`
	RegionCode                types.String `tfsdk:"region_code"`
	GenerationCode            types.String `tfsdk:"generation_code"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	PostgresqlConfigList      types.List   `tfsdk:"postgresql_config_list"`
	PostgresqlServerList      types.List   `tfsdk:"postgresql_server_list"`
}

type postgresqlServer struct {
	ServerInstanceNo    types.String `tfsdk:"server_instance_no"`
	ServerName          types.String `tfsdk:"server_name"`
	ServerRole          types.String `tfsdk:"server_role"`
	ProductCode         types.String `tfsdk:"product_code"`
	ZoneCode            types.String `tfsdk:"zone_code"`
	SubnetNo            types.String `tfsdk:"subnet_no"`
	PublicSubnet        types.Bool   `tfsdk:"public_subnet"`
	PublicDomain        types.String `tfsdk:"public_domain"`
	PrivateDomain       types.String `tfsdk:"private_domain"`
	PrivateIp           types.String `tfsdk:"private_ip"`
	DataStorageSize     types.Int64  `tfsdk:"data_storage_size"`
	UsedDataStorageSize types.Int64  `tfsdk:"used_data_storage_size"`
	CpuCount            types.Int64  `tfsdk:"cpu_count"`
	MemorySize          types.Int64  `tfsdk:"memory_size"`
	Uptime              types.String `tfsdk:"uptime"`
	CreateDate          types.String `tfsdk:"create_date"`
}

func (r postgresqlServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_instance_no":     types.StringType,
		"server_name":            types.StringType,
		"server_role":            types.StringType,
		"product_code":           types.StringType,
		"zone_code":              types.StringType,
		"subnet_no":              types.StringType,
		"public_subnet":          types.BoolType,
		"public_domain":          types.StringType,
		"private_domain":         types.StringType,
		"private_ip":             types.StringType,
		"data_storage_size":      types.Int64Type,
		"used_data_storage_size": types.Int64Type,
		"cpu_count":              types.Int64Type,
		"memory_size":            types.Int64Type,
		"uptime":                 types.StringType,
		"create_date":            types.StringType,
	}
}

func (r *postgresqlResourceModel) refreshFromOutput(ctx context.Context, output *vpostgresql.CloudPostgresqlInstance) diag.Diagnostics {
	r.ID = types.StringPointerValue(output.CloudPostgresqlInstanceNo)
	r.ServiceName = types.StringPointerValue(output.CloudPostgresqlServiceName)
	r.VpcNo = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].VpcNo)
	r.ImageProductCode = types.StringPointerValue(output.CloudPostgresqlImageProductCode)
	r.EngineVersionCode = types.StringValue(common.ExtractEngineVersion(*output.EngineVersion))
	r.DataStorageType = types.StringPointerValue(common.GetCodePtrByCommonCode(output.CloudPostgresqlServerInstanceList[0].DataStorageType))
	r.StorageEncryption = types.BoolPointerValue(output.CloudPostgresqlServerInstanceList[0].IsStorageEncryption)
	r.Ha = types.BoolPointerValue(output.IsHa)
	r.MultiZone = types.BoolPointerValue(output.IsMultiZone)
	r.Backup = types.BoolPointerValue(output.IsBackup)
	r.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	r.BackupTime = types.StringPointerValue(output.BackupTime)
	r.Port = common.Int64ValueFromInt32(output.CloudPostgresqlPort)
	r.RegionCode = types.StringPointerValue(output.CloudPostgresqlServerInstanceList[0].RegionCode)
	r.GenerationCode = types.StringPointerValue(output.GenerationCode)

	acgList, diags := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	if diags.HasError() {
		return diags
	}
	r.AccessControlGroupNoList = acgList
	configList, diags := types.ListValueFrom(ctx, types.StringType, output.CloudPostgresqlConfigList)
	if diags.HasError() {
		return diags
	}
	r.PostgresqlConfigList = configList

	r.PostgresqlServerList, diags = listValueFromPostgresqlServerInatanceList(ctx, output.CloudPostgresqlServerInstanceList)

	return diags
}

func listValueFromPostgresqlServerInatanceList(ctx context.Context, serverInatances []*vpostgresql.CloudPostgresqlServerInstance) (basetypes.ListValue, diag.Diagnostics) {
	var serverList []postgresqlServer
	for _, server := range serverInatances {
		postgresqlServerInstance := postgresqlServer{
			ServerInstanceNo:    types.StringPointerValue(server.CloudPostgresqlServerInstanceNo),
			ServerName:          types.StringPointerValue(server.CloudPostgresqlServerName),
			ServerRole:          types.StringPointerValue(common.GetCodePtrByCommonCode(server.CloudPostgresqlServerRole)),
			ProductCode:         types.StringPointerValue(server.CloudPostgresqlProductCode),
			ZoneCode:            types.StringPointerValue(server.ZoneCode),
			SubnetNo:            types.StringPointerValue(server.SubnetNo),
			PublicSubnet:        types.BoolPointerValue(server.IsPublicSubnet),
			PublicDomain:        types.StringPointerValue(server.PublicDomain),
			PrivateDomain:       types.StringPointerValue(server.PrivateDomain),
			PrivateIp:           types.StringPointerValue(server.PrivateIp),
			DataStorageSize:     types.Int64PointerValue(server.DataStorageSize),
			UsedDataStorageSize: types.Int64PointerValue(server.UsedDataStorageSize),
			CpuCount:            common.Int64ValueFromInt32(server.CpuCount),
			MemorySize:          types.Int64PointerValue(server.MemorySize),
			Uptime:              types.StringPointerValue(server.Uptime),
			CreateDate:          types.StringPointerValue(server.CreateDate),
		}

		serverList = append(serverList, postgresqlServerInstance)
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: postgresqlServer{}.attrTypes()}, serverList)
}

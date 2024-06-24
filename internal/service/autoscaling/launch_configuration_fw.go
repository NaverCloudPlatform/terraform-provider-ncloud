package autoscaling

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	resource.Resource = &launchConfigurationResource{}
	resource.ResourceWithConfigure = &launchConfigurationResource{}
	resource.ResourceWithImportState = &launchConfigurationResource{}
)

func NewLaunchConfigResource() resource.Resource{
	return &launchConfigurationResource{}
}

type launchConfigurationResource struct {
	config *conn.ProviderConfig
}

func (l *launchConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse){
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resq)
}

func(l *launchConfigurationResource) Metadata(_ context.Contextm req resource.MetadataRequest, resp *resource.MetadataResponse){
	resp.TypeName = req.ProviderTypeName + "_launch_configuration"
}

func (l *launchConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse){
	resp.Schema=schema.Schema{
		Attributes: map[string]schema.Attribute{
			"launch_configuration_no": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(255),

				}
			},
			"server_image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
				// [TODO] ConflictsWith: []string{"member_server_image_no"},
			},
			"server_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
			},
			"member_server_image_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
				// [TODO] ConflictsWith: []string{"server_image_product_code"},
			},
			"login_key_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
			},
			"init_script_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
			},
			"user_data": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
			},
			//[TODO] "schema.Typelist : access_control_group_no_list"
			"is_encrypted_volume": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifiers.String{
					stringplanmodifier.RequiresReplace()
				},
			},

		},
	}
}

func (l *launchConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	v.config = config
}

func (l *launchConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse){
	var plan launchConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !v.config.SupportVPC {
		// [TODO] createclassicLaunchConfiguration
	}else{
		//[TODO] createVPClaucnhCOnfiguration
	}
	
	}


}
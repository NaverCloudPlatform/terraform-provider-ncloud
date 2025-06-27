package redis

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vredis"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &redisConfigGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &redisConfigGroupDataSource{}
)

func NewRedisConfigGroupDataSource() datasource.DataSource {
	return &redisConfigGroupDataSource{}
}

type redisConfigGroupDataSource struct {
	config *conn.ProviderConfig
}

func (r *redisConfigGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redis_config_group"
}

func (r *redisConfigGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 15),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`),
						"Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number.",
					),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"redis_version": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *redisConfigGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.config = config
}

func (r *redisConfigGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data redisConfigGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetRedisConfigGroup(ctx, r.config, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	data.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type redisConfigGroupDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	RedisVersion types.String `tfsdk:"redis_version"`
	Description  types.String `tfsdk:"description"`
}

func (r *redisConfigGroupDataSourceModel) refreshFromOutput(ctx context.Context, output *vredis.CloudRedisConfigGroup) {
	r.ID = types.StringPointerValue(output.ConfigGroupNo)
	r.Name = types.StringPointerValue(output.ConfigGroupName)
	r.RedisVersion = types.StringPointerValue(output.CloudRedisVersion)
	r.Description = types.StringPointerValue(output.ConfigGroupDescription)
}

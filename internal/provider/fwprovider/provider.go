package fwprovider

import (
	"context"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/mongodb"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/mysql"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
)

func New(primary interface{ Meta() interface{} }) provider.Provider {
	return &fwprovider{
		Primary: primary,
	}
}

type fwprovider struct {
	Primary interface{ Meta() interface{} }
}

func (p *fwprovider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ncloud"
}

func (p *fwprovider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"access_key": schema.StringAttribute{
				Optional:    true,
				Description: "Access key of ncloud",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "Region of ncloud",
			},
			"secret_key": schema.StringAttribute{
				Optional:    true,
				Description: "Secret key of ncloud",
			},
			"site": schema.StringAttribute{
				Optional:    true,
				Description: "Site of ncloud (public / gov / fin)",
			},
			"support_vpc": schema.BoolAttribute{
				Optional:    true,
				Description: "Support VPC platform",
			},
		},
	}
}

func (p *fwprovider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	providerConfig := p.Primary.Meta().(*conn.ProviderConfig)

	resp.DataSourceData = providerConfig
	resp.ResourceData = providerConfig
}

func (p *fwprovider) DataSources(ctx context.Context) []func() datasource.DataSource {
	var errs *multierror.Error
	var dataSources []func() datasource.DataSource

	dataSources = append(dataSources, vpc.NewVpcDataSource)
	dataSources = append(dataSources, vpc.NewVpcsDataSource)
	dataSources = append(dataSources, vpc.NewSubnetDataSource)
	dataSources = append(dataSources, vpc.NewSubnetsDataSource)
	dataSources = append(dataSources, vpc.NewNatGatewayDataSource)
	dataSources = append(dataSources, vpc.NewVpcPeeringDataSource)
	dataSources = append(dataSources, server.NewInitScriptDataSource)
	dataSources = append(dataSources, mysql.NewMysqlDataSource)
	dataSources = append(dataSources, mysql.NewMysqlProductDataSource)
	dataSources = append(dataSources, mysql.NewMysqlImageProductDataSource)
	dataSources = append(dataSources, mysql.NewMysqlImageProductsDataSource)
	dataSources = append(dataSources, mysql.NewMysqlProductsDataSource)
	dataSources = append(dataSources, mongodb.NewMongoDbDataSource)
	dataSources = append(dataSources, mongodb.NewMongoDbProductDataSource)
	dataSources = append(dataSources, mongodb.NewMongoDbProductsDataSource)
	dataSources = append(dataSources, mongodb.NewMongoDbImageProductDataSource)
	dataSources = append(dataSources, mongodb.NewMongoDbImageProductsDataSource)

	if err := errs.ErrorOrNil(); err != nil {
		tflog.Warn(ctx, "registering resources", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return dataSources
}

func (p *fwprovider) Resources(ctx context.Context) []func() resource.Resource {
	var errs *multierror.Error
	var resources []func() resource.Resource

	resources = append(resources, vpc.NewVpcResource)
	resources = append(resources, vpc.NewSubnetResource)
	resources = append(resources, vpc.NewNatGatewayResource)
	resources = append(resources, vpc.NewVpcPeeringResource)
	resources = append(resources, server.NewLoginKeyResource)
	resources = append(resources, server.NewInitScriptResource)
	resources = append(resources, mysql.NewMysqlResource)
	resources = append(resources, mongodb.NewMongoDbResource)

	if err := errs.ErrorOrNil(); err != nil {
		tflog.Warn(ctx, "registering resources", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return resources
}

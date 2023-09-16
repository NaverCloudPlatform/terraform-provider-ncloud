package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/region"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/autoscaling"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/cdss"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/classicloadbalancer"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/devtools"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/loadbalancer"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/memberserverimage"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/nasvolume"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/nks"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/ses"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func New(ctx context.Context) *schema.Provider {
	dataSourceMap := map[string]*schema.Resource{
		"ncloud_access_control_group":                    server.DataSourceNcloudAccessControlGroup(),
		"ncloud_access_control_groups":                   server.DataSourceNcloudAccessControlGroups(),
		"ncloud_access_control_rule":                     server.DataSourceNcloudAccessControlRule(),
		"ncloud_access_control_rules":                    server.DataSourceNcloudAccessControlRules(),
		"ncloud_auto_scaling_group":                      autoscaling.DataSourceNcloudAutoScalingGroup(),
		"ncloud_auto_scaling_policy":                     autoscaling.DataSourceNcloudAutoScalingPolicy(),
		"ncloud_auto_scaling_schedule":                   autoscaling.DataSourceNcloudAutoScalingSchedule(),
		"ncloud_auto_scaling_adjustment_types":           autoscaling.DataSourceNcloudAutoScalingAdjustmentTypes(),
		"ncloud_block_storage":                           server.DataSourceNcloudBlockStorage(),
		"ncloud_block_storage_snapshot":                  server.DataSourceNcloudBlockStorageSnapshot(),
		"ncloud_cdss_cluster":                            cdss.DataSourceNcloudCDSSCluster(),
		"ncloud_cdss_config_group":                       cdss.DataSourceNcloudCDSSConfigGroup(),
		"ncloud_cdss_kafka_version":                      cdss.DataSourceNcloudCDSSKafkaVersion(),
		"ncloud_cdss_kafka_versions":                     cdss.DataSourceNcloudCDSSKafkaVersions(),
		"ncloud_cdss_node_product":                       cdss.DataSourceNcloudCDSSNodeProduct(),
		"ncloud_cdss_node_products":                      cdss.DataSourceNcloudCDSSNodeProducts(),
		"ncloud_cdss_os_image":                           cdss.DataSourceNcloudCDSSOsImage(),
		"ncloud_cdss_os_images":                          cdss.DataSourceNcloudCDSSOsImages(),
		"ncloud_init_script":                             server.DataSourceNcloudInitScript(),
		"ncloud_launch_configuration":                    autoscaling.DataSourceNcloudLaunchConfiguration(),
		"ncloud_lb":                                      loadbalancer.DataSourceNcloudLb(),
		"ncloud_lb_listener":                             loadbalancer.DataSourceNcloudLbListener(),
		"ncloud_lb_target_group":                         loadbalancer.DataSourceNcloudLbTargetGroup(),
		"ncloud_member_server_image":                     memberserverimage.DataSourceNcloudMemberServerImage(),
		"ncloud_member_server_images":                    memberserverimage.DataSourceNcloudMemberServerImages(),
		"ncloud_nas_volume":                              nasvolume.DataSourceNcloudNasVolume(),
		"ncloud_nas_volumes":                             nasvolume.DataSourceNcloudNasVolumes(),
		"ncloud_network_acls":                            vpc.DataSourceNcloudNetworkAcls(),
		"ncloud_network_acl_deny_allow_groups":           vpc.DataSourceNcloudNetworkACLDenyAllowGroups(),
		"ncloud_network_interface":                       server.DataSourceNcloudNetworkInterface(),
		"ncloud_network_interfaces":                      server.DataSourceNcloudNetworkInterfaces(),
		"ncloud_nks_cluster":                             nks.DataSourceNcloudNKSCluster(),
		"ncloud_nks_clusters":                            nks.DataSourceNcloudNKSClusters(),
		"ncloud_nks_kube_config":                         nks.DataSourceNcloudNKSKubeConfig(),
		"ncloud_nks_node_pool":                           nks.DataSourceNcloudNKSNodePool(),
		"ncloud_nks_node_pools":                          nks.DataSourceNcloudNKSNodePools(),
		"ncloud_nks_server_images":                       nks.DataSourceNcloudNKSServerImages(),
		"ncloud_nks_server_products":                     nks.DataSourceNcloudNKSServerProducts(),
		"ncloud_nks_versions":                            nks.DataSourceNcloudNKSVersions(),
		"ncloud_placement_group":                         server.DataSourceNcloudPlacementGroup(),
		"ncloud_port_forwarding_rule":                    server.DataSourceNcloudPortForwardingRule(),
		"ncloud_port_forwarding_rules":                   server.DataSourceNcloudPortForwardingRules(),
		"ncloud_public_ip":                               server.DataSourceNcloudPublicIp(),
		"ncloud_regions":                                 region.DataSourceNcloudRegions(),
		"ncloud_root_password":                           server.DataSourceNcloudRootPassword(),
		"ncloud_route_table":                             vpc.DataSourceNcloudRouteTable(),
		"ncloud_route_tables":                            vpc.DataSourceNcloudRouteTables(),
		"ncloud_server":                                  server.DataSourceNcloudServer(),
		"ncloud_server_image":                            server.DataSourceNcloudServerImage(),
		"ncloud_server_images":                           server.DataSourceNcloudServerImages(),
		"ncloud_server_product":                          server.DataSourceNcloudServerProduct(),
		"ncloud_server_products":                         server.DataSourceNcloudServerProducts(),
		"ncloud_servers":                                 server.DataSourceNcloudServers(),
		"ncloud_ses_cluster":                             ses.DataSourceNcloudSESCluster(),
		"ncloud_ses_clusters":                            ses.DataSourceNcloudSESClusters(),
		"ncloud_ses_node_os_images":                      ses.DataSourceNcloudSESNodeOsImage(),
		"ncloud_ses_node_products":                       ses.DataSourceNcloudSESNodeProduct(),
		"ncloud_ses_versions":                            ses.DataSourceNcloudSESVersions(),
		"ncloud_sourcebuild_project_computes":            devtools.DataSourceNcloudSourceBuildComputes(),
		"ncloud_sourcebuild_project":                     devtools.DataSourceNcloudSourceBuildProject(),
		"ncloud_sourcebuild_project_docker_engines":      devtools.DataSourceNcloudSourceBuildDockerEngines(),
		"ncloud_sourcebuild_project_os":                  devtools.DataSourceNcloudSourceBuildOs(),
		"ncloud_sourcebuild_project_os_runtime_versions": devtools.DataSourceNcloudSourceBuildRuntimeVersions(),
		"ncloud_sourcebuild_project_os_runtimes":         devtools.DataSourceNcloudSourceBuildRuntimes(),
		"ncloud_sourcebuild_projects":                    devtools.DataSourceNcloudSourceBuildProjects(),
		"ncloud_sourcecommit_repositories":               devtools.DataSourceNcloudSourceCommitRepositories(),
		"ncloud_sourcecommit_repository":                 devtools.DataSourceNcloudSourceCommitRepository(),
		"ncloud_sourcedeploy_project_stage":              devtools.DataSourceNcloudSourceDeployStageContext(),
		"ncloud_sourcedeploy_project_stage_scenario":     devtools.DataSourceNcloudSourceDeployScenarioContext(),
		"ncloud_sourcedeploy_project_stage_scenarios":    devtools.DataSourceNcloudSourceDeployscenariosContext(),
		"ncloud_sourcedeploy_project_stages":             devtools.DataSourceNcloudSourceDeployStagesContext(),
		"ncloud_sourcedeploy_projects":                   devtools.DataSourceNcloudSourceDeployProjectsContext(),
		"ncloud_sourcepipeline_project":                  devtools.DataSourceNcloudSourcePipelineProject(),
		"ncloud_sourcepipeline_projects":                 devtools.DataSourceNcloudSourcePipelineProjects(),
		"ncloud_sourcepipeline_trigger_timezone":         devtools.DataSourceNcloudSourcePipelineTimeZone(),
		"ncloud_zones":                                   zone.DataSourceNcloudZones(),
	}

	resourceMap := map[string]*schema.Resource{
		"ncloud_access_control_group_rule":           server.ResourceNcloudAccessControlGroupRule(),
		"ncloud_access_control_group":                server.ResourceNcloudAccessControlGroup(),
		"ncloud_auto_scaling_group":                  autoscaling.ResourceNcloudAutoScalingGroup(),
		"ncloud_auto_scaling_policy":                 autoscaling.ResourceNcloudAutoScalingPolicy(),
		"ncloud_auto_scaling_schedule":               autoscaling.ResourceNcloudAutoScalingSchedule(),
		"ncloud_block_storage_snapshot":              server.ResourceNcloudBlockStorageSnapshot(),
		"ncloud_block_storage":                       server.ResourceNcloudBlockStorage(),
		"ncloud_cdss_cluster":                        cdss.ResourceNcloudCDSSCluster(),
		"ncloud_cdss_config_group":                   cdss.ResourceNcloudCDSSConfigGroup(),
		"ncloud_init_script":                         server.ResourceNcloudInitScript(),
		"ncloud_launch_configuration":                autoscaling.ResourceNcloudLaunchConfiguration(),
		"ncloud_lb_listener":                         loadbalancer.ResourceNcloudLbListener(),
		"ncloud_lb_target_group_attachment":          loadbalancer.ResourceNcloudLbTargetGroupAttachment(),
		"ncloud_lb_target_group":                     loadbalancer.ResourceNcloudLbTargetGroup(),
		"ncloud_lb":                                  loadbalancer.ResourceNcloudLb(),
		"ncloud_load_balancer_ssl_certificate":       classicloadbalancer.ResourceNcloudLoadBalancerSSLCertificate(),
		"ncloud_load_balancer":                       classicloadbalancer.ResourceNcloudLoadBalancer(),
		"ncloud_nas_volume":                          nasvolume.ResourceNcloudNasVolume(),
		"ncloud_network_acl":                         vpc.ResourceNcloudNetworkACL(),
		"ncloud_network_acl_deny_allow_group":        vpc.ResourceNcloudNetworkACLDenyAllowGroup(),
		"ncloud_network_acl_rule":                    vpc.ResourceNcloudNetworkACLRule(),
		"ncloud_network_interface":                   server.ResourceNcloudNetworkInterface(),
		"ncloud_nks_cluster":                         nks.ResourceNcloudNKSCluster(),
		"ncloud_nks_node_pool":                       nks.ResourceNcloudNKSNodePool(),
		"ncloud_placement_group":                     server.ResourceNcloudPlacementGroup(),
		"ncloud_port_forwarding_rule":                server.ResourceNcloudPortForwadingRule(),
		"ncloud_public_ip":                           server.ResourceNcloudPublicIpInstance(),
		"ncloud_route":                               vpc.ResourceNcloudRoute(),
		"ncloud_route_table":                         vpc.ResourceNcloudRouteTable(),
		"ncloud_route_table_association":             vpc.ResourceNcloudRouteTableAssociation(),
		"ncloud_server":                              server.ResourceNcloudServer(),
		"ncloud_ses_cluster":                         ses.ResourceNcloudSESCluster(),
		"ncloud_sourcebuild_project":                 devtools.ResourceNcloudSourceBuildProject(),
		"ncloud_sourcecommit_repository":             devtools.ResourceNcloudSourceCommitRepository(),
		"ncloud_sourcedeploy_project_stage_scenario": devtools.ResourceNcloudSourceDeployScenario(),
		"ncloud_sourcedeploy_project_stage":          devtools.ResourceNcloudSourceDeployStage(),
		"ncloud_sourcedeploy_project":                devtools.ResourceNcloudSourceDeployProject(),
		"ncloud_sourcepipeline_project":              devtools.ResourceNcloudSourcePipeline(),
	}

	return &schema.Provider{
		Schema:               SchemaMap(),
		DataSourcesMap:       dataSourceMap,
		ResourcesMap:         resourceMap,
		ConfigureContextFunc: ProviderConfigure,
	}
}

func SchemaMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access key of ncloud",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Region of ncloud",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Secret key of ncloud",
		},
		"site": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Site of ncloud (public / gov / fin)",
		},
		"support_vpc": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Support VPC platform",
		},
	}
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	providerConfig := conn.ProviderConfig{
		SupportVPC: true,
	}

	if _, ok := getOrFromEnv(d, "support_vpc", "NCLOUD_SUPPORT_VPC"); !ok {
		providerConfig.SupportVPC = false
	}

	// Set site
	if site, ok := getOrFromEnv(d, "site", "NCLOUD_SITE"); ok {
		providerConfig.Site = site.(string)

		switch site {
		case "gov":
			os.Setenv("NCLOUD_API_GW", "https://ncloud.apigw.gov-ntruss.com")
		case "fin":
			os.Setenv("NCLOUD_API_GW", "https://fin-ncloud.apigw.fin-ntruss.com")
		}
	}

	// Fin only supports VPC
	if providerConfig.Site == "fin" {
		providerConfig.SupportVPC = true
	}

	accessKey, ok := getOrFromEnv(d, "access_key", "NCLOUD_ACCESS_KEY")
	if !ok {
		return nil, diag.Errorf("missing provider configuration: ACCESS_KEY")
	}
	secretKey, ok := getOrFromEnv(d, "secret_key", "NCLOUD_SECRET_KEY")
	if !ok {
		return nil, diag.Errorf("missing provider configuration: SECRET_KEY")
	}
	region, ok := getOrFromEnv(d, "region", "NCLOUD_REGION")
	if !ok {
		return nil, diag.Errorf("missing provider configuration: REGION")
	}

	// Set client
	config := conn.Config{
		AccessKey: accessKey.(string),
		SecretKey: secretKey.(string),
		Region:    region.(string),
	}

	if client, err := config.Client(); err != nil {
		return nil, diag.FromErr(err)
	} else {
		providerConfig.Client = client
	}

	// Set region
	if err := conn.SetRegionCache(providerConfig.Client, providerConfig.SupportVPC); err != nil {
		return nil, diag.FromErr(err)
	}

	if conn.IsValidRegionCode(region.(string)) {
		os.Setenv("NCLOUD_REGION", region.(string))
		providerConfig.RegionCode = region.(string)
		if !providerConfig.SupportVPC {
			providerConfig.RegionNo = *conn.GetRegionNoByCode(region.(string))
		}
	} else {
		return nil, []diag.Diagnostic{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("no region data for region_code `%s`. please change region_code and try again", region),
			},
		}
	}

	return &providerConfig, nil
}

func getOrFromEnv(d *schema.ResourceData, name, env string) (any, bool) {
	if v, ok := d.GetOk(name); ok {
		return v, true
	}

	if v := os.Getenv(env); v != "" {
		return v, true
	}
	return nil, false
}

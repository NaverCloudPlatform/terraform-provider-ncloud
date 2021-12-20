package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v3"
)

func init() {
	RegisterDataSource("ncloud_nks_kube_config", dataSourceNcloudNKSKubeConfig())
}

func dataSourceNcloudNKSKubeConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSKubeConfigRead,
		Schema: map[string]*schema.Schema{
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_ca_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudNKSKubeConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_nks_kube_config`"))
	}
	clusterUuid := d.Get("cluster_uuid").(string)

	kubeConfig, err := getNKSKubeConfig(ctx, config, clusterUuid)
	if err != nil {
		return diag.FromErr(err)
	}

	if kubeConfig == nil {
		d.SetId("")
		return nil
	}

	d.SetId(clusterUuid)

	d.Set("host", kubeConfig.Clusters[0].Cluster.Server)
	d.Set("client_certificate", kubeConfig.Users[0].User.ClientCertificateData)
	d.Set("client_key", kubeConfig.Users[0].User.ClientKeyData)
	d.Set("cluster_ca_certificate", kubeConfig.Clusters[0].Cluster.ClusterCaCertificate)

	return nil
}

func getNKSKubeConfig(ctx context.Context, config *ProviderConfig, uuid string) (kc *KubeConfig, err error) {
	resp, err := config.Client.vnks.V2Api.ClustersUuidKubeconfigGet(ctx, ncloud.String(uuid))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(ncloud.StringValue(resp.Kubeconfig)), &kc)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return kc, nil
}

type KubeConfig struct {
	Clusters []struct {
		Cluster struct {
			Server               string `yaml:"server"`
			ClusterCaCertificate string `yaml:"certificate-authority-data"`
		}
	}
	Users []struct {
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		}
	}
}

package classicloadbalancer

import (
	"fmt"
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudLoadBalancerSSLCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLoadBalancerSSLCertificateCreate,
		Read:   resourceNcloudLoadBalancerSSLCertificateRead,
		Update: resourceNcloudLoadBalancerSSLCertificateUpdate,
		Delete: resourceNcloudLoadBalancerSSLCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"certificate_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of a certificate to add",
			},
			"privatekey": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private key for a certificate",
			},
			"publickey_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public key for a certificate",
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chainca certificate (Required if the certificate is issued with a chainca)",
			},
		},
	}
}

func resourceNcloudLoadBalancerSSLCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client
	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		return NotSupportVpc("resource `ncloud_load_balancer_ssl_certificate`")
	}

	reqParams, err := buildCreateLoadBalancerSSLCertificateParams(d)
	if err != nil {
		LogErrorResponse("AddLoadBalancerSslCertificate", err, reqParams)
		return err
	}

	LogCommonRequest("AddLoadBalancerSslCertificate", reqParams)

	resp, err := client.Loadbalancer.V2Api.AddLoadBalancerSslCertificate(reqParams)
	if err != nil {
		LogErrorResponse("AddLoadBalancerSslCertificate", err, reqParams)
		return err
	}

	LogCommonResponse("AddLoadBalancerSslCertificate", GetCommonResponse(resp))

	if len(resp.SslCertificateList) == 0 {
		return fmt.Errorf("no SSL certificate found in the API response")
	}

	cert := resp.SslCertificateList[0]
	d.SetId(*cert.CertificateName)

	return resourceNcloudLoadBalancerSSLCertificateRead(d, meta)
}

func resourceNcloudLoadBalancerSSLCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client

	lb, err := GetLoadBalancerSslCertificateList(client, d.Id())
	if err != nil {
		return err
	}
	if lb != nil {
		d.Set("certificate_name", lb.CertificateName)
		d.Set("privatekey", lb.PrivateKey)
		d.Set("publickey_certificate", lb.PublicKeyCertificate)
		d.Set("certificate_chain", lb.CertificateChain)
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudLoadBalancerSSLCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client
	if err := deleteLoadBalancerSSLCertificate(client, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudLoadBalancerSSLCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudLoadBalancerSSLCertificateRead(d, meta)
}

func buildCreateLoadBalancerSSLCertificateParams(d *schema.ResourceData) (*loadbalancer.AddLoadBalancerSslCertificateRequest, error) {
	reqParams := &loadbalancer.AddLoadBalancerSslCertificateRequest{
		CertificateName:      ncloud.String(d.Get("certificate_name").(string)),
		PrivateKey:           ncloud.String(d.Get("privatekey").(string)),
		PublicKeyCertificate: ncloud.String(d.Get("publickey_certificate").(string)),
	}

	if certificateChain, ok := d.GetOk("certificate_chain"); ok {
		reqParams.CertificateChain = ncloud.String(certificateChain.(string))
	}

	return reqParams, nil
}

func GetLoadBalancerSslCertificateList(client *conn.NcloudAPIClient, certificateName string) (*loadbalancer.SslCertificate, error) {
	reqParams := loadbalancer.GetLoadBalancerSslCertificateListRequest{CertificateName: ncloud.String(certificateName)}
	LogCommonRequest("GetLoadBalancerSslCertificateList", reqParams)
	resp, err := client.Loadbalancer.V2Api.GetLoadBalancerSslCertificateList(&reqParams)
	if err != nil {
		LogErrorResponse("GetLoadBalancerSslCertificateList", err, certificateName)
		return nil, err
	}
	LogCommonResponse("GetLoadBalancerSslCertificateList", GetCommonResponse(resp))

	for _, cert := range resp.SslCertificateList {
		if certificateName == ncloud.StringValue(cert.CertificateName) {
			log.Printf("[DEBUG] %s CertificateName: %s,", "GetLoadBalancerSslCertificateList", ncloud.StringValue(cert.CertificateName))
			return cert, nil
		}
	}
	return nil, nil
}

func deleteLoadBalancerSSLCertificate(client *conn.NcloudAPIClient, certificateName string) error {
	reqParams := loadbalancer.DeleteLoadBalancerSslCertificateRequest{CertificateName: ncloud.String(certificateName)}
	LogCommonRequest("DeleteLoadBalancerSslCertificate", reqParams)
	resp, err := client.Loadbalancer.V2Api.DeleteLoadBalancerSslCertificate(&reqParams)
	if err != nil {
		LogErrorResponse("DeleteLoadBalancerSslCertificate", err, certificateName)
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}
	LogCommonResponse("DeleteLoadBalancerSslCertificate", commonResponse)

	return nil
}

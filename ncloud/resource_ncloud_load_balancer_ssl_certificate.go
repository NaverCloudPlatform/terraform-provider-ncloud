package ncloud

import (
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterResource("ncloud_load_balancer_ssl_certificate", resourceNcloudLoadBalancerSSLCertificate())
}

func resourceNcloudLoadBalancerSSLCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLoadBalancerSSLCertificateCreate,
		Read:   resourceNcloudLoadBalancerSSLCertificateRead,
		Delete: resourceNcloudLoadBalancerSSLCertificateDelete,
		Update: resourceNcloudLoadBalancerSSLCertificateUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
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
	client := meta.(*ProviderConfig).Client
	config := meta.(*ProviderConfig)

	if config.SupportVPC {
		return NotSupportVpc("resource `ncloud_load_balancer_ssl_certificate`")
	}

	reqParams, err := buildCreateLoadBalancerSSLCertificateParams(d)
	if err != nil {
		logErrorResponse("AddLoadBalancerSslCertificate", err, reqParams)
		return err
	}

	logCommonRequest("AddLoadBalancerSslCertificate", reqParams)

	resp, err := client.loadbalancer.V2Api.AddLoadBalancerSslCertificate(reqParams)
	if err != nil {
		logErrorResponse("AddLoadBalancerSslCertificate", err, reqParams)
		return err
	}

	logCommonResponse("AddLoadBalancerSslCertificate", GetCommonResponse(resp))

	cert := resp.SslCertificateList[0]
	d.SetId(*cert.CertificateName)

	return resourceNcloudLoadBalancerSSLCertificateRead(d, meta)
}

func resourceNcloudLoadBalancerSSLCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	lb, err := getLoadBalancerSslCertificateList(client, d.Id())
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
	client := meta.(*ProviderConfig).Client
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

func getLoadBalancerSslCertificateList(client *NcloudAPIClient, certificateName string) (*loadbalancer.SslCertificate, error) {
	reqParams := loadbalancer.GetLoadBalancerSslCertificateListRequest{CertificateName: ncloud.String(certificateName)}
	logCommonRequest("GetLoadBalancerSslCertificateList", reqParams)
	resp, err := client.loadbalancer.V2Api.GetLoadBalancerSslCertificateList(&reqParams)
	if err != nil {
		logErrorResponse("GetLoadBalancerSslCertificateList", err, certificateName)
		return nil, err
	}
	logCommonResponse("GetLoadBalancerSslCertificateList", GetCommonResponse(resp))

	for _, cert := range resp.SslCertificateList {
		if certificateName == ncloud.StringValue(cert.CertificateName) {
			log.Printf("[DEBUG] %s CertificateName: %s,", "GetLoadBalancerSslCertificateList", ncloud.StringValue(cert.CertificateName))
			return cert, nil
		}
	}
	return nil, nil
}

func deleteLoadBalancerSSLCertificate(client *NcloudAPIClient, certificateName string) error {
	reqParams := loadbalancer.DeleteLoadBalancerSslCertificateRequest{CertificateName: ncloud.String(certificateName)}
	logCommonRequest("DeleteLoadBalancerSslCertificate", reqParams)
	resp, err := client.loadbalancer.V2Api.DeleteLoadBalancerSslCertificate(&reqParams)
	if err != nil {
		logErrorResponse("DeleteLoadBalancerSslCertificate", err, certificateName)
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}
	logCommonResponse("DeleteLoadBalancerSslCertificate", commonResponse)

	return nil
}

package ncloud

import (
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

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
				Description: "Public key certificate",
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
	log.Println("[DEBUG] resourceNcloudLoadBalancerSSLCertificateCreate")
	conn := meta.(*NcloudSdk).conn

	reqParams, err := buildCreateLoadBalancerSSLCertificateParams(d)
	if err != nil {
		logErrorResponse("AddLoadBalancerSslCertificate", err, reqParams)
		return err
	}

	resp, err := conn.AddLoadBalancerSslCertificate(reqParams)
	if err != nil {
		logErrorResponse("AddLoadBalancerSslCertificate", err, reqParams)
		return err
	}

	logCommonResponse("AddLoadBalancerSslCertificate", reqParams, resp.CommonResponse)

	cert := &resp.SslCertificateList[0]
	d.SetId(cert.CertificateName)

	return resourceNcloudLoadBalancerSSLCertificateRead(d, meta)
}

func resourceNcloudLoadBalancerSSLCertificateRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudLoadBalancerSSLCertificateRead")
	conn := meta.(*NcloudSdk).conn

	lb, err := getLoadBalancerSslCertificateList(conn, d.Id())
	if err != nil {
		return err
	}
	if lb != nil {
		d.Set("certificate_name", lb.CertificateName)
		d.Set("privatekey", lb.PrivateKey)
		d.Set("publickey_certificate", lb.PublicKeyCertificate)
		d.Set("certificate_chain", lb.CertificateChain)
	}

	return nil
}

func resourceNcloudLoadBalancerSSLCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudLoadBalancerSSLCertificateDelete")
	conn := meta.(*NcloudSdk).conn
	return deleteLoadBalancerSSLCertificate(conn, d.Id())
}

func resourceNcloudLoadBalancerSSLCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudLoadBalancerSSLCertificateUpdate")
	return resourceNcloudLoadBalancerSSLCertificateRead(d, meta)
}

func buildCreateLoadBalancerSSLCertificateParams(d *schema.ResourceData) (*sdk.RequestAddSslCertificate, error) {
	reqParams := &sdk.RequestAddSslCertificate{
		CertificateName:      d.Get("certificate_name").(string),
		PrivateKey:           d.Get("privatekey").(string),
		PublicKeyCertificate: d.Get("publickey_certificate").(string),
	}

	if certificateChain, ok := d.GetOk("certificate_chain"); ok {
		reqParams.CertificateChain = certificateChain.(string)
	}

	return reqParams, nil
}

func getLoadBalancerSslCertificateList(conn *sdk.Conn, certificateName string) (*sdk.SslCertificate, error) {
	resp, err := conn.GetLoadBalancerSslCertificateList(certificateName)
	if err != nil {
		logErrorResponse("GetLoadBalancerSslCertificateList", err, certificateName)
		return nil, err
	}
	logCommonResponse("GetLoadBalancerSslCertificateList", certificateName, resp.CommonResponse)

	for _, cert := range resp.SslCertificateList {
		if certificateName == cert.CertificateName {
			log.Printf("[DEBUG] %s CertificateName: %s,", "GetLoadBalancerSslCertificateList", cert.CertificateName)
			return &cert, nil
		}
	}
	return nil, nil
}

func deleteLoadBalancerSSLCertificate(conn *sdk.Conn, certificateName string) error {
	resp, err := conn.DeleteLoadBalancerSslCertificate(certificateName)
	if err != nil {
		logErrorResponse("DeleteLoadBalancerSslCertificate", err, certificateName)
		return err
	}
	var commonResponse = common.CommonResponse{}
	if resp != nil {
		commonResponse = resp.CommonResponse
	}
	logCommonResponse("DeleteLoadBalancerSslCertificate", certificateName, commonResponse)

	return nil
}

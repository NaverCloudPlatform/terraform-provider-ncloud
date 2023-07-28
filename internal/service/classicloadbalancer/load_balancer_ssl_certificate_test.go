package classicloadbalancer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/classicloadbalancer"
)

func TestAccNcloudLoadBalancerSSLCertificateBasic(t *testing.T) {
	var sc loadbalancer.SslCertificate
	prefix := GetTestPrefix()
	testSSLCertificateName := prefix + "_cert"
	testLoadBalancerName := prefix + "_lb"
	testCertPEM := `-----BEGIN CERTIFICATE-----
MIIDGDCCAgACCQDGDiYiQixnsTANBgkqhkiG9w0BAQsFADBOMQswCQYDVQQGEwJL
UjEOMAwGA1UECAwFc2VvdWwxDjAMBgNVBAcMBXNlb3VsMQwwCgYDVQQKDANuYnAx
ETAPBgNVBAsMCG5jbG91ZGV2MB4XDTE4MDYwODA5NTEyOVoXDTE4MDcwODA5NTEy
OVowTjELMAkGA1UEBhMCS1IxDjAMBgNVBAgMBXNlb3VsMQ4wDAYDVQQHDAVzZW91
bDEMMAoGA1UECgwDbmJwMREwDwYDVQQLDAhuY2xvdWRldjCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAMV0paXrbjzipw875D6ZKABd7KQFvHH46fWAxwRb
wz/jrQisPcopwJTutSd19fDCdLCsL62+S/oAJFrFK32BMxgK/Feamepj9SS35VZR
yWO5rKrI6a5HkEFMexzz+qr5jN7me/pihqxMPsintEqx6I7ajvXAQGOSr9qiDI1I
T6XJ2++atgWqlDok37HsMyIeMJx0fRkmRr9z5fzfjgqxpbEpoXvXjwXGwLS/aGJZ
/ie6fCRiUnDCujbVCXePCGI4AtQHjXrcmWKwthZ4UEwsZAyv8qtpjEjvpfVlpmkf
LlOM8R6mowcnB3L9csH7aWTpwVXbKfiJNzysxueF+y8sV4kCAwEAATANBgkqhkiG
9w0BAQsFAAOCAQEABDFl73C8ta9zYfyQXIbtv2tXt6oIhphjD5sV5KO6lgVSw7db
XoiDlpQb5/LIXVYEwf8GLlSPORLsU36DQk0EyE6veTDYq4Cexkp5U1ca+jAGSum3
Rtl02Dj6w7pz44sFnZc1IwKWUI7nTc0rQoloKRVyWnb5EPoNEe5QI1R5HJh7vV2M
OUpdTkOueLjy+4BfUAt+LqeNz5u9WhvVgkFul1V9e3UXOaf0KvMYaRUYRj+IThaK
7/UsUgXKmbhVOZrgNPjQT1cTC1lSD+xfCeA5UT/6xqxq91LL6cSNBOqEoDwyWGt5
Kp6GlcmqTBU7RXAwPE4MeQ20yMmdipoPMb8nug==
-----END CERTIFICATE-----`
	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAxXSlpetuPOKnDzvkPpkoAF3spAW8cfjp9YDHBFvDP+OtCKw9
yinAlO61J3X18MJ0sKwvrb5L+gAkWsUrfYEzGAr8V5qZ6mP1JLflVlHJY7msqsjp
rkeQQUx7HPP6qvmM3uZ7+mKGrEw+yKe0SrHojtqO9cBAY5Kv2qIMjUhPpcnb75q2
BaqUOiTfsewzIh4wnHR9GSZGv3Pl/N+OCrGlsSmhe9ePBcbAtL9oYln+J7p8JGJS
cMK6NtUJd48IYjgC1AeNetyZYrC2FnhQTCxkDK/yq2mMSO+l9WWmaR8uU4zxHqaj
BycHcv1ywftpZOnBVdsp+Ik3PKzG54X7LyxXiQIDAQABAoIBAEOQY3H/uivZPmLH
EpWc4IQnn2aMk+vHyX54/yBtqcS9yiKSlV4MpVoQyCnlgi9MypL9iB8CY4r663Wn
y/bY87vBXpE3VH1QkLxstGux9qBKE1wo/VTmJeVCH0pL7bT9SQeohDmr5vsj58PP
JrD8aWAgRxSuIRoxQj0kf/kECkTm0IVD0YqMe4b956aul5noet8wbQ+2xEjsnBCy
Z4rI8JKMq+V2brgOn7XK4BLfL0YbIFnS3s4A5zu0xvi7aAK5gb7igz0zMVLp/Rhc
u8HNnxsDkC7Fj8XeWEw4I4IGfS1y7Pz83LJfqM9kfdhjwn4h25KWOPxS6bOwpVqz
wzZV1fUCgYEA6RTWBW8h15FZ7jlUG5JAzqQGxtKCFIFHZmkIiVfV4c/45XVfD3Kz
5xg6RsQCUMF3sBTcLNZYuRcS7es2eL6aCXXTkfT4I5KUoMuSN5YoDvvRZLz+JZMy
XmhYzzy/Uh7ZsOTc9Fdo/YXJIXD0Cl1P0nGiunOH2Qhz3eno4gJtJBcCgYEA2N8I
sjRy9GW+Ha2Jt13RHPoXhQbFk39zZAn9B8Rx2z4FaAO//l1S8PDL417yivdLY1UW
7Tu2MOGE3YgYTefAlN61Rtv/QjNjvYA2c4GjA2yNvVCKoPLV5sudhyKUZ6mQGLdu
oiBa/jDFbtbWl7u9NiHyvEX6zbfmYnN+AKn0hV8CgYBXTCiEzITeWmBWaz5nPTXs
r16iZQG3cFwvrTM3TaCb/Or59iXugUWETny1OICtgmizmHyGhpmgaVX7qlcyjiDf
XjQpvJibqjDksJpJG4JRaluY4XhG1oTM+0QYCmaV+VwLdwySr5JxMgSM8+NTZnOZ
HFqYfuDoltPez9cbn1EFbQKBgQCdu8olYsRhQUa/ayKI/XFEdBl7JWu6Va5limZA
qf5tiXSBLIkNxm6200xXuQ0LScXJH3AnZ5ChiMUMIxoaP37wR/Ls8MF9MsdOYtw3
sogPy3pjwRqy6SvuSxXt3Za2trsZXwDWZlYIHwzaCuPVRDTgFFzp1rQNv72OyZVR
gktYXQKBgQCegh5CycXxo+l2AZiok8qEYythpSMoWgbYIXOl3ewDYPBXc7czwKm0
tvFCSA63gNUikyzOqW0MT7cuvZvR1/7HfvqIB6ZKN2icWndTdHj7TOflaJr7TspN
GTfhUTV7jTQ0dt9U1E+oxRkjqC2HFYlpewXP0rcQxhtK7p6kiaUDIw==
-----END RSA PRIVATE KEY-----`

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if *sc.CertificateName != testSSLCertificateName {
				return fmt.Errorf("not found: CertificateName [%s]", testSSLCertificateName)
			}
			if *sc.PrivateKey != testPrivateKey {
				return fmt.Errorf("not found: PrivateKey [%s]", testPrivateKey)
			}
			if *sc.PublicKeyCertificate != testCertPEM {
				return fmt.Errorf("not found: PublicKeyCertificate [%s]", testCertPEM)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(false),
		CheckDestroy: testAccCheckLoadBalancerSSLCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerSSLCertificateConfig(testSSLCertificateName, testPrivateKey, testCertPEM, testLoadBalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerSSLCertificateExists("ncloud_load_balancer_ssl_certificate.cert", &sc),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer_ssl_certificate.cert",
						"certificate_name",
						testSSLCertificateName),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer_ssl_certificate.cert",
						"privatekey",
						testPrivateKey),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer_ssl_certificate.cert",
						"publickey_certificate",
						testCertPEM),
				),
			},
			{
				ResourceName:      "ncloud_load_balancer_ssl_certificate.cert",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoadBalancerSSLCertificateExists(n string, i *loadbalancer.SslCertificate) resource.TestCheckFunc {
	return testAccCheckLoadBalancerSSLCertificateExistsWithProvider(n, i, func() *schema.Provider { return GetTestProvider(false) })
}

func testAccCheckLoadBalancerSSLCertificateExistsWithProvider(n string, i *loadbalancer.SslCertificate, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		client := provider.Meta().(*conn.ProviderConfig).Client
		sc, err := classicloadbalancer.GetLoadBalancerSslCertificateList(client, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if sc != nil {
			*i = *sc
			return nil
		}

		return fmt.Errorf("SSL Certificate not found")
	}
}

func testAccCheckLoadBalancerSSLCertificateDestroy(s *terraform.State) error {
	return testAccCheckLoadBalancerSSLCertificateDestroyWithProvider(s, GetTestProvider(false))
}

func testAccCheckLoadBalancerSSLCertificateDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*conn.ProviderConfig).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_load_balancer_ssl_certificate" {
			continue
		}
		sc, err := classicloadbalancer.GetLoadBalancerSslCertificateList(client, rs.Primary.ID)
		if sc == nil {
			return nil
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("failed to delete SSL Certificate: %s", *sc.CertificateName)
	}

	return nil
}

func testAccLoadBalancerSSLCertificateConfig(certificateName string, privatekey string, publickeyCertificate string, lbName string) string {
	return fmt.Sprintf(`
		resource "ncloud_load_balancer_ssl_certificate" "cert" {
			certificate_name      = "%s"
			privatekey            = "%s"
			publickey_certificate = "%s"
		}

		resource "ncloud_load_balancer" "lb" {
			name           = "%s"
			algorithm_type = "SIPHS"
			description    = "tftest_lb description"

			rule_list {
				protocol_type        = "HTTP"
				load_balancer_port   = 80
				server_port          = 80
				l7_health_check_path = "/monitor/l7check"
			}
			
			rule_list {
				protocol_type        = "HTTPS"
				load_balancer_port   = 443
				server_port          = 443
				l7_health_check_path = "/monitor/l7check"
				certificate_name     = "${ncloud_load_balancer_ssl_certificate.cert.certificate_name}"
			}

			region = "KR"
		}
		`, certificateName, strings.Replace(privatekey, "\n", "\\n", -1), strings.Replace(publickeyCertificate, "\n", "\\n", -1), lbName)
}

package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

// DeleteLoadBalancerSslCertificate deletes SSL Certificate
func (s *Conn) DeleteLoadBalancerSslCertificate(certificateName string) (*SslCertificateList, error) {
	params := make(map[string]string)

	if certificateName == "" {
		return nil, fmt.Errorf("CertificateName is required field")
	}

	params["certificateName"] = certificateName
	params["action"] = "deleteLoadBalancerSslCertificate"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"loadbalancer/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := SslCertificateList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var SslCertificateList = SslCertificateList{}
	if err := xml.Unmarshal([]byte(bytes), &SslCertificateList); err != nil {
		return nil, err
	}

	return &SslCertificateList, nil
}

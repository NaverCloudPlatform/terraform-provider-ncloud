package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processAddSslCertificateParams(reqParams *RequestAddSslCertificate) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return nil, errors.New("CertificateName, PrivateKey and PublicKeyCertificate are required fields")
	}

	if err := validateRequiredField("CertificateName", reqParams.CertificateName); err != nil {
		return nil, err
	}
	params["certificateName"] = reqParams.CertificateName

	if err := validateRequiredField("PrivateKey", reqParams.PrivateKey); err != nil {
		return nil, err
	}
	params["privateKey"] = reqParams.PrivateKey

	if err := validateRequiredField("PublicKeyCertificate", reqParams.PublicKeyCertificate); err != nil {
		return nil, err
	}
	params["publicKeyCertificate"] = reqParams.PublicKeyCertificate

	if reqParams.CertificateChain != "" {
		params["certificateChain"] = reqParams.CertificateChain
	}

	return params, nil
}

// AddLoadBalancerSslCertificate get SSL Certificate
func (s *Conn) AddLoadBalancerSslCertificate(reqParams *RequestAddSslCertificate) (*SslCertificateList, error) {
	params, err := processAddSslCertificateParams(reqParams)

	if err != nil {
		return nil, err
	}

	params["action"] = "addLoadBalancerSslCertificate"

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

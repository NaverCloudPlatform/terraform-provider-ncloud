package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processChangeLoadBalancerInstanceConfigurationParams(reqParams *RequestChangeLoadBalancerInstanceConfiguration) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return nil, errors.New("LoadBalancerInstanceNo field is required")
	}

	if err := validateRequiredField("LoadBalancerInstanceNo", reqParams.LoadBalancerInstanceNo); err != nil {
		return nil, err
	}
	params["loadBalancerInstanceNo"] = reqParams.LoadBalancerInstanceNo

	if err := validateRequiredField("LoadBalancerAlgorithmTypeCode", reqParams.LoadBalancerAlgorithmTypeCode); err != nil {
		return nil, err
	}

	if err := validateIncludeValues("LoadBalancerAlgorithmTypeCode", reqParams.LoadBalancerAlgorithmTypeCode, []string{"RR", "LC", "SIPHS"}); err != nil {
		return nil, err
	}
	params["loadBalancerAlgorithmTypeCode"] = reqParams.LoadBalancerAlgorithmTypeCode

	if reqParams.LoadBalancerDescription != "" {
		if err := validateStringLenBetween("LoadBalancerDescription", reqParams.LoadBalancerDescription, 1, 1000); err != nil {
			return nil, err
		}
		params["loadBalancerDescription"] = reqParams.LoadBalancerDescription
	}

	for k, v := range reqParams.LoadBalancerRuleList {
		// ProtocolTypeCode
		protocolTypeCode := fmt.Sprintf("loadBalancerRuleList.%d.protocolTypeCode", k+1)
		if err := validateRequiredField(protocolTypeCode, v.ProtocolTypeCode); err != nil {
			return params, err
		}
		if err := validateIncludeValues(protocolTypeCode, v.ProtocolTypeCode, []string{"HTTP", "HTTPS", "TCP", "SSL"}); err != nil {
			return params, err
		}
		params[protocolTypeCode] = v.ProtocolTypeCode

		// LoadBalancerPort
		loadBalancerPort := fmt.Sprintf("loadBalancerRuleList.%d.loadBalancerPort", k+1)
		if err := validateRequiredField(loadBalancerPort, v.LoadBalancerPort); err != nil {
			return params, err
		}
		if err := validateIntegerInRange(loadBalancerPort, v.LoadBalancerPort, 1, 65534); err != nil {
			return nil, err
		}
		params[loadBalancerPort] = strconv.Itoa(v.LoadBalancerPort)

		// ServerPort
		serverPort := fmt.Sprintf("loadBalancerRuleList.%d.serverPort", k+1)
		if err := validateRequiredField(serverPort, v.ServerPort); err != nil {
			return params, err
		}
		if err := validateIntegerInRange(serverPort, v.ServerPort, 1, 65534); err != nil {
			return nil, err
		}
		params[serverPort] = strconv.Itoa(v.ServerPort)

		// L7HealthCheckPath
		l7HealthCheckPath := fmt.Sprintf("loadBalancerRuleList.%d.l7HealthCheckPath", k+1)
		if v.ProtocolTypeCode == "HTTP" || v.ProtocolTypeCode == "HTTPS" {
			if err := validateRequiredField(l7HealthCheckPath, v.L7HealthCheckPath); err != nil {
				return params, err
			}

			if err := validateStringLenBetween("L7HealthCheckPath", v.L7HealthCheckPath, 1, 600); err != nil {
				return nil, err
			}
			params[l7HealthCheckPath] = v.L7HealthCheckPath
		}

		// CertificateName
		certificateName := fmt.Sprintf("loadBalancerRuleList.%d.certificateName", k+1)
		if v.ProtocolTypeCode == "HTTPS" || v.ProtocolTypeCode == "SSL" {
			if err := validateRequiredField(certificateName, v.CertificateName); err != nil {
				return params, err
			}

			if err := validateStringLenBetween(certificateName, v.CertificateName, 1, 300); err != nil {
				return nil, err
			}
			params[certificateName] = v.CertificateName
		}

		// ProxyProtocolUseYn
		if v.ProxyProtocolUseYn != "" {
			proxyProtocolUseYn := fmt.Sprintf("loadBalancerRuleList.%d.proxyProtocolUseYn", k+1)
			if err := validateIncludeValues(proxyProtocolUseYn, v.ProxyProtocolUseYn, []string{"Y", "N"}); err != nil {
				return nil, err
			}
			params[proxyProtocolUseYn] = v.ProxyProtocolUseYn
		}
	}

	return params, nil
}

// ChangeLoadBalancerInstanceConfiguration change public ip instance list
func (s *Conn) ChangeLoadBalancerInstanceConfiguration(reqParams *RequestChangeLoadBalancerInstanceConfiguration) (*LoadBalancerInstanceList, error) {
	params, err := processChangeLoadBalancerInstanceConfigurationParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "changeLoadBalancerInstanceConfiguration"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"loadbalancer/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := LoadBalancerInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var LoadBalancerInstanceList = LoadBalancerInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &LoadBalancerInstanceList); err != nil {
		return nil, err
	}

	return &LoadBalancerInstanceList, nil
}

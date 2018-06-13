package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processChangeLoadBalancedServerInstancesParams(reqParams *RequestChangeLoadBalancedServerInstances) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil || reqParams.LoadBalancerInstanceNo == "" {
		return nil, errors.New("LoadBalancerInstanceNo field is required")
	}

	params["loadBalancerInstanceNo"] = reqParams.LoadBalancerInstanceNo

	if len(reqParams.ServerInstanceNoList) > 0 {
		for k, v := range reqParams.ServerInstanceNoList {
			params[fmt.Sprintf("serverInstanceNoList.%d", k+1)] = v
		}
	}

	return params, nil
}

// ChangeLoadBalancedServerInstances changes load balancer server instance
func (s *Conn) ChangeLoadBalancedServerInstances(reqParams *RequestChangeLoadBalancedServerInstances) (*LoadBalancerInstanceList, error) {
	params, err := processChangeLoadBalancedServerInstancesParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "changeLoadBalancedServerInstances"

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

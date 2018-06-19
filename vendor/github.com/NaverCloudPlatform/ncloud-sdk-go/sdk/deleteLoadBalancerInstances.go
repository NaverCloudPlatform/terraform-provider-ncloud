package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processDeleteLoadBalancerInstancesParams(reqParams *RequestDeleteLoadBalancerInstances) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return nil, errors.New("LoadBalancerInstanceNoList is required field")
	}

	if len(reqParams.LoadBalancerInstanceNoList) == 0 {
		return nil, errors.New("LoadBalancerInstanceNoList is required field")
	}

	for k, v := range reqParams.LoadBalancerInstanceNoList {
		params[fmt.Sprintf("loadBalancerInstanceNoList.%d", k+1)] = v
	}

	return params, nil
}

// DeleteLoadBalancerInstances delete load balancer instances
func (s *Conn) DeleteLoadBalancerInstances(reqParams *RequestDeleteLoadBalancerInstances) (*LoadBalancerInstanceList, error) {
	params, err := processDeleteLoadBalancerInstancesParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "deleteLoadBalancerInstances"

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

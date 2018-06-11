package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processGetLoadBalancerTargetServerInstanceListParams(reqParams *RequestGetLoadBalancerTargetServerInstanceList) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return params, nil
	}

	if reqParams.InternetLineTypeCode != "" {
		if err := validateIncludeValues("InternetLineTypeCode", reqParams.InternetLineTypeCode, []string{"PUBLC", "GLBL"}); err != nil {
			return nil, err
		}
		params["internetLineTypeCode"] = reqParams.InternetLineTypeCode
	}

	if reqParams.NetworkUsageTypeCode != "" {
		if err := validateIncludeValues("NetworkUsageTypeCode", reqParams.NetworkUsageTypeCode, []string{"PBLIP", "PRVT"}); err != nil {
			return nil, err
		}
		params["networkUsageTypeCode"] = reqParams.NetworkUsageTypeCode
	}

	if reqParams.RegionNo != "" {
		params["regionNo"] = reqParams.RegionNo
	}

	return params, nil
}

// GetLoadBalancerTargetServerInstanceList get load balancer target server instance list
func (s *Conn) GetLoadBalancerTargetServerInstanceList(reqParams *RequestGetLoadBalancerTargetServerInstanceList) (*ServerInstanceList, error) {
	params, err := processGetLoadBalancerTargetServerInstanceListParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "getLoadBalancerTargetServerInstanceList"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "GET", s.apiURL+"loadbalancer/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := ServerInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var serverInstanceList = ServerInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &serverInstanceList); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &serverInstanceList, nil
}

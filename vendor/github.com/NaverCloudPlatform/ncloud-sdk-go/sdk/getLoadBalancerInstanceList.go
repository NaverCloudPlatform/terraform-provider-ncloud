package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processGetLoadBalancerInstanceListParams(reqParams *RequestLoadBalancerInstanceList) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return params, nil
	}

	if len(reqParams.LoadBalancerInstanceNoList) > 0 {
		for k, v := range reqParams.LoadBalancerInstanceNoList {
			params[fmt.Sprintf("loadBalancerInstanceNoList.%d", k+1)] = v
		}
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

	if reqParams.PageNo != 0 {
		if err := validateIntegerInRange("PageNo", reqParams.PageNo, 0, 2147483647); err != nil {
			return nil, err
		}

		params["pageNo"] = strconv.Itoa(reqParams.PageNo)
	}

	if reqParams.PageSize != 0 {
		if err := validateIntegerInRange("PageSize", reqParams.PageSize, 0, 2147483647); err != nil {
			return nil, err
		}

		params["pageSize"] = strconv.Itoa(reqParams.PageSize)
	}

	if reqParams.SortedBy != "" {
		if err := validateIncludeValuesIgnoreCase("SortedBy", reqParams.SortedBy, []string{"loadBalancerName", "loadBalancerInstanceNo"}); err != nil {
			return nil, err
		}

		params["sortedBy"] = reqParams.SortedBy
	}

	if reqParams.SortingOrder != "" {
		if err := validateIncludeValuesIgnoreCase("SortingOrder", reqParams.SortingOrder, []string{"ascending", "descending"}); err != nil {
			return nil, err
		}

		params["sortingOrder"] = reqParams.SortingOrder
	}

	return params, nil
}

// GetLoadBalancerInstanceList get load balancer instance list
func (s *Conn) GetLoadBalancerInstanceList(reqParams *RequestLoadBalancerInstanceList) (*LoadBalancerInstanceList, error) {
	params, err := processGetLoadBalancerInstanceListParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "getLoadBalancerInstanceList"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "GET", s.apiURL+"loadbalancer/", params)
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

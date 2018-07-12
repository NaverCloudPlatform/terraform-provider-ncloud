package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processGetLaunchConfigurationListParams(reqParams *RequestGetLaunchConfigurationList) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return params, nil
	}

	if len(reqParams.LaunchConfigurationNameList) > 0 {
		for k, v := range reqParams.LaunchConfigurationNameList {
			params[fmt.Sprintf("launchConfigurationNameList.%d", k+1)] = v
		}
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
		if err := validateIncludeValuesIgnoreCase("SortedBy", reqParams.SortedBy, []string{"launchConfigurationName", "createDate"}); err != nil {
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

// GetLaunchConfigurationList get launch configuration list
func (s *Conn) GetLaunchConfigurationList(reqParams *RequestGetLaunchConfigurationList) (*LaunchConfigurationList, error) {
	params, err := processGetLaunchConfigurationListParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "getLaunchConfigurationList"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "GET", s.apiURL+"autoscaling/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := LaunchConfigurationList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var launchConfigurationList = LaunchConfigurationList{}
	if err := xml.Unmarshal([]byte(bytes), &launchConfigurationList); err != nil {
		return nil, err
	}

	return &launchConfigurationList, nil
}

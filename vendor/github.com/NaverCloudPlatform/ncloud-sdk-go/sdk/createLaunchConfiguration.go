package sdk

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processCreateLaunchConfigurationParams(reqParams *RequestCreateLaunchConfiguration) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams.LaunchConfigurationName != "" {
		if err := validateStringMaxLen("LaunchConfigurationName", reqParams.LaunchConfigurationName, 255); err != nil {
			return nil, err
		}
		params["launchConfigurationName"] = reqParams.LaunchConfigurationName
	}

	if reqParams.ServerImageProductCode == "" && reqParams.MemberServerImageNo == "" {
		return nil, errors.New("Required field is not specified. location : serverImageProductCode or memberServerImageNo")
	}

	if reqParams.ServerImageProductCode != "" && reqParams.MemberServerImageNo != "" {
		return nil, errors.New("Only one field is required. location : serverImageProductCode or memberServerImageNo")
	}

	if reqParams.ServerImageProductCode != "" {
		if err := validateStringMaxLen("ServerImageProductCode", reqParams.ServerImageProductCode, 20); err != nil {
			return nil, err
		}
		params["serverImageProductCode"] = reqParams.ServerImageProductCode
	}

	if reqParams.ServerProductCode != "" {
		if err := validateStringMaxLen("ServerProductCode", reqParams.ServerProductCode, 20); err != nil {
			return nil, err
		}
		params["serverProductCode"] = reqParams.ServerProductCode
	}

	if reqParams.MemberServerImageNo != "" {
		params["memberServerImageNo"] = reqParams.MemberServerImageNo
	}

	if len(reqParams.AccessControlGroupConfigurationNoList) > 0 {
		for k, v := range reqParams.AccessControlGroupConfigurationNoList {
			params[fmt.Sprintf("accessControlGroupConfigurationNoList.%d", k+1)] = v
		}
	}

	if reqParams.LoginKeyName != "" {
		if err := validateStringLenBetween("LoginKeyName", reqParams.LoginKeyName, 3, 30); err != nil {
			return nil, err
		}
		params["loginKeyName"] = reqParams.LoginKeyName
	}

	if reqParams.UserData != "" {
		if err := validateStringMaxLen("UserData", reqParams.UserData, 21847); err != nil {
			return nil, err
		}
		params["userData"] = base64.StdEncoding.EncodeToString([]byte(reqParams.UserData))
	}

	return params, nil
}

// CreateLaunchConfiguration creates launch configuration list
func (s *Conn) CreateLaunchConfiguration(reqParams *RequestCreateLaunchConfiguration) (*LaunchConfigurationList, error) {
	params, err := processCreateLaunchConfigurationParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "createLaunchConfiguration"

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

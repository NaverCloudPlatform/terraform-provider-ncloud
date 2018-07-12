package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processDeleteLaunchConfigurationParams(launchConfigurationName string) error {
	if err := validateRequiredField("LaunchConfigurationName", launchConfigurationName); err != nil {
		return err
	}

	if err := validateStringMaxLen("LaunchConfigurationName", launchConfigurationName, 255); err != nil {
		return err
	}

	return nil
}

// DeleteLaunchConfiguration delete Launch Configuration
func (s *Conn) DeleteLaunchConfiguration(launchConfigurationName string) (*common.CommonResponse, error) {
	if err := processDeleteLaunchConfigurationParams(launchConfigurationName); err != nil {
		return nil, err
	}

	params := make(map[string]string)

	params["launchConfigurationName"] = launchConfigurationName
	params["action"] = "deleteAutoScalingLaunchConfiguration"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"autoscaling/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := common.CommonResponse{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var responseDeleteLoginKey = common.CommonResponse{}
	if err := xml.Unmarshal([]byte(bytes), &responseDeleteLoginKey); err != nil {
		return nil, err
	}

	return &responseDeleteLoginKey, nil
}

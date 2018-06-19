package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processRecreateServerInstanceParams(reqParams *RequestRecreateServerInstance) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return nil, errors.New("ServerInstanceNo and ChangeServerImageProductCode are required fields")
	}
	if err := validateRequiredField("ServerInstanceNo", reqParams.ServerInstanceNo); err != nil {
		return nil, err
	}
	params["serverInstanceNo"] = reqParams.ServerInstanceNo

	if reqParams.ServerInstanceName != "" {
		if err := validateStringLenBetween("ServerInstanceName", reqParams.ServerInstanceName, 3, 30); err != nil {
			return nil, err
		}
		params["serverInstanceName"] = reqParams.ServerInstanceName
	}

	if err := validateRequiredField("ChangeServerImageProductCode", reqParams.ChangeServerImageProductCode); err != nil {
		return nil, err
	}
	if err := validateStringMaxLen("ChangeServerImageProductCode", reqParams.ChangeServerImageProductCode, 20); err != nil {
		return nil, err
	}
	params["changeServerImageProductCode"] = reqParams.ChangeServerImageProductCode

	return params, nil
}

// RecreateServerInstance recreate server instance
func (s *Conn) RecreateServerInstance(reqParams *RequestRecreateServerInstance) (*ServerInstanceList, error) {
	params, err := processRecreateServerInstanceParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "recreateServerInstance"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"server/", params)
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

	var responseRecreateServerInstance = ServerInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &responseRecreateServerInstance); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &responseRecreateServerInstance, nil
}

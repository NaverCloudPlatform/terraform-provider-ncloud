package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processAttachBlockStorageInstanceParams(reqParams *RequestAttachBlockStorageInstance) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return nil, errors.New("ServerInstanceNo and BlockStorageInstanceNo are required fields")
	}

	if err := validateRequiredField("ServerInstanceNo", reqParams.ServerInstanceNo); err != nil {
		return nil, err
	}
	params["serverInstanceNo"] = reqParams.ServerInstanceNo

	if err := validateRequiredField("BlockStorageInstanceNo", reqParams.BlockStorageInstanceNo); err != nil {
		return nil, err
	}
	params["blockStorageInstanceNo"] = reqParams.BlockStorageInstanceNo

	return params, nil
}

// AttachBlockStorageInstance attaches block storage instance to a server instance
func (s *Conn) AttachBlockStorageInstance(reqParams *RequestAttachBlockStorageInstance) (*BlockStorageInstanceList, error) {
	params, err := processAttachBlockStorageInstanceParams(reqParams)

	if err != nil {
		return nil, err
	}

	params["action"] = "attachBlockStorageInstance"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"server/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := BlockStorageInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var list = BlockStorageInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &list); err != nil {
		return nil, err
	}

	return &list, nil
}

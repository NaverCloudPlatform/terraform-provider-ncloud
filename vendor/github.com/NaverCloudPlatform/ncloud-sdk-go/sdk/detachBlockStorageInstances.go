package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processDetachBlockStorageInstanceParams(reqParams *RequestDetachBlockStorageInstance) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil || len(reqParams.BlockStorageInstanceNoList) == 0 {
		return nil, errors.New("BlockStorageInstanceNoList is required field")
	}

	if len(reqParams.BlockStorageInstanceNoList) > 0 {
		for k, v := range reqParams.BlockStorageInstanceNoList {
			params[fmt.Sprintf("blockStorageInstanceNoList.%d", k+1)] = v
		}
	}

	return params, nil
}

// DetachBlockStorageInstance detaches block storage instance from server instance
func (s *Conn) DetachBlockStorageInstance(reqParams *RequestDetachBlockStorageInstance) (*BlockStorageInstanceList, error) {
	params, err := processDetachBlockStorageInstanceParams(reqParams)

	if err != nil {
		return nil, err
	}

	params["action"] = "detachBlockStorageInstances"

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

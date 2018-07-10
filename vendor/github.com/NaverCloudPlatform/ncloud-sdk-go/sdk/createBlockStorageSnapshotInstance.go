package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processCreateBlockStorageSnapshotInstanceParams(reqParams *RequestCreateBlockStorageSnapshotInstance) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil || reqParams.BlockStorageInstanceNo == "" {
		return nil, errors.New("BlockStorageInstanceNo field is required")
	}
	params["blockStorageInstanceNo"] = reqParams.BlockStorageInstanceNo

	if reqParams.BlockStorageSnapshotName != "" {
		params["blockStorageSnapshotName"] = reqParams.BlockStorageSnapshotName
	}

	if reqParams.BlockStorageSnapshotDescription != "" {
		params["blockStorageSnapshotDescription"] = reqParams.BlockStorageSnapshotDescription
	}

	return params, nil
}

// CreateBlockStorageSnapshotInstance create block storage snapshot instance
func (s *Conn) CreateBlockStorageSnapshotInstance(reqParams *RequestCreateBlockStorageSnapshotInstance) (*BlockStorageSnapshotInstanceList, error) {

	params, err := processCreateBlockStorageSnapshotInstanceParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "createBlockStorageSnapshotInstance"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"server/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := BlockStorageSnapshotInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var blockStorageSnapshotInstanceList = BlockStorageSnapshotInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &blockStorageSnapshotInstanceList); err != nil {
		return nil, err
	}

	return &blockStorageSnapshotInstanceList, nil
}

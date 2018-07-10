package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processDeleteBlockStorageSnapshotInstancesParams(blockStorageSnapshotInstanceNoList []string) (map[string]string, error) {
	params := make(map[string]string)

	if len(blockStorageSnapshotInstanceNoList) == 0 {
		return params, errors.New("BlockStorageSnapshotInstanceNoList is required field")
	}

	for k, v := range blockStorageSnapshotInstanceNoList {
		params[fmt.Sprintf("blockStorageSnapshotInstanceNoList.%d", k+1)] = v
	}

	return params, nil
}

// DeleteBlockStorageSnapshotInstances delete block storage snapshot instances
func (s *Conn) DeleteBlockStorageSnapshotInstances(blockStorageSnapshotInstanceNoList []string) (*BlockStorageSnapshotInstanceList, error) {
	params, err := processDeleteBlockStorageSnapshotInstancesParams(blockStorageSnapshotInstanceNoList)
	if err != nil {
		return nil, err
	}

	params["action"] = "deleteBlockStorageSnapshotInstances"

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
		fmt.Println(err)
		return nil, err
	}

	return &blockStorageSnapshotInstanceList, nil
}

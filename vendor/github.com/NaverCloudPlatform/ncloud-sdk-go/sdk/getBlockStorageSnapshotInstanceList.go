package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processGetBlockStorageSnapshotInstanceListParams(reqParams *RequestGetBlockStorageSnapshotInstanceList) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return params, nil
	}

	if len(reqParams.BlockStorageSnapshotInstanceNoList) > 0 {
		for k, v := range reqParams.BlockStorageSnapshotInstanceNoList {
			params[fmt.Sprintf("blockStorageSnapshotInstanceNoList.%d", k+1)] = v
		}
	}

	if len(reqParams.OriginalBlockStorageInstanceNoList) > 0 {
		for k, v := range reqParams.OriginalBlockStorageInstanceNoList {
			params[fmt.Sprintf("originalBlockStorageInstanceNoList.%d", k+1)] = v
		}
	}

	if reqParams.RegionNo != "" {
		params["regionNo"] = reqParams.RegionNo
	}

	if reqParams.PageNo != 0 {
		if reqParams.PageNo > 2147483647 {
			return nil, errors.New("PageNo should be up to 2147483647")
		}

		params["pageNo"] = strconv.Itoa(reqParams.PageNo)
	}

	if reqParams.PageSize != 0 {
		if reqParams.PageSize > 2147483647 {
			return nil, errors.New("PageSize should be up to 2147483647")
		}

		params["pageSize"] = strconv.Itoa(reqParams.PageSize)
	}

	return params, nil
}

// GetBlockStorageSnapshotInstanceList Get block storage snapshot instance list
func (s *Conn) GetBlockStorageSnapshotInstanceList(reqParams *RequestGetBlockStorageSnapshotInstanceList) (*BlockStorageSnapshotInstanceList, error) {
	params, err := processGetBlockStorageSnapshotInstanceListParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "getBlockStorageSnapshotInstanceList"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "GET", s.apiURL+"server/", params)
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

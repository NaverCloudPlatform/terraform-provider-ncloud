package request

import (
	"io/ioutil"
	"net/http"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/oauth"
)

// NewRequest is http request with oauth
func NewRequest(accessKey string, secretKey string, method string, url string, params map[string]string) ([]byte, *http.Response, error) {
	c := oauth.NewConsumer(accessKey, secretKey, method, url)

	for k, v := range params {
		c.AdditionalParams[k] = v
	}

	reqURL, body, err := c.GetRequest()
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, nil, err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	return bytes, resp, nil
}

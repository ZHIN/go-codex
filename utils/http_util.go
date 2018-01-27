package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	nurl "net/url"
	"strings"
)

type httpUtil struct {
	client *http.Client
}

func (c *httpUtil) GetJSON(url string, queryParams map[string]string, responseData interface{}, headers map[string]string) error {

	u, err := nurl.Parse(url)
	if err != nil {
		return err
	}
	params := u.Query()
	if queryParams != nil {
		for key, value := range queryParams {
			params.Set(key, value)
		}
	}
	u.RawQuery = params.Encode()

	return c.RequestJSON(http.MethodGet, u.String(), nil, responseData, headers)
}

func (c *httpUtil) PostJSON(url string, formParams map[string]interface{}, responseData interface{}, headers map[string]string) error {
	return c.RequestJSON(http.MethodPost, url, formParams, responseData, headers)
}

type httpRequestEncodeType int

const (
	UrlEncoded  httpRequestEncodeType = 1
	JSONEncoded httpRequestEncodeType = 2
)

func (c *httpUtil) RequestJSON(method string, url string, formParams map[string]interface{}, responseData interface{}, headers map[string]string) error {

	var byteBuff *bytes.Buffer
	var err error

	encodeType := UrlEncoded

	if headers != nil {
		if contentType := headers["Content-Type"]; strings.HasPrefix(contentType,
			"application/x-www-form-urlencoded") {
			encodeType = UrlEncoded
		} else {
			encodeType = JSONEncoded
		}
	}

	if formParams != nil {
		if encodeType == UrlEncoded {
			urlValues := nurl.Values{}

			for key, value := range formParams {
				urlValues.Add(key, fmt.Sprintf("%v", value))
			}
			byteBuff = bytes.NewBuffer([]byte(urlValues.Encode()))
		} else if encodeType == JSONEncoded {
			byteBuff, err = JSONUtil.MapToByteBuffer(formParams)
			if err != nil {
				return err
			}
		}
	}

	var req *http.Request

	log.Println(method, url, byteBuff)
	if byteBuff == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, byteBuff)
	}
	if err != nil {
		return err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}
	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(responseData)

	return err
}

var HttpUtil = httpUtil{
	client: http.DefaultClient,
}

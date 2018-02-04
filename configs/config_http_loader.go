package configs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func requestEncryptTOML(url string, params map[string]interface{}, headers map[string]string) ([]byte, error) {

	var byteBuff *bytes.Buffer
	var err error

	if params != nil {
		buff, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		byteBuff = bytes.NewBuffer(buff)
		if err != nil {
			return nil, err
		}
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, byteBuff)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buff, err := ioutil.ReadAll(resp.Body)
	return buff, err
}

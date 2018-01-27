package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type jsonUtil struct {
}

func (*jsonUtil) LintJSON(jsonText string) string {

	var data map[string]interface{}
	err := json.NewDecoder(bytes.NewBuffer([]byte(jsonText))).Decode(&data)
	if err != nil {
		log.Fatalln(err)
	}

	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	err = encoder.Encode(data)
	return string(b.Bytes())
}

func (*jsonUtil) LintJSONFromData(data interface{}) string {

	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	encoder.Encode(data)
	return string(b.Bytes())
}

func (*jsonUtil) ParseStreamToMap(r io.Reader) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

func (c *jsonUtil) MapToString(data map[string]interface{}) (string, error) {

	buff, err := c.MapToBytes(data)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func (c *jsonUtil) MapToByteBuffer(data map[string]interface{}) (*bytes.Buffer, error) {

	buff, err := c.MapToBytes(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buff), nil
}

func (*jsonUtil) MapToBytes(data map[string]interface{}) ([]byte, error) {

	buff, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

var JSONUtil = jsonUtil{}

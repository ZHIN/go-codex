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

func (*jsonUtil) DataToJSONString(data interface{}) (string, error) {
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	err := encoder.Encode(data)
	return string(b.Bytes()), err
}

func (*jsonUtil) LintJSONFromData(data interface{}) string {

	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	encoder.Encode(data)
	return string(b.Bytes())
}

func (j *jsonUtil) ParseStreamToMap(r io.Reader) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := j.ParseFromStream(r, &data)
	return data, err
}

func (*jsonUtil) ParseFromStream(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	err := decoder.Decode(v)
	return err
}
func (*jsonUtil) ParseFromBytes(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err := decoder.Decode(v)
	return err
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

func (c *jsonUtil) MustStringify(data interface{}) string {
	text, e := c.DataToJSONString(data)
	if e != nil {
		return ""
	}
	return text
}

var JSONUtil = jsonUtil{}

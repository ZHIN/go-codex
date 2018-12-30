package exthttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	nurl "net/url"
	"strings"
	"time"

	"github.com/zhin/go-codex/configs"
)

type HttpClient struct {
	client         *http.Client
	hookBeforeSend func(*http.Request)
}

func (c *HttpClient) GetJSON(url string, queryParams map[string]string, responseData interface{}, options *RequestOptions) error {

	return c.RequestJSON(http.MethodGet, url, queryParams, nil, responseData, options)
}

func (c *HttpClient) PostJSON(url string, formParams map[string]interface{}, responseData interface{}, options *RequestOptions) error {
	return c.RequestJSON(http.MethodPost, url, nil, formParams, responseData, options)
}

func (c *HttpClient) RequestJSON(method string, url string, queryParams map[string]string, formParams map[string]interface{}, responseData interface{}, options *RequestOptions) error {

	var err error
	responseBuff, err := c.Request(method, url, queryParams, formParams, options)

	err = json.NewDecoder(bytes.NewBuffer(responseBuff)).Decode(responseData)

	if err != nil && configs.Settings.GetBool(httpDebugErrorJSON) {
		log.Println(fmt.Sprintf("Error:%s", err.Error()))
		log.Println(fmt.Sprintf("url:%s\ncontent:%s", url, string(responseBuff)))
	}
	return err
}

func (c *HttpClient) Request(method string, url string, queryParams map[string]string, formParams map[string]interface{}, options *RequestOptions) ([]byte, error) {
	var byteBuff *bytes.Buffer
	var err error
	encodeType := JSONEncoded

	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}
	params := u.Query()
	if queryParams != nil {
		for key, value := range queryParams {
			params.Set(key, value)
		}
	}
	u.RawQuery = params.Encode()
	url = u.String()

	if options != nil && options.Headers != nil {
		if contentType := options.Headers["Content-Type"]; strings.HasPrefix(strings.ToLower(contentType),
			"application/x-www-form-urlencoded") {
			encodeType = URLEncoded
		} else if contentType := options.Headers["Content-Type"]; strings.HasPrefix(strings.ToLower(contentType),
			"application/json") {
			encodeType = JSONEncoded
		} else {
			return nil, fmt.Errorf("unknow http content-type \"%s\"", contentType)
		}
	}

	if formParams != nil {
		if encodeType == URLEncoded {
			urlValues := nurl.Values{}
			for key, value := range formParams {
				urlValues.Add(key, fmt.Sprintf("%v", value))
			}
			byteBuff = bytes.NewBuffer([]byte(urlValues.Encode()))
		} else if encodeType == JSONEncoded {
			byteBuff, err = mapToByteBuffer(formParams)
			if err != nil {
				return nil, err
			}
		}
	}

	var req *http.Request

	if byteBuff == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, byteBuff)
	}
	if err != nil {
		return nil, err
	}

	if options != nil && options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Add(key, value)
		}
	}

	if c.hookBeforeSend != nil {
		c.hookBeforeSend(req)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseBuff []byte

	responseBuff, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBuff, nil
}

func (c *HttpClient) SetProxy(host string, port int) {
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%d", host, port))
	if err != nil {
		// log
	}
	transport := c.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxyURL)
}

func (c *HttpClient) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

func (c *HttpClient) SetBeforeSendHook(func(r *http.Request)) {

}
func (c *HttpClient) UseCookieJar(use bool) {
	if use {
		if c.client.Jar == nil {
			jar, _ := cookiejar.New(nil)
			c.client.Jar = jar
		}
	} else {
		c.client.Jar = nil
	}
}

func newClient(option ClientOption) *http.Client {

	transport := &http.Transport{}

	if option.Proxy != "" {
		proxyURL, err := url.Parse(option.Proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		} else {
			// TODO: LOG ERROR
		}
	}
	client := &http.Client{Transport: transport, Timeout: time.Duration(option.TimeoutSecond) * time.Second}
	return client
}

func NewHttpClient(option ClientOption) *HttpClient {
	proxyValue := configs.Settings.GetString(httpProxySettingKey)
	transport := &http.Transport{}

	if proxyValue != "" {
		proxyURL, err := url.Parse(proxyValue)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		} else {
			// TODO: LOG ERROR
		}
	}

	client := &HttpClient{
		client: newClient(ClientOption{
			Proxy:         proxyValue,
			TimeoutSecond: configs.Settings.GetInt(httpTimeoutSecondSettingKey),
		}),
	}

	if option.UseCookieJar {
		client.UseCookieJar(true)
	}
	return client
}

type ClientOption struct {
	TimeoutSecond int
	Proxy         string
	UseCookieJar  bool
}

var DefaultClient *HttpClient

const (
	httpTimeoutSecondSettingKey = "http_timeout_seoncd"
	httpProxySettingKey         = "http_proxy"
	httpUseCookieJarSettingKey  = "http_use_cookie_jar"
	httpDebugErrorJSON          = "http_debug_error_json"
)

type httpRequestEncodeType int

const (
	URLEncoded  httpRequestEncodeType = 1
	JSONEncoded httpRequestEncodeType = 2
)

func init() {
	configs.Settings.SetDefault(httpTimeoutSecondSettingKey, 90)
	configs.Settings.SetDefault(httpUseCookieJarSettingKey, false)
	configs.Settings.SetDefault(httpDebugErrorJSON, true)

	DefaultClient = NewHttpClient(ClientOption{
		Proxy:         configs.Settings.GetString(httpProxySettingKey),
		TimeoutSecond: configs.Settings.GetInt(httpTimeoutSecondSettingKey),
		UseCookieJar:  configs.Settings.GetBool(httpUseCookieJarSettingKey),
	})
}

type RequestOptions struct {
	Headers map[string]string
}

func mapToByteBuffer(data map[string]interface{}) (*bytes.Buffer, error) {

	buff, err := mapToBytes(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buff), nil
}

func mapToBytes(data map[string]interface{}) ([]byte, error) {

	buff, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

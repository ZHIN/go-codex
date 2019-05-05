package exthttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	httpcookiejar "net/http/cookiejar"
	"net/url"
	nurl "net/url"
	"strings"
	"time"

	cookiejar "github.com/juju/persistent-cookiejar"

	"github.com/zhin/go-codex/configs"
)

type HttpClient struct {
	client         *http.Client
	hookBeforeSend func(*http.Request)
}

type HttpResponseError struct {
	RequestURL     string
	Status         int
	ResponseHeader http.Header
	ResponseData   []byte
}

func (s *HttpResponseError) Error() string {
	if s.ResponseData != nil {
		return fmt.Sprintf("http request error:%s status:%d response:%s", s.RequestURL, s.Status, string(s.ResponseData))

	}
	return fmt.Sprintf("http request error:%s status:%d", s.RequestURL, s.Status)
}

func IsHttpResponseError(err error) bool {
	_, ok := err.(*HttpResponseError)
	return ok
}

func (c *HttpClient) GetJSON(url string, queryParams map[string]string, responseData interface{}, options *RequestOptions) error {

	return c.RequestJSON(http.MethodGet, url, queryParams, nil, responseData, options)
}

func (c *HttpClient) PostJSON(url string, formParams map[string]interface{}, responseData interface{}, options *RequestOptions) error {
	return c.RequestJSON(http.MethodPost, url, nil, formParams, responseData, options)
}

func (c *HttpClient) RequestJSON(method string, url string, queryParams map[string]string, formParams map[string]interface{}, responseData interface{}, options *RequestOptions) error {

	var err error

	if options == nil {
		options = &RequestOptions{}
	}
	if options.Headers["Content-Type"] == "" && options.ContentType == 0 {
		options.ContentType = JSONEncoded
	}

	responseBuff, err := c.Request(method, url, queryParams, formParams, options)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewBuffer(responseBuff))
	decoder.UseNumber()
	err = decoder.Decode(responseData)
	if err != nil && configs.Settings.GetBool(httpDebugErrorJSON) {
		log.Println(fmt.Sprintf("Error:%s", err.Error()))
		log.Println(fmt.Sprintf("url:%s\ncontent:%s", url, string(responseBuff)))
	}
	return err
}

func (c *HttpClient) RawRequest(method string, url string, queryParams map[string]string, body []byte, options *RequestOptions) ([]byte, error) {

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

	if options == nil {
		options = &RequestOptions{}
	}
	if options.Headers == nil {
		options.Headers = map[string]string{}
	}

	if method == "POST" {

		if options.ContentType == 0 {
			if contentType := options.Headers["Content-Type"]; strings.HasPrefix(strings.ToLower(contentType),
				"application/x-www-form-urlencoded") {
				encodeType = URLEncoded
			} else if contentType := options.Headers["Content-Type"]; strings.HasPrefix(strings.ToLower(contentType),
				"application/json") {
				encodeType = JSONEncoded
			} else {
				return nil, fmt.Errorf("unknow http content-type \"%s\"", contentType)
			}
		} else if options.ContentType != 0 {
			encodeType = options.ContentType
			if encodeType == URLEncoded {
				options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
			} else if encodeType == JSONEncoded {
				options.Headers["Content-Type"] = "application/json"
			}
		}
		byteBuff = bytes.NewBuffer(body)
	}

	var req *http.Request

	if byteBuff == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, byteBuff)
	}
	if err != nil {
		return nil, fmt.Errorf("request content error:%s", err.Error())
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
		return nil, fmt.Errorf("request error:%s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var buff []byte
		buff, _ = ioutil.ReadAll(resp.Body)
		return nil, &HttpResponseError{
			Status:         resp.StatusCode,
			ResponseHeader: resp.Header,
			ResponseData:   buff,
		}
	}

	var responseBuff []byte
	responseBuff, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response error:%s", err.Error())
	}

	if options != nil && options.ResponseHeaders != nil {
		options.ResponseHeaders = resp.Header
	}

	return responseBuff, nil

}

func (c *HttpClient) Request(method string, url string, queryParams map[string]string, formParams map[string]interface{}, options *RequestOptions) ([]byte, error) {
	var byteBuff *bytes.Buffer
	var err error
	encodeType := JSONEncoded
	if method == "POST" {
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
					return nil, fmt.Errorf("request json content error:%s", err.Error())
				}
			}
		}
	}
	if byteBuff != nil {
		return c.RawRequest(method, url, queryParams, byteBuff.Bytes(), options)
	}

	return c.RawRequest(method, url, queryParams, nil, options)

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

type CookieJarOption struct {
	Filename string
}

func (c *HttpClient) UseCookieJar(use bool, option *CookieJarOption) {
	if use {
		if c.client.Jar == nil {

			opt := &cookiejar.Options{}

			if option == nil {
			} else {
				opt.Filename = option.Filename
			}

			if opt.Filename == "" {
				opt.NoPersist = true
				// opt.Filename = path.Join(os.TempDir(), fmt.Sprintf("gocodex_tmp_cookie_jar_%s", randString(40)))
			}

			jar, _ := cookiejar.New(opt)
			c.client.Jar = jar
		}
	} else {
		c.client.Jar = nil
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (c *HttpClient) Jar() *httpcookiejar.Jar {
	return c.Jar()
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

	return client
}

type ClientOption struct {
	TimeoutSecond int
	Proxy         string
}

var DefaultClient *HttpClient

const (
	httpTimeoutSecondSettingKey = "http_timeout_seoncd"
	httpProxySettingKey         = "http_proxy"
	httpUseCookieJarSettingKey  = "http_use_cookie_jar"
	httpDebugErrorJSON          = "http_debug_error_json"
)

type HttpRequestEncodeType int

const (
	URLEncoded  HttpRequestEncodeType = 1
	JSONEncoded HttpRequestEncodeType = 2
)

func init() {
	configs.Settings.SetDefault(httpTimeoutSecondSettingKey, 90)
	configs.Settings.SetDefault(httpUseCookieJarSettingKey, false)
	configs.Settings.SetDefault(httpDebugErrorJSON, true)

	DefaultClient = NewHttpClient(ClientOption{
		Proxy:         configs.Settings.GetString(httpProxySettingKey),
		TimeoutSecond: configs.Settings.GetInt(httpTimeoutSecondSettingKey),
	})

	DefaultClient.UseCookieJar(configs.Settings.GetBool(httpUseCookieJarSettingKey), nil)
}

type RequestOptions struct {
	Headers         map[string]string
	ContentType     HttpRequestEncodeType
	ResponseHeaders http.Header
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

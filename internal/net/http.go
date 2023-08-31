package net

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type HttpClient struct {
	httpClient  *http.Client
	method      string
	params      url.Values // 查询参数
	formData    url.Values // 表单数据
	body        io.Reader
	uri         string
	header      map[string]string
	cookies     []*http.Cookie
	contentType string // 指定Content-Type
	userAgent   string // 指定User-Agent
	debug       bool
}

func NewHttpClient(client *http.Client) *HttpClient {
	return &HttpClient{
		httpClient: client,
		method:     http.MethodGet,
		params:     make(url.Values),
		formData:   make(url.Values),
		header:     make(map[string]string),
		cookies:    make([]*http.Cookie, 0),
	}
}

func (c *HttpClient) Clone() *HttpClient {
	clonedParams := make(url.Values)
	for key, values := range c.params {
		clonedParams[key] = append([]string(nil), values...)
	}

	clonedFormData := make(url.Values)
	for key, values := range c.formData {
		clonedFormData[key] = append([]string(nil), values...)
	}

	clonedHeaders := make(map[string]string)
	for key, value := range c.header {
		clonedHeaders[key] = value
	}

	clonedCookies := make([]*http.Cookie, len(c.cookies))
	copy(clonedCookies, c.cookies)

	return &HttpClient{
		httpClient:  c.httpClient,
		method:      c.method,
		params:      clonedParams,
		formData:    clonedFormData,
		body:        c.body,
		uri:         c.uri,
		header:      clonedHeaders,
		cookies:     clonedCookies,
		contentType: c.contentType,
		userAgent:   c.userAgent,
		debug:       c.debug,
	}
}

func (c *HttpClient) Get(uri string) *HttpClient {
	c.method = http.MethodGet
	c.uri = uri

	return c
}

func (c *HttpClient) Post(uri string) *HttpClient {
	c.method = http.MethodPost
	c.uri = uri

	return c
}

func (c *HttpClient) Put(uri string) *HttpClient {
	c.method = http.MethodPut
	c.uri = uri

	return c
}

func (c *HttpClient) SetContentType(contentType string) *HttpClient {
	c.contentType = contentType

	return c
}

func (c *HttpClient) SetUserAgent(ua string) *HttpClient {
	c.userAgent = ua

	return c
}

func (c *HttpClient) AddParams(key string, value string) *HttpClient {
	c.params.Add(key, value)

	return c
}

func (c *HttpClient) AddFormData(key string, value string) *HttpClient {
	c.formData.Add(key, value)

	return c
}

func (c *HttpClient) SetParams(key string, value string) *HttpClient {
	c.params.Set(key, value)

	return c
}

func (c *HttpClient) SetHeader(key string, value string) *HttpClient {
	c.header[key] = value

	return c
}

func (c *HttpClient) SendBody(body io.Reader) *HttpClient {
	c.body = body

	return c
}

func (c *HttpClient) SetCookies(cookies []*http.Cookie) *HttpClient {
	c.cookies = cookies

	return c
}

func (c *HttpClient) Debug() *HttpClient {
	c.debug = true

	return c
}

func (c *HttpClient) End() (resp *http.Response, body []byte, err error) {
	if len(c.formData) > 0 {
		c.body = strings.NewReader(c.formData.Encode())
		c.contentType = "application/x-www-form-urlencoded"
	}

	request, err := http.NewRequest(c.method, c.uri, c.body)
	if err != nil {
		return nil, nil, err
	}

	if c.params != nil {
		request.URL.RawQuery = c.params.Encode()
	}

	for k, v := range c.header {
		request.Header.Set(k, v)
	}

	for _, cookie := range c.cookies {
		request.AddCookie(cookie)
	}

	if c.contentType != "" {
		request.Header.Add("Content-Type", c.contentType)
	}

	if c.userAgent != "" {
		request.Header.Add("User-Agent", c.userAgent)
	}

	if c.debug {
		bts, _ := httputil.DumpRequest(request, true)
		fmt.Printf("\n[REQUEST]\n%s\n", string(bts))
	}

	resp, err = c.httpClient.Do(request)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if c.debug {
		bts, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("\n[RESPONSE]\n%s\n", string(bts))
	}

	body, err = io.ReadAll(resp.Body)

	return
}

func (c *HttpClient) EndStruct(value any, callbacks ...func(*http.Response) error) error {
	resp, body, err := c.End()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &value); err != nil {
		return err
	}

	for _, f := range callbacks {
		if err := f(resp); err != nil {
			return err
		}
	}

	return nil
}
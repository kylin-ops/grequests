package grequests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

var version = "0.1"

type Header map[string]string
type Data map[string]interface{}
type Param map[string]string

type RequestOptions struct {
	Header  Header
	Data    Data
	Json    bool
	Params  Param
	Timeout time.Duration
}

type Response struct {
	Response *http.Response
}

func request(url, method string, options ...*RequestOptions) (resp *Response, err error) {
	var option = RequestOptions{}
	if len(options) > 0 {
		option = *options[0]
	}
	var r *http.Request
	var response Response
	var params []string
	client := http.Client{Timeout: option.Timeout}
	data, _ := json.Marshal(option.Data)
	for k, v := range option.Params {
		params = append(params, k+"="+v)
	}
	if p := strings.Join(params, "&"); p != "" {
		url = url + "?" + p
	}
	r, err = http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return resp, err
	}
	r.Header.Set("User-Agent", "go-request"+version)
	for k, v := range option.Header {
		r.Header.Set(k, v)
	}
	if option.Json {
		r.Header.Set("Content-Type", "application/json")
	}
	response.Response, err = client.Do(r)
	return &response, err
}

func Get(url string, options ...*RequestOptions) (resp *Response, err error) {
	return request(url, "GET", options...)
}

func Post(url string, options ...*RequestOptions) (resp *Response, err error) {
	return request(url, "POST", options...)
}

func Put(url string, options ...*RequestOptions) (resp *Response, err error) {
	return request(url, "PUT", options...)
}

func Patch(url string, options ...*RequestOptions) (resp *Response, err error) {
	return request(url, "PATCH", options...)
}

func Delete(url string, options ...*RequestOptions) (resp *Response, err error) {
	return request(url, "DELETE", options...)
}

func Head(url string, options ...*RequestOptions) (resp *Response, err error) {
	return request(url, "HEAD", options...)
}

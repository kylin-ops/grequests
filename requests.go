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
	Url     string
	Header  Header
	Data    Data
	Json    bool
	Params  Param
	Timeout time.Duration
}

type Response struct {
	Response *http.Response
}

func request(method string, req *RequestOptions) (resp *Response, err error) {
	var r *http.Request
	var response Response
	var params []string
	client := http.Client{Timeout: req.Timeout}
	data, _ := json.Marshal(req.Data)
	for k, v := range req.Params {
		params = append(params, k+"="+v)
	}
	if p := strings.Join(params, "&"); p != "" {
		req.Url = req.Url + "?" + p
	}
	r, err = http.NewRequest(method, req.Url, bytes.NewReader(data))
	if err != nil {
		return resp, err
	}
	r.Header.Set("User-Agent", "go-request"+version)
	for k, v := range req.Header {
		r.Header.Set(k, v)
	}
	if req.Json {
		r.Header.Set("Content-Type", "application/json")
	}
	response.Response, err = client.Do(r)
	return &response, err
}

func Get(options *RequestOptions) (resp *Response, err error) {
	return request("GET", options)
}

func Post(options *RequestOptions) (resp *Response, err error) {
	return request("POST", options)
}

func Put(options *RequestOptions) (resp *Response, err error) {
	return request("PUT", options)
}

func Patch(options *RequestOptions) (resp *Response, err error) {
	return request("PATCH", options)
}

func Delete(options *RequestOptions) (resp *Response, err error) {
	return request("DELETE", options)
}

func Head(options *RequestOptions) (resp *Response, err error) {
	return request("HEAD", options)
}

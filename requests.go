package grequests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var version = "0.1"

type Header map[string]string
type Data map[string]interface{}
type Param map[string]string

type baseAuth struct {
	UserName string
	Password string
}

type RequestOptions struct {
	Header   Header
	Data     Data
	Json     bool
	Params   Param
	Timeout  time.Duration
	BashAuth baseAuth
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
	if option.BashAuth.UserName != "" && option.BashAuth.Password != "" {
		r.SetBasicAuth(option.BashAuth.UserName, option.BashAuth.Password)
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

func destValidator(dest, filename string) (destFile string) {
	if dest == "" {
		dest = filename
	} else {
		f, err := os.Stat(dest)
		if err == nil {
			if f.IsDir() {
				fmt.Println("a")
				dest = path.Join(dest, filename)
			} else {
				dest = dest + strconv.Itoa(int(time.Now().Unix()))
			}
		}
	}
	return dest
}

// dest=""时，文件下载在本地
func DownloadFile(url, dest string) error {
	var buf = make([]byte, 32*1024)
	var written int
	fileName := path.Base(url)
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.Body == nil {
		return errors.New("下载的类容为nil")
	}
	fsize, err := strconv.ParseInt(req.Header.Get("Content-Length"), 10, 32)
	destFile := destValidator(dest, fileName)
	destFileTemp := destFile + ".tmp"
	f, err := os.Create(destFileTemp)
	if err != nil {
		return nil
	}
	defer func() {
		f.Sync()
		f.Close()
		err = os.Rename(destFileTemp, destFile)
	}()
	for {
		nr, err := req.Body.Read(buf)
		written += nr
		if err != nil {
			if err == io.EOF {
				if fsize != int64(written) {
					return errors.New("下载文件大小不一致")
				}
				return nil
			}
			return err
		}
		if nr > 0 {
			if _, er := f.Write(buf[0:nr]); er != nil {
				return er
			}
		}
	}
}

type Response struct {
	Response *http.Response
}

func (r *Response) Text() (string, error) {
	d, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func (r *Response) Json(data interface{}) error {
	d, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, data)
}

func (r *Response) Close() error {
	return r.Response.Body.Close()
}

func (r *Response) StatusCode() int {
	return r.Response.StatusCode
}

func (r *Response) Header() http.Header {
	return r.Response.Header
}

func (r *Response) Proto() string {
	return r.Response.Proto
}

func (r *Response) Body() io.Reader {
	return r.Response.Body
}

func (r *Response) Request() *http.Request {
	return r.Response.Request
}

package ziva

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/zituocn/ziva/logx"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	defaultContentType = "text/html; charset=utf-8"
	defaultUserAgent   = "Go-http-client/ziva/1.0"
)

func HttpGet(url string, vs ...interface{}) *Context {
	return DoRequest(url, "GET", vs...)
}

func HttpPost(url string, vs ...interface{}) *Context {
	return DoRequest(url, "POST", vs...)
}

func HttpPut(url string, vs ...interface{}) *Context {
	return DoRequest(url, "PUT", vs...)
}

func DoTask(task *Task) *Context {
	return DoRequest(task.Url, task.Method, task.Header, task.FormData, task.Payload, task)
}

func DoRequest(url, method string, vs ...interface{}) *Context {
	ctx, err := NewRequest(url, method, vs...)
	if err != nil {
		logx.Errorf("DoRequest error : %s ", err.Error())
		return nil
	}
	return ctx
}

func NewRequest(urlStr, method string, vs ...interface{}) (*Context, error) {
	u, errU := validUrl(urlStr)
	if errU != nil {
		return nil, errU
	}
	var task *Task
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range vs {
		switch vv := v.(type) {
		case *Task:
			task = vv
		case FormData:
			if len(vv) > 0 {
				formData := url.Values{}
				for k, v := range vv {
					formData.Set(k, v)
				}
				req, err = http.NewRequest(method, u, strings.NewReader(formData.Encode()))
				if err != nil {
					return nil, err
				}
			}
		case []byte:
			if len(vv) > 0 {
				req, err = http.NewRequest(method, u, bytes.NewReader(vv))
				if err != nil {
					return nil, err
				}
				req.ContentLength = int64(len(vv))
			}
		default:

		}
	}
	req.Header = http.Header{}
	ctx := NewContext(req, vs...)
	ctx.Task = task
	return ctx, nil
}

func NewContext(req *http.Request, vs ...interface{}) *Context {
	var (
		client *http.Client
	)
	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case *http.Header:
			for key, values := range *vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case *http.Client:
			client = vv
		case *http.Cookie:
			req.AddCookie(vv)
		case []http.Cookie:
			for _, cookie := range vv {
				req.AddCookie(&cookie)
			}
		case []*http.Cookie:
			for _, cookie := range vv {
				req.AddCookie(cookie)
			}
		case FormData:
			if len(vv) > 0 {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		}
	}

	if client == nil {
		client = getDefaultClient()
	}
	client.Transport = getDefaultTransport()

	// cookie jar
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		client.Jar = jar
	}
	if length := req.Header.Get("Content-Length"); length != "" {
		req.ContentLength = Str2Int64(length)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", defaultContentType)
	}

	if req.Header.Get("User-Agent") == "" || req.Header.Get("User-Agent") == "Go-http-client/2.0" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}

	return &Context{
		client:  client,
		Request: req,
		Data:    make(map[string]interface{}),
	}
}

/*
ziva.Cookie
*/

type Cookie struct {
	Name   string
	Value  string
	Domain string
	Path   string

	HttpOnly bool
}

/*
ziva.FormData
*/

type FormData map[string]string

func (f FormData) Set(k, v string) {
	f[k] = v
}

/*
private
*/

func getDefaultClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func getDefaultTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:    100,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func validUrl(urlStr string) (string, error) {
	length := len(urlStr)
	if length < 7 {
		return "", fmt.Errorf("error request url : %s", urlStr)
	}

	if urlStr[:7] == "http://" || urlStr[:8] == "https://" {
		return urlStr, nil
	}

	return "http://" + urlStr, nil
}

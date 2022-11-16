package ziva

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zituocn/ziva/logx"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CallbackFunc func(ctx *Context)

type Context struct {
	client *http.Client

	Request *http.Request

	Response *http.Response

	Err error

	RespBody []byte

	Task *Task

	Data map[string]interface{}

	execTime time.Duration

	Options Options
}

func (c *Context) Do() {
	var (
		bodyBytes []byte
		err       error
	)
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	if c.Options.SheepTime > 0 {
		time.Sleep(time.Duration(c.Options.SheepTime) * time.Millisecond)
	}
	if c.Options.StartFunc != nil {
		c.Options.StartFunc(c)
	}
	startTime := time.Now()
	c.Response, c.Err = c.client.Do(c.Request)
	if c.Err != nil {
		if c.Options.RetryFunc != nil {
			logx.Warnf("[%s] callback -> %s", "deadline", GetFuncName(c.Options.RetryFunc))
			c.Options.RetryFunc(c)
			return
		} else {
			logx.Errorf("request error :%s", c.Err.Error())
			return
		}
	}
	defer func(c *Context) {
		if c.Response != nil {
			err = c.Response.Body.Close()
			if err != nil {
				logx.Errorf("response body close error : %s", err.Error())
			}
		}
	}(c)
	c.execTime = time.Now().Sub(startTime)
	if c.Response.Header.Get("Content-Encoding") == "gzip" {
		c.Response.Body, err = gzip.NewReader(c.Response.Body)
		if err != nil {
			logx.Errorf("unzip failed: %s", err.Error())
			return
		}
	}

	if c.Response != nil {
		code := c.Response.StatusCode
		status := GetStatusByCode(code)
		body, err := ioutil.ReadAll(c.Response.Body)
		if err != nil {
			logx.Errorf("read response body error : %s", err.Error())
			logx.Debugf("task %v", c.Task)
			return
		}
		c.RespBody = body
		switch status {
		case "success":
			if c.Options.SucceedFunc != nil {
				logx.Infof("[%s] callback -> %s", status, GetFuncName(c.Options.SucceedFunc))
				c.Options.SucceedFunc(c)
			}
		case "retry":
			if c.Options.RetryFunc != nil {
				logx.Warnf("[%s] callback -> %s", status, GetFuncName(c.Options.RetryFunc))
				c.Options.RetryFunc(c)
			}
		case "fail":
			if c.Options.FailedFunc != nil {
				logx.Errorf("[%s] callback -> %s", status, GetFuncName(c.Options.FailedFunc))
				c.Options.FailedFunc(c)
			}
		default:
			logx.Warnf("unhandled status code :%d", code)
		}
	}

	if c.Options.CompleteFunc != nil {
		c.Options.CompleteFunc(c)
	}

	if c.Options.IsDebug {
		c.debugPrint()
	}
}

func (c *Context) SetProxy(httpProxy string) *Context {
	if httpProxy == "" {
		return c
	}

	proxy, _ := url.Parse(httpProxy)
	transport := getDefaultTransport()
	transport.Proxy = http.ProxyURL(proxy)
	c.client.Transport = transport
	return c
}

func (c *Context) SetTransport(f func() *http.Transport) *Context {
	c.client.Transport = f()
	return c
}

/*
response to
*/

func (c *Context) ToByte() []byte {
	if c.RespBody != nil {
		return c.RespBody
	}
	return []byte("")
}

func (c *Context) ToString() string {
	if c.RespBody != nil {
		return string(c.RespBody)
	}
	return ""
}

func (c *Context) ToSection(path string) string {
	s := c.ToString()
	if s != "" {
		return gjson.Get(s, path).String()
	}
	return ""
}

func (c *Context) ToJSON(v interface{}) error {
	if c.RespBody != nil {
		return json.Unmarshal(c.RespBody, &v)
	}
	return errors.New("response body is nil")
}

func (c *Context) ToHTML() string {
	s := c.ToString()
	return strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&#34;", `"`,
		"&#39;", "'",
	).Replace(s)
}

/*
private
*/
func (c *Context) debugPrint() {
	fmt.Printf("%s %v \n", leftText("URL:"), c.Request.URL)
	fmt.Printf("%s %v \n", leftText("Method:"), c.Request.Method)
	fmt.Printf("%s %v \n", leftText("Request Header:"), c.Request.Header)
	fmt.Printf("%s %v \n", leftText("Response code:"), c.Response.StatusCode)
	fmt.Printf("%s %v \n", leftText("Response Header:"), c.Response.Header)
}

func leftText(s string) string {
	return fmt.Sprintf("%15s", s)
}
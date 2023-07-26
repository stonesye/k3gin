package httpx

import "C"
import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	URL         string
	Method      string
	ContentType string
	Cookies     []Cookie
	Header      []Header
	Body        []Body
}

type Cookie struct {
	k, v string
}

type Header struct {
	k, v string
}

type Body struct {
	k, v string
}

func WithURL(u string) func(*Request) {
	return func(r *Request) {
		r.URL = u
	}
}

func WithBody(bodys ...Body) func(*Request) {
	return func(r *Request) {
		for _, b := range bodys {
			r.Body = append(r.Body, b)
		}
	}
}

func WithMethod(meth string) func(*Request) {
	return func(r *Request) {
		r.Method = meth
	}
}

func WithContentType(ctype string) func(*Request) {
	return func(r *Request) {
		r.ContentType = ctype
	}
}

func WithCookies(cookies ...Cookie) func(*Request) {
	return func(r *Request) {
		for _, cookie := range cookies {
			r.Cookies = append(r.Cookies, cookie)
		}
	}
}

func WithHeaders(headers ...Header) func(*Request) {
	return func(r *Request) {
		for _, header := range headers {
			r.Header = append(r.Header, header)
		}
	}
}

func NewRequest(options ...func(*Request)) *Request {
	r := &Request{}

	for _, option := range options {
		option(r)
	}

	return r
}

type Response struct {
	Status  int
	Body    string
	Message string
	Err     error
}

func NewSuccessResponse(body string, msg string) Response {
	return Response{
		Status:  0,
		Body:    body,
		Message: msg,
		Err:     nil,
	}
}

func New500Response(msg string, err ...error) Response {
	var e error

	if len(err) > 0 {
		e = err[0]
	} else {
		e = errors.New(msg)
	}

	return Response{
		Status:  500,
		Body:    "",
		Message: msg,
		Err:     e,
	}
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATH   = "PATH"
	DELETE = "DELETE"
)

const (
	FormData              = "multipart/form-data"
	XWwwFormUrlEncode     = "application/x-www-form-urlencoded"
	TextPlain             = "text/plain"
	ApplicationJson       = "application/json"
	ApplicationJavaScript = "application/javascript"
	ApplicationXml        = "application/xml"
	TextXml               = "text/xml"
	TextHtml              = "text/html"
)

type Client struct {
	C *http.Client
}

func InitHttp() (*Client, func(), error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,                                 // 是否开始keepalive
			Proxy:             http.ProxyFromEnvironment,             // 使用系统默认代理， nil表示不用代理
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // 是否校验服务端ssl证书
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(30) * time.Second,
				KeepAlive: time.Duration(30) * time.Second,
			}).DialContext, // 即将创建的网络连接的超时时间以及保持长链的时间
			MaxIdleConns:        10000,                           // 最大空闲连接数
			MaxIdleConnsPerHost: 2000,                            // 每一个 host 的最大连接数
			IdleConnTimeout:     time.Duration(90) * time.Second, // 一个连接在空闲多久之后关闭

		},
		Timeout: time.Duration(10) * time.Second, // 链接超时时间
	}

	cleanFunc := func() {

	}

	return &Client{
		C: client,
	}, cleanFunc, nil
}

func (c *Client) Get(URL string, cookies []Cookie, headers []Header) Response {
	r := NewRequest(WithURL(URL), WithCookies(cookies...), WithHeaders(headers...), WithMethod(GET), WithContentType(FormData))
	return httpRequest(r, c.C)
}

func (c *Client) Post(URL string, cookies []Cookie, headers []Header, bodys []Body) Response {
	r := NewRequest(WithURL(URL), WithCookies(cookies...), WithHeaders(headers...), WithMethod(POST), WithBody(bodys...), WithContentType(FormData))
	return httpRequest(r, c.C)
}

func httpRequest(r *Request, c *http.Client) Response {

	var bodyStr = ""

	if len(r.Body) > 0 {
		for _, b := range r.Body {
			bodyStr = b.k + "=" + b.v + "&"
		}

		bodyStr = bodyStr[:len(bodyStr)-1]
	}

	request, err := http.NewRequest(r.Method, r.URL, strings.NewReader(bodyStr))

	if err != nil {
		return New500Response("Create new request fail ")
	}

	defer request.Body.Close()

	request.Header.Add("content-type", r.ContentType)

	for _, h := range r.Header {
		request.Header.Add(h.k, h.v)
	}

	for _, c := range r.Cookies {
		request.AddCookie(&http.Cookie{Name: c.k, Value: c.v, Expires: time.Time{}, HttpOnly: true})
	}

	res, err := c.Do(request)

	if err != nil {
		return New500Response("Execute request fail ")
	}

	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return New500Response("Read response body fail ")
	}

	return NewSuccessResponse(string(buf), "ok")
}

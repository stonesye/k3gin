package http

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Url         string
	Method      string
	ContentType string
	Body        string
	Cookies     map[string]string
	Header      map[string]string
}

type Response struct {
	StatusCode int
	Body       string
}

type HttpClient struct {
	client *http.Client
	contentType
	method
}

type method struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
}

type contentType struct {
	FormData              string
	XWwwFormUrlEncode     string
	TextPlain             string
	ApplicationJson       string
	ApplicationJavaScript string
	ApplicationXml        string
	TextXml               string
	TextHtml              string
}

var httpClient *HttpClient

func InitHttp() {
	httpClient = new(HttpClient)

	httpClient.client = &http.Client{
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

	httpClient.contentType = contentType{
		FormData:              "multipart/form-data",
		XWwwFormUrlEncode:     "application/x-www-form-urlencoded",
		TextPlain:             "text/plain",
		ApplicationJson:       "application/json",
		ApplicationJavaScript: "application/javascript",
		ApplicationXml:        "application/xml",
		TextXml:               "text/xml",
		TextHtml:              "text/html",
	}

	httpClient.method = method{
		GET:    "GET",
		POST:   "POST",
		PUT:    "PUT",
		PATCH:  "PATCH",
		DELETE: "DELETE",
	}
}

func HttpGet(url string, header map[string]string, cookies map[string]string) string {
	reqs := Request{
		Url:         url,
		Method:      httpClient.GET,
		ContentType: httpClient.XWwwFormUrlEncode,
		Body:        "",
		Cookies:     cookies,
		Header:      header,
	}
	resp := new(Response)
	httpRequest(reqs, resp)
	return resp.Body
}

func HttpPost(url string, data map[string]string, header map[string]string, cookies map[string]string) string {

	var query string
	for k, v := range data {
		query += k + "=" + v + "&"
	}
	if query != "" {
		query = query[:len(query)-1]
	}

	reqs := Request{
		Url:         url,
		Method:      httpClient.POST,
		ContentType: httpClient.XWwwFormUrlEncode,
		Body:        query,
		Cookies:     cookies,
		Header:      header,
	}

	resp := new(Response)

	httpRequest(reqs, resp)

	return resp.Body
}

func httpRequest(req Request, resp *Response) (err error) {

	request, err := http.NewRequest(req.Method, req.Url, strings.NewReader(req.Body))
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()

	request.Header.Add("content-type", req.ContentType)
	for key, val := range req.Header {
		request.Header.Add(key, val)
	}
	for k, v := range req.Cookies {
		request.AddCookie(&http.Cookie{
			Name:     k,
			Value:    v,
			Expires:  time.Time{},
			HttpOnly: true,
		})
	}

	response, err := httpClient.client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	resp.StatusCode = response.StatusCode
	resp.Body = string(body)

	return err
}

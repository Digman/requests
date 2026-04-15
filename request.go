package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	tls_client "github.com/Digman/tls-client"
	JSON "github.com/tidwall/gjson"

	"io"
	"maps"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

type Request struct {
	client      tls_client.HttpClient
	request     *http.Request
	response    *http.Response
	method      string
	url         string
	header      map[string]string
	headerOrder []string
	cookies     *[]*http.Cookie
	dataType    any
	data        url.Values
	jsonData    string
	fileData    map[bool]map[string]string
	debug       bool
	err         error
	username    string
	password    string
}

func NewRequest(client tls_client.HttpClient) *Request {
	return &Request{
		client:   client,
		method:   "GET",
		header:   make(map[string]string),
		cookies:  new([]*http.Cookie),
		data:     url.Values{},
		fileData: make(map[bool]map[string]string),
	}
}

func (r *Request) Get(url string) *Request {
	return r.SetUrl(url).SetMethod(http.MethodGet)
}

func (r *Request) Post(url string) *Request {
	return r.SetUrl(url).SetMethod(http.MethodPost)
}

func (r *Request) Head(url string) *Request {
	return r.SetUrl(url).SetMethod(http.MethodHead)
}

func (r *Request) Put(url string) *Request {
	return r.SetUrl(url).SetMethod(http.MethodPut)
}

func (r *Request) Options(url string) *Request {
	return r.SetUrl(url).SetMethod(http.MethodOptions)
}

func (r *Request) SetBasicAuth(userName, password string) *Request {
	r.username = userName
	r.password = password
	return r
}

func (r *Request) SetMethod(name string) *Request {
	r.method = strings.ToUpper(name)
	return r
}

func (r *Request) SetUrl(url string) *Request {
	r.url = url
	return r
}

func (r *Request) SetHeader(name, value string) *Request {
	r.header[name] = value
	return r
}

func (r *Request) SetHeaders(values map[string]string) *Request {
	maps.Copy(r.header, values)
	return r
}

func (r *Request) SetHeaderOrder(order []string) *Request {
	r.headerOrder = order
	return r
}

func (r *Request) SetCookies(cookies *[]*http.Cookie) *Request {
	r.cookies = cookies
	return r
}

func (r *Request) SetDebug(d bool) *Request {
	r.debug = d
	return r
}

func (r *Request) SetData(name, value string) *Request {
	r.data.Set(name, value)
	return r
}

func (r *Request) GetAllData() url.Values {
	return r.data
}

func (r *Request) SetAllData(data url.Values) *Request {
	r.data = data
	return r
}

func (r *Request) SetJsonData(s string) *Request {
	r.jsonData = s
	r.dataType = "json"
	return r
}

func (r *Request) SetJson(data any) *Request {
	jsonData, err := json.Marshal(data)
	if err == nil {
		r.jsonData = string(jsonData)
	}
	r.dataType = "json"
	return r
}

func (r *Request) SetFileData(name, value string, isFile bool) *Request {
	if _, exists := r.fileData[isFile]; exists {
		r.fileData[isFile][name] = value
	} else {
		r.fileData[isFile] = map[string]string{name: value}
	}
	r.dataType = "file"
	return r
}

func (r *Request) log(t string) {
	if r.debug == true {
		fmt.Printf("[Request Debug]\n")
		fmt.Printf("-------------------------------------------------------------------\n")
		fmt.Printf("Request: %s %s\nHeader: %v\nCookies: %v\n", r.method, r.url, r.request.Header, r.request.Cookies())
		switch t {
		case "url":
			if r.method == "GET" {
				fmt.Printf("Query: %v\n", r.request.URL.RawQuery)
			} else {
				fmt.Printf("Body: %v\n", r.data)
			}
		case "json":
			fmt.Printf("Body: %v\n", r.jsonData)
		default:
			fmt.Printf("Body: %v\n", r.fileData)
		}
		fmt.Printf("-------------------------------------------------------------------\n")
	}
}

func (r *Request) Send(a ...any) *Request {
	var err error
	if len(a) > 0 {
		r.dataType = a[0]
	}
	r.err = nil
	switch r.dataType {
	case nil, "url":
		var body io.Reader
		if r.method != "GET" {
			body = strings.NewReader(r.data.Encode())
		}
		r.request, err = http.NewRequest(r.method, r.url, body)

		defer r.log("url")
		if err != nil {
			r.err = err
			return r
		}

		if r.username != "" && r.password != "" {
			r.request.SetBasicAuth(r.username, r.password)
		}

		if r.method == "POST" {
			r.request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		} else if r.method == "GET" && len(r.data) > 0 {
			r.request.URL.RawQuery = r.data.Encode()
		}
	case "json":
		r.request, err = http.NewRequest(r.method, r.url, strings.NewReader(r.jsonData))
		defer r.log("json")
		if err != nil {
			r.err = err
			return r
		}
		r.request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	case "file":
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)
		for h, m := range r.fileData {
			for k, v := range m {
				if h {
					fd, err := os.Open(v)
					if err != nil {
						r.err = err
						return r
					}
					fileWriter, _ := bodyWriter.CreateFormFile(k, filepath.Base(v))
					_, _ = io.Copy(fileWriter, fd)
					_ = fd.Close()
				} else {
					_ = bodyWriter.WriteField(k, v)
				}
			}
		}

		contentType := bodyWriter.FormDataContentType()
		_ = bodyWriter.Close()
		r.request, err = http.NewRequest(r.method, r.url, io.NopCloser(bodyBuf))
		defer r.log("file")
		if err != nil {
			r.err = err
			return r
		}

		r.request.Header.Set("Content-Type", contentType)
	default:
		r.err = errors.New("unsupported data type")
		return r
	}
	for k, v := range r.header {
		r.request.Header.Set(k, v)
	}

	if len(r.headerOrder) > 0 {
		r.request.Header[http.HeaderOrderKey] = r.headerOrder
	}

	for _, v := range *r.cookies {
		r.request.AddCookie(v)
	}

	r.response, err = r.client.Do(r.request)
	if err != nil {
		r.err = err
	}

	return r
}

func (r *Request) Close() {
	if r.response == nil {
		return
	}
	_ = r.response.Body.Close()
}

func (r *Request) End() (*http.Response, string, error) {
	response, bodyByte, err := r.EndByte()
	if err != nil {
		return nil, "", err
	}
	return response, string(bodyByte), nil
}

func (r *Request) EndJson() (*http.Response, JSON.Result, error) {
	response, body, err := r.EndByte()

	if err != nil {
		return nil, JSON.Result{}, err
	}

	return response, JSON.ParseBytes(body), nil
}

func (r *Request) EndResponse() (*http.Response, error) {
	defer r.Close()

	if r.err != nil {
		return nil, r.err
	}

	if r.response == nil {
		return nil, errors.New("response empty")
	}

	// HTTP/2 多路復用：Body.Close() 即可釋放串流，無需 drain
	// HTTP/1.1 keep-alive：需讀完 body 才能將連線歸還連線池，限額避免超大 body 拖慢
	if r.response.ProtoMajor < 2 {
		const maxDrain = 64 << 10
		_, _ = io.CopyN(io.Discard, r.response.Body, maxDrain)
	}

	return r.response, nil
}

func (r *Request) EndByte() (*http.Response, []byte, error) {
	defer r.Close()

	if r.err != nil {
		return nil, nil, r.err
	}

	if r.response == nil {
		return nil, nil, errors.New("response empty")
	}

	body := r.response.Body
	if r.response.ProtoMajor < 2 {
		body = http.DecompressBody(r.response)
	}

	buf := bytes.Buffer{}
	if n := r.response.ContentLength; n > 0 {
		buf.Grow(int(n))
	}
	if _, err := io.Copy(&buf, body); err != nil {
		return nil, nil, err
	}

	return r.response, buf.Bytes(), nil
}

func (r *Request) EndFile(savePath, saveFileName string) (*http.Response, error) {
	defer r.Close()

	if r.err != nil {
		return nil, r.err
	}

	if r.response == nil {
		return nil, errors.New("response empty")
	}

	if r.response.StatusCode != http.StatusOK {
		return nil, errors.New("not written")
	}

	if saveFileName == "" {
		path := strings.Split(r.request.URL.String(), "/")
		if len(path) > 1 {
			saveFileName = path[len(path)-1]
		}
	}

	f, err := os.OpenFile(savePath+saveFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err = io.Copy(f, r.response.Body); err != nil {
		return nil, err
	}

	return r.response, nil
}

package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
)

type Response struct {
	Text       string
	StatusCode int
}

func (resp *Response) Json() map[string]interface{} {
	var mapData map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Text), &mapData); err != nil {
		return nil
	}
	return mapData
}

type Session struct {
	core *http.Client
}

func (client *Session) request(url, method, body string) (*Response, error) {
	var req *http.Request
	var res *http.Response
	var err error
	var respTextBytes []byte
	var bodyReader io.Reader = nil

	if body != "" {
		var bodyBytes = []byte(body)
		bodyReader = bytes.NewBuffer(bodyBytes)
	}
	if req, err = http.NewRequest(method, url, bodyReader); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if res, err = client.core.Do(req); err != nil {
		return nil, err
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	if respTextBytes, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	} else {
		return &Response{Text: string(respTextBytes),
			StatusCode: res.StatusCode}, nil
	}
}

func (client *Session) Post(url, body string) (*Response, error) {
	return client.request(url, "POST", body)
}

func (client *Session) Get(url string) (*Response, error) {
	return client.request(url, "GET", "")
}

type MyTransport struct {
	Transport http.RoundTripper
}

func (t *MyTransport) transport() http.RoundTripper {
	if nil != t.Transport {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *MyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729)")
	return t.transport().RoundTrip(req)
}

func NewSession() *Session {
	t := &MyTransport{}
	jar, err := cookiejar.New(nil)
	if nil != err {
		log.Fatal(err)
	}
	client := http.DefaultClient
	client.Transport = t
	client.Jar = jar
	return &Session{
		core: client,
	}
}

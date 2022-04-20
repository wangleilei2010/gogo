package http

import (
	"bytes"
	"encoding/json"
	"github.com/oliveagle/jsonpath"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Response[T any] struct {
	Text       string
	StatusCode int
	Value      T
}

func (resp *Response[T]) Json() map[string]interface{} {
	var mapData map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Text), &mapData); err != nil {
		return nil
	}
	return mapData
}

func (resp *Response[T]) J(expression string) interface{} {
	if j := resp.Json(); j != nil {
		if r, e := jsonpath.JsonPathLookup(j, expression); e != nil {
			return nil
		} else {
			return r
		}
	} else {
		return nil
	}
}

type Session struct {
	core *http.Client
}

func request[T any](client *Session, url, method, body string) (*Response[T], error) {
	var (
		req           *http.Request
		res           *http.Response
		err           error
		respTextBytes []byte
		bodyReader    io.Reader = nil
	)

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
		var value T
		err = json.Unmarshal(respTextBytes, &value)
		if err != nil {
			return nil, err
		}
		return &Response[T]{Text: string(respTextBytes), Value: value,
			StatusCode: res.StatusCode}, nil
	}
}

func Post[T any](client *Session, url, body string) (*Response[T], error) {
	var session *Session
	if client == nil {
		session = NewSession(time.Second * 120)
	} else {
		session = client
	}
	return request[T](session, url, "POST", body)
}

func Get[T any](client *Session, url string) (*Response[T], error) {
	var session *Session
	if client == nil {
		session = NewSession(time.Second * 120)
	} else {
		session = client
	}
	return request[T](session, url, "GET", "")
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

func NewSession(timeOut time.Duration) *Session {
	t := &MyTransport{}
	jar, err := cookiejar.New(nil)
	if nil != err {
		log.Fatal(err)
	}
	client := &http.Client{}
	client.Transport = t
	client.Jar = jar
	if timeOut != 0 {
		client.Timeout = timeOut
	}
	return &Session{
		core: client,
	}
}

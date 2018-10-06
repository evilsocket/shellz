package plugins

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpPackage struct {
}

var hp = httpPackage{}

func getHTTP() httpPackage {
	return hp
}

type httpResponse struct {
	Error    error
	Response *http.Response
	Raw      []byte
	Body     string
}

func (c httpPackage) Request(method string, uri string, headers map[string]string, form map[string]string) httpResponse {
	var reader io.Reader
	if form != nil {
		data := url.Values{}
		for k, v := range form {
			data.Set(k, v)
		}
		reader = bytes.NewBufferString(data.Encode())
	}

	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		return httpResponse{Error: err}
	}

	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return httpResponse{Error: err}
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return httpResponse{Error: err}
	}

	return httpResponse{
		Error:    nil,
		Response: resp,
		Raw:      raw,
		Body:     string(raw),
	}
}

func (c httpPackage) Get(url string, headers map[string]string) httpResponse {
	return c.Request("GET", url, headers, nil)
}

func (c httpPackage) Post(url string, headers map[string]string, form map[string]string) httpResponse {
	return c.Request("POST", url, headers, form)
}

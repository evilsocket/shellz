package session

import (
	"io/ioutil"
	"net/http"
)

type httpClient struct {
}

func newHttpClient() httpClient {
	return httpClient{}
}

type response struct {
	Error    error
	Response *http.Response
	Raw      []byte
	Body     string
}

func (c httpClient) Get(url string, headers map[string]string) response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response{Error: err}
	}

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response{Error: err}
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response{Error: err}
	}

	return response{
		Error:    nil,
		Response: resp,
		Raw:      raw,
		Body:     string(raw),
	}
}

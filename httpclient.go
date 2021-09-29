package sendhttp

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type HttpClient struct {
	httpClient *http.Client
	debug      bool
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		httpClient: &http.Client{},
		debug:      false,
	}
}

var _ Client = (*HttpClient)(nil)

func (c *HttpClient) SetDebug(flag bool) {
	c.debug = flag
}

func (c *HttpClient) Send(request Request, response Response) error {
	if err := request.CheckValid(); err != nil {
		return err
	}

	bodyReader, err := GetBody(request)
	if err != nil {
		return err
	}

	httpRequest, err := http.NewRequest(request.GetMethod(), request.GetUrl(), bodyReader)
	if err != nil {
		return err
	}
	for k, v := range request.GetHeaders() {
		httpRequest.Header.Set(k, v)
	}

	if c.debug {
		outbytes, err := httputil.DumpRequest(httpRequest, true)
		if err != nil {
			log.Printf("[ERROR] dump request failed because %s", err)
			return err
		}
		log.Printf("[DEBUG] http request:\n%s", outbytes)
	}

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		msg := fmt.Sprintf("Fail to get response because %s", err)
		return errors.New(msg)
	}

	if c.debug {
		//  but is does not contain response body
		outbytes, err := httputil.DumpResponse(httpResponse, false)
		if err != nil {
			log.Printf("[ERROR] dump response failed because %s", err)
			return err
		}
		log.Printf("[DEBUG] http response:\n%s", outbytes)
	}

	err = ParseHttpResponse(httpResponse, response)
	return err
}

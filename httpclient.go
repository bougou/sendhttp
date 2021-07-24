package sendhttp

import (
	"bytes"
	"encoding/json"
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

	reqBody, err := json.Marshal(request)
	if err != nil {
		msg := fmt.Sprintf("marshal request failed, %s", err)
		return errors.New(msg)
	}

	httpRequest, err := http.NewRequest(request.GetMethod(), request.GetUrl(), bytes.NewReader(reqBody))
	if err != nil {
		return err
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
		outbytes, err := httputil.DumpResponse(httpResponse, false)
		if err != nil {
			log.Printf("[ERROR] dump response failed because %s", err)
			return err
		}
		log.Printf("[DEBUG] http response:\n%s", outbytes)
	}

	response.SetRaw(reqBody)
	err = ParseHttpResponse(httpResponse, response)

	return err
}

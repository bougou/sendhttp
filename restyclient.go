package sendhttp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http/httputil"
	"strings"

	"github.com/go-resty/resty/v2"
)

type RestyClient struct {
	restyClient *resty.Client
	debug       bool
}

func NewRestyClient() *RestyClient {
	return &RestyClient{
		restyClient: resty.New(),
		debug:       false,
	}
}

var _ Client = (*RestyClient)(nil)

func (c *RestyClient) SetDebug(flag bool) {
	c.debug = flag
}

func (c *RestyClient) Send(request Request, response Response) error {
	// Note, this is a must.
	c.restyClient.SetDoNotParseResponse(true)

	restyReq := c.restyClient.R()
	restyReq.SetHeaders(request.GetHeaders())

	bodyReader, err := GetBody(request)
	if err != nil {
		return err
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, bodyReader)
	if err != nil {
		log.Printf("read bodyReader failed because %s", err)
		return err
	}
	bodyStr := buf.String()
	restyReq.SetBody(bodyStr)

	// the restyReq.RawRequest is ONLY set after restyReq.Execute
	restyResp, err := restyReq.Execute(request.GetMethod(), request.GetUrl())
	if err != nil {
		msg := fmt.Sprintf("request failed, err = %s", err.Error())
		return errors.New(msg)
	}

	if c.debug {
		// DumpRequest of restyReq.RawRequest does not contain body bytes
		reqbytes, err := httputil.DumpRequest(restyReq.RawRequest, true)
		if err != nil {
			log.Printf("[ERROR] dump request failed because %s", err)
			return err
		}
		log.Printf("[DEBUG] http request:\n%s%s", reqbytes, bodyStr)
	}

	if c.debug {
		resBytes, err := httputil.DumpResponse(restyResp.RawResponse, false)
		if err != nil {
			log.Printf("[ERROR] dump response failed because %s", err)
			return err
		}
		log.Printf("[DEBUG] http response:\n%s", resBytes)
	}

	err = ParseHttpResponse(restyResp.RawResponse, response)
	return err
}

package sendhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http/httputil"

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
	if c.debug {
		c.restyClient.EnableTrace()
	}

	// Note, this is a must.
	c.restyClient.SetDoNotParseResponse(true)

	restyReq := c.restyClient.R()
	restyReq.SetHeaders(request.GetHeaders())

	reqBody, err := json.Marshal(request)
	if err != nil {
		msg := fmt.Sprintf("marshal request failed, %s", err)
		return errors.New(msg)
	}
	restyReq.SetBody(reqBody)

	restyResp, err := restyReq.Execute(request.GetMethod(), request.GetUrl())

	if c.debug {
		outbytes, err := httputil.DumpRequest(restyReq.RawRequest, true)
		if err != nil {
			log.Printf("[ERROR] dump request failed because %s", err)
			return err
		}
		log.Printf("[DEBUG] http request = %s", outbytes)
	}

	if err != nil {
		msg := fmt.Sprintf("request failed, err = %s", err.Error())
		return errors.New(msg)
	}

	err = ParseHttpResponse(restyResp.RawResponse, response)
	return err
}

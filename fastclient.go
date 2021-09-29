package sendhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type FastClient struct {
	client *fasthttp.Client
	debug  bool
}

func NewFastClient() *FastClient {
	return &FastClient{
		client: &fasthttp.Client{},
		debug:  false,
	}
}

var _ Client = (*FastClient)(nil)

func (c *FastClient) SetDebug(flag bool) {
	c.debug = flag
}

func (c *FastClient) Send(request Request, response Response) error {
	if err := request.CheckValid(); err != nil {
		return err
	}

	fastRequest := new(fasthttp.Request)
	bodyReader, err := GetBody(request)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, bodyReader)
	if err != nil {
		log.Printf("read bodyReader failed because %s", err)
		return err
	}

	fastRequest.SetBody(buf.Bytes())

	fastRequest.Header.SetMethod(request.GetMethod())

	fastRequest.SetRequestURI(request.GetUrl())

	for k, v := range request.GetHeaders() {
		fastRequest.Header.Set(k, v)
	}

	if c.debug {
		requestCtx := fasthttp.RequestCtx{Request: *fastRequest}
		httpRequest := new(http.Request)
		if err := fasthttpadaptor.ConvertRequest(&requestCtx, httpRequest, false); err != nil {
			return fmt.Errorf("convert request failed, err: %s", err)
		}

		outbytes, err := httputil.DumpRequest(httpRequest, true)
		if err != nil {
			log.Printf("[ERROR] dump request failed because %s", err)
			return err
		}
		log.Printf("[DEBUG] http request:\n%s", outbytes)
	}

	var fastResponse = new(fasthttp.Response)
	err = c.client.Do(fastRequest, fastResponse)
	if err != nil {
		msg := fmt.Sprintf("Fail to get response because %s", err)
		return errors.New(msg)
	}

	headers := make(http.Header)
	fastResponse.Header.VisitAll(func(key []byte, value []byte) {
		headers.Add(string(key), string(value))
	})

	httpResponse := http.Response{
		StatusCode:    fastResponse.StatusCode(),
		ContentLength: int64(fastResponse.Header.ContentLength()),
		Header:        headers,
	}

	body := fastResponse.Body()

	err = json.Unmarshal(body, &response)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors.New(msg)
	}

	response.Fill(&httpResponse)

	if c.debug {
		var out bytes.Buffer
		fastResponse.Header.VisitAll(func(key []byte, value []byte) {
			out.Write(key)
			out.WriteString(": ")
			out.Write(value)
			out.WriteByte('\n')
		})
		out.Write(body)

		log.Printf("[DEBUG] http response:\n%s", out.Bytes())
	}

	// r.status = response.Status
	// r.statusCode = response.StatusCode
	// r.proto = response.Proto
	// r.protoMajor = response.ProtoMajor
	// r.protoMinor = response.ProtoMinor
	// r.contentLength = response.ContentLength
	// r.header = response.Header

	response.SetRaw(body)

	// err = ParseHttpResponse(httpResponse, response)

	return err
}

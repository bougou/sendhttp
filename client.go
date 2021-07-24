package sendhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client is a http client which has a Send(request Request, response Response) method.
// Client does not contain URL.
// All HTTP request informations like method, sheme, URL, body and others are ALL fetched from the Request interface.
type Client interface {
	Send(request Request, response Response) error
	SetDebug(debug bool)
}

// ParseFromHttpResponse parses the native http.Response hr to response which implements the Response interface.
func ParseHttpResponse(hr *http.Response, response Response) (err error) {
	defer hr.Body.Close()
	body, err := ioutil.ReadAll(hr.Body)
	if err != nil {
		msg := fmt.Sprintf("Fail to read response body because %s", err)
		return errors.New(msg)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors.New(msg)
	}

	response.Fill(hr)
	response.SetRaw(body)

	return
}

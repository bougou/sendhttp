package sendhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeRequest struct {
	*BaseRequest
}

type fakeResponse struct {
	*BaseResponse `json:"-"`
	MSG           string `json:"msg"`
}

// fakeHttpClient test HttpClient
type fakeHttpClient struct {
	client *HttpClient
	url    string
}

// fakeRestyClient test RestyClient
type fakeRestyClient struct {
	client *RestyClient
	url    string
}

func newFakeRequest() *fakeRequest {
	return &fakeRequest{
		BaseRequest: NewBaseRequest(),
	}

}

func newFakeResponse() *fakeResponse {
	return &fakeResponse{
		BaseResponse: NewBaseResponse(),
	}
}

func newFakeHttpClient(url string) *fakeHttpClient {
	return &fakeHttpClient{
		client: NewHttpClient(),
		url:    url,
	}
}

func newFakeRestyClient(url string) *fakeRestyClient {
	return &fakeRestyClient{
		client: NewRestyClient(),
		url:    url,
	}
}

func (c *fakeHttpClient) Complete(r Request) error {
	r.SetMethod(GET)
	r.SetUrl(c.url)
	r.SetPath("/")
	return nil
}

func (c *fakeRestyClient) Complete(r Request) error {
	r.SetMethod(GET)
	r.SetUrl(c.url)
	r.SetPath("/")
	return nil
}

// Fake is normally your exported method for your client
func (c *fakeHttpClient) Fake(request *fakeRequest) (response *fakeResponse, err error) {
	response = newFakeResponse()
	err = c.client.Send(request, response)
	return
}

func (c *fakeRestyClient) Fake(request *fakeRequest) (response *fakeResponse, err error) {
	response = newFakeResponse()
	err = c.client.Send(request, response)
	return
}

var mockGetResposne = `{"msg": "Hello, World"}`

func TestHttpClient(t *testing.T) {
	// server.URL returns the randomly listened address, like "http://127.0.0.1:58713"
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note, the output contains an extra line break character
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeHttpClient(server.URL)
	c.client.SetDebug(true)
	request := newFakeRequest()
	c.Complete(request)

	response, err := c.Fake(request)
	if err != nil {
		t.Error(err)
	}

	got := string(response.GetRaw())
	expected := mockGetResposne + "\n"

	if got != expected {
		t.Errorf("response not matched, expected: %s, got: %s", expected, got)
	}
}

func TestRestyClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeRestyClient(server.URL)
	c.client.SetDebug(true)
	request := newFakeRequest()
	c.Complete(request)

	response, err := c.Fake(request)
	if err != nil {
		t.Error(err)
	}
	got := string(response.GetRaw())
	expected := mockGetResposne + "\n"

	if got != expected {
		t.Errorf("response not matched, expected: %s, got: %s", expected, got)
	}
}

func TestHttpClientForm(t *testing.T) {
	// server.URL returns the randomly listened address, like "http://127.0.0.1:58713"
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note, the output contains an extra line break character
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeHttpClient(server.URL)
	c.client.SetDebug(true)
	request := newFakeRequest()
	request.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	request.SetParams(map[string]string{
		"limit":  "10",
		"offset": "20",
	})
	request.SetFormParams(map[string]string{
		"a": "hello",
		"b": "100",
	})
	c.Complete(request)

	response, err := c.Fake(request)
	if err != nil {
		t.Error(err)
	}

	got := string(response.GetRaw())
	expected := mockGetResposne + "\n"

	if got != expected {
		t.Errorf("response not matched, expected: %s, got: %s", expected, got)
	}
}

func Test_RestyClientForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeRestyClient(server.URL)
	c.client.SetDebug(true)
	request := newFakeRequest()
	request.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	request.SetParams(map[string]string{
		"limit":  "10",
		"offset": "20",
	})
	request.SetFormParams(map[string]string{
		"a": "hello",
		"b": "100:200",
	})
	c.Complete(request)

	response, err := c.Fake(request)
	if err != nil {
		t.Error(err)
	}
	got := string(response.GetRaw())
	expected := mockGetResposne + "\n"

	if got != expected {
		t.Errorf("response not matched, expected: %s, got: %s", expected, got)
	}
}

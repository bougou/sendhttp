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

// fakeClient test HttpClient
type fakeClient struct {
	client *HttpClient
	url    string
}

// fakeClient2 test RestyClient
type fakeClient2 struct {
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

func newFakeClient(url string) *fakeClient {
	return &fakeClient{
		client: NewHttpClient(),
		url:    url,
	}
}

func newFakeClient2(url string) *fakeClient2 {
	return &fakeClient2{
		client: NewRestyClient(),
		url:    url,
	}
}

func (c *fakeClient) Complete(r Request) error {
	r.SetMethod(GET)
	r.SetUrl(c.url)
	r.SetPath("/")
	return nil
}

func (c *fakeClient2) Complete(r Request) error {
	r.SetMethod(GET)
	r.SetUrl(c.url)
	r.SetPath("/")
	return nil
}

// Fake is normally your exported method for your client
func (c *fakeClient) Fake(request *fakeRequest) (response *fakeResponse, err error) {
	response = newFakeResponse()
	err = c.client.Send(request, response)
	return
}

func (c *fakeClient2) Fake(request *fakeRequest) (response *fakeResponse, err error) {
	response = newFakeResponse()
	err = c.client.Send(request, response)
	return
}

var mockGetResposne = `{"msg": "Hello, World"}`

func TestHttpClient(t *testing.T) {
	var (
		server *httptest.Server
	)

	// server.URL returns the randomly listened address, like "http://127.0.0.1:58713"
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note, the output contains an extra line break character
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeClient(server.URL)
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

	c := newFakeClient2(server.URL)
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

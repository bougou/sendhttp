package sendhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeFastClient test FastClient
type fakeFastClient struct {
	client *FastClient
	url    string
}

func newFakeFastClient(url string) *fakeFastClient {
	return &fakeFastClient{
		client: NewFastClient(),
		url:    url,
	}
}

func (c *fakeFastClient) Complete(r Request) error {
	r.SetMethod(GET)
	r.SetUrl(c.url)
	r.SetPath("/")
	return nil
}

func (c *fakeFastClient) Fake(request *fakeRequest) (response *fakeResponse, err error) {
	response = newFakeResponse()
	err = c.client.Send(request, response)
	return
}

func TestFastClient(t *testing.T) {
	var mockGetResposne = `{"msg": "Hello, World"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeFastClient(server.URL)
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

func Test_FastClientForm(t *testing.T) {
	var mockGetResposne = `{"msg": "Hello, World"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeFastClient(server.URL)
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

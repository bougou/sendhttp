package sendhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeRestyClient test RestyClient
type fakeRestyClient struct {
	client *RestyClient
	url    string
}

func newFakeRestyClient(url string) *fakeRestyClient {
	return &fakeRestyClient{
		client: NewRestyClient(),
		url:    url,
	}
}

func (c *fakeRestyClient) Complete(r Request) error {
	r.SetMethod(GET)
	r.SetUrl(c.url)
	r.SetPath("/")
	return nil
}

func (c *fakeRestyClient) Fake(request *fakeRequest) (response *fakeResponse, err error) {
	response = newFakeResponse()
	err = c.client.Send(request, response)
	return
}

func TestRestyClient(t *testing.T) {
	var mockGetResposne = `{"msg": "Hello, World"}`

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

func Test_RestyClientForm(t *testing.T) {
	var mockGetResposne = `{"msg": "Hello, World"}`

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

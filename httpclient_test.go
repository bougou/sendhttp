package sendhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
)

// fakeHttpClient test HttpClient
type fakeHttpClient struct {
	client *HttpClient
	url    string
}

func newFakeHttpClient(url string) *fakeHttpClient {
	return &fakeHttpClient{
		client: NewHttpClient(),
		url:    url,
	}
}

func (c *fakeHttpClient) Complete(r Request) error {
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

func TestHttpClient(t *testing.T) {
	var mockGetResposne = `{"msg": "Hello, World"}`

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

func TestHttpClientForm(t *testing.T) {
	var mockGetResposne = `{"msg": "Hello, World"}`

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

func Test_ClientMultipart(t *testing.T) {
	var mockGetResposne = `{"msg": "Hello, World"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockGetResposne)
	}))
	defer server.Close()

	c := newFakeHttpClient(server.URL)
	c.client.SetDebug(true)
	request := newFakeRequest()
	request.SetHeaders(map[string]string{
		"Content-Type": "multipart/form-data",
	})
	request.SetParams(map[string]string{
		"limit":  "10",
		"offset": "20",
	})
	request.SetFormParams(map[string]string{
		"a": "hello",
		"b": "100:200",
	})

	home, _ := homedir.Dir()
	file, _ := os.Open(path.Join(home, "file1.txt"))
	request.AddMultipart("key1", "file1.txt", file)
	request.AddMultipart("key2", "", strings.NewReader("hello world"))
	request.AddMultipart("key3", "", strings.NewReader("test test"))
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

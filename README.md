# sendhttp

`sendhttp` can make the process of creating golang sdk for http service more easier.

This package defines three interface.

- `Request`
- `Response`
- `Client`

And four struct:

- `BaseRequest`
- `BaseResponse`
- `HttpClient`
- `RestyClient`

The `BaseRequest` struct implemented the `Request` interface.
So all business request struct just embed `*BaseRequest` into itself.

The `BaseResponse` struct implemented the `Response` interface.
So all business response struct just embed `*BaseResponse` into itself.

This package also provides two struct `RestyClient` and `HttpClient` which implemented the `Client` interface.

The `HttpClient` is based on the `http.Client` of the standard library `net/http`.

The `RestyClient` is based on the `resty.Client` of the `github.com/go-resty/resty/v2`.

In your business code, you can directly embeds `RestyClient` or `HttpClient` in your own client struct.

**Warning**: Currently, the `HttpClient` and `RestyClient` CAN ONLY handle **json** compatible request and response. So your business request and response struct MUST can be marshaled to or unmarshaled from plain json.

## example

The `client_test.go` source code demonstrates the basic usage.

Below is the common procedures to use this package.

Suppose you have a service named `Foo`, which provides informations about movies and others.
You want to provide a golang sdk for accessing the service.

1. Firstly, create a Client which holds information used to connect to your service.

```go
type FooClient struct {
  // embed a client
  // Note here we use the Client interface
  client sendhttp.Client

  url string
  // and many other informations, like user, pass etc
}

func NewFooClient() *FooClient {
  return &FooClient{
    client: sendhttp.NewHttpClient()
    // and initialize other fields if necessary
  }
}
```

2. Then, define the request and response struct which describe the REST API's input and output.

```go
type ListMoviesRequest struct {
  *sendhttp.BaseRequest

  // define other fields if the request contains json body
}

func NewListMoviesRequest() *ListMoviesRequest{
  return &ListMoviesRequest{
    BaseRequest: sendhttp.NewBaseRequest(),

    // init other fields if necessary
  }
}

type ListMoviesResonse struct {
  *sendhttp.BaseResponse

  // define fields which holds the response
  Items []Movie `json:"items"`
}

func NewListMoviesResponse() *ListMoviesResonse{
  return &ListMoviesResonse{
    BaseResponse: sendhttp.NewBaseResponse(),

    // init other fields if necessary
  }
}

// define business model
type Movie struct {
  Name string
  // ...
}
```

3. Define a method

```go
func (c *Client) ListMovies(request *ListMoviesRequest) (response *ListMoviesResonse, err error) {

  // use the client information to fill the request
  // normally your will
  // c.Complete(request)

  resonse := NewListMoviesResponse()
  err c.client.Send(request, response)

  return
}

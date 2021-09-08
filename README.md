# sendhttp

`sendhttp` can make the process of creating golang sdk for http service more easier.

This package defines three interfaces.

- `Request`
- `Response`
- `Client`

And four structs:

- `BaseRequest`
- `BaseResponse`
- `HttpClient`
- `RestyClient`

The `BaseRequest` struct have implemented the `Request` interface.
So all business request struct just needs to have the `*BaseRequest` nested.

The `BaseResponse` struct have implemented the `Response` interface.
So all business response struct just needs to have the `*BaseResponse` nested.

The `RestyClient` and `HttpClient` structs have implemented the `Client` interface which
defines a `Send(request Request, response Response) error` method.

- The `HttpClient` is based on the `http.Client` of the standard library `net/http`.
- The `RestyClient` is based on the `resty.Client` of the `github.com/go-resty/resty/v2`.

In your business code, you can directly embeds `RestyClient` or `HttpClient` in your own client struct.

## Features

Currently, supported content types of the HTTP Request are:

- "application/json"
- "application/x-www-form-urlencoded"
- "multipart/form-data"

And, only JSON format HTTP Response are supported. So your business response struct MUST can be unmarshaled from plain json.

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
// define the business request
type ListMoviesRequest struct {
  *sendhttp.BaseRequest

  // define other fields if the request contains json body
}

// define a function to create a business request object
func NewListMoviesRequest() *ListMoviesRequest{
  return &ListMoviesRequest{
    // Note, MUST init the BaseRequest, or else you might encounter runtime nil pointer panic
    BaseRequest: sendhttp.NewBaseRequest(),

    // init other fields if necessary
  }
}

// define the business response
type ListMoviesResonse struct {
  *sendhttp.BaseResponse

  // define fields which holds the response
  Items []Movie `json:"items"`
}

// define a function to create a business response object
func NewListMoviesResponse() *ListMoviesResonse{
  return &ListMoviesResonse{
    // Note, MUST init the BaseResponse, or else you might encounter runtime nil pointer panic
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

3. Define a method as your SDK exported method.

```go
func (c *Client) ListMovies(request *ListMoviesRequest) (response *ListMoviesResonse, err error) {

  // use the client information to fill the request
  // normally your will
  // c.Complete(request)

  resonse = NewListMoviesResponse()
  // The main logic of this function is to call the Send method provided by sendhttp.Client
  err = c.client.Send(request, response)

  return
}

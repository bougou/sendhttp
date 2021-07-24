package sendhttp

import "net/http"

type Response interface {
	SetRaw(body []byte)
	GetRaw() []byte
	Fill(response *http.Response)
}

// BaseResponse implements Response interface, and should be embeded into other concrete response structs.
type BaseResponse struct {
	status        string
	statusCode    int
	proto         string
	protoMajor    int
	protoMinor    int
	contentLength int64
	header        map[string][]string

	raw []byte
}

func NewBaseResponse() *BaseResponse {
	return &BaseResponse{
		raw:    make([]byte, 0),
		header: make(map[string][]string),
	}
}

func (r *BaseResponse) SetRaw(body []byte) {
	r.raw = body
}

func (r *BaseResponse) GetRaw() []byte {
	return r.raw
}

// Fill MUST not read response Body field
func (r *BaseResponse) Fill(response *http.Response) {
	r.status = response.Status
	r.statusCode = response.StatusCode
	r.proto = response.Proto
	r.protoMajor = response.ProtoMajor
	r.protoMinor = response.ProtoMinor
	r.contentLength = response.ContentLength
	r.header = response.Header
}

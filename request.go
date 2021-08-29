package sendhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	//"log"

	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	POST = "POST"
	GET  = "GET"

	HTTP  = "http"
	HTTPS = "https"

	Path = "/"
)

type Request interface {
	// GET params to POST body
	// GetBodyReader() io.Reader

	GetMethod() string
	GetUrl() string
	GetScheme() string
	GetDomain() string
	GetPath() string
	GetParams() map[string]string
	GetFormParams() map[string]string
	GetHeaders() map[string]string
	GetMultipart() []*multipartValue

	SetMethod(string)
	SetUrl(url string)
	SetScheme(string)
	SetDomain(string)
	SetPath(string)
	SetParams(map[string]string)
	SetFormParams(map[string]string)
	SetHeaders(map[string]string)
	AddMultipart(fieldname string, filename string, value io.Reader)

	CheckValid() error

	IsContentTypeForm() bool
	IsContentTypeJSON() bool
	IsContentTypeMultipart() bool
}

func GetBody(request Request) (io.Reader, error) {
	if request.IsContentTypeForm() {
		return strings.NewReader(GetUrlQueriesEncoded(request.GetFormParams())), nil
	}

	if request.IsContentTypeMultipart() {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		for _, m := range request.GetMultipart() {
			var formWriter io.Writer
			var err error
			if m.FileName != "" {
				formWriter, err = w.CreateFormFile(m.FieldName, m.FileName)
			} else {
				formWriter, err = w.CreateFormField(m.FieldName)
			}
			if err != nil {
				return nil, err
			}
			io.Copy(formWriter, m.Value)
		}
		// Reset request Content-Type Header, fill the boundary
		request.SetHeaders(map[string]string{
			"Content-Type": w.FormDataContentType(),
		})
		return buf, nil
	}

	// default json
	reqBody, err := json.Marshal(request)
	if err != nil {
		msg := fmt.Sprintf("marshal request failed, %s", err)
		return nil, errors.New(msg)
	}
	return bytes.NewReader(reqBody), nil
}

type BaseRequest struct {
	method string

	// url is splitted into scheme/domain/path
	scheme string
	domain string
	path   string
	params map[string]string

	formParams map[string]string

	headers map[string]string

	multipart []*multipartValue
}

type multipartValue struct {
	FieldName string // must
	FileName  string // pass a empty value if not a file
	Value     io.Reader
}

var _ Request = (*BaseRequest)(nil)

func NewBaseRequest() *BaseRequest {
	r := &BaseRequest{
		path:       Path,
		params:     make(map[string]string),
		formParams: make(map[string]string),
		headers:    make(map[string]string),
		multipart:  make([]*multipartValue, 0),
	}

	return r
}

func (r *BaseRequest) CheckValid() error {
	if r.method == "" {
		return errors.New("invalid request, no method set")
	}

	if r.scheme == "" {
		return errors.New("invalid request, no scheme set")
	}

	if r.domain == "" {
		return errors.New("invalid request, no domain set")
	}

	if r.path == "" {
		return errors.New("invalid request, no path set")
	}

	return nil
}

func (r *BaseRequest) GetMethod() string {
	return r.method
}

func (r *BaseRequest) GetUrl() string {
	s := r.scheme + "://" + r.domain + r.path
	if len(r.params) > 0 {
		s += "?" + GetUrlQueriesEncoded(r.params)
	}
	return s
}

func (r *BaseRequest) GetScheme() string {
	return r.scheme
}

func (r *BaseRequest) GetDomain() string {
	return r.domain
}

func (r *BaseRequest) GetPath() string {
	return r.path
}

func (r *BaseRequest) GetParams() map[string]string {
	return r.params
}

func (r *BaseRequest) GetFormParams() map[string]string {
	return r.formParams
}

func (r *BaseRequest) GetHeaders() map[string]string {
	return r.headers
}

func (r *BaseRequest) HasHeader(name string) bool {
	_, ok := r.GetHeaders()[name]
	return ok
}

func (r *BaseRequest) GetHeader(name string) string {
	s := r.GetHeaders()[name]
	return s
}

func (r *BaseRequest) GetMultipart() []*multipartValue {
	return r.multipart
}

// SetUrl set scheme/domain/path of r according to urlstr
func (r *BaseRequest) SetUrl(urlstr string) {
	u, err := url.Parse(urlstr)
	if err != nil {
		// do nothing if err
		return
	}
	r.scheme = u.Scheme
	r.domain = u.Host
	r.path = u.Path
}

func (r *BaseRequest) SetDomain(domain string) {
	r.domain = domain
}

func (r *BaseRequest) SetPath(path string) {
	r.path = path
}

func (r *BaseRequest) SetScheme(scheme string) {
	scheme = strings.ToLower(scheme)
	switch scheme {
	case HTTP:
		r.scheme = HTTP
	default:
		r.scheme = HTTPS
	}
}

func (r *BaseRequest) SetMethod(method string) {
	switch strings.ToUpper(method) {
	case POST:
		{
			r.method = POST
		}
	case GET:
		{
			r.method = GET
		}
	default:
		{
			r.method = GET
		}
	}
}

func (r *BaseRequest) SetParams(params map[string]string) {
	for k, v := range params {
		r.params[k] = v
	}
}

func (r *BaseRequest) SetFormParams(params map[string]string) {
	for k, v := range params {
		r.formParams[k] = v
	}
}

func (r *BaseRequest) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		r.headers[k] = v
	}
}

func (r *BaseRequest) AddMultipart(fieldname string, filename string, value io.Reader) {
	v := &multipartValue{
		FieldName: fieldname,
		FileName:  filename,
		Value:     value,
	}
	r.multipart = append(r.multipart, v)
}

func (r *BaseRequest) IsContentTypeForm() bool {
	return r.GetHeader("Content-Type") == "application/x-www-form-urlencoded"
}

func (r *BaseRequest) IsContentTypeJSON() bool {
	return r.GetHeader("Content-Type") == "application/json"
}

func (r *BaseRequest) IsContentTypeMultipart() bool {
	return r.GetHeader("Content-Type") == "multipart/form-data"
}

// todo
func ConstructParams(req Request) (err error) {
	value := reflect.ValueOf(req).Elem()
	err = flatStructure(value, req, "")
	return
}

func flatStructure(value reflect.Value, request Request, prefix string) (err error) {
	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		tag := valueType.Field(i).Tag
		nameTag, hasNameTag := tag.Lookup("name")
		if !hasNameTag {
			continue
		}
		field := value.Field(i)
		kind := field.Kind()
		if kind == reflect.Ptr && field.IsNil() {
			continue
		}
		if kind == reflect.Ptr {
			field = field.Elem()
			kind = field.Kind()
		}
		key := prefix + nameTag
		if kind == reflect.String {
			s := field.String()
			if s != "" {
				request.GetParams()[key] = s
			}
		} else if kind == reflect.Bool {
			request.GetParams()[key] = strconv.FormatBool(field.Bool())
		} else if kind == reflect.Int || kind == reflect.Int64 {
			request.GetParams()[key] = strconv.FormatInt(field.Int(), 10)
		} else if kind == reflect.Uint || kind == reflect.Uint64 {
			request.GetParams()[key] = strconv.FormatUint(field.Uint(), 10)
		} else if kind == reflect.Float64 {
			request.GetParams()[key] = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		} else if kind == reflect.Slice {
			list := value.Field(i)
			for j := 0; j < list.Len(); j++ {
				vj := list.Index(j)
				key := prefix + nameTag + "." + strconv.Itoa(j)
				kind = vj.Kind()
				if kind == reflect.Ptr && vj.IsNil() {
					continue
				}
				if kind == reflect.Ptr {
					vj = vj.Elem()
					kind = vj.Kind()
				}
				if kind == reflect.String {
					request.GetParams()[key] = vj.String()
				} else if kind == reflect.Bool {
					request.GetParams()[key] = strconv.FormatBool(vj.Bool())
				} else if kind == reflect.Int || kind == reflect.Int64 {
					request.GetParams()[key] = strconv.FormatInt(vj.Int(), 10)
				} else if kind == reflect.Uint || kind == reflect.Uint64 {
					request.GetParams()[key] = strconv.FormatUint(vj.Uint(), 10)
				} else if kind == reflect.Float64 {
					request.GetParams()[key] = strconv.FormatFloat(vj.Float(), 'f', -1, 64)
				} else {
					if err = flatStructure(vj, request, key+"."); err != nil {
						return
					}
				}
			}
		} else {
			if err = flatStructure(reflect.ValueOf(field.Interface()), request, prefix+nameTag+"."); err != nil {
				return
			}
		}
	}
	return
}

func GetUrlQueriesEncoded(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			// values.Set(key, value)
			values.Add(key, value)
		}
	}
	return values.Encode()
}

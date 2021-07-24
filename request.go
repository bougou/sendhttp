package sendhttp

import (
	"errors"

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

	GetHeaders() map[string]string

	SetUrl(url string)
	SetScheme(string)
	SetMethod(string)
	SetDomain(string)
	SetPath(string)

	SetHeaders(map[string]string)

	CheckValid() error
}

type BaseRequest struct {
	method string

	// url is splitted into scheme/domain/path
	scheme string
	domain string
	path   string

	params     map[string]string
	formParams map[string]string

	headers map[string]string
}

var _ Request = (*BaseRequest)(nil)

func NewBaseRequest() *BaseRequest {
	r := &BaseRequest{
		path:       Path,
		params:     make(map[string]string),
		formParams: make(map[string]string),
		headers:    make(map[string]string),
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
	if r.method == GET {
		return r.GetScheme() + "://" + r.domain + r.path + "?" + GetUrlQueriesEncoded(r.params)
	} else if r.method == POST {
		return r.GetScheme() + "://" + r.domain + r.path
	} else {
		return ""
	}
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

func (r *BaseRequest) GetHeaders() map[string]string {
	return r.headers
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

func (r *BaseRequest) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		r.headers[k] = v
	}
}

// todo
func ConstructParams(req Request) (err error) {
	value := reflect.ValueOf(req).Elem()
	err = flatStructure(value, req, "")
	//log.Printf("[DEBUG] params=%s", req.GetParams())
	return
}

func flatStructure(value reflect.Value, request Request, prefix string) (err error) {
	//log.Printf("[DEBUG] reflect value: %v", value.Type())
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
			values.Add(key, value)
		}
	}
	return values.Encode()
}

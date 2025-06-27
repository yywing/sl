package types

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/native"
	"github.com/yywing/sl/util"
)

const (
	TypeKindURL          = "url"
	TypeKindHTTPRequest  = "http_request"
	TypeKindHTTPResponse = "http_response"
)

var (
	URLType          = native.NewNativeSelectorType[*URL](TypeKindURL)
	HTTPRequestType  = native.NewNativeSelectorType[*HTTPRequest](TypeKindHTTPRequest)
	HTTPResponseType = native.NewNativeSelectorType[*HTTPResponse](TypeKindHTTPResponse)
)

type URL struct {
	Scheme   string `sl:"scheme"`
	Domain   string `sl:"domain"`
	Host     string `sl:"host"`
	Port     string `sl:"port"`
	Path     string `sl:"path"`
	Query    string `sl:"query"`
	Fragment string `sl:"fragment"`

	URL string
}

func NewURL(u string) *URL {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil
	}

	// 处理一下默认端口
	port := parsedURL.Port()
	if port == "" {
		port = strconv.Itoa(util.HttpSchemeToPort(parsedURL.Scheme))
		parsedURL.Host = parsedURL.Host + ":" + port
	}

	return &URL{
		Scheme:   parsedURL.Scheme,
		Domain:   parsedURL.Hostname(),
		Host:     parsedURL.Host,
		Port:     port,
		Path:     parsedURL.Path,
		Query:    parsedURL.RawQuery,
		Fragment: parsedURL.Fragment,
		URL:      u,
	}
}

func (v *URL) Type() ast.ValueType {
	return URLType
}

func (v *URL) String() string {
	return v.URL
}

func (v *URL) Equal(other ast.Value) bool {
	return false
}

func (v *URL) Get(key ast.Value) (ast.Value, bool) {
	switch key.Type().Kind() {
	case ast.TypeKindString:
		return URLType.Get(v, key.(*ast.StringValue).StringValue)
	default:
		return nil, false
	}
}

type HTTPRequest struct {
	URL       *URL   `sl:"url"`
	Raw       []byte `sl:"raw"`
	Method    string `sl:"method"`
	RawHeader []byte `sl:"raw_header"`
	Body      []byte `sl:"body"`

	ContentType string            `sl:"content_type"`
	Headers     map[string]string `sl:"headers"`
}

func NewHTTPRequestFromRequest(req *http.Request) (*HTTPRequest, error) {
	headers := make(map[string]string)
	var contentType string
	for k, v := range req.Header {
		if k == "Content-Type" {
			contentType = strings.Join(v, ",")
		}
		headers[k] = strings.Join(v, ",")
	}

	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	raw, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}

	rawHeader, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, err
	}

	return &HTTPRequest{
		URL:         NewURL(req.URL.String()),
		Raw:         raw,
		Method:      req.Method,
		Headers:     headers,
		Body:        body,
		ContentType: contentType,
		RawHeader:   rawHeader,
	}, nil
}

func (v *HTTPRequest) Type() ast.ValueType {
	return HTTPRequestType
}

func (v *HTTPRequest) String() string {
	return v.URL.URL
}

func (v *HTTPRequest) Equal(other ast.Value) bool {
	return false
}

func (v *HTTPRequest) Get(key ast.Value) (ast.Value, bool) {
	switch key.Type().Kind() {
	case ast.TypeKindString:
		return HTTPRequestType.Get(v, key.(*ast.StringValue).StringValue)
	default:
		return nil, false
	}
}

type HTTPResponse struct {
	URL       *URL   `sl:"url"`
	Raw       []byte `sl:"raw"`
	Status    int    `sl:"status"`
	RawHeader []byte `sl:"raw_header"`
	Body      []byte `sl:"body"`
	RawCert   []byte `sl:"raw_cert"`
	Cert      Cert   `sl:"cert"`

	Latency     int               `sl:"latency"`
	ContentType string            `sl:"content_type"`
	Headers     map[string]string `sl:"headers"`
	BodyString  string            `sl:"body_string"`
	Title       []byte            `sl:"title"`
	TitleString string            `sl:"title_string"`
}

func NewHTTPResponseFromResponse(resp *http.Response) (*HTTPResponse, error) {
	req := util.GetFristRequest(resp)

	return &HTTPResponse{
		// TODO: 空指针
		URL: NewURL(req.URL.String()),
		// TODO: 其它属性
	}, nil
}

func (v *HTTPResponse) Type() ast.ValueType {
	return HTTPResponseType
}

func (v *HTTPResponse) String() string {
	return v.URL.URL
}

func (v *HTTPResponse) Equal(other ast.Value) bool {
	return false
}

func (v *HTTPResponse) Get(key ast.Value) (ast.Value, bool) {
	switch key.Type().Kind() {
	case ast.TypeKindString:
		return HTTPResponseType.Get(v, key.(*ast.StringValue).StringValue)
	default:
		return nil, false
	}
}

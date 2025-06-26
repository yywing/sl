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
	TypeKindURL         = "url"
	TypeKindHTTPRequest = "http_request"
)

var (
	URLType         = native.NewNativeSelectorType[*URL](TypeKindURL)
	HTTPRequestType = native.NewNativeSelectorType[*HTTPRequestValue](TypeKindHTTPRequest)
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

type HTTPRequestValue struct {
	URL     *URL              `sl:"url"`
	Raw     []byte            `sl:"raw"`
	Method  string            `sl:"method"`
	Headers map[string]string `sl:"headers"`
	Body    []byte            `sl:"body"`

	ContentType string `sl:"content_type"`
	RawHeader   []byte `sl:"raw_header"`
}

func NewHTTPRequestValueFromRequest(req *http.Request) (*HTTPRequestValue, error) {
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

	return &HTTPRequestValue{
		URL:         NewURL(req.URL.String()),
		Raw:         raw,
		Method:      req.Method,
		Headers:     headers,
		Body:        body,
		ContentType: contentType,
		RawHeader:   rawHeader,
	}, nil
}

func (v *HTTPRequestValue) Type() ast.ValueType {
	return HTTPRequestType
}

func (v *HTTPRequestValue) String() string {
	return v.URL.URL
}

func (v *HTTPRequestValue) Equal(other ast.Value) bool {
	return false
}

func (v *HTTPRequestValue) Get(key ast.Value) (ast.Value, bool) {
	switch key.Type().Kind() {
	case ast.TypeKindString:
		return HTTPRequestType.Get(v, key.(*ast.StringValue).StringValue)
	default:
		return nil, false
	}
}

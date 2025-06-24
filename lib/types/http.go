package types

import (
	"io"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/native"
)

const (
	TypeKindHTTPRequest = "http_request"
)

var (
	HTTPRequestType = native.NewNativeSelectorType[*HTTPRequestValue](TypeKindHTTPRequest)
)

type HTTPRequestValue struct {
	URL     string            `sl:"url"`
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

	raw, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	rawHeader, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, err
	}

	return &HTTPRequestValue{
		URL:         req.URL.String(),
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
	return v.URL
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

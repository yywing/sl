package types

import (
	"crypto/x509"
	"strings"

	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/native"
)

const (
	TypeKindCert = "cert"
)

var (
	CertType = native.NewNativeSelectorType[*Cert](TypeKindCert)
)

type Cert struct {
	Issuer       string          `sl:"issuer"`
	Subject      string          `sl:"subject"`
	DNSNames     string          `sl:"dnsnames"`
	NotAfter     *TimestampValue `sl:"not_after"`
	NotBefore    *TimestampValue `sl:"not_before"`
	SerialNumber string          `sl:"serial_number"`
}

func NewCert(c *x509.Certificate) *Cert {
	return &Cert{
		Issuer:       c.Issuer.String(),
		Subject:      c.Subject.String(),
		DNSNames:     strings.Join(c.DNSNames, ","),
		NotAfter:     NewTimestampValueWithTime(&c.NotAfter),
		NotBefore:    NewTimestampValueWithTime(&c.NotBefore),
		SerialNumber: c.SerialNumber.String(),
	}
}

func (v *Cert) Type() ast.ValueType {
	return CertType
}

func (v *Cert) String() string {
	return ""
}

func (v *Cert) Equal(other ast.Value) bool {
	return false
}

func (v *Cert) Get(key ast.Value) (ast.Value, bool) {
	switch key.Type().Kind() {
	case ast.TypeKindString:
		return CertType.Get(v, key.(*ast.StringValue).StringValue)
	default:
		return nil, false
	}
}

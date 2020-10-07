package conflate

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func init() {
	// annoyingly the format checker list is a global variable
	gojsonschema.FormatCheckers.Add(newXMLFormatChecker("xml-template"))
	gojsonschema.FormatCheckers.Add(newHTMLFormatChecker("html-template"))
	gojsonschema.FormatCheckers.Add(newRegexFormatChecker("regex"))
	gojsonschema.FormatCheckers.Add(newCryptoFormatChecker("pkcs1-private-key", pkcs1PrivateKey))
	gojsonschema.FormatCheckers.Add(newCryptoFormatChecker("pkcs1-public-key", pkcs1PublicKey))
	gojsonschema.FormatCheckers.Add(newCryptoFormatChecker("pkcs8-private-key", pkcs8PrivateKey))
	gojsonschema.FormatCheckers.Add(newCryptoFormatChecker("pkcs8-public-key", pkixPublicKey)) // deprecated, use pkix-public-key
	gojsonschema.FormatCheckers.Add(newCryptoFormatChecker("pkix-public-key", pkixPublicKey))
	gojsonschema.FormatCheckers.Add(newCryptoFormatChecker("x509-certificate", x509Certificate))
}

// ----------------

type formatErrors map[string]error

var formatErrs = formatErrors{}

func (errs formatErrors) clear() {
	formatErrs = formatErrors{}
}

func (errs formatErrors) add(name interface{}, value interface{}, err error) {
	errs[errs.key(name, value)] = err
}

func (errs formatErrors) get(name interface{}, value interface{}) error {
	return errs[errs.key(name, value)]
}

func (errs formatErrors) key(name interface{}, value interface{}) string {
	return fmt.Sprintf("%v#%v", name, value)
}

// ----------------

type xmlFormatChecker struct {
	tags *regexp.Regexp
	name string
}

func newXMLFormatChecker(name string) (string, gojsonschema.FormatChecker) {
	return name, xmlFormatChecker{name: name, tags: regexp.MustCompile(`{{[^{}]*}}`)}
}

func (f xmlFormatChecker) IsFormat(input interface{}) bool {
	var err error

	if s, ok := input.(string); ok {
		s = f.tags.ReplaceAllString(s, "")
		if len(s) > 0 {
			var v interface{}
			err = xml.Unmarshal([]byte(s), &v)
			err = wrapError(err, "Failed to parse xml")
		}
	} else {
		err = makeError("The value is not a string")
	}
	if err != nil {
		formatErrs.add(f.name, input, err)
		return false
	}
	return true
}

// ----------------

type htmlFormatChecker struct {
	tags *regexp.Regexp
	name string
}

func newHTMLFormatChecker(name string) (string, gojsonschema.FormatChecker) {
	return name, htmlFormatChecker{name: name, tags: regexp.MustCompile(`{{[^{}]*}}`)}
}

func (f htmlFormatChecker) IsFormat(input interface{}) bool {
	var err error

	if s, ok := input.(string); ok {
		s = f.tags.ReplaceAllString(s, "")
		_, err = html.Parse(strings.NewReader(s))
		err = wrapError(err, "Failed to parse html")
	} else {
		err = makeError("The value is not a string")
	}
	if err != nil {
		formatErrs.add(f.name, input, err)
		return false
	}
	return true
}

// ----------------

type cryptoFormatChecker struct {
	name  string
	cType cryptoType
}

type cryptoType int

const (
	pkcs1PrivateKey cryptoType = 1 + iota
	pkcs1PublicKey
	pkcs8PrivateKey
	pkixPublicKey
	x509Certificate
)

func newCryptoFormatChecker(name string, cType cryptoType) (string, gojsonschema.FormatChecker) {
	return name, cryptoFormatChecker{
		name:  name,
		cType: cType,
	}
}

func (f cryptoFormatChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		formatErrs.add(f.name, input, makeError("The value is not a string"))
		return false
	}

	var err error

	var data []byte

	block, _ := pem.Decode([]byte(s))
	if block != nil {
		data = block.Bytes
	} else {
		// Try to directly base64 decode if not valid PEM
		data, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			formatErrs.add(f.name, input, wrapError(err, "Failed to decode the data"))
			return false
		}
	}

	switch f.cType {
	case pkcs1PrivateKey:
		_, err = x509.ParsePKCS1PrivateKey(data)
	case pkcs1PublicKey:
		_, err = x509.ParsePKCS1PublicKey(data)
	case pkcs8PrivateKey:
		_, err = x509.ParsePKCS8PrivateKey(data)
	case pkixPublicKey:
		_, err = x509.ParsePKIXPublicKey(data)
	case x509Certificate:
		_, err = x509.ParseCertificate(data)
	default:
		err = makeError(f.name + " called with unsupported type")
	}

	if err != nil {
		formatErrs.add(f.name, input, wrapError(err, "Failed to parse key"))
		return false
	}

	return true
}

// ----------------

type regexFormatChecker struct{ name string }

func newRegexFormatChecker(name string) (string, gojsonschema.FormatChecker) {
	return name, regexFormatChecker{name: name}
}

func (f regexFormatChecker) IsFormat(input interface{}) bool {
	var err error

	if s, ok := input.(string); ok {
		_, err = regexp.Compile(s)
		err = wrapError(err, "Failed to parse regular expression")
	} else {
		err = makeError("The value is not a string")
	}

	if err != nil {
		formatErrs.add(f.name, input, err)
		return false
	}

	return true
}

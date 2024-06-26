package conflate

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/net/html"
)

var (
	errRequiredString  = errors.New("the value is not a string")
	errUnsupportedType = errors.New("called with unsupported type")
)

func initFormatCheckers() {
	// annoyingly the format checker list is a global variable
	gojsonschema.FormatCheckers.Add(newXMLFormatChecker("xml"))
	gojsonschema.FormatCheckers.Add(newXMLTemplateFormatChecker("xml-template"))
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

func (errs formatErrors) add(name, value interface{}, err error) {
	errs[errs.key(name, value)] = err
}

func (errs formatErrors) get(name, value interface{}) error {
	return errs[errs.key(name, value)]
}

func (errs formatErrors) key(name, value interface{}) string {
	return fmt.Sprintf("%v#%v", name, value)
}

// ----------------

type xmlFormatChecker struct{ name string }

//nolint:unparam // left for extensibility
func newXMLFormatChecker(name string) (string, gojsonschema.FormatChecker) {
	return name, xmlFormatChecker{name: name}
}

func (f xmlFormatChecker) IsFormat(input interface{}) bool {
	var ferr error

	if s, ok := input.(string); ok {
		if err := xml.Unmarshal([]byte(s), new(interface{})); err != nil {
			ferr = fmt.Errorf("failed to parse xml: %w", err)
		}
	} else {
		ferr = errRequiredString
	}

	if ferr != nil {
		formatErrs.add(f.name, input, ferr)

		return false
	}

	return true
}

// ----------------

type xmlTemplateFormatChecker struct {
	tags *regexp.Regexp
	name string
}

func newXMLTemplateFormatChecker(name string) (string, gojsonschema.FormatChecker) {
	return name, xmlTemplateFormatChecker{name: name, tags: regexp.MustCompile(`{{[^{}]*}}`)}
}

func (f xmlTemplateFormatChecker) IsFormat(input interface{}) bool {
	var ferr error

	if s, ok := input.(string); ok {
		s = f.tags.ReplaceAllString(s, "")
		if s != "" {
			var v interface{}
			if err := xml.Unmarshal([]byte(s), &v); err != nil {
				ferr = fmt.Errorf("failed to parse xml: %w", err)
			}
		}
	} else {
		ferr = errRequiredString
	}

	if ferr != nil {
		formatErrs.add(f.name, input, ferr)

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
	var ferr error

	if s, ok := input.(string); ok {
		s = f.tags.ReplaceAllString(s, "")

		if _, err := html.Parse(strings.NewReader(s)); err != nil {
			ferr = fmt.Errorf("failed to parse html: %w", err)
		}
	} else {
		ferr = errRequiredString
	}

	if ferr != nil {
		formatErrs.add(f.name, input, ferr)

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
		formatErrs.add(f.name, input, errRequiredString)

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
			formatErrs.add(f.name, input, fmt.Errorf("failed to decode the data: %w", err))

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
		err = fmt.Errorf("%v %w", f.name, errUnsupportedType)
	}

	if err != nil {
		formatErrs.add(f.name, input, fmt.Errorf("failed to parse key: %w", err))

		return false
	}

	return true
}

// ----------------

type regexFormatChecker struct{ name string }

//nolint:unparam // left for extensibility
func newRegexFormatChecker(name string) (string, gojsonschema.FormatChecker) {
	return name, regexFormatChecker{name: name}
}

func (f regexFormatChecker) IsFormat(input interface{}) bool {
	var ferr error

	if s, ok := input.(string); ok {
		if _, err := regexp.Compile(s); err != nil {
			ferr = fmt.Errorf("failed to parse regular expression: %w", err)
		}
	} else {
		ferr = errRequiredString
	}

	if ferr != nil {
		formatErrs.add(f.name, input, ferr)

		return false
	}

	return true
}

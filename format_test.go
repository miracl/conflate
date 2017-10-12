package conflate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatErrors_Get(t *testing.T) {
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	errExpected := makeError("An error")
	formatErrs["name#value"] = errExpected
	err := formatErrs.get("name", "value")
	assert.Equal(t, errExpected, err)
}

func TestFormatErrors_Add(t *testing.T) {
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	errExpected := makeError("An error")
	formatErrs.add("name", "value", errExpected)
	assert.Equal(t, errExpected, formatErrs["name#value"])
}

func TestFormatErrors_Clear(t *testing.T) {
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	formatErrs["stuff"] = nil
	formatErrs.clear()
	assert.Empty(t, formatErrs)
}

// --------

func TestXmlFormatCheckerIsFormat_NotString(t *testing.T) {
	givenName := "xml"
	givenValue := 1
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newXMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The value is not a string")
}

func TestXmlFormatCheckerIsFormat_Valid(t *testing.T) {
	givenName := "xml"
	givenValue := "<test>{{.Value}}</test>"
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newXMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.True(t, result)
	err := formatErrs.get(name, givenValue)
	assert.Nil(t, err)
}

func TestXmlFormatCheckerIsFormat_NotValid(t *testing.T) {
	givenName := "xml"
	givenValue := "<test>"
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newXMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to parse xml")
}

// --------

func TestHtmlFormatCheckerIsFormat_NotString(t *testing.T) {
	givenName := "html"
	givenValue := 1
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newHTMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The value is not a string")
}

func TestHtmlFormatCheckerIsFormat_Valid(t *testing.T) {
	givenName := "html"
	givenValue := "<html>{{.Value}}</html>"
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newHTMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.True(t, result)
	err := formatErrs.get(name, givenValue)
	assert.Nil(t, err)
}

/*
func TestHtmlFormatCheckerIsFormat_NotValid(t *testing.T) {
	givenName := "html"
	givenValue := "<!DOCTYPE wibble> </html/>"
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newHTMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to parse html")
}
*/
// --------

var testCryptoFormatCheckerKeys = map[cryptoType]string{
	pkcs1PrivateKey: `MIIEogIBAAKCAQEA+4pTx7LDtTdXrMCYYZUHbLkIyOV/DPwJVJo4V7CPjv9CxKT/+1s75Zf0Ek64I6UFoK39Adlneqnhw6OShkM9UlBXOWVN+eVucKo4WeuiEwYgmyPpFq3yYt1nVIZMeNcho9AramnCS2D0wSXdSQLCtpGTeKbPklDS8CH6Br/fXTgcMQs4PlRuJt6utggjTHfk4cuQUMB1xuvqCdNPvlz/IQRb+pkEcnE5tw41WTOoSySV7RKThsIy3fP14UZrse+f6rZ5rqm5XreHV7FR4eN04gY06Kf2kPPPoqDX5tsRS7MGsK1tuSabahSRQzpmRUkXQBx0H5FusmyUZvkBxH0KRQIDAQABAoIBAQCmJUAhb5qFYACxQqVkDyPQVIfQ6oV49iDTmOoOQbkRItnfOX07cY+mny+/x+8o3B9TepjbV9ZZ10wgBTMLK+6dKCP07U0L+tdv439FJbbwCeQPlTCookyvsL5Hvk1UpAS8RwvmReToFSqVSYVYbFJgdNeLoeMJKgmqQ0b6nlHTJ6V6anwjWHN9yWB/2T1q0VALFMNeQS3YHjrI3ldzQ37KuyPWWahh+OTXcFP22e45aJ68cXqaX4WPbDTEzAUuHYXl/vfas/2+qsnwIlGpBnEhm4EMT8DRzel9IW/i8ggz7nfc0167bK931vQGIR6ofFhizsZEot57wLEHjCgvdOUhAoGBAP+nVJtEcR/aH+p1Fb4E7srTTUh9vr60qpFqm95YlDGFzKrx5FqN83kf5PiN5oAPzKyL5xQBITJ+1Evp6ptJn+Q6a1tqEgp+4vhe4mqAuefgUZ/zVZwsGPO9WPYyuQxsJF/4eO+XEPgctZGVoUmw9R3UO3TrhMtBn04ZNDnxrVEfAoGBAPvhkfSoAz9oQFyWhbDn4VHnq3jF+CJRWBXrMQG765ymXxwvR24R5vqjnuFAjPMhC5kcMGyTJbWUPlgQnjoijxswVer5D7VjNbtVFjz+I8QfrFWD6ZN9OkQOeK1N7gFufeF12xnjYI+8a18TDMJAzmj8XytG63R4C5ovUSdq7wQbAoGAVT8ZojB/PCzCqa5jEpqYeX9d7McpPHQH0DdpRAnrWFwSNAo2p89fmUr+UbcXXipmiD6aTfxWcn0CE1IrjZ5ON48XT0MeQuuuiT0yaGsGEoClFx3PtSSrKVNA/89IDxFcS7gRs2p/GQyRqrrnLihYg5rKPKFwBUqbcTJMTOn+becCfzYdQD5P2mLFAw7hR9e5+a6fFzfbUHQPBd2vqde1h+kf1U842R/MuNMgVAIQ3Ddf4h4C8GRjOLbKIprI9zWGNzigKdVRrQ5LQj/9B8oSz5xTMbMtpAEL5ni4ozSYiYnFM0Y9C9WEBDrdQDQs45DYW4AyuD+T/QIIyHVXtRfC4ZcCgYEA+dgnLtWWBojqexyEiCt9vCVMcdtek+GrkBWvyGF2EoFynoLcTMz0xcSISQjYn47TxKhIBbt3RQCoiBI10gBnm2p+7yD2zdG6XQF7ckQji1LvuidEcGbKGvVg/WmSNZdPbcASLO7r0w9BGXpbiLJCSR7b9jTbIo/2yUOYoDTSLGQ=`,
	pkcs1PublicKey:  `MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwYoff7zCqefmQXkLDSN/9lDSB0i/IAPKAEL0hNR8Zr6UtfFhW0zvzs4eEAyOJVFNjkmV/NE7pLuwnFMj0neeWsw/e7H8mQwcDtnCBc4/a3V6izk+IQ+WE1Cydhq+NDjmG1swxpiJm/NWLmAx4ySYxgL6JFpRMZ0ZDYtAh5dVUgBX+OpKqrKS1EAlKGmYFPdrsEuspevHquFCYUUB8KvAhaG4pCXTV2xYzFq9LgmurzLHJP/rtd5SLMGwuDOLSIpUvziNNyV3hbHiPAlG6QGifp6g0ioNy6dAqlBN8Goe3ZY5+lJTQL1MM9bGv53bOHyxiFpsdQRMhcTfRAxpxKePt1eOvt1i034rII4CsDfNduPIB1goKSTxr/d4e+FmKo07ZpV9Z1WsM7UrU4R8PE4ZRsnBrcIo6OeKpNMJ4zpUQlk/f3xptIzH32EKYGRqT2yBiy8b0fuk+UKDw+J0otXZVUKtnwLwdhU8bal0+pYvffPW25M/MGGMS7zgHMQrrYL31oIN/bXxA60YF+eDsi8ffVgQT6Xm+qftMmqtgDHTSO8kT/bqsoboPRMU6jtKY/+zJdj3SI6Ji4PHNpBlWO2nrRQ2DWcab08Hb78oHVEGnrGszqBNBHZuFnuhKCuGBfP3kBV5h0xda1fbJkEQZRaidD85sNcfnv6QyECwhGeQIq8CAwEAAQ==`,
	pkcs8PrivateKey: `MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAN/eySAfwBAPSocWY3f6ugcZ0zCTKO825kLQ5AUNan976RZOg+lI1ZBAbcOnp25+r+tEmdcWdjkl76pY8miyfqCyByOwzAS3fZm0H2ZxaDTpSraupG449HmP9hOjT+MEqd94WpGxVDUVv2nuN9so2Giq+Ua9JKOTVZeDFFcEyM/pAgMBAAECgYBR0mpuLEyVUhFqODuYsXHmgIDYwyB19fnIt/JvCk0/VPmHJqE91VoBBBtaJF+kmfcQJR2ZKGODVaM3/pRywbJXuAQBXn5J+Dyo352C6oe51GwZr/gYYDLl9JLE8tEhZRS4gUHYLS6WqMZQiQDkj5PljJ2l6F541MA9QolmBOEjgQJBAP7Dcp4dXOX+2ajgpTgJeExpPtrf+wsEe8sCJilaEbytDla6VH+rJXrx8AqNGrwz85cv/C3gOwHi6gw4+XLRKNECQQDg9PO5528NB2wGmkRZiyaL88vPQnE7LS0Wx6+kdAknWK5lmyLscHjJw6PmHCFd1c+xIz4/UeeesyffKSpPO3uZAkANVmAexPzDEbeFbtVXfK9umKfwv38eEYsFksQ6r/tzvD+e7LxVvrkEmbvuYXq/1ZxwEtQJ9s5ACAgmwajViPuxAkAM5H0DbowixwYd6sF4EE2JP9OycTBkH9axs4gReMT9iEuWbym1O0zw41GFYp9W7WYV1NtAbcVEhygF1ioimFohAkEA2KbvdSARygUfZqRrDiC6FipRaVrsO18ZOMqteXfxzthEYVtLJFk/1gO77gVsKfoeIQf8GrznROIh9bPpRUlhEA==`,
	pkcs8PublicKey:  `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1eSZbr9I6ZRts+ByRmKnYYcSm/S/IvvxWZKuRL1t4kQns/nurhAM6rAu1E4OH75yAj/SF7FnpEzZ617oak6zSNz7Ow/QZ9APCtA8GpDEkg4p1O2hldIbm3ww8T1cQ6wyy77pQbSWMB70FSNV0CQX57v9e6TxHalMJ5PYvSucTwSiLkBye5uk2ktwnoqP/NNGC3/nr0NXGjeYHWqTIa/gf4GX2jngrm7NsZK4l30IL0AN7wQ07+FFwgJYuQCYCgmMb1Wa1o7hI6zIpCSf9bIlVgPNck/QV4lP/H40zj5S8VeSPysJtV2LhscH7erY9RWEb01/kYVbJpHWhAbTPxolJQIDAQAB`,
	x509Certificate: `MIIDXTCCAkWgAwIBAgIJAISKTtZx8oAqMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMTcxMjEzMTAyODE5WhcNMTgwMTEyMTAyODE5WjBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+eKfL1lMFOVP4L3LfNcKueYX0Spl2x1+B9UFufFyF/kC7hI6KjevMA5Y7HgSImSbajpa6B8c4iAdtGvG2p/DBddClv5Xq9wCTZffDxQbdb6syUfUo6/22hvCFDx/kpYmyLpq8r2AMAAeDsOmiO9pybhkMi4VHDs5DC/ptKR4h3w6T2c7RGb65kOO94rFU9lzITlIQOrgRPVUt38n5V2zVqWuUjU214CdSNHLiFhMaTH7WI/F/7JDQ2KI44Bd5htoZfgqNRC453KdDuUJO3kK/tSW2BXECTfR5ZstSJF34sMe4PBKOvg9xS5sr6CFe4VdTjHRzt/ELe05JYGoTjU3fwIDAQABo1AwTjAdBgNVHQ4EFgQUT4MTeTRAnnhjjQkeFAh1on6ZqCowHwYDVR0jBBgwFoAUT4MTeTRAnnhjjQkeFAh1on6ZqCowDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAadVMg8QHAPTpx+MXbXOTXWQUYIWZ4CJz0+rXQjUwgRygYd65CNO4cI9yk9w0Fx7GYvEbY6pL0iG3PzZs9nvG5SgxRrBxk7jbBbYaBLL0tCL+c0DryR2DOfbiPOcKBQGX6SULNcnAJGxkTwbhk7tiisbOxnW0XtAezcyQvCDWGvo55rHNTfZ8iBvIjh5aUofwYplG1xrJtPOY5OhQDzdsV4d/A9rdUmaamES0eA1EASRPs4hs0HWYq/dNu22VREyRQLuepkoYvrFYQFgwM8qU3kE3/4R4ZGTK2JuHfnM75hr6y6X9RHwJVyUMw45KgfwPCITUYOva6kdfudrlHYRR6A==`,
}

func testCryptoFormatCheckerIsFormatNotString(t *testing.T, cryptoType cryptoType) {
	givenName := "crypto"
	givenValue := 1
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The value is not a string")
}

func testCryptoFormatCheckerIsFormatValid(t *testing.T, cryptoType cryptoType) {
	givenName := "crypto"
	givenValue := testCryptoFormatCheckerKeys[cryptoType]
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.True(t, result)
	err := formatErrs.get(name, givenValue)
	assert.Nil(t, err)
}

func testCryptoFormatCheckerIsFormatNotValid(t *testing.T, cryptoType cryptoType) {
	givenName := "crypto"
	givenValue := testCryptoFormatCheckerKeys[cryptoType][0:16]
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to parse")
}

func TestCryptoFormatCheckerIsFormat_NotString(t *testing.T) {
	for cryptoType := range testCryptoFormatCheckerKeys {
		testCryptoFormatCheckerIsFormatNotString(t, cryptoType)
	}
}

func TestCryptoFormatCheckerIsFormat_Valid(t *testing.T) {
	for cryptoType := range testCryptoFormatCheckerKeys {
		testCryptoFormatCheckerIsFormatValid(t, cryptoType)
	}
}

func TestCryptoFormatCheckerIsFormat_NotValid(t *testing.T) {
	for cryptoType := range testCryptoFormatCheckerKeys {
		testCryptoFormatCheckerIsFormatNotValid(t, cryptoType)
	}
}

func TestCryptoFormatCheckerIsFormat_NotBase64(t *testing.T) {
	givenName := "crypto"
	givenValue := "not base-64"
	cryptoType := cryptoType(9999)
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to base-64 decode the data")
}

func TestCryptoFormatCheckerIsFormat_UnsupportedType(t *testing.T) {
	givenName := "crypto"
	givenValue := ""
	cryptoType := cryptoType(9999)
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

// --------

func TestRegexFormatCheckerIsFormat_NotString(t *testing.T) {
	givenName := "regex"
	givenValue := 1
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newRegexFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The value is not a string")
}

func TestRegexFormatCheckerIsFormat_Valid(t *testing.T) {
	givenName := "regex"
	givenValue := "^.*$"
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newRegexFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.True(t, result)
	err := formatErrs.get(name, givenValue)
	assert.Nil(t, err)
}

func TestRegexFormatCheckerIsFormat_NotValid(t *testing.T) {
	givenName := "regex"
	givenValue := "^(.*$"
	formatErrs.clear()
	defer func() { formatErrs.clear() }()
	name, checker := newRegexFormatChecker(givenName)
	assert.Equal(t, givenName, name)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)
	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to parse regular expression")
}

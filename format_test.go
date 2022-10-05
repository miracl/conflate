package conflate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	cryptoName = "crypto"
	htmlName   = "html"
	regexName  = "regex"
	xmlName    = "xml"
)

func TestFormatErrors_Get(t *testing.T) {
	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	formatErrs["name#value"] = errTest

	err := formatErrs.get("name", "value")
	assert.Equal(t, errTest, err)
}

func TestFormatErrors_Add(t *testing.T) {
	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	formatErrs.add("name", "value", errTest)
	assert.Equal(t, errTest, formatErrs["name#value"])
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
	givenName := xmlName
	givenValue := 1

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newXMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the value is not a string")
}

func TestXmlFormatCheckerIsFormat_Valid(t *testing.T) {
	givenName := xmlName
	givenValue := "<test>Value</test>"

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
	givenName := xmlName
	givenValue := "<test1>"

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newXMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to parse xml")
}

// --------

func TestXmlTemplateFormatCheckerIsFormat_NotString(t *testing.T) {
	givenName := xmlName
	givenValue := 1

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newXMLTemplateFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the value is not a string")
}

func TestXmlTemplateFormatCheckerIsFormat_Valid(t *testing.T) {
	givenName := xmlName
	givenValue := "<test>{{.Value}}</test>"

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newXMLTemplateFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.True(t, result)

	err := formatErrs.get(name, givenValue)
	assert.Nil(t, err)
}

func TestXmlTemplateFormatCheckerIsFormat_NotValid(t *testing.T) {
	givenName := xmlName
	givenValue := "<test>"

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newXMLTemplateFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to parse xml")
}

// --------

func TestHtmlFormatCheckerIsFormat_NotString(t *testing.T) {
	givenName := htmlName
	givenValue := 1

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newHTMLFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the value is not a string")
}

func TestHtmlFormatCheckerIsFormat_Valid(t *testing.T) {
	givenName := htmlName
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

func testCryptoFormatCheckerIsFormatNotString(t *testing.T, cryptoType cryptoType) {
	givenName := cryptoName
	givenValue := 1

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the value is not a string")
}

func TestCryptoFormatCheckerIsFormat_NotString(t *testing.T) {
	testCryptoFormatCheckerIsFormatNotString(t, pkcs1PrivateKey)
	testCryptoFormatCheckerIsFormatNotString(t, pkcs1PublicKey)
	testCryptoFormatCheckerIsFormatNotString(t, pkcs8PrivateKey)
	testCryptoFormatCheckerIsFormatNotString(t, pkixPublicKey)
	testCryptoFormatCheckerIsFormatNotString(t, x509Certificate)
}

func testCryptoFormatCheckerIsFormatNotValid(t *testing.T, cryptoType cryptoType) {
	givenName := cryptoName
	givenValue := "dGhpcyBpcyBub3QgYSB2YWxpZCBjZXJ0aWZpY2F0ZQo="

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestCryptoFormatCheckerIsFormat_NotValid(t *testing.T) {
	testCryptoFormatCheckerIsFormatNotValid(t, pkcs1PrivateKey)
	testCryptoFormatCheckerIsFormatNotValid(t, pkcs1PublicKey)
	testCryptoFormatCheckerIsFormatNotValid(t, pkcs8PrivateKey)
	testCryptoFormatCheckerIsFormatNotValid(t, pkixPublicKey)
	testCryptoFormatCheckerIsFormatNotValid(t, x509Certificate)
}

func testCryptoFormatCheckerIsFormatValid(t *testing.T, cryptoType cryptoType, givenValue string) {
	givenName := cryptoName

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.True(t, result)

	err := formatErrs.get(name, givenValue)
	assert.Nil(t, err)
}

//nolint:funlen // test case
func TestCryptoFormatCheckerIsFormat_Valid(t *testing.T) {
	testCryptoFormatCheckerKeys := []struct {
		ct  cryptoType
		val string
	}{
		{
			ct:  pkcs1PrivateKey,
			val: `MIIEogIBAAKCAQEA+4pTx7LDtTdXrMCYYZUHbLkIyOV/DPwJVJo4V7CPjv9CxKT/+1s75Zf0Ek64I6UFoK39Adlneqnhw6OShkM9UlBXOWVN+eVucKo4WeuiEwYgmyPpFq3yYt1nVIZMeNcho9AramnCS2D0wSXdSQLCtpGTeKbPklDS8CH6Br/fXTgcMQs4PlRuJt6utggjTHfk4cuQUMB1xuvqCdNPvlz/IQRb+pkEcnE5tw41WTOoSySV7RKThsIy3fP14UZrse+f6rZ5rqm5XreHV7FR4eN04gY06Kf2kPPPoqDX5tsRS7MGsK1tuSabahSRQzpmRUkXQBx0H5FusmyUZvkBxH0KRQIDAQABAoIBAQCmJUAhb5qFYACxQqVkDyPQVIfQ6oV49iDTmOoOQbkRItnfOX07cY+mny+/x+8o3B9TepjbV9ZZ10wgBTMLK+6dKCP07U0L+tdv439FJbbwCeQPlTCookyvsL5Hvk1UpAS8RwvmReToFSqVSYVYbFJgdNeLoeMJKgmqQ0b6nlHTJ6V6anwjWHN9yWB/2T1q0VALFMNeQS3YHjrI3ldzQ37KuyPWWahh+OTXcFP22e45aJ68cXqaX4WPbDTEzAUuHYXl/vfas/2+qsnwIlGpBnEhm4EMT8DRzel9IW/i8ggz7nfc0167bK931vQGIR6ofFhizsZEot57wLEHjCgvdOUhAoGBAP+nVJtEcR/aH+p1Fb4E7srTTUh9vr60qpFqm95YlDGFzKrx5FqN83kf5PiN5oAPzKyL5xQBITJ+1Evp6ptJn+Q6a1tqEgp+4vhe4mqAuefgUZ/zVZwsGPO9WPYyuQxsJF/4eO+XEPgctZGVoUmw9R3UO3TrhMtBn04ZNDnxrVEfAoGBAPvhkfSoAz9oQFyWhbDn4VHnq3jF+CJRWBXrMQG765ymXxwvR24R5vqjnuFAjPMhC5kcMGyTJbWUPlgQnjoijxswVer5D7VjNbtVFjz+I8QfrFWD6ZN9OkQOeK1N7gFufeF12xnjYI+8a18TDMJAzmj8XytG63R4C5ovUSdq7wQbAoGAVT8ZojB/PCzCqa5jEpqYeX9d7McpPHQH0DdpRAnrWFwSNAo2p89fmUr+UbcXXipmiD6aTfxWcn0CE1IrjZ5ON48XT0MeQuuuiT0yaGsGEoClFx3PtSSrKVNA/89IDxFcS7gRs2p/GQyRqrrnLihYg5rKPKFwBUqbcTJMTOn+becCfzYdQD5P2mLFAw7hR9e5+a6fFzfbUHQPBd2vqde1h+kf1U842R/MuNMgVAIQ3Ddf4h4C8GRjOLbKIprI9zWGNzigKdVRrQ5LQj/9B8oSz5xTMbMtpAEL5ni4ozSYiYnFM0Y9C9WEBDrdQDQs45DYW4AyuD+T/QIIyHVXtRfC4ZcCgYEA+dgnLtWWBojqexyEiCt9vCVMcdtek+GrkBWvyGF2EoFynoLcTMz0xcSISQjYn47TxKhIBbt3RQCoiBI10gBnm2p+7yD2zdG6XQF7ckQji1LvuidEcGbKGvVg/WmSNZdPbcASLO7r0w9BGXpbiLJCSR7b9jTbIo/2yUOYoDTSLGQ=`,
		},
		{
			ct: pkcs1PrivateKey,
			val: `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAvkgrPq510pZVOcG63aYS9txkHegzMTQTjWfw22/5RZ40mtxr
M/fxWCxPUOq7uPy4HeiGE8Pj038jwEbeQ5qW+WmecsL8KBRYwbU0v21q/sjhybz0
hqhPj1qET6pxst3iBp8BFms90jaexq9B3t6MDMP9oVbPlZA83QLJhM3amG74Do+g
fZMKT4uO8MTuAWGXSNsuhFViZKJE2nJbV2ifib1U5WVrUGL5uRbusiyWOecBkKjj
8jUpXevQ9a3hsA6GUNfZ/L8OIxIQB060vrJHIZGy6d4TT1qxwmxQlbymvoj3Gbnr
yR73vN2MFeURGUqm14XCXlCU6wGlHEhImP9QMQIDAQABAoIBAQCQh5qpUrRVzzBO
3XLFdqaC8WXuPXfc0LRlPOT6mHObSI6mbfPTlmmv0oAwcwtBeFuVBBALJUiAVdre
6jVa3in8qhfbDXWgL8w23h+Bo2eYvRabODX3rhI4TcWgCVOxS82pv86Aq9ZaPHif
a3E6CooQDDIzARBNhzBjowFUKVP3YUsUmaHot+qFK8ZyzVVGhjrCW+wi2MdXIu5T
2qQnm2KCTIQC0mgYE7tklO/9+ZN06uE/pkNm6ypvvwsK0Pqzpt8WgOepYyT2cxnw
3q05ih1PPkxHtJnfd5VIPrtN61yGv6o/Xo2DkFCs3fA9IoJpTcxs9zfnlVZsWM3M
Oy+O/OjRAoGBAOGyl17nyXpjyH8vSz2tx3nQir++ibgJD445axgHsuMI2KpNmRLT
8jbNO5Oby5gKuxJlj7arb2U+ZMtd2vC4kGj5nCtPqWNgnuqSDt/ioIhGohDBQr+x
oTSRsnUaIJOTyMLO9NtOxwQPr/iI1xrNhmtO3L7L2NeLBDVcFR4IEk2dAoGBANfU
T6YuJZRw4pqXUk4rAgHR0zcG0MgRVMw0MWMSVCHBVN5vYm55Uyrvw32XJ7Ne09vm
pV4Qih0eWjguzd4tWChO6tkY1LutFNeQ/G4scs4Yr9XNSbK++aFlYzMTXPNhE98S
WuSuiwx+jldpWZAErsC0wCweUrOm1Haxe0DFg1KlAoGAfaJNUp4R8FgVJn8sEfRn
Qq7MXXnx7YjVqOTbcW/vqyOkgABcAjgK72iFDmC+Dy+B/Pad7iA2DRSTRQVEt5T1
hgnUXeOlNdV2ALs3HndnxxQSaOM7hbuaMcocncTid2PfcFmFwYJzlBYrbVy26IuZ
lKg8htSwKyOOPym385Soo50CgYAGwfsMdP1wPGib9oj5MZeKfwth+bCn0wMYsbmq
JHTF6cvCezJVyy6zdXZlhEoV764qgHpFC7eHWd/xSmXfwwOzn2TzDzf5+F1isoN0
36dolJVM2HSqSBiA2S/V9ZE/fZalsWlvJ5fq+Dt0uTO4sqzWE9LAjuKABYU5gi0d
xhFqkQKBgEoQs7kkywxFi8wSLzC58aCpJCd3KvhC6WWuxfNuaXvW+NonQhGjGCk+
QqAsSLzL4iWSrZjzc7ypYbdaL/sr0chAHixxoJwmGgont1SSf9qhHzvoAQLU3xGH
x4jdnZi5canjhM6v9ZOEsyy8OEGtZ2EXoS/lQUKn/70AdnQdbWCw
-----END RSA PRIVATE KEY-----`,
		},
		{
			ct:  pkcs1PublicKey,
			val: `MIIBCgKCAQEAvkgrPq510pZVOcG63aYS9txkHegzMTQTjWfw22/5RZ40mtxrM/fxWCxPUOq7uPy4HeiGE8Pj038jwEbeQ5qW+WmecsL8KBRYwbU0v21q/sjhybz0hqhPj1qET6pxst3iBp8BFms90jaexq9B3t6MDMP9oVbPlZA83QLJhM3amG74Do+gfZMKT4uO8MTuAWGXSNsuhFViZKJE2nJbV2ifib1U5WVrUGL5uRbusiyWOecBkKjj8jUpXevQ9a3hsA6GUNfZ/L8OIxIQB060vrJHIZGy6d4TT1qxwmxQlbymvoj3GbnryR73vN2MFeURGUqm14XCXlCU6wGlHEhImP9QMQIDAQAB`,
		},
		{
			ct: pkcs1PublicKey,
			val: `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAvkgrPq510pZVOcG63aYS9txkHegzMTQTjWfw22/5RZ40mtxrM/fx
WCxPUOq7uPy4HeiGE8Pj038jwEbeQ5qW+WmecsL8KBRYwbU0v21q/sjhybz0hqhP
j1qET6pxst3iBp8BFms90jaexq9B3t6MDMP9oVbPlZA83QLJhM3amG74Do+gfZMK
T4uO8MTuAWGXSNsuhFViZKJE2nJbV2ifib1U5WVrUGL5uRbusiyWOecBkKjj8jUp
XevQ9a3hsA6GUNfZ/L8OIxIQB060vrJHIZGy6d4TT1qxwmxQlbymvoj3GbnryR73
vN2MFeURGUqm14XCXlCU6wGlHEhImP9QMQIDAQAB
-----END RSA PUBLIC KEY-----`,
		},
		{
			ct:  pkcs8PrivateKey,
			val: `MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAN/eySAfwBAPSocWY3f6ugcZ0zCTKO825kLQ5AUNan976RZOg+lI1ZBAbcOnp25+r+tEmdcWdjkl76pY8miyfqCyByOwzAS3fZm0H2ZxaDTpSraupG449HmP9hOjT+MEqd94WpGxVDUVv2nuN9so2Giq+Ua9JKOTVZeDFFcEyM/pAgMBAAECgYBR0mpuLEyVUhFqODuYsXHmgIDYwyB19fnIt/JvCk0/VPmHJqE91VoBBBtaJF+kmfcQJR2ZKGODVaM3/pRywbJXuAQBXn5J+Dyo352C6oe51GwZr/gYYDLl9JLE8tEhZRS4gUHYLS6WqMZQiQDkj5PljJ2l6F541MA9QolmBOEjgQJBAP7Dcp4dXOX+2ajgpTgJeExpPtrf+wsEe8sCJilaEbytDla6VH+rJXrx8AqNGrwz85cv/C3gOwHi6gw4+XLRKNECQQDg9PO5528NB2wGmkRZiyaL88vPQnE7LS0Wx6+kdAknWK5lmyLscHjJw6PmHCFd1c+xIz4/UeeesyffKSpPO3uZAkANVmAexPzDEbeFbtVXfK9umKfwv38eEYsFksQ6r/tzvD+e7LxVvrkEmbvuYXq/1ZxwEtQJ9s5ACAgmwajViPuxAkAM5H0DbowixwYd6sF4EE2JP9OycTBkH9axs4gReMT9iEuWbym1O0zw41GFYp9W7WYV1NtAbcVEhygF1ioimFohAkEA2KbvdSARygUfZqRrDiC6FipRaVrsO18ZOMqteXfxzthEYVtLJFk/1gO77gVsKfoeIQf8GrznROIh9bPpRUlhEA==`,
		},
		{
			ct: pkcs8PrivateKey,
			val: `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDAd90/zUaEouif
L6ShL0p1DMnwL6tUj57gMG9mQkDLPzDOnuPxKcyeBSFnApADPF5ZVKcsKCbk0vSG
ibWMNlLdr2ECOIVvINx8iGVFQd4jx53bk4tEUM5VVlCSCXENWUGgNyKKdwtYKDG0
tEIDEX4VVryRpf5Zcx3zzGa2PFbU4DGKscdQnDN5TXPI6Y8iqg17OtaodG3zB+74
UkaSYSDP69PI9ZXw6/bv29Z65kr3cCQMBF1EuHpsKsEDzEsvtJWtj7RUDFF5T+I3
qov3YrUz8aqumD6wyvZP4tSuVmi5aaxCrmkxDwYPe/cZfDV73k0krVBgR6KLVPlL
R4ccf3s9AgMBAAECggEBALwVZpTu2TbMrF7DhDIvsKdp8/1P2DIul6emLXbUz9TO
z7da9o25t0fE47tyeaFv4ROS8jrokfmTXXoMIOoAPEJ/HGG7MkpC5rSb5bskfxyf
+deV+8aq4LfsPZg4uc29p7AVsR927hMcVFauwMOW2Iup02TRvhlTsbWJzeXQWVp0
dtxLtvv5RkClDffYjCgBajogzpyvrsrYo2kxPW0Fs7lMJeTGRpG/URH/ACTmDFeU
6x2nG+7OnVwCWkKgnN9klFrTT3+0gkaC2Fykjkhs1o5GhCHGtmMzoVhp/pez48U9
WF6/0jzcQNqSQH/mHlETzvnB/9YuYdDTWs810xioooECgYEA/S8gINV93zhDMKag
PZNFs4fEpuo3vTLU9Zyg8oM8HrInBtVzAlKwCTQ4hfz+sg57evFakUXM7bYmKOzO
WE8mQERzEjrgnryOxeSc6hCINPslkPZ+zR18YBXcTAB5wbpDwrUNsXCFOY/+cEHB
e9bg6wOlA82Np+up5AAG3YlGSqUCgYEAwpvdzw/xcGctOEH/oBBeeXeZQ2+DlniB
tf06Hq3GZfi9HPTeUv2NksARDe4TJXqUG6/lUU3xj5QbPc/chIl5W2iqXdAuS+Zi
nr+x5EzZ1VnThpZNoZEXqZK8NHlodlqFVAomQQfvawihspR8MwUv6IrYAYUggCKU
5FNK4tqRQrkCgYEA5ZZnVv7h7ppRa3ud0ViMGznhnM7Fjr7amILY/DD/QoKgmzTR
3uhmk2IUY4RA3ev+E0VrsFKQe3rybagXkcLsV9j6VCyp5afs/AzMMgCd0xVvQl4U
LCIx19va8dx4jLFAov6VlTMIzGMEAn3OW2NGgDbE24b5jq1IWOWhVCEYabUCgYEA
pZUdS1sTYUJItX9iUvzahZt/amNtoQ/zvbcyRnwxPP5BWmv25sIaPWzyldmlrNP5
RP2Krn0VNccczqEXziVyfpY1rxC885OZAd21LL0+80s0sWUdtITRj1TskfFjMqCe
pPzlw5tO2NFFU78HVhnpw4CvfcuZ9ax25zb/lrnFsvECgYEA0Wb2Jkhe1o/rTEXg
yTLbfcBAbT+1ftxxZIZTbrbpdNpi0WoJuTxfLfTq95HYOpP3xyVID02U6xCOFCIg
IzKaKH2y0r2nzJAX9LNWPsBjnPTLcMCPOnl1ZByQRByUsjnT74M1/SNDe+Dg5Gv3
8Q5FzNEiFnzu/0ivRYPiDEhv75s=
-----END PRIVATE KEY-----`,
		},
		{
			ct:  pkixPublicKey,
			val: `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1eSZbr9I6ZRts+ByRmKnYYcSm/S/IvvxWZKuRL1t4kQns/nurhAM6rAu1E4OH75yAj/SF7FnpEzZ617oak6zSNz7Ow/QZ9APCtA8GpDEkg4p1O2hldIbm3ww8T1cQ6wyy77pQbSWMB70FSNV0CQX57v9e6TxHalMJ5PYvSucTwSiLkBye5uk2ktwnoqP/NNGC3/nr0NXGjeYHWqTIa/gf4GX2jngrm7NsZK4l30IL0AN7wQ07+FFwgJYuQCYCgmMb1Wa1o7hI6zIpCSf9bIlVgPNck/QV4lP/H40zj5S8VeSPysJtV2LhscH7erY9RWEb01/kYVbJpHWhAbTPxolJQIDAQAB`,
		},
		{
			ct: pkixPublicKey,
			val: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvLdok1vcwFjg4NbLWxpX
409ijMEHquPGaokuy+nekajRePjnzBf+yV3fCKey4tsS4vpk2q/PzHi5aiRBfvq5
dVAS8oegZ0qLE44g037LMLCMxwNak9LktJLk/B+Cc5gXLlg7hRCDB1XeiMQYiMOf
OPhwMLWtuhvWpHyH7dyivPm3XY+Bp1z7fhwbQ2EcLAhDVdw95yQ9vdINzsGJjr5Y
UGZZy53IVsyYKTZ59/NNtf5qSOc2uxzRFmPO43SVlGeRvtq6MjHqdcgsnRLSnvI7
qdRoDyQyEz2EPWx7N+yksbWJLmljHo+g8uM9EszvSnKMyJuSp997Zo63xL67MutD
6wIDAQAB
-----END PUBLIC KEY-----`,
		},
		{
			ct: pkcs8PrivateKey,
			val: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQ+oPbjugJroh7FavDAR2McYtj
d+vCuocfTEQ1TrHKF9U7nVdl5hW3TmWzGrIiM9zinovu/cXTq60xVSo6ERF7Y4KnTWzOssB8hbge
yQLr3Sx1+yb9mljmHV3wt6rhyHTYc6KnhUSbHIY60yvLVxqa+RE65tupvNCtZ/0IRo5/kLRIC0tu
kNPBQ2FHN4jDBp4NleTL2kB2fPZwMG7CjWBbBYM5jPCSWAPZi1/IpxnaAd1askYM8HGNYrT3Vz8B
uk8Fw2eGs+ftDnO7oiGo0mJh7cCsQxCd67ItMkCmL71M4UchniykQ0rRdiQeSmITtL6G6vqN+HH+
q3xJK5a4UV53AgMBAAECggEAMxqYBSiYyMlUGdluU/UhXzdjEVOBpXCM25QAeBLC+ntCi14KQeJ2
vAlhaF+mYSPSp96CtNJ9PqqUY2SCu4lx/30Rtc0CxzdCSBFBOdIJM1m2ZjGhUtIcXEWhM90OXkFx
hX8drx6xbMcYciky4/psiKpQ0tkirYr9cvJjl9L9ROdSV2CrYvwr+CriA+JDv/4pY8x4QP7UNnrF
PWENQUdu0WT/D0IWy2HT3POtNO9EIpK2YMO7wF7WwcI7nvNGW8yR9jJ/738xx0rYbf/Hq0fNME7/
tAWJK02P2yxAZlNyoToAQ0O92HTbldYQOQgM5BHhe0x+NX8lrXwiDr1NQhmGjQKBgQDyVkSQm5Md
MuO8ta/VMZJAl0i0CEDfq5LepN2m5XSJGUUFdOoEdQI777R4U4PHWuooJ5/LQ8bs1lHEbcVFQuHh
MUI7RPtMh3RRfz6bNV6TKMdg3bUp0+fjEB2Tm0w5K1PLvkCkmGdZyYmRvMH0F3zR347OyGuqyUWX
dBLwH0OTowKBgQDcwsYicS0XY+0gJspXRuQ2tf+O55ooUDCg7Lo1D78DF0wVtl6ic/tLe6lMEK5I
cZTYZ9Kyb5rlh3khA6Av5QX5MeshmjhwWzu38050EoyiAiRaeKsxptZivs2d0YnE/SWZJIJnk7Hm
A5vfXSz/kPGbqj5hVMVgrBvCjiLVfWIXHQKBgQDgWuHLh1zh0XVqBkMte2FNj0Ht+x4kdXHZ0oSq
uQ/0xYJTFPR6/+D7oGZSZ+8+p3rVhim4Q51tMtYspvvVrZ/1nmcU/D4zkcwsj0Nk6joOv9gmY9wP
R3INk6PuPf6JhwVjQVYTjE1SoPVOCZT+6KfUncZWxtJ3ITPejcXirO8hRQKBgQCwH/zvcZfl92Ux
p4D7DKX6OE2Bd6l3zDJf0T3mI3/jOW0MTYlG1n2AhVJWS2Cgj22PEZX5oizUPrcW7cuZKoEPhRHw
pxesHD2SztyiokHs5wSV4XvDizWzZkKpTIk7zjN28LfRZvYhanOrSq0h4EPCS5qlEHrAW89x8vA7
n2LoyQKBgH43xkaO9dxQjqvHMh5iTzn7pcXPI9RABVG8Un+WkG90l04xRilYg/1RFznQfQXtxyDv
3XAzfTAJHPR3GXCtcZIfsTx12OOweCMbphojYy1yXuo1B48jfXNJG0G+aEQ3s6Wh1eJ5kXi34JMp
HbpMNg17/p8b9Za7IROLV4Ypd7sb
-----END PRIVATE KEY-----`,
		},
		{
			ct:  x509Certificate,
			val: `MIIDXTCCAkWgAwIBAgIJAISKTtZx8oAqMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMTcxMjEzMTAyODE5WhcNMTgwMTEyMTAyODE5WjBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+eKfL1lMFOVP4L3LfNcKueYX0Spl2x1+B9UFufFyF/kC7hI6KjevMA5Y7HgSImSbajpa6B8c4iAdtGvG2p/DBddClv5Xq9wCTZffDxQbdb6syUfUo6/22hvCFDx/kpYmyLpq8r2AMAAeDsOmiO9pybhkMi4VHDs5DC/ptKR4h3w6T2c7RGb65kOO94rFU9lzITlIQOrgRPVUt38n5V2zVqWuUjU214CdSNHLiFhMaTH7WI/F/7JDQ2KI44Bd5htoZfgqNRC453KdDuUJO3kK/tSW2BXECTfR5ZstSJF34sMe4PBKOvg9xS5sr6CFe4VdTjHRzt/ELe05JYGoTjU3fwIDAQABo1AwTjAdBgNVHQ4EFgQUT4MTeTRAnnhjjQkeFAh1on6ZqCowHwYDVR0jBBgwFoAUT4MTeTRAnnhjjQkeFAh1on6ZqCowDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAadVMg8QHAPTpx+MXbXOTXWQUYIWZ4CJz0+rXQjUwgRygYd65CNO4cI9yk9w0Fx7GYvEbY6pL0iG3PzZs9nvG5SgxRrBxk7jbBbYaBLL0tCL+c0DryR2DOfbiPOcKBQGX6SULNcnAJGxkTwbhk7tiisbOxnW0XtAezcyQvCDWGvo55rHNTfZ8iBvIjh5aUofwYplG1xrJtPOY5OhQDzdsV4d/A9rdUmaamES0eA1EASRPs4hs0HWYq/dNu22VREyRQLuepkoYvrFYQFgwM8qU3kE3/4R4ZGTK2JuHfnM75hr6y6X9RHwJVyUMw45KgfwPCITUYOva6kdfudrlHYRR6A==`,
		},
		{
			ct: x509Certificate,
			val: `-----BEGIN CERTIFICATE-----
MIIDbDCCAlSgAwIBAgIGAWkz3Fw1MA0GCSqGSIb3DQEBCwUAMHcxFDASBgNVBAoTC0dvb2dsZSBJ
bmMuMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MRQwEgYDVQQDEwtMREFQIENsaWVudDEPMA0GA1UE
CxMGR1N1aXRlMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTAeFw0xOTAyMjgxMTI3
NTJaFw0yMjAyMjcxMTI3NTJaMHcxFDASBgNVBAoTC0dvb2dsZSBJbmMuMRYwFAYDVQQHEw1Nb3Vu
dGFpbiBWaWV3MRQwEgYDVQQDEwtMREFQIENsaWVudDEPMA0GA1UECxMGR1N1aXRlMQswCQYDVQQG
EwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AND6g9uO6AmuiHsVq8MBHYxxi2N368K6hx9MRDVOscoX1TudV2XmFbdOZbMasiIz3OKei+79xdOr
rTFVKjoREXtjgqdNbM6ywHyFuB7JAuvdLHX7Jv2aWOYdXfC3quHIdNhzoqeFRJschjrTK8tXGpr5
ETrm26m80K1n/QhGjn+QtEgLS26Q08FDYUc3iMMGng2V5MvaQHZ89nAwbsKNYFsFgzmM8JJYA9mL
X8inGdoB3VqyRgzwcY1itPdXPwG6TwXDZ4az5+0Oc7uiIajSYmHtwKxDEJ3rsi0yQKYvvUzhRyGe
LKRDStF2JB5KYhO0vobq+o34cf6rfEkrlrhRXncCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAwhFn
hnqpmzGwHSvptRsCHBn1uUrxMaQDitIc0mMJbIANjBduF9prU0PYqZkXK4PrJZ9CJoMQCcvDf0dy
ABbXee2b/9+HlBE78Fye+BGTi2cK//Vi4i+j6UFGtF8Ougsde1J8q3Dir1m/u6TXI/N1etMOifNi
aw3GIiLIpfoGAAaUVS8qcB9G94PUWIKSze9cHb0/FF1J3uTdu5vqz5CWaSuHtmj7nT0tWxqT0AiW
boU2Ixx0QOjn5XoHskWdY065eOc14bSGCHutiiAYBwJMnGol18X7D/hUKEiAfggG8927rnRGzGpN
otero2Uksx3F/9g6PtyQx2WRah4A8eA9kQ==
-----END CERTIFICATE-----`,
		},
	}

	for _, c := range testCryptoFormatCheckerKeys {
		testCryptoFormatCheckerIsFormatValid(t, c.ct, c.val)
	}
}

func TestCryptoFormatCheckerIsFormat_NotBase64(t *testing.T) {
	givenName := cryptoName
	givenValue := "not base-64"
	cryptoType := cryptoType(9999)

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newCryptoFormatChecker(givenName, cryptoType)
	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to decode the data")
}

func TestCryptoFormatCheckerIsFormat_UnsupportedType(t *testing.T) {
	givenName := cryptoName
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
	givenName := regexName
	givenValue := 1

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newRegexFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the value is not a string")
}

func TestRegexFormatCheckerIsFormat_Valid(t *testing.T) {
	givenName := regexName
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
	givenName := regexName
	givenValue := "^(.*$"

	formatErrs.clear()

	defer func() { formatErrs.clear() }()

	name, checker := newRegexFormatChecker(givenName)
	assert.Equal(t, givenName, name)

	result := checker.IsFormat(givenValue)
	assert.False(t, result)

	err := formatErrs.get(name, givenValue)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to parse regular expression")
}

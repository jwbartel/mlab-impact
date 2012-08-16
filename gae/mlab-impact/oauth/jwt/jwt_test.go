// Copyright 2012 The goauth2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// For package documentation please see jwt.go.
//
package jwt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	StdHeaderStr = "{\"alg\":\"RS256\",\"typ\":\"JWT\"}"
	Iss          = "761326798069-r5mljlln1rd4lrbhg75efgigp36m78j5@developer.gserviceaccount.com"
	Scope        = "https://www.googleapis.com/auth/prediction"
	Aud          = "https://accounts.google.com/o/oauth2/token"
	Exp          = 1328554385
	Iat          = 1328550785 // Exp + 1 hour
)

// Base64url encoded Header
const HeaderEnc = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9"

// Base64url encoded ClaimSet
const ClaimSetEnc = "eyJpc3MiOiI3NjEzMjY3OTgwNjktcjVtbGpsbG4xcmQ0bHJiaGc3NWVmZ2lncDM2bTc4ajVAZGV2ZWxvcGVyLmdzZXJ2aWNlYWNjb3VudC5jb20iLCJzY29wZSI6Imh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL2F1dGgvcHJlZGljdGlvbiIsImF1ZCI6Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi90b2tlbiIsImV4cCI6MTMyODU1NDM4NSwiaWF0IjoxMzI4NTUwNzg1fQ"

// Base64url encoded Signature
const SigEnc = "olukbHreNiYrgiGCTEmY3eWGeTvYDSUHYoE84Jz3BRPBSaMdZMNOn_0CYK7UHPO7OdvUofjwft1dH59UxE9GWS02pjFti1uAQoImaqjLZoTXr8qiF6O_kDa9JNoykklWlRAIwGIZkDupCS-8cTAnM_ksSymiH1coKJrLDUX_BM0x2f4iMFQzhL5vT1ll-ZipJ0lNlxb5QsyXxDYcxtHYguF12-vpv3ItgT0STfcXoWzIGQoEbhwB9SBp9JYcQ8Ygz6pYDjm0rWX9LrchmTyDArCodpKLFtutNgcIFUP9fWxvwd1C2dNw5GjLcKr9a_SAERyoJ2WnCR1_j9N0wD2o0g"

// Base64url encoded Token
const TokEnc = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI3NjEzMjY3OTgwNjktcjVtbGpsbG4xcmQ0bHJiaGc3NWVmZ2lncDM2bTc4ajVAZGV2ZWxvcGVyLmdzZXJ2aWNlYWNjb3VudC5jb20iLCJzY29wZSI6Imh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL2F1dGgvcHJlZGljdGlvbiIsImF1ZCI6Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi90b2tlbiIsImV4cCI6MTMyODU1NDM4NSwiaWF0IjoxMzI4NTUwNzg1fQ.olukbHreNiYrgiGCTEmY3eWGeTvYDSUHYoE84Jz3BRPBSaMdZMNOn_0CYK7UHPO7OdvUofjwft1dH59UxE9GWS02pjFti1uAQoImaqjLZoTXr8qiF6O_kDa9JNoykklWlRAIwGIZkDupCS-8cTAnM_ksSymiH1coKJrLDUX_BM0x2f4iMFQzhL5vT1ll-ZipJ0lNlxb5QsyXxDYcxtHYguF12-vpv3ItgT0STfcXoWzIGQoEbhwB9SBp9JYcQ8Ygz6pYDjm0rWX9LrchmTyDArCodpKLFtutNgcIFUP9fWxvwd1C2dNw5GjLcKr9a_SAERyoJ2WnCR1_j9N0wD2o0g"

// Private key for testing
const PrivateKeyPem = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA4ej0p7bQ7L/r4rVGUz9RN4VQWoej1Bg1mYWIDYslvKrk1gpj
7wZgkdmM7oVK2OfgrSj/FCTkInKPqaCR0gD7K80q+mLBrN3PUkDrJQZpvRZIff3/
xmVU1WeruQLFJjnFb2dqu0s/FY/2kWiJtBCakXvXEOb7zfbINuayL+MSsCGSdVYs
SliS5qQpgyDap+8b5fpXZVJkq92hrcNtbkg7hCYUJczt8n9hcCTJCfUpApvaFQ18
pe+zpyl4+WzkP66I28hniMQyUlA1hBiskT7qiouq0m8IOodhv2fagSZKjOTTU2xk
SBc//fy3ZpsL7WqgsZS7Q+0VRK8gKfqkxg5OYQIDAQABAoIBAQDGGHzQxGKX+ANk
nQi53v/c6632dJKYXVJC+PDAz4+bzU800Y+n/bOYsWf/kCp94XcG4Lgsdd0Gx+Zq
HD9CI1IcqqBRR2AFscsmmX6YzPLTuEKBGMW8twaYy3utlFxElMwoUEsrSWRcCA1y
nHSDzTt871c7nxCXHxuZ6Nm/XCL7Bg8uidRTSC1sQrQyKgTPhtQdYrPQ4WZ1A4J9
IisyDYmZodSNZe5P+LTJ6M1SCgH8KH9ZGIxv3diMwzNNpk3kxJc9yCnja4mjiGE2
YCNusSycU5IhZwVeCTlhQGcNeV/skfg64xkiJE34c2y2ttFbdwBTPixStGaF09nU
Z422D40BAoGBAPvVyRRsC3BF+qZdaSMFwI1yiXY7vQw5+JZh01tD28NuYdRFzjcJ
vzT2n8LFpj5ZfZFvSMLMVEFVMgQvWnN0O6xdXvGov6qlRUSGaH9u+TCPNnIldjMP
B8+xTwFMqI7uQr54wBB+Poq7dVRP+0oHb0NYAwUBXoEuvYo3c/nDoRcZAoGBAOWl
aLHjMv4CJbArzT8sPfic/8waSiLV9Ixs3Re5YREUTtnLq7LoymqB57UXJB3BNz/2
eCueuW71avlWlRtE/wXASj5jx6y5mIrlV4nZbVuyYff0QlcG+fgb6pcJQuO9DxMI
aqFGrWP3zye+LK87a6iR76dS9vRU+bHZpSVvGMKJAoGAFGt3TIKeQtJJyqeUWNSk
klORNdcOMymYMIlqG+JatXQD1rR6ThgqOt8sgRyJqFCVT++YFMOAqXOBBLnaObZZ
CFbh1fJ66BlSjoXff0W+SuOx5HuJJAa5+WtFHrPajwxeuRcNa8jwxUsB7n41wADu
UqWWSRedVBg4Ijbw3nWwYDECgYB0pLew4z4bVuvdt+HgnJA9n0EuYowVdadpTEJg
soBjNHV4msLzdNqbjrAqgz6M/n8Ztg8D2PNHMNDNJPVHjJwcR7duSTA6w2p/4k28
bvvk/45Ta3XmzlxZcZSOct3O31Cw0i2XDVc018IY5be8qendDYM08icNo7vQYkRH
504kQQKBgQDjx60zpz8ozvm1XAj0wVhi7GwXe+5lTxiLi9Fxq721WDxPMiHDW2XL
YXfFVy/9/GIMvEiGYdmarK1NW+VhWl1DC5xhDg0kvMfxplt4tynoq1uTsQTY31Mx
BeF5CT/JuNYk3bEBF0H/Q3VGO1/ggVS+YezdFbLWIRoMnLj6XCFEGg==
-----END RSA PRIVATE KEY-----`

// Public key to go with the private key for testing
const PublicKeyPem = `-----BEGIN CERTIFICATE-----
MIIDIzCCAgugAwIBAgIJAMfISuBQ5m+5MA0GCSqGSIb3DQEBBQUAMBUxEzARBgNV
BAMTCnVuaXQtdGVzdHMwHhcNMTExMjA2MTYyNjAyWhcNMjExMjAzMTYyNjAyWjAV
MRMwEQYDVQQDEwp1bml0LXRlc3RzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEA4ej0p7bQ7L/r4rVGUz9RN4VQWoej1Bg1mYWIDYslvKrk1gpj7wZgkdmM
7oVK2OfgrSj/FCTkInKPqaCR0gD7K80q+mLBrN3PUkDrJQZpvRZIff3/xmVU1Wer
uQLFJjnFb2dqu0s/FY/2kWiJtBCakXvXEOb7zfbINuayL+MSsCGSdVYsSliS5qQp
gyDap+8b5fpXZVJkq92hrcNtbkg7hCYUJczt8n9hcCTJCfUpApvaFQ18pe+zpyl4
+WzkP66I28hniMQyUlA1hBiskT7qiouq0m8IOodhv2fagSZKjOTTU2xkSBc//fy3
ZpsL7WqgsZS7Q+0VRK8gKfqkxg5OYQIDAQABo3YwdDAdBgNVHQ4EFgQU2RQ8yO+O
gN8oVW2SW7RLrfYd9jEwRQYDVR0jBD4wPIAU2RQ8yO+OgN8oVW2SW7RLrfYd9jGh
GaQXMBUxEzARBgNVBAMTCnVuaXQtdGVzdHOCCQDHyErgUOZvuTAMBgNVHRMEBTAD
AQH/MA0GCSqGSIb3DQEBBQUAA4IBAQBRv+M/6+FiVu7KXNjFI5pSN17OcW5QUtPr
odJMlWrJBtynn/TA1oJlYu3yV5clc/71Vr/AxuX5xGP+IXL32YDF9lTUJXG/uUGk
+JETpKmQviPbRsvzYhz4pf6ZIOZMc3/GIcNq92ECbseGO+yAgyWUVKMmZM0HqXC9
ovNslqe0M8C1sLm1zAR5z/h/litE7/8O2ietija3Q/qtl2TOXJdCA6sgjJX2WUql
ybrC55ct18NKf3qhpcEkGQvFU40rVYApJpi98DiZPYFdx1oBDp/f4uZ3ojpxRVFT
cDwcJLfNRCPUhormsY7fDS9xSyThiHsW9mjJYdcaKQkwYZ0F11yB
-----END CERTIFICATE-----`

var (
	PrivateKeyPemBytes = []byte(PrivateKeyPem)
	PublicKeyPemBytes  = []byte(PublicKeyPem)
)

// The signature in bytes
var SigBytes = Signature{162, 91, 164, 108, 122, 222, 54, 38, 43, 130, 33, 130, 76, 73, 152, 221, 229, 134, 121, 59, 216, 13, 37, 7, 98, 129, 60, 224, 156, 247, 5, 19, 193, 73, 163, 29, 100, 195, 78, 159, 253, 2, 96, 174, 212, 28, 243, 187, 57, 219, 212, 161, 248, 240, 126, 221, 93, 31, 159, 84, 196, 79, 70, 89, 45, 54, 166, 49, 109, 139, 91, 128, 66, 130, 38, 106, 168, 203, 102, 132, 215, 175, 202, 162, 23, 163, 191, 144, 54, 189, 36, 218, 50, 146, 73, 86, 149, 16, 8, 192, 98, 25, 144, 59, 169, 9, 47, 188, 113, 48, 39, 51, 249, 44, 75, 41, 162, 31, 87, 40, 40, 154, 203, 13, 69, 255, 4, 205, 49, 217, 254, 34, 48, 84, 51, 132, 190, 111, 79, 89, 101, 249, 152, 169, 39, 73, 77, 151, 22, 249, 66, 204, 151, 196, 54, 28, 198, 209, 216, 130, 225, 117, 219, 235, 233, 191, 114, 45, 129, 61, 18, 77, 247, 23, 161, 108, 200, 25, 10, 4, 110, 28, 1, 245, 32, 105, 244, 150, 28, 67, 198, 32, 207, 170, 88, 14, 57, 180, 173, 101, 253, 46, 183, 33, 153, 60, 131, 2, 176, 168, 118, 146, 139, 22, 219, 173, 54, 7, 8, 21, 67, 253, 125, 108, 111, 193, 221, 66, 217, 211, 112, 228, 104, 203, 112, 170, 253, 107, 244, 128, 17, 28, 168, 39, 101, 167, 9, 29, 127, 143, 211, 116, 192, 61, 168, 210}

// Testing the urlEncode function.
func TestUrlEncode(t *testing.T) {
	enc := urlEncode([]byte(StdHeaderStr))
	b := []byte(enc)
	if b[len(b)-1] == 61 {
		t.Error("TestUrlEncode: last chat == \"=\"")
	}
	if enc != HeaderEnc {
		t.Error("TestUrlEncode: enc != HeaderEnc")
		t.Errorf("        enc = %s", enc)
		t.Errorf("  HeaderEnc = %s", HeaderEnc)
	}
}

// Given a well formed Header, test for proper encoding.
func TestHeaderEncode(t *testing.T) {
	h := &Header{
		Algorithm: "RS256",
		Type:      "JWT",
	}
	enc, err := h.Encode()
	if err != nil {
		t.Errorf("TestHeaderEncode:h.Encode: %v", err)
	}
	if enc != HeaderEnc {
		t.Error("TestHeaderEncode: enc != HeaderEnc")
		t.Errorf("        enc = %s", enc)
		t.Errorf("  HeaderEnc = %s", HeaderEnc)
	}
}

// Test that the times are set properly.
func TestClaimSetSetTimes(t *testing.T) {
	c := &ClaimSet{
		Iss:   Iss,
		Scope: Scope,
		Aud:   Aud,
	}
	iatTime := time.Unix(Iat, 0)
	c.setTimes(iatTime)
	if c.Exp != Exp {
		t.Error("TestClaimSetSetTimes: c.Exp != Exp")
		t.Errorf("  c.Exp = %d", c.Exp)
		t.Errorf("    Exp = %d", Exp)
	}
}

// Given a well formed ClaimSet, test for proper encoding.
func TestClaimSetEncode(t *testing.T) {
	c := &ClaimSet{
		Iss:   Iss,
		Scope: Scope,
		Aud:   Aud,
		Exp:   Exp,
		Iat:   Iat,
	}
	enc, err := c.Encode()
	if err != nil {
		t.Errorf("TestClaimSetEncode:c.Encode: %v", err)
	}
	if enc != ClaimSetEnc {
		t.Error("TestClaimSetEncode: enc != ClaimSetEnc")
		t.Errorf("        enc = %s", enc)
		t.Errorf("  HeaderEnc = %s", ClaimSetEnc)
	}
}

// Given a well formed Signature, test for proper encoding.
func TestSignatureEncode(t *testing.T) {
	enc, err := SigBytes.Encode()
	if err != nil {
		t.Errorf("TestSignatureEncode:SigBytes.Encode: %v", err)
	}
	if enc != SigEnc {
		t.Error("TestSignatureEncode: enc != SigEnc")
		t.Errorf("     enc = %s", enc)
		t.Errorf("  SigEnc = %s", SigEnc)
	}
}

// Make sure the private key parsing functions work.
func TestParsePrivateKey(t *testing.T) {
	tok := &Token{
		Key: PrivateKeyPemBytes,
	}
	err := tok.parsePrivateKey()
	if err != nil {
		t.Errorf("TestParsePrivateKey:tok.parsePrivateKey: %v", err)
	}
}

// Test that the token signature generated matches the golden standard.
func TestTokenSign(t *testing.T) {
	tok := &Token{
		Key:   PrivateKeyPemBytes,
		head:  HeaderEnc,
		claim: ClaimSetEnc,
	}
	err := tok.parsePrivateKey()
	if err != nil {
		t.Errorf("TestTokenSign:tok.parsePrivateKey: %v", err)
	}
	err = tok.sign()
	if err != nil {
		t.Errorf("TestTokenSign:tok.sign: %v", err)
	}
	if len(tok.Signature) != len(SigBytes) {
		t.Errorf("TestTokenSign: len(tok.Signature) != len(SigBytes)")
	}
	for i := range tok.Signature {
		if tok.Signature[i] != SigBytes[i] {
			t.Errorf("TestTokenSign: tok.Signature[%d] != SigBytes[%d]", i, i)
		}
	}
}

// Test that the token expiration function is working.
func TestTokenExpired(t *testing.T) {
	c := &ClaimSet{}
	tok := &Token{
		ClaimSet: c,
	}
	now := time.Now()
	c.setTimes(now)
	if tok.Expired() != false {
		t.Error("TestTokenExpired: tok.Expired != false")
	}
	// Set the times as if they were set 2 hours ago.
	c.setTimes(now.Add(-2 * time.Hour))
	if tok.Expired() != true {
		t.Error("TestTokenExpired: tok.Expired != true")
	}
}

// Given a well formed Token, test for proper encoding.
func TestTokenEncode(t *testing.T) {
	c := &ClaimSet{
		Iss:   Iss,
		Scope: Scope,
		Aud:   Aud,
		Exp:   Exp,
		Iat:   Iat,
	}
	tok := &Token{
		Header:   StdHeader,
		ClaimSet: c,
		Key:      PrivateKeyPemBytes,
	}
	enc, err := tok.Encode()
	if err != nil {
		t.Errorf("TestTokenEncode:tok.Assertion: %v", err)
	}
	if enc != TokEnc {
		t.Error("TestTokenEncode: enc != TokEnc")
		t.Errorf("     enc = %s", enc)
		t.Errorf("  TokEnc = %s", TokEnc)
	}
}

// Given a well formed Token we should get back a well formed request.
func TestBuildRequest(t *testing.T) {
	c := &ClaimSet{
		Iss:   Iss,
		Scope: Scope,
		Aud:   Aud,
		Exp:   Exp,
		Iat:   Iat,
	}
	tok := &Token{
		Header:   StdHeader,
		ClaimSet: c,
		Key:      PrivateKeyPemBytes,
	}
	r, err := tok.BuildRequest()
	if err != nil {
		t.Errorf("TestBuildRequest:BuildRequest: %v", err)
	}
	if r.Method != "POST" {
		t.Error("TestBuildRequest: r.Method != POST")
		t.Errorf("  r.Method = %s", r.Method)
	}
	v := r.URL.Query()
	if v.Get("grant_type") != StdGrantType {
		t.Error("TestBuildRequest: grant_type != StdGrantType")
		t.Errorf("    grant_type = %s", v.Get("grant_type"))
		t.Errorf("  StdGrantType = %s", StdGrantType)
	}
	if v.Get("assertion_type") != tok.AssertionType() {
		t.Error("TestBuildRequest: assertion_type != t.AssertionType")
		t.Errorf("   assertion_type = %s", v.Get("assertion_type"))
		t.Errorf("  t.AssertionType = %s", tok.AssertionType())
	}
	if v.Get("assertion") != TokEnc {
		t.Error("TestBuildRequest: assertion != TokEnc")
		t.Errorf("  assertion = %s", v.Get("assertion"))
		t.Errorf("     TokEnc = %s", TokEnc)
	}
}

// Given a well formed access request response we should get back a oauth.Token.
func TestHandleResponse(t *testing.T) {
	rb := &respBody{
		Access:    "1/8xbJqaOZXSUZbHLl5EOtu1pxz3fmmetKx9W8CV4t79M",
		Type:      "Bearer",
		ExpiresIn: 3600,
	}
	b, err := json.Marshal(rb)
	if err != nil {
		t.Errorf("TestHandleResponse:json.Marshal: %v", err)
	}
	r := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}
	tok := &Token{}
	o, err := tok.HandleResponse(r)
	if err != nil {
		t.Errorf("TestHandleResponse:HandleResponse: %v", err)
	}
	if o.AccessToken != rb.Access {
		t.Error("TestHandleResponse: o.AccessToken != rb.Access")
		t.Errorf("  o.AccessToken = %s", o.AccessToken)
		t.Errorf("       rb.Access = %s", rb.Access)
	}
	if o.Expired() {
		t.Error("TestHandleResponse: o.Expired == true")
	}
}

// Placeholder for future Assert tests.
func TestAssert(t *testing.T) {
	// Since this method makes a call to BuildRequest, an htttp.Client, and
	// finally HandleResponse there is not much more to test.  This is here
	// as a placeholder if that changes.
}

// Benchmark for the end-to-end encoding of a well formed token.
func BenchmarkTokenEncode(b *testing.B) {
	b.StopTimer()
	c := &ClaimSet{
		Iss:   Iss,
		Scope: Scope,
		Aud:   Aud,
		Exp:   Exp,
		Iat:   Iat,
	}
	tok := &Token{
		Header:   StdHeader,
		ClaimSet: c,
		Key:      PrivateKeyPemBytes,
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Exp = 0
		tok.Encode()
	}
}

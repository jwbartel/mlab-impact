// Copyright 2012 The goauth2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The jwt package provides support for creating credentials for OAuth2 service
// account requests.
//
// For examples of the package usage please see jwt_test.go.
//
// For info on OAuth2 service accounts please see the online documentation.
// https://developers.google.com/accounts/docs/OAuth2ServiceAccount
//
package jwt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.google.com/p/goauth2/oauth"
)

const (
	StdAlgorithm     = "RS256"
	StdType          = "JWT"
	StdAssertionType = "http://oauth.net/grant_type/jwt/1.0/bearer"
	StdGrantType     = "assertion"
)

var (
	ErrInvKey = errors.New("Invalid Key")

	StdHeader = &Header{
		Algorithm: StdAlgorithm,
		Type:      StdType,
	}
)

// urlEncode returns and Base64url encoded version of the input string with any
// trailing "=" stripped.
func urlEncode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

// The header consists of two fields that indicate the signing algorithm and the
// format of the assertion. Both fields are mandatory, and each field has only
// one value. As additional algorithms and formats are introduced, this header
// will change accordingly.
//
// Service Accounts rely on the RSA SHA256 algorithm and the JWT token format
// (see StdAlgorithm and StdType).
type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// Encode returns the Base64url encoded form of the Header.
func (h *Header) Encode() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	return urlEncode(b), nil
}

// The JWT claim set contains information about the JWT including the
// permissions being requested (scopes), the target of the token, the issuer,
// the time the token was issued, and the lifetime of the token. Most of the
// fields are mandatory.
//
// Aud is usually https://accounts.google.com/o/oauth2/token
type ClaimSet struct {
	Iss     string `json:"iss"`           // email address of the client_id of the application making the access token request
	Scope   string `json:"scope"`         // space-delimited list of the permissions the application requests
	Aud     string `json:"aud"`           // a descriptor of the intended target of the assertion
	Exp     int64  `json:"exp"`           // the expiration time of the assertion (maximum of 1 hour from the JWT encoding time)
	Iat     int64  `json:"iat"`           // the time the assertion was issued (the time that the JWT was encoded)
	Prn     string `json:"prn,omitempty"` // the email for which the application is requesting delegated access (Optional).
	expTime time.Time
	iatTime time.Time
}

// setTimes sets Iat and Exp to time.Now() and Iat.Add(time.Hour) respectively.
//
// Note that these times have nothing to do with the expiration time for the
// access_token returned by the server.  These have to do with the lifetime of
// the encoded JWT.
//
// A JWT can be re-used for up to one hour after it was encoded.  The access
// token that is granted will also be good for one hour so there is little point
// in trying to use the JWT a second time.
func (c *ClaimSet) setTimes(t time.Time) {
	c.iatTime = t
	c.expTime = c.iatTime.Add(time.Hour)
	c.Iat = c.iatTime.Unix() // The time that the JWT was encoded.
	c.Exp = c.expTime.Unix() // The time that the encoded JWT will expire.
}

// Encode returns the Base64url encoded form of the Signature.  If either of Iat
// or Exp are 0, then they will be set to time.Now() and Iat.Add(time.Hour)
// respectively.
func (c *ClaimSet) Encode() (string, error) {
	if c.Exp == 0 || c.Iat == 0 {
		c.setTimes(time.Now())
	}
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return urlEncode(b), nil
}

// JSON Web Signature (JWS) is the specification that guides the mechanics of
// generating the signature for the JWT.  The contents of a Signature will be
// generated based on a well formed Header and ClaimSet.
type Signature []byte

// Encode returns the Base64url encoded form of the Signature.  The error
// returned from this function is always nil and can be safely ignored.
func (s Signature) Encode() (string, error) {
	return urlEncode([]byte(s)), nil
}

// A JWT is composed of three parts: a header, a claim set, and a signature.
// The well formed and encoded JWT can then be exchanged for an access token.
//
// The Token is not a JWT, but is is encoded to produce a well formed JWT.
type Token struct {
	Header    *Header
	ClaimSet  *ClaimSet
	Signature Signature
	Key       []byte
	pKey      *rsa.PrivateKey
	head      string
	claim     string
	sig       string
}

// Expired returns a boolean value letting us know if the token has expired.
func (t *Token) Expired() bool {
	// The truth is that if the sis is empty, then it was never encoded
	// and was never sent to the server, so you should re-encode and send it
	// anyways.
	if t.ClaimSet.Exp == 0 {
		return true
	}
	return t.ClaimSet.expTime.Before(time.Now())
}

// AssertionType returns the standard assertion type for a JWT.
func (t *Token) AssertionType() string {
	return StdAssertionType
}

// Encode constructs and signs a Token returning a JWT ready to use for
// requesting an access token.
func (t *Token) Encode() (string, error) {
	var tok string
	var err error
	t.head, err = t.Header.Encode()
	if err != nil {
		return tok, err
	}
	t.claim, err = t.ClaimSet.Encode()
	if err != nil {
		return tok, err
	}
	err = t.sign()
	if err != nil {
		return tok, err
	}
	t.sig, _ = t.Signature.Encode()
	tok = fmt.Sprintf("%s.%s.%s", t.head, t.claim, t.sig)
	return tok, nil
}

// sign computes the signature for a Token.  The details for this can be found
// in the OAuth2 Service Account documentation.
// https://developers.google.com/accounts/docs/OAuth2ServiceAccount#computingsignature
func (t *Token) sign() error {
	ss := fmt.Sprintf("%s.%s", t.head, t.claim)
	if t.pKey == nil {
		err := t.parsePrivateKey()
		if err != nil {
			return err
		}
	}
	h := sha256.New()
	h.Write([]byte(ss))
	b, err := rsa.SignPKCS1v15(rand.Reader, t.pKey, crypto.SHA256, h.Sum(nil))
	t.Signature = Signature(b)
	return err
}

// parsePrivateKey converts the Token's Key ([]byte) into a parsed
// rsa.PrivateKey.  If the key is not well formed this method will return an
// ErrInvKey error.
func (t *Token) parsePrivateKey() error {
	block, _ := pem.Decode(t.Key)
	if block == nil {
		return ErrInvKey
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
	}
	var ok bool
	t.pKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return ErrInvKey
	}
	return nil
}

// Assert returns an *oauth.Token that can be used to make future requests.  The
// access_token will expire in one hour (3600 seconds) and cannot be refreshed
// (no refresh_token is returned with the response).  Once this token expires
// call this method again to get a fresh one.
func (t *Token) Assert(client *http.Client) (*oauth.Token, error) {
	
	j, err := t.Encode()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("grant_type", StdGrantType)
	v.Set("assertion_type", t.AssertionType())
	v.Set("assertion", j)
	/*u, err := url.Parse(t.ClaimSet.Aud)
	if err != nil {
		return nil, err
	}*/
	//u.RawQuery = v.Encode()
	//return http.NewRequest("POST", u.String(), bytes.NewBufferString(v.Encode()))

	var o *oauth.Token
	//req, err := t.BuildRequest()
	if err != nil {
		return nil, errors.New("r error")
		return o, err
	}
	resp, err := client.PostForm(t.ClaimSet.Aud, v)
	if err != nil {
		return o, err
	}
	o, err = t.HandleResponse(resp)
	return o, err
}

// BuildRequest returns and *http.Request that can be used to make an access
// token request.  Per the Assertion profile in OAuth 2.0 specification, this
// access token request is an HTTPs POST.
//
// Most users should simply use the Assert method.
func (t *Token) BuildRequest() (*http.Request, error) {
	j, err := t.Encode()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("grant_type", StdGrantType)
	v.Set("assertion_type", t.AssertionType())
	v.Set("assertion", j)
	u, err := url.Parse(t.ClaimSet.Aud)
	if err != nil {
		return nil, err
	}
	//u.RawQuery = v.Encode()
	return http.NewRequest("POST", u.String(), bytes.NewBufferString(v.Encode()))
}

// Used for decoding the response body.
type respBody struct {
	Access    string        `json:"access_token"`
	Type      string        `json:"token_type"`
	ExpiresIn time.Duration `json:"expires_in"`
}

// HandleResponse returns an *oauth.Token given the *http.Response from a
// *http.Request created by BuildRequest.
//
// Most users should simply use the Assert method.
func (t *Token) HandleResponse(r *http.Response) (*oauth.Token, error) {
	o := &oauth.Token{}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return o, errors.New("invalid response: " + r.Status)
	}
	b := &respBody{}
	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		return o, err
	}
	o.AccessToken = b.Access
	o.Expiry = time.Now().Add(b.ExpiresIn * time.Second)
	return o, nil
}

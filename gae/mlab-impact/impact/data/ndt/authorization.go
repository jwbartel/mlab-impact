package ndt

import (
	"appengine"
	"appengine/urlfetch"
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"impact/data/secrets"
	"net/http"
	"oauth/jwt"
	"os"
	"time"
)

type JWTTransport struct {
	Token   *jwt.Token
	Oauth   *oauth.Token
	Request *http.Request
	Client  *http.Client
}

func (t *JWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Token == nil {
		t.CreateToken()
	}
	if t.Oauth == nil || t.Oauth.Expired() {
		err := t.UpdateOauth()
		if err != nil {
			return nil, err
		}
	}
	req.Header.Set("Authorization", "OAuth "+t.Oauth.AccessToken)
	if false {
		req.Write(os.Stdout)
	}
	transport, err := t.transport()
	if err != nil {
		return nil, err
	}
	return transport.RoundTrip(req)
}

func (t *JWTTransport) CreateToken() {
	c := &jwt.ClaimSet{
		Iss:   secrets.Keys().ServiceAccountEmailAddress,
		Scope: "https://www.googleapis.com/auth/bigquery",
		Aud:   "https://accounts.google.com/o/oauth2/token",
	}
	t.Token = &jwt.Token{
		Header:   jwt.StdHeader,
		ClaimSet: c,
		Key:      []byte(secrets.Keys().ServiceAccountPrivateKey),
	}
}

func (t *JWTTransport) UpdateOauth() error {
	oauth, err := t.Token.Assert(t.Client)
	if err != nil {
		return err
	}
	t.Oauth = oauth
	return nil
}

func (t *JWTTransport) transport() (*urlfetch.Transport, error) {

	deadline, err := time.ParseDuration("60s")
	if err != nil {
		return nil, err
	}

	transport := &urlfetch.Transport{
		Context:  appengine.NewContext(t.Request),
		Deadline: deadline,
	}

	return transport, nil
	/*return &oauth.Transport{
		Config: &oauth.Config{},
		Token: t.Oauth,
	}*/
}

func toJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}
	buffer := bytes.NewBuffer(b)
	json := buffer.String()

	return json
}

func getJWTClient(r *http.Request) *http.Client {
	context := appengine.NewContext(r)
	client := urlfetch.Client(context)
	transport := &JWTTransport{
		Request: r,
		Client:  client,
	}
	jwtClient := &http.Client{Transport: transport}
	return jwtClient
}

var (
	DefaultTransport      = &JWTTransport{}
	DefaultBigqueryClient = &http.Client{Transport: DefaultTransport}
)

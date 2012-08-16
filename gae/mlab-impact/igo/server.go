/* The igo package is intended to provide the core m-lab impact functionality.  */
package igo

/*
 */

import ( // for docs http://golang.org/pkg/ pkgname
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"bytes"
	"encoding/json"
	"fmt"
	"impact/data/ndt"
	"impact/queryHandler"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const (
	api_key = "AIzaSyCxu5gJmbAFcJp7mRoIeTBg555M9bxoOOQ"
)

var baseConfig = map[string]string{
	"trackSetTrackingCode": "'UA-16051354-8'",
	"trackSetRemoteAddr":   "''",
	"trackSetReferer":      "''",
	"pageSetLoggedInAdmin": "false",
	"pageSetLoggedIn":      "false",
	"pageSetLoginUrl":      "''",
	"pageSetLogoutUrl":     "''",
	"pageSetUserName":      "''",
	"pageSetDisplayWidth":  "640",
	"pageSetDisplayHeight": "360",
}

var rootConfig = map[string]string{}

func newConfig(tempConfig map[string]string, c appengine.Context, r *http.Request, gt int64) (map[string]string, error) {
	config := make(map[string]string, len(tempConfig)+len(baseConfig))
	for k, v := range baseConfig {
		config[k] = v
	}
	for k, v := range tempConfig {
		config[k] = v
	}
	u := user.Current(c)
	if u != nil {
		go StoreUser(c, u)
		if user.IsAdmin(c) {
			config["pageSetLoggedInAdmin"] = "true"
		} else {
			config["pageSetLoggedInAdmin"] = "false"
		}
		config["pageSetLoggedIn"] = "true"
		config["pageSetUserName"] = fmt.Sprintf("'%s'", u.String())
	} else {
		config["pageSetLoggedIn"] = "false"
		config["pageSetLoggedIn"] = "false"
		config["pageSetUserName"] = "''"
		config["pageSetLoggedInAdmin"] = "false"
	}
	if loginUrl, err := user.LoginURL(c, r.URL.Path); err == nil {
		config["pageSetLoginUrl"] = fmt.Sprintf("'%s'", loginUrl)
	} else {
		c.Errorf("newConfig:user.LoginURL err = %v", err)
		return config, err
	}
	if logoutUrl, err := user.LogoutURL(c, r.URL.Path); err == nil {
		config["pageSetLogoutUrl"] = fmt.Sprintf("'%s'", logoutUrl)
	} else {
		c.Errorf("newConfig:user.LogoutURL err = %v", err)
		return config, err
	}

	config["trackSetRemoteAddr"] = fmt.Sprintf("'%s'", r.RemoteAddr)
	config["trackSetReferer"] = fmt.Sprintf("'%s'", r.Referer())
	config["trackSetRequestTime"] = fmt.Sprintf("%v", gt)
	return config, nil
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/query", query)
	http.HandleFunc("/bq_job", bigqueryJob)
	//http.HandleFunc("/oauth2callback", oauth2callback)
	http.HandleFunc("/admin/logout", logout)
	http.HandleFunc("/user/logout", logout)
	http.HandleFunc("/admin/login", login)
	http.HandleFunc("/user/login", login)
	staticSet.Parse(headHTML)
	staticSet.Parse(pageHTML)
}

type HTMLHeader struct {
	Title      string
	ScriptList []string
	StyleList  []string
	JsConfig   map[string]string
}

type HTMLPage struct {
	Head   HTMLHeader
	OnLoad string
}

func VersionFile(path string, fileName string, c appengine.Context) string {
	v := appengine.VersionID(c)
	sv := strings.Split(v, ".")[0]
	return fmt.Sprintf("%s%s.%s", path, sv, fileName)
}

func MapURL(key string, sensorOn string) string {
	return fmt.Sprintf("http://maps.googleapis.com/maps/api/js?key=%s&sensor=%s", key, sensorOn)
}

const headHTML = `{{define "header"}}<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN""http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"><html xmlns="http://www.w3.org/1999/xhtml" lang="en" xml:lang="en"><head><title>{{.Title}}</title><meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1" /><meta name="keywords" content="broadband,performance,network,Internet,measurement" /><meta name="description" content="a tool for measuring and analyzing broadband performance" />{{range .ScriptList}}<script type="text/javascript" src="{{html .}}"></script>{{end}}{{range .StyleList}}<link rel="stylesheet" href="{{html .}}" />{{end}}<script type="text/javascript">{{range $k, $v := .JsConfig}}{{printf "%s(%s);" $k $v}}{{end}}</script></head>{{end}}`

const pageHTML = `{{define "page"}}{{template "header" .Head}}<body onload="{{js .OnLoad}}"><noscript>To use M-Lab Impact you need to allow JavaScript to run.</noscript></body></html>{{end}}`

var staticSet = new(template.Template)

/*Writes the User entity to our datastore.

This never needs to be put in memcache because we never need to read it.
There may be a small gain to be made by putting the key in memcache, but that
only saves 1 Small Read on a hit.
*/
func StoreUser(c appengine.Context, u *user.User) {
	key := datastore.NewKey(c, "User", u.ID, 0, nil)
	eu := new(user.User)
	err := datastore.Get(c, key, eu)
	if err == datastore.ErrNoSuchEntity {
		if _, err := datastore.Put(c, key, u); err != nil {
			c.Errorf("PUT user (%s) err = %v", u, err)
		}
		c.Infof("WELCOME NEW USER = %s", u)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		//already logged out send them home
		http.Redirect(w, r, "/", 307)
		return
	}
	//log them out and then home
	logoutUrl, err := user.LogoutURL(c, "/")
	if err != nil {
		c.Errorf("logout:user.LogoutURL err = %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, logoutUrl, 307)
}

func login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u != nil {
		//already logged in send them home
		http.Redirect(w, r, "/", 307)
		return
	}
	//send them to the login page and then home
	loginUrl, err := user.LoginURL(c, "/")
	if err != nil {
		c.Errorf("login:user.LoginURL err = %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, loginUrl, 307)
}

func query(w http.ResponseWriter, r *http.Request) {

	result, err := queryHandler.GetResult(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(result)

	if err == nil {
		buffer := bytes.NewBuffer(b)
		json := buffer.String()

		fmt.Fprint(w, json)
	} else {
		fmt.Fprint(w, err.Error())
		//http.Error(w, "Error marshalling JSON values", 400)
	}
}

func bigqueryJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.FormValue("jobID")

	result, err := ndt.NDT_Source().JobResult(r, jobID)

	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	b, err := json.Marshal(result)

	if err == nil {
		buffer := bytes.NewBuffer(b)
		json := buffer.String()

		fmt.Fprint(w, json)
	} else {
		fmt.Fprint(w, err.Error())
		//http.Error(w, "Error marshalling JSON values", 400)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	gt := time.Now().Unix()
	c := appengine.NewContext(r)
	//get the user if we have one.
	fc, err := newConfig(rootConfig, c, r, gt)
	if err != nil {
		c.Errorf("root:newConfig err = %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h := HTMLHeader{
		Title: "M-Lab Impact",
		ScriptList: []string{
			"https://www.google.com/jsapi",
			VersionFile("/", "impact-comp.js", c),
			MapURL(api_key, "false"),
		},
		StyleList: []string{VersionFile("/", "impact-comp.css", c)},
		JsConfig:  fc,
	}
	p := &HTMLPage{
		Head:   h,
		OnLoad: "window.impactStart();",
	}
	err = staticSet.ExecuteTemplate(w, "page", p)
	if err != nil {
		c.Errorf("root:staticSet.ExecuteTemplate err = %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

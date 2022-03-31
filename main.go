package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

var ProxyToken Token
var TargetURL *url.URL

func main() {
	expPort, protocol, authPath, usr, conHost, conPort := readEnv()
	jsonValue, err := json.Marshal(usr)
	if err != nil {
		panic(err)
	}
	u := protocol + "://" + conHost + ":" + conPort
	wp, err := http.Post(u+authPath, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	defer wp.Body.Close()
	body, err := ioutil.ReadAll(wp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &ProxyToken); err != nil {
		panic(err)
	}

	TargetURL, err = url.Parse(u)
	if err != nil {
		panic(err)
	}

	log.Println("Starting proxy on port:", expPort, "for", u)
	http.HandleFunc("/", handleRequestAndRedirect)
	err = http.ListenAndServe(":"+expPort, nil)
	if err != nil {
		panic(err)
	}
}

func handleRequestAndRedirect(w http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(TargetURL)
	req.URL.Host = TargetURL.Host
	req.URL.Scheme = TargetURL.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = TargetURL.Host
	proxy.ModifyResponse = func(r *http.Response) error {
		if strings.Contains(r.Header.Get("Content-Type"), "html") {
			b, _ := ioutil.ReadAll(r.Body)
			// inject localStorage
			buf := bytes.NewBufferString(
				"<script>" +
					"localStorage.setItem('loggedIn', 'true');" +
					"localStorage.setItem('token', JSON.stringify('" + ProxyToken.Token + "'));" +
					"</script>")
			buf.Write(b)
			r.Body = ioutil.NopCloser(buf)
			r.Header["Content-Length"] = []string{fmt.Sprint(buf.Len())}
		}
		return nil
	}
	proxy.ServeHTTP(w, req)
}
func readEnv() (expPort string, protocol string, authPath string, usr User, conHost string, conPort string) {
	var ok bool
	protocol, ok = os.LookupEnv("PROTOCOL")
	if !ok {
		protocol = "http"
	}
	authPath, ok = os.LookupEnv("AUTH_PATH")
	if !ok {
		authPath = "/auth/login"
	}
	expPort, ok = os.LookupEnv("EXPOSE_PORT")
	if !ok {
		expPort = "9999"
	}
	conHost, ok = os.LookupEnv("CONNECT_HOST")
	if !ok {
		conHost = "192.168.2.128"
	}
	conPort, ok = os.LookupEnv("CONNECT_PORT")
	if !ok {
		conPort = "4000"
	}
	usr.Username, ok = os.LookupEnv("USER")
	if !ok {
		usr.Username = "admin"
	}
	usr.Password, ok = os.LookupEnv("PASSWORD")
	if !ok {
		usr.Password = "zero-ui"
	}
	return
}

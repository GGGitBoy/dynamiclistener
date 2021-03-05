package dynamiclistener

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Approach taken from letsencrypt, except manglePort is specific to us
func HTTPRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Printf("jiandao === HTTPRedirect Content-Type %+v\n", r.Header)
			fmt.Printf("jiandao === HTTPRedirect ProtoMajor %+v\n", r.ProtoMajor)
			fmt.Printf("jiandao === HTTPRedirect TLS %+v\n", r.TLS)
			fmt.Printf("jiandao === HTTPRedirect Path %+v\n", r.URL.Path)
			fmt.Printf("jiandao === HTTPRedirect Method %+v\n", r.Method)
			if r.TLS != nil ||
				r.Header.Get("x-Forwarded-Proto") == "https" ||
				r.Header.Get("x-Forwarded-Proto") == "wss" ||
				strings.HasPrefix(r.URL.Path, "/.well-known/") ||
				strings.HasPrefix(r.URL.Path, "/ping") ||
				strings.HasPrefix(r.URL.Path, "/health") {
				next.ServeHTTP(rw, r)
				return
			}
			if r.Method != "GET" && r.Method != "HEAD" {
				http.Error(rw, "Use HTTPS", http.StatusBadRequest)
				return
			}
			target := "https://" + manglePort(r.Host) + r.URL.RequestURI()
			http.Redirect(rw, r, target, http.StatusFound)
		})
}

func manglePort(hostport string) string {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return hostport
	}

	portInt = ((portInt / 1000) * 1000) + 443

	return net.JoinHostPort(host, strconv.Itoa(portInt))
}

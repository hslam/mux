package proxy
import (
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"net/url"
)
var httpTransport *http.Transport
func init() {
	httpTransport = &http.Transport{
		MaxIdleConnsPerHost: 65536,
		MaxConnsPerHost:65536,
	}
}
func Proxy(w http.ResponseWriter, r *http.Request, remote_url string) {
	remote_url_parse, _ := url.Parse(remote_url)
	remote, err := url.Parse("http://"+remote_url_parse.Host)
	if err != nil {
		panic(err)
	}
	r.URL.Path = remote_url_parse.Path
	proxy := httputil.NewSingleHostReverseProxy(remote)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		actualHost, err := actualRemoteHost(req)
		if err == nil {
			req.Header.Set("HTTP_X_FORWARDED_FOR", actualHost)
		}
		director(req)
	}
	proxy.Transport = httpTransport
	proxy.ServeHTTP(w, r)
}
func actualRemoteHost(r *http.Request) (host string, err error) {
	host = r.Header.Get("HTTP_X_FORWARDED_FOR")
	if host == "" {
		host = r.Header.Get("X-FORWARDED-FOR")
	}
	if strings.Contains(host, ",") {
		host = host[0:strings.Index(host, ",")]
	}
	if host == "" {
		host, _, err = net.SplitHostPort(r.RemoteAddr)
	}
	return
}
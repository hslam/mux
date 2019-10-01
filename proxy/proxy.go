package proxy
import (
	"net/http"
	"net/http/httputil"
	"net/url"
)
var httpTransport *http.Transport
func init() {
	httpTransport = &http.Transport{
		Proxy:http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: 65536,
		MaxConnsPerHost:65536,
	}
}
func Proxy(w http.ResponseWriter, r *http.Request, target_url string) {
	target_url_parse,err := url.Parse(target_url)
	if err != nil {
		panic(err)
	}
	target, _ := url.Parse("http://"+target_url_parse.Host)
	r.URL.Path = target_url_parse.Path
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = httpTransport
	proxy.ServeHTTP(w, r)
}

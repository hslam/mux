package main
import (
	"log"
	"net/http"
	"hslam.com/mgit/Mort/mux"
	"hslam.com/mgit/Mort/mux/gzip"
	"fmt"
)
func main() {
	router := mux.New()
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found : "+r.URL.String(), http.StatusNotFound)
	})
	router.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})
	router.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Host:%s Path:%s Method:%s\n",r.Host,r.URL.Path,r.Method)
	})
	router.HandleFunc("/hello/:key/mort/:value/huang", func(w http.ResponseWriter, r *http.Request) {
		params:=router.Params(r)
		w.Write([]byte(fmt.Sprintf("hello Method:%s key:%s value:%s\n",r.Method,params["key"], params["value"])))
	}).GET().POST()
	router.Group("/group", func(router *mux.Router) {
		router.HandleFunc("/:key/mort/:value/huang", func(w http.ResponseWriter, r *http.Request) {
			params:=router.Params(r)
			w.Write([]byte(fmt.Sprintf("group Method:%s key:%s value:%s\n",r.Method,params["key"], params["value"])))
		}).GET().POST()
		router.HandleFunc("/:foo/:bar", func(w http.ResponseWriter, r *http.Request) {
			params:=router.Params(r)
			gzip.WriteGzip(w,r,http.StatusOK,[]byte(fmt.Sprintf("group Method:%s foo:%s bar:%s\n",r.Method,params["foo"], params["bar"])))
		}).All()
	})
	router.Once()//before listen
	log.Fatal(http.ListenAndServe(":8080", router))
}
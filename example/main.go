package main
import (
	"log"
	"net/http"
	"hslam.com/mgit/Mort/mux"
	"fmt"
)
func main() {
	router := mux.New()
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found : "+r.URL.String(), http.StatusNotFound)
	})
	router.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Host:%s Path:%s Method:%s\n",r.Host,r.URL.Path,r.Method)
	})
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello world Method:%s\n",r.Method)))
	}).All()
	router.HandleFunc("/hello/:key/mort/:value/huang", func(w http.ResponseWriter, r *http.Request) {
		params:=router.Params(r)
		w.Write([]byte(fmt.Sprintf("hello key:%s value:%s\n",params["key"], params["value"])))
	}).GET().POST().PUT().DELETE().End()
	router.Group("/group", func(router *mux.Router) {
		router.HandleFunc("/foo/:id", func(w http.ResponseWriter, r *http.Request) {
			params:=router.Params(r)
			w.Write([]byte(fmt.Sprintf("group/foo id:%s\n",params["id"])))
		}).GET()
		router.HandleFunc("/bar/:id", func(w http.ResponseWriter, r *http.Request) {
			params:=router.Params(r)
			w.Write([]byte(fmt.Sprintf("group/bar id:%s\n",params["id"])))
		}).GET()
	})
	router.Once()//before listen
	log.Fatal(http.ListenAndServe(":8080", router))
}
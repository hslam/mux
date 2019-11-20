package main
import (
	"log"
	"fmt"
	"net/http"
	"hslam.com/git/x/mux"
)
func main() {
	m := mux.New()
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found : "+r.URL.String(), http.StatusNotFound)
	})
	m.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Host:%s Path:%s Method:%s\n",r.Host,r.URL.Path,r.Method)
	})
	m.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello world Method:%s\n",r.Method)))
	}).All()
	m.HandleFunc("/hello/:key/mort/:value/huang", func(w http.ResponseWriter, r *http.Request) {
		params:=m.Params(r)
		w.Write([]byte(fmt.Sprintf("hello key:%s value:%s\n",params["key"], params["value"])))
	}).GET().POST().PUT().DELETE()
	m.Group("/group", func(m *mux.Mux) {
		m.HandleFunc("/foo/:id", func(w http.ResponseWriter, r *http.Request) {
			params:=m.Params(r)
			w.Write([]byte(fmt.Sprintf("group/foo id:%s\n",params["id"])))
		}).GET()
		m.HandleFunc("/bar/:id", func(w http.ResponseWriter, r *http.Request) {
			params:=m.Params(r)
			w.Write([]byte(fmt.Sprintf("group/bar id:%s\n",params["id"])))
		}).GET()
	})
	log.Fatal(http.ListenAndServe(":8080", m))
}
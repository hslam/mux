package main
import (
	"log"
	"net/http"
	"hslam.com/mgit/Mort/mux"
	"fmt"
)
func main() {
	router := mux.New()
	router.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.Host)
	})
	router.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
	})
	router.HandleFunc("/hello/:key/mort/:value/huang", func(w http.ResponseWriter, r *http.Request) {
		params:=router.Params(r)
		w.Write([]byte(fmt.Sprintf("hello world Method:%s key:%s value:%s",r.Method,params["key"], params["value"])))
	}).GET().POST()
	log.Fatal(http.ListenAndServe(":8080", router))
}
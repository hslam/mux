package main

import (
	"log"
	"net/http"
	"hslam.com/mgit/Mort/mux"
	"fmt"
)
func main() {
	router := mux.New()
	router.Middleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
	})
	router.HandleFunc("/hello/:key/meng/:value/huang", func(w http.ResponseWriter, r *http.Request) {
		params:=router.Params(r)
		w.Write([]byte(fmt.Sprintf("hello world Method:%s key:%s value:%s",r.Method,params["key"], params["value"])))
	}).GET().POST()
	log.Fatal(http.ListenAndServe(":8080", router))
}
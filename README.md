# mux
## A url path router implementation written in Golang.

## Features

* Middleware
* Group
* Path matching and routing
* Fully compatible with the http.HandlerFunc
* Not found HandlerFunc setting

## Get started

### Install
```
go get hslam.com/mgit/Mort/mux
```
### Import
```
import "hslam.com/mgit/Mort/mux"
```
### Usage
#### Example
```
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
		w.Write([]byte(fmt.Sprintf("hello Method:%s key:%s value:%s\n",r.Method,params["key"], params["value"])))
	}).GET().POST()
	router.Group("/group", func(router *mux.Router) {
		router.HandleFunc("/:key/mort/:value/huang", func(w http.ResponseWriter, r *http.Request) {
			params:=router.Params(r)
			w.Write([]byte(fmt.Sprintf("group Method:%s key:%s value:%s\n",r.Method,params["key"], params["value"])))
		}).GET().POST()
		router.HandleFunc("/:foo/:bar", func(w http.ResponseWriter, r *http.Request) {
			params:=router.Params(r)
			w.Write([]byte(fmt.Sprintf("group Method:%s foo:%s bar:%s\n",r.Method,params["foo"], params["bar"])))
		}).GET().POST()
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}
```

curl http://localhost:8080/hello/123/mort/456/huang
#### Output
```
hello Method:GET key:123 value:456
```
curl -XPOST http://localhost:8080/hello/123/mort/456/huang
#### Output
```
hello Method:POST key:123 value:456
```
curl http://localhost:8080/group/123/mort/456/huang
#### Output
```
group Method:GET key:123 value:456
```
curl -XPOST http://localhost:8080/group/123/mort/456/huang
#### Output
```
group Method:POST key:123 value:456
```

curl http://localhost:8080/group/123/456
#### Output
```
group Method:GET foo:123 bar:456
```
curl -XPOST http://localhost:8080/group/123/456
#### Output
```
group Method:POST foo:123 bar:456
```
### Licence
This package is licenced under a MIT licence (Copyright (c) 2017 Mort Huang)


### Authors
mux was written by Mort Huang.



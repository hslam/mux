# mux
## A url path router implementation written in Golang.

## Features

* Middleware
* Group
* Path matching and routing
* Fully compatible with the http.HandlerFunc
* Not found
* Gzip

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
			gz:=gzip.NewGzipWriter(w,r)
			defer gz.Close()
			gz.Write([]byte(fmt.Sprintf("group Method:%s foo:%s bar:%s\n",r.Method,params["foo"], params["bar"])))
		}).All()
	})
	router.Once()//before listen
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
curl -H "Accept-Encoding: gzip,deflate" --compressed http://localhost:8080/group/123/456
#### gzip Output
```
group Method:GET foo:123 bar:456
```
curl -XPOST http://localhost:8080/group/123/456
#### Output
```
group Method:POST foo:123 bar:456
```
curl --HEAD http://localhost:8080/group/123/456

or

curl -I http://localhost:8080/group/123/456
#### Output
```
HTTP/1.1 200 OK
Date: Mon, 30 Sep 2019 10:01:11 GMT
Content-Length: 34
Content-Type: text/plain; charset=utf-8
```
curl -XPATCH http://localhost:8080/group/123/456
#### Output
```
group Method:PATCH foo:123 bar:456
```
curl -XOPTIONS http://localhost:8080/group/123/456
#### Output
```
group Method:OPTIONS foo:123 bar:456
```
### Licence
This package is licenced under a MIT licence (Copyright (c) 2017 Mort Huang)


### Authors
mux was written by Mort Huang.



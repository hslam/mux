# mux
## A url path router implementation written in Golang.

## Features

* Middleware
* Group
* Path matching and routing
* Fully compatible with the http.HandlerFunc
* Not found
* [Mux Handler](https://hslam.com/git/x/handler "handler")

## Get started

### Install
```
go get hslam.com/git/x/mux
```
### Import
```
import "hslam.com/git/x/mux"
```
### Usage
#### Example
```
package main
import (
	"log"
	"net/http"
	"hslam.com/git/x/mux"
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
	router.Once()//before listening
	log.Fatal(http.ListenAndServe(":8080", router))
}
```

curl -XGET http://localhost:8080/hello
```
hello world Method:GET
```

curl -XPOST http://localhost:8080/hello
```
hello world Method:POST
```

curl --HEAD http://localhost:8080/hello

or

curl -I http://localhost:8080/hello
```
HTTP/1.1 200 OK
Date: Tue, 01 Oct 2019 20:28:42 GMT
Content-Length: 24
Content-Type: text/plain; charset=utf-8
```

curl -XPATCH http://localhost:8080/hello
```
hello world Method:PATCH
```

curl -XOPTIONS http://localhost:8080/hello
```
hello world Method:OPTIONS
```

curl http://localhost:8080/hello/123/mort/456/huang
```
hello key:123 value:456
```
curl http://localhost:8080/group/foo/123
```
group/foo id:123
```
curl http://localhost:8080/group/bar/123
```
group/bar id:123
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)


### Authors
mux was written by Mort Huang.



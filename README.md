# mux
An implementation of url path router written in Golang.

## Features

* Middleware
* Group
* Path matching and routing
* Fully compatible with the http.HandlerFunc
* Not found
* [Mux Handler](https://github.com/hslam/handler "handler")

## Get started

### Install
```
go get github.com/hslam/mux
```
### Import
```
import "github.com/hslam/mux"
```
### Usage
#### Example
```
package main
import (
	"log"
	"fmt"
	"net/http"
	"github.com/hslam/mux"
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
```

curl -XGET http://localhost:8080/hello
```
hello world Method:GET
```

curl -XPOST http://localhost:8080/hello
```
hello world Method:POST
```

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
This package is licenced under a MIT licence (Copyright (c) 2019 Meng Huang)


### Authors
mux was written by Meng Huang.



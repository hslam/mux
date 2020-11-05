# mux
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/mux)](https://pkg.go.dev/github.com/hslam/mux)
[![Build Status](https://travis-ci.org/hslam/mux.svg?branch=master)](https://travis-ci.org/hslam/mux)
[![codecov](https://codecov.io/gh/hslam/mux/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/mux)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/mux?v=7e100)](https://goreportcard.com/report/github.com/hslam/mux)
[![LICENSE](https://img.shields.io/github/license/hslam/mux.svg?style=flat-square)](https://github.com/hslam/mux/blob/master/LICENSE)

Package mux implements an HTTP request multiplexer.

## Features

* Middleware
* Group
* Path matching and routing
* Fully compatible with the http.HandlerFunc
* Not found
* [HTTP Handler](https://github.com/hslam/handler "handler")

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
	"fmt"
	"github.com/hslam/mux"
	"log"
	"net/http"
)

func main() {
	m := mux.New()
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found : "+r.URL.String(), http.StatusNotFound)
	})
	m.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Host:%s Path:%s Method:%s\n", r.Host, r.URL.Path, r.Method)
	})
	m.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello world Method:%s\n", r.Method)))
	}).All()
	m.HandleFunc("/hello/:key/world/:value", func(w http.ResponseWriter, r *http.Request) {
		params := m.Params(r)
		w.Write([]byte(fmt.Sprintf("hello key:%s value:%s\n", params["key"], params["value"])))
	}).GET().POST().PUT().DELETE()
	m.Group("/group", func(m *mux.Mux) {
		m.HandleFunc("/foo/:id", func(w http.ResponseWriter, r *http.Request) {
			params := m.Params(r)
			w.Write([]byte(fmt.Sprintf("group/foo id:%s\n", params["id"])))
		}).GET()
		m.HandleFunc("/bar/:id", func(w http.ResponseWriter, r *http.Request) {
			params := m.Params(r)
			w.Write([]byte(fmt.Sprintf("group/bar id:%s\n", params["id"])))
		}).GET()
	})
	log.Fatal(http.ListenAndServe(":8080", m))
}
```

curl -XGET http://localhost:8080/favicon.ico
```
Not Found : /favicon.ico
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

curl http://localhost:8080/hello/123/world/456
```
hello key:123 value:456
```
curl http://localhost:8080/group/foo/1
```
group/foo id:1
```
curl http://localhost:8080/group/bar/2
```
group/bar id:2
```

### License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)


### Author
mux was written by Meng Huang.



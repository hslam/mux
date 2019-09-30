# mux
## A worker pool implementation written in Golang.

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
	router.Middleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
	})
	router.HandleFunc("/hello/:key/meng/:value/huang", func(w http.ResponseWriter, r *http.Request) {
		params:=router.Params(r)
		w.Write([]byte(fmt.Sprintf("hello world Method:%s key:%s value:%s",r.Method,params["key"], params["value"])))
	}).GET().POST()
	log.Fatal(http.ListenAndServe(":8080", router))
}
```


http://127.0.0.1:8080/count
```
curl http://localhost:8080/hello/123/meng/456/huang
#### Output
```
hello world Method:GET key:123 value:456
```
curl -XPOST http://localhost:8080/hello/123/meng/456/huang
#### Output
```
hello world Method:POST key:123 value:456
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2017 Mort Huang)


### Authors
workerpool was written by Mort Huang.



package mux

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestParseMatch(t *testing.T) {
	pattern := "/db/:key/meng/:value/huang"
	i := strings.Index(pattern, ":")
	prefix := pattern[:i]
	match := strings.Split(pattern[i:], "/")
	params := make(map[string]string)
	key := ""
	for i := 0; i < len(match); i++ {
		if strings.Contains(match[i], ":") {
			match[i] = strings.Trim(match[i], ":")
			params[match[i]] = ""
			if i > 0 {
				key += "/"
			}
		} else {
			key += "/" + match[i]
			match[i] = ""
		}
	}
	path := "/db/123/meng/456/huang"
	strs := strings.Split(strings.Trim(path, prefix), "/")
	if len(strs) == len(match) {
		for i := 0; i < len(strs); i++ {
			if match[i] != "" {
				if _, ok := params[match[i]]; ok {
					params[match[i]] = strs[i]
				}
			}
		}
	}
	if params["key"] != "123" || params["value"] != "456" {
		t.Error(params)
	}
}

func TestMux(t *testing.T) {
	m := New()
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found : "+r.URL.String(), http.StatusNotFound)
	})
	m.Use(func(w http.ResponseWriter, r *http.Request) {
	})
	m.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello world Method:%s\n", r.Method)))
	}).All()
	m.HandleFunc("/hello/:key/world/:value", func(w http.ResponseWriter, r *http.Request) {
		params := m.Params(r)
		w.Write([]byte(fmt.Sprintf("hello key:%s value:%s\n", params["key"], params["value"])))
	}).GET().POST().PUT().DELETE()
	m.Group("/group", func(m *Mux) {
		m.HandleFunc("/foo/:id", func(w http.ResponseWriter, r *http.Request) {
			params := m.Params(r)
			w.Write([]byte(fmt.Sprintf("group/foo id:%s\n", params["id"])))
		}).GET()
		m.HandleFunc("/bar/:id", func(w http.ResponseWriter, r *http.Request) {
			params := m.Params(r)
			w.Write([]byte(fmt.Sprintf("group/bar id:%s\n", params["id"])))
		}).GET()
	})
	addr := ":8080"
	httpServer := &http.Server{
		Addr:    addr,
		Handler: m,
	}
	l, _ := net.Listen("tcp", addr)
	go httpServer.Serve(l)
	if resp, err := http.Get("http://" + addr + "/hello"); err != nil {
		t.Error(err)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if string(body) != "hello world Method:GET\n" {
		t.Error(string(body))
	}
	if resp, err := http.Post("http://"+addr+"/hello", "", nil); err != nil {
		t.Error(err)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if string(body) != "hello world Method:POST\n" {
		t.Error(string(body))
	}
	if resp, err := http.Head("http://" + addr + "/hello"); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error(resp.ContentLength)
	}
	if resp, err := http.Get("http://" + addr + "/hello/123/world/456"); err != nil {
		t.Error(err)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if string(body) != "hello key:123 value:456\n" {
		t.Error(string(body))
	}
	if resp, err := http.Get("http://" + addr + "/group/foo/1"); err != nil {
		t.Error(err)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if string(body) != "group/foo id:1\n" {
		t.Error(string(body))
	}
	if resp, err := http.Get("http://" + addr + "/group/bar/2"); err != nil {
		t.Error(err)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if string(body) != "group/bar id:2\n" {
		t.Error(string(body))
	}
	httpServer.Close()
}

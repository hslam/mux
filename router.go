package mux

import (
	"net/http"
	"sync"
	"strings"
)
const (
	GET         = iota
	POST
	PUT
	DELETE
	PATCH
	HEAD
	OPTIONS
)


type Router struct {
	mut    		sync.RWMutex
	mux  		map[string]*Prefix
	middlewares []http.HandlerFunc
	notFound 	http.HandlerFunc
}
type Prefix struct {
	m 			map[string]*Entry
	pattern 	string
	prefix 		string
}
type Entry struct {
	handler		http.HandlerFunc
	key 		string
	match 		[]string
	params		map[string]string
	get 		http.HandlerFunc
	post 		http.HandlerFunc
	put 		http.HandlerFunc
	delete 		http.HandlerFunc
	patch 		http.HandlerFunc
	head     	http.HandlerFunc
	options	 	http.HandlerFunc
}

func New() *Router {
	router := &Router{
		mux: make(map[string]*Prefix),
	}
	return router
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mut.RLock()
	defer router.mut.RUnlock()
	if entry:=router.getHandlerFunc(r.URL.Path);entry!=nil{
		if r.Method=="GET"&&entry.get!=nil{
			router.middleware(w,r)
			entry.get(w,r)
			return
		}else if r.Method=="POST"&&entry.post!=nil{
			router.middleware(w,r)
			entry.post(w,r)
			return
		}else if r.Method=="PUT"&&entry.put!=nil{
			router.middleware(w,r)
			entry.put(w,r)
			return
		}else if r.Method=="DELETE"&&entry.delete!=nil{
			router.middleware(w,r)
			entry.delete(w,r)
			return
		}else if r.Method=="PATCH"&&entry.patch!=nil{
			router.middleware(w,r)
			entry.patch(w,r)
			return
		}else if r.Method=="HEAD"&&entry.head!=nil{
			router.middleware(w,r)
			entry.head(w,r)
			return
		}else if r.Method=="OPTIONS"&&entry.options!=nil{
			router.middleware(w,r)
			entry.options(w,r)
			return
		}
	}
	if router.notFound!=nil{
		router.notFound(w,r)
		return
	}
	http.Error(w, "Not Found : "+r.URL.String(), http.StatusNotFound)
}

func (router *Router) getHandlerFunc(path string) *Entry{
	if prefix,key,ok:=router.matchParams(path);ok {
		if entry, ok := router.mux[prefix].m[key]; ok {
			return entry
		}
	}
	return nil
}

func (router *Router) HandleFunc(pattern string, handler http.HandlerFunc) *Entry{
	router.mut.RLock()
	defer router.mut.RUnlock()
	prefix,key,match,params:=router.parseParams(pattern)
	if v, ok := router.mux[prefix]; ok {
		if entry, ok := v.m[key]; ok {
			entry.handler=handler
			entry.key=key
			entry.match=match
			entry.params=params
			router.mux[prefix].m[key] = entry
			return entry
		}
	}
	router.mux[prefix]=&Prefix{m:make(map[string]*Entry),pattern:pattern,prefix:prefix}
	entry:=&Entry{}
	entry.handler=handler
	entry.key=key
	entry.match=match
	entry.params=params
	router.mux[prefix].m[key] = entry
	return entry
}
func (router *Router) SetNotFound(handler http.HandlerFunc){
	router.mut.RLock()
	defer router.mut.RUnlock()
	router.notFound=handler
}
func (router *Router) Middleware(handler http.HandlerFunc){
	router.mut.RLock()
	defer router.mut.RUnlock()
	router.middlewares=append(router.middlewares,handler)
}
func (router *Router) middleware(w http.ResponseWriter, r *http.Request){
	for _,handler:=range router.middlewares{
		handler(w,r)
	}
}
func (router *Router) Params(r *http.Request)(map[string]string){
	router.mut.RLock()
	defer router.mut.RUnlock()
	params:=make(map[string]string)
	if prefix,key,ok:=router.matchParams(r.URL.Path);ok{
		if entry,ok:=router.mux[prefix].m[key];ok{
			strs := strings.Split(strings.Trim(r.URL.Path,prefix), "/")
			if len(strs)==len(entry.match){
				for i:=0;i<len(strs);i++{
					if entry.match[i]!=""{
						params[entry.match[i]]=strs[i]
					}
				}
			}
		}
	}
	return params
}
func (router *Router) matchParams(path string)(string,string,bool){
	for _,p:=range router.mux{
		if strings.HasPrefix(path,p.prefix){
			for _,v:=range p.m{
				form:=strings.Split(strings.Trim(path,p.prefix),"/")
				if len(form)==len(v.match){
					key:=""
					for i:=0;i<len(form);i++ {
						if v.match[i]!=""{
							if i>0{
								key+="/"
							}
						}else {
							if i>0{
								key+="/"+form[i]
							}else {
								key+=form[i]
							}
						}
					}
					if key==v.key{
						return p.prefix,v.key,true
					}
				}
			}
		}
	}
	return "","",false
}
func (router *Router) parseParams(pattern string)(string,string,[]string,map[string]string){
	prefix:=""
	var match []string
	key:=""
	params:=make(map[string]string)
	if strings.Contains(pattern,":"){
		idx:=strings.Index(pattern,":")
		prefix=pattern[:idx]
		match = strings.Split(pattern[idx:], "/")
		for i:=0;i<len(match);i++{
			if strings.Contains(match[i],":"){
				match[i]=strings.Trim(match[i],":")
				params[match[i]]=""
				if i>0{
					key+="/"
				}
			}else {
				if i>0{
					key+="/"+match[i]
				}else {
					key+=match[i]
				}
				match[i]=""
			}
		}
	}else{
		prefix=pattern
	}
	return prefix,key,match,params

}
func (entry *Entry) GET() *Entry{
	entry.get=entry.handler
	return entry
}
func (entry *Entry) POST() *Entry{
	entry.post=entry.handler
	return entry
}
func (entry *Entry) PUT()*Entry {
	entry.put=entry.handler
	return entry
}
func (entry *Entry) DELETE() *Entry{
	entry.delete=entry.handler
	return entry
}
func (entry *Entry) PATCH()*Entry {
	entry.patch=entry.handler
	return entry
}
func (entry *Entry) HEAD() *Entry{
	entry.head=entry.handler
	return entry
}
func (entry *Entry) OPTIONS() *Entry{
	entry.options=entry.handler
	return entry
}
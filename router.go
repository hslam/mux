package mux

import (
	"net/http"
	"sync"
	"strings"
	"fmt"
)
const (
	GET         = iota
	POST
	PUT
	DELETE
	PATCH
	HEAD
	OPTIONS
	TRACE
	CONNECT
)

type Router struct {
	mut    		sync.RWMutex
	prefixes  		map[string]*Prefix
	middlewares []http.HandlerFunc
	notFound 	http.HandlerFunc
	groups 		map[string]*Router
	group		string
}
type Prefix struct {
	m 			map[string]*Entry
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
	trace	 	http.HandlerFunc
	connect	 	http.HandlerFunc

}

func New() *Router {
	router := &Router{
		prefixes: make(map[string]*Prefix),
		groups: make(map[string]*Router),
	}
	return router
}

func newGroup(group string) *Router {
	router := &Router{
		prefixes: make(map[string]*Prefix),
		groups: make(map[string]*Router),
		group:group,
	}
	return router
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w,fmt.Sprint(err),http.StatusBadRequest)
		}
	}()
	router.mut.RLock()
	defer router.mut.RUnlock()
	if router.serve(w,r){
		return
	}else {
		for _,groupRouter:=range router.groups{
			if groupRouter.serve(w,r){
				return
			}
		}
	}
	if router.notFound!=nil{
		router.notFound(w,r)
		return
	}
	http.Error(w, "404 Not Found : "+r.URL.String(), http.StatusNotFound)
}
func (router *Router) serve(w http.ResponseWriter, r *http.Request)bool {
	path:=router.replace(r.URL.Path)
	if entry:=router.getHandlerFunc(path);entry!=nil{
		if r.Method=="GET"&&entry.get!=nil{
			router.serveEntry(entry.get,w,r)
			return true
		}else if r.Method=="POST"&&entry.post!=nil{
			router.serveEntry(entry.post,w,r)
			return true
		}else if r.Method=="PUT"&&entry.put!=nil{
			router.serveEntry(entry.put,w,r)
			return true
		}else if r.Method=="DELETE"&&entry.delete!=nil{
			router.serveEntry(entry.delete,w,r)
			return true
		}else if r.Method=="PATCH"&&entry.patch!=nil{
			router.serveEntry(entry.patch,w,r)
			return true
		}else if r.Method=="HEAD"&&entry.head!=nil{
			router.serveEntry(entry.head,w,r)
			return true
		}else if r.Method=="OPTIONS"&&entry.options!=nil{
			router.serveEntry(entry.options,w,r)
			return true
		}else if r.Method=="TRACE"&&entry.trace!=nil{
			router.serveEntry(entry.trace,w,r)
			return true
		}else if r.Method=="CONNECT"&&entry.connect!=nil{
			router.serveEntry(entry.connect,w,r)
			return true
		}else if entry.handler!=nil{
			router.serveEntry(entry.handler,w,r)
			return true
		}
	}
	return false
}
func (router *Router) serveEntry(handler http.HandlerFunc,w http.ResponseWriter, r *http.Request) {
	router.middleware(w,r)
	handler(w,r)
}
func (router *Router) getHandlerFunc(path string) *Entry{
	if prefix,key,ok:=router.matchParams(path);ok {
		if entry, ok := router.prefixes[prefix].m[key]; ok {
			return entry
		}
	}
	return nil
}

func (router *Router) HandleFunc(pattern string, handler http.HandlerFunc) *Entry{
	router.mut.RLock()
	defer router.mut.RUnlock()
	pattern=router.replace(pattern)
	prefix,key,match,params:=router.parseParams(router.group+pattern)
	if v, ok := router.prefixes[prefix]; ok {
		if entry, ok := v.m[key]; ok {
			entry.handler=handler
			entry.key=key
			entry.match=match
			entry.params=params
			router.prefixes[prefix].m[key] = entry
			return entry
		}else {
			entry:=&Entry{}
			entry.handler=handler
			entry.key=key
			entry.match=match
			entry.params=params
			router.prefixes[prefix].m[key] = entry
			return entry
		}
	}else {
		router.prefixes[prefix]=&Prefix{m:make(map[string]*Entry),prefix:prefix}
		entry:=&Entry{}
		entry.handler=handler
		entry.key=key
		entry.match=match
		entry.params=params
		router.prefixes[prefix].m[key] = entry
		return entry
	}
}

func (router *Router) Group(group string,f func(router *Router)){
	router.mut.RLock()
	defer router.mut.RUnlock()
	group=router.replace(group)
	groupRouter:=newGroup(group)
	f(groupRouter)
	for _,p:=range groupRouter.prefixes{
		for _,v:=range p.m{
			v.End()
		}
	}
	if _,ok:=router.groups[group];ok{
		panic("Group Existed")
	}
	groupRouter.middlewares=router.middlewares
	router.groups[group]=groupRouter
}
func (router *Router) NotFound(handler http.HandlerFunc){
	router.mut.RLock()
	defer router.mut.RUnlock()
	router.notFound=handler
}
func (router *Router) Use(handler http.HandlerFunc){
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
	path:=router.replace(r.URL.Path)
	if prefix,key,ok:=router.matchParams(path);ok{
		if entry,ok:=router.prefixes[prefix].m[key];ok{
			strs := strings.Split(strings.Trim(path,prefix), "/")
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
	for _,p:=range router.prefixes{
		if strings.HasPrefix(path,p.prefix){
			for _,v:=range p.m{
				r:=strings.TrimLeft(path,p.prefix)
				if r==""{
					return p.prefix,v.key,true
				}
				form:=strings.Split(r,"/")
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
func (router *Router) Once(){
	router.mut.RLock()
	defer router.mut.RUnlock()
	for _,p:=range router.prefixes{
		for _,v:=range p.m{
			v.End()
		}
	}
	for _,groupRouter:=range router.groups{
		groupRouter.Once()
	}
}
func (router *Router) replace(s string) string {
	for strings.Contains(s,"//"){
		s=strings.ReplaceAll(s,"//","/")
	}
	return s
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
func (entry *Entry) TRACE() *Entry{
	entry.trace=entry.handler
	return entry
}
func (entry *Entry) CONNECT() *Entry{
	entry.connect=entry.handler
	return entry
}
func (entry *Entry) End(){
	entry.handler=nil
}
func (entry *Entry) All() {
	entry.GET()
	entry.POST()
	entry.HEAD()
	entry.OPTIONS()
	entry.PUT()
	entry.PATCH()
	entry.DELETE()
	entry.TRACE()
	entry.CONNECT()
	entry.End()
}


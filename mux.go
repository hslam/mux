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

type Mux struct {
	mut    		sync.RWMutex
	prefixes	map[string]*Prefix
	middlewares []http.HandlerFunc
	notFound 	http.HandlerFunc
	groups 		map[string]*Mux
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

func New() *Mux {
	m := &Mux{
		prefixes: make(map[string]*Prefix),
		groups: make(map[string]*Mux),
	}
	return m
}

func newGroup(group string) *Mux {
	m := &Mux{
		prefixes: make(map[string]*Prefix),
		groups: make(map[string]*Mux),
		group:group,
	}
	return m
}
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w,fmt.Sprint(err),http.StatusBadRequest)
		}
	}()
	m.mut.RLock()
	defer m.mut.RUnlock()
	if m.serve(w,r){
		return
	}else {
		for _,groupMux:=range m.groups{
			if groupMux.serve(w,r){
				return
			}
		}
	}
	if m.notFound!=nil{
		m.notFound(w,r)
		return
	}
	http.Error(w, "404 Not Found : "+r.URL.String(), http.StatusNotFound)
}
func (m *Mux) serve(w http.ResponseWriter, r *http.Request)bool {
	path:=m.replace(r.URL.Path)
	if entry:=m.getHandlerFunc(path);entry!=nil{
		if r.Method=="GET"&&entry.get!=nil{
			m.serveEntry(entry.get,w,r)
			return true
		}else if r.Method=="POST"&&entry.post!=nil{
			m.serveEntry(entry.post,w,r)
			return true
		}else if r.Method=="PUT"&&entry.put!=nil{
			m.serveEntry(entry.put,w,r)
			return true
		}else if r.Method=="DELETE"&&entry.delete!=nil{
			m.serveEntry(entry.delete,w,r)
			return true
		}else if r.Method=="PATCH"&&entry.patch!=nil{
			m.serveEntry(entry.patch,w,r)
			return true
		}else if r.Method=="HEAD"&&entry.head!=nil{
			m.serveEntry(entry.head,w,r)
			return true
		}else if r.Method=="OPTIONS"&&entry.options!=nil{
			m.serveEntry(entry.options,w,r)
			return true
		}else if r.Method=="TRACE"&&entry.trace!=nil{
			m.serveEntry(entry.trace,w,r)
			return true
		}else if r.Method=="CONNECT"&&entry.connect!=nil{
			m.serveEntry(entry.connect,w,r)
			return true
		}else if entry.isEmpty()&&entry.handler!=nil{
			m.serveEntry(entry.handler,w,r)
			return true
		}
	}
	return false
}
func (m *Mux) serveEntry(handler http.HandlerFunc,w http.ResponseWriter, r *http.Request) {
	m.middleware(w,r)
	handler(w,r)
}
func (m *Mux) getHandlerFunc(path string) *Entry{
	if prefix,key,ok:=m.matchParams(path);ok {
		if entry, ok := m.prefixes[prefix].m[key]; ok {
			return entry
		}
	}
	return nil
}

func (m *Mux) HandleFunc(pattern string, handler http.HandlerFunc) *Entry{
	m.mut.RLock()
	defer m.mut.RUnlock()
	pattern=m.replace(pattern)
	prefix,key,match,params:=m.parseParams(m.group+pattern)
	if v, ok := m.prefixes[prefix]; ok {
		if entry, ok := v.m[key]; ok {
			entry.handler=handler
			entry.key=key
			entry.match=match
			entry.params=params
			m.prefixes[prefix].m[key] = entry
			return entry
		}else {
			entry:=&Entry{}
			entry.handler=handler
			entry.key=key
			entry.match=match
			entry.params=params
			m.prefixes[prefix].m[key] = entry
			return entry
		}
	}else {
		m.prefixes[prefix]=&Prefix{m:make(map[string]*Entry),prefix:prefix}
		entry:=&Entry{}
		entry.handler=handler
		entry.key=key
		entry.match=match
		entry.params=params
		m.prefixes[prefix].m[key] = entry
		return entry
	}
}

func (m *Mux) Group(group string,f func(m *Mux)){
	m.mut.RLock()
	defer m.mut.RUnlock()
	group=m.replace(group)
	groupMux:=newGroup(group)
	f(groupMux)
	for _,p:=range groupMux.prefixes{
		for _,v:=range p.m{
			v.FIN()
		}
	}
	if _,ok:=m.groups[group];ok{
		panic("Group Existed")
	}
	groupMux.middlewares=m.middlewares
	m.groups[group]=groupMux
}
func (m *Mux) NotFound(handler http.HandlerFunc){
	m.mut.RLock()
	defer m.mut.RUnlock()
	m.notFound=handler
}
func (m *Mux) Use(handler http.HandlerFunc){
	m.mut.RLock()
	defer m.mut.RUnlock()
	m.middlewares=append(m.middlewares,handler)
}
func (m *Mux) middleware(w http.ResponseWriter, r *http.Request){
	for _,handler:=range m.middlewares{
		handler(w,r)
	}
}

func (m *Mux) Params(r *http.Request)(map[string]string){
	m.mut.RLock()
	defer m.mut.RUnlock()
	params:=make(map[string]string)
	path:=m.replace(r.URL.Path)
	if prefix,key,ok:=m.matchParams(path);ok{
		if entry,ok:=m.prefixes[prefix].m[key];ok{
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
func (m *Mux) matchParams(path string)(string,string,bool){
	for _,p:=range m.prefixes{
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
func (m *Mux) parseParams(pattern string)(string,string,[]string,map[string]string){
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
func (m *Mux) Once(){
	m.mut.RLock()
	defer m.mut.RUnlock()
	for _,p:=range m.prefixes{
		for _,v:=range p.m{
			v.FIN()
		}
	}
	for _,groupMux:=range m.groups{
		groupMux.Once()
	}
}
func (m *Mux) replace(s string) string {
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
func (entry *Entry) FIN(){
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
	entry.FIN()
}
func (entry *Entry) isEmpty()bool{
	if entry.get==nil&&
		entry.post==nil&&
			entry.head==nil&&
				entry.options==nil&&
					entry.put==nil&&
						entry.patch==nil&&
							entry.delete==nil&&
								entry.trace==nil&&
									entry.connect==nil{
		return true
	}
	return false
}


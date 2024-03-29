// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package mux implements an HTTP request multiplexer.
package mux

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (
	options = iota
	get
	head
	post
	put
	delete
	trace
	connect
	patch
)

// ErrGroupExisted is the error returned by Group when registers a existed group.
var ErrGroupExisted = errors.New("Group Existed")

// ErrParamsKeyEmpty is the error returned by HandleFunc when the params key is empty.
var ErrParamsKeyEmpty = errors.New("Params key must be not empty")

// contextKey is a key for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

// String returns a context key.
func (k *contextKey) String() string { return "github.com/hslam/mux context key " + k.name }

// RecoveryContextKey is a context key.
var RecoveryContextKey = &contextKey{"recovery"}

// Mux is an HTTP request multiplexer.
type Mux struct {
	mut      sync.RWMutex
	prefixes map[string]*prefix
	group    string
	groups   map[string]*Mux
	context  struct {
		middlewares []http.Handler
		recovery    http.Handler
		notFound    http.Handler
	}
}

type prefix struct {
	prefix string
	m      map[string]*Entry
}

// Entry represents an HTTP HandlerFunc entry.
type Entry struct {
	handler  http.Handler
	handlers [9]http.Handler
	key      string
	match    []string
	params   map[string]string
}

// New returns a new Mux.
func New() *Mux {
	m := &Mux{
		prefixes: make(map[string]*prefix),
		groups:   make(map[string]*Mux),
	}
	return m
}

func newGroup(group string) *Mux {
	m := &Mux{
		prefixes: make(map[string]*prefix),
		groups:   make(map[string]*Mux),
		group:    group,
	}
	return m
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := m.replace(r.URL.Path)
	m.mut.RLock()
	entry := m.searchEntry(path, w, r)
	m.mut.RUnlock()
	if entry != nil {
		m.serveEntry(entry, w, r)
		return
	}
	if m.context.notFound != nil {
		m.context.notFound.ServeHTTP(w, r)
		return
	}
	http.Error(w, "404 Not Found : "+r.URL.String(), http.StatusNotFound)
}

func (m *Mux) searchEntry(path string, w http.ResponseWriter, r *http.Request) *Entry {
	if entry := m.getHandlerFunc(path); entry != nil {
		return entry
	}
	for _, groupMux := range m.groups {
		if entry := groupMux.searchEntry(path, w, r); entry != nil {
			return entry
		}
	}
	return nil
}

func (m *Mux) serveEntry(entry *Entry, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && entry.handlers[get] != nil {
		m.serveHandler(entry.handlers[get], w, r)
	} else if r.Method == "POST" && entry.handlers[post] != nil {
		m.serveHandler(entry.handlers[post], w, r)
	} else if r.Method == "PUT" && entry.handlers[put] != nil {
		m.serveHandler(entry.handlers[put], w, r)
	} else if r.Method == "DELETE" && entry.handlers[delete] != nil {
		m.serveHandler(entry.handlers[delete], w, r)
	} else if r.Method == "PATCH" && entry.handlers[patch] != nil {
		m.serveHandler(entry.handlers[patch], w, r)
	} else if r.Method == "HEAD" && entry.handlers[head] != nil {
		m.serveHandler(entry.handlers[head], w, r)
	} else if r.Method == "OPTIONS" && entry.handlers[options] != nil {
		m.serveHandler(entry.handlers[options], w, r)
	} else if r.Method == "TRACE" && entry.handlers[trace] != nil {
		m.serveHandler(entry.handlers[trace], w, r)
	} else if r.Method == "CONNECT" && entry.handlers[connect] != nil {
		m.serveHandler(entry.handlers[connect], w, r)
	} else {
		m.serveHandler(entry.handler, w, r)
	}
}

// Recovery returns a recovery handler function that recovers from any panics and writes a 500 status code.
func Recovery(w http.ResponseWriter, r *http.Request) {
	err := r.Context().Value(RecoveryContextKey)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "500 Internal Server Error : %v\n", err)
}

func (m *Mux) serveHandler(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	if m.context.recovery != nil {
		defer func() {
			if err := recover(); err != nil {
				ctx := context.WithValue(r.Context(), RecoveryContextKey, err)
				m.context.recovery.ServeHTTP(w, r.WithContext(ctx))
			}
		}()
	}
	m.middleware(w, r)
	if handler != nil {
		handler.ServeHTTP(w, r)
	}
}

func (m *Mux) getHandlerFunc(path string) *Entry {
	if prefix, key, ok := m.matchParams(path); ok {
		if entry, ok := m.prefixes[prefix].m[key]; ok {
			return entry
		}
	}
	return nil
}

// HandleFunc registers a handler function with the given pattern to the Mux.
func (m *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Entry {
	return m.Handle(pattern, http.HandlerFunc(handler))
}

// Handle registers a handler with the given pattern to the Mux.
func (m *Mux) Handle(pattern string, handler http.Handler) *Entry {
	m.mut.Lock()
	defer m.mut.Unlock()
	pattern = m.replace(pattern)
	pre, key, match, params := m.parseParams(m.group + pattern)
	if v, ok := m.prefixes[pre]; ok {
		if entry, ok := v.m[key]; ok {
			entry.handler = handler
			entry.key = key
			entry.match = match
			entry.params = params
			m.prefixes[pre].m[key] = entry
			return entry
		}
		entry := &Entry{}
		entry.handler = handler
		entry.key = key
		entry.match = match
		entry.params = params
		m.prefixes[pre].m[key] = entry
		return entry
	}
	m.prefixes[pre] = &prefix{m: make(map[string]*Entry), prefix: pre}
	entry := &Entry{}
	entry.handler = handler
	entry.key = key
	entry.match = match
	entry.params = params
	m.prefixes[pre].m[key] = entry
	return entry
}

// Group registers a group with the given pattern to the Mux.
func (m *Mux) Group(group string, f func(m *Mux)) {
	m.mut.Lock()
	defer m.mut.Unlock()
	group = m.replace(group)
	groupMux := newGroup(group)
	f(groupMux)
	if _, ok := m.groups[group]; ok {
		panic(ErrGroupExisted)
	}
	groupMux.context = m.context
	m.groups[group] = groupMux
}

// NotFound registers a not found handler function to the Mux.
func (m *Mux) NotFound(handler http.HandlerFunc) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.context.notFound = handler
}

// Recovery registers a recovery handler function to the Mux.
func (m *Mux) Recovery(handler http.HandlerFunc) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.context.recovery = handler
}

// Use uses middleware.
func (m *Mux) Use(handler http.HandlerFunc) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.context.middlewares = append(m.context.middlewares, handler)
}

func (m *Mux) middleware(w http.ResponseWriter, r *http.Request) {
	for _, handler := range m.context.middlewares {
		handler.ServeHTTP(w, r)
	}
}

// Params returns http request params.
func (m *Mux) Params(r *http.Request) map[string]string {
	params := make(map[string]string)
	path := m.replace(r.URL.Path)
	m.mut.RLock()
	if prefix, key, ok := m.matchParams(path); ok {
		if entry, ok := m.prefixes[prefix].m[key]; ok &&
			len(entry.match) > 0 && len(path) > len(prefix) {
			strs := strings.Split(path[len(prefix):], "/")
			if len(strs) == len(entry.match) {
				for i := 0; i < len(strs); i++ {
					if entry.match[i] != "" {
						params[entry.match[i]] = strs[i]
					}
				}
			}
		}
	}
	m.mut.RUnlock()
	return params
}

func (m *Mux) matchParams(path string) (string, string, bool) {
	for _, p := range m.prefixes {
		if strings.HasPrefix(path, p.prefix) {
			r := path[len(p.prefix):]
			if r == "" {
				return p.prefix, "", true
			}
			for _, v := range p.m {
				count := strings.Count(r, "/")
				if count+1 == len(v.match) {
					form := strings.Split(r, "/")
					key := ""
					for i := 0; i < len(form); i++ {
						if v.match[i] != "" {
							if i > 0 {
								key += "/:"
							} else {
								key += ":"
							}
						} else {
							key += "/" + form[i]
						}
					}
					if key == v.key {
						return p.prefix, v.key, true
					}
				}
			}
		}
	}
	return "", "", false
}

func (m *Mux) parseParams(pattern string) (string, string, []string, map[string]string) {
	prefix := ""
	var match []string
	key := ""
	params := make(map[string]string)
	if strings.Contains(pattern, ":") {
		idx := strings.Index(pattern, ":")
		prefix = pattern[:idx]
		if idx+1 == len(pattern) || strings.Contains(pattern, ":/") {
			panic(ErrParamsKeyEmpty)
		}
		match = strings.Split(pattern[idx:], "/")
		for i := 0; i < len(match); i++ {
			if strings.Contains(match[i], ":") {
				match[i] = strings.Trim(match[i], ":")
				params[match[i]] = ""
				if i > 0 {
					key += "/:"
				} else {
					key += ":"
				}
			} else {
				key += "/" + match[i]
				match[i] = ""
			}
		}
	} else {
		prefix = pattern
	}
	return prefix, key, match, params
}

func (m *Mux) replace(s string) string {
	for strings.Contains(s, "//") {
		s = strings.ReplaceAll(s, "//", "/")
	}
	return s
}

// GET adds a GET HTTP method to the entry.
func (entry *Entry) GET() *Entry {
	entry.handlers[get] = entry.handler
	return entry
}

// POST adds a POST HTTP method to the entry.
func (entry *Entry) POST() *Entry {
	entry.handlers[post] = entry.handler
	return entry
}

// PUT adds a PUT HTTP method to the entry.
func (entry *Entry) PUT() *Entry {
	entry.handlers[put] = entry.handler
	return entry
}

// DELETE adds a DELETE HTTP method to the entry.
func (entry *Entry) DELETE() *Entry {
	entry.handlers[delete] = entry.handler
	return entry
}

// PATCH adds a PATCH HTTP method to the entry.
func (entry *Entry) PATCH() *Entry {
	entry.handlers[patch] = entry.handler
	return entry
}

// HEAD adds a HEAD HTTP method to the entry.
func (entry *Entry) HEAD() *Entry {
	entry.handlers[head] = entry.handler
	return entry
}

// OPTIONS adds a OPTIONS HTTP method to the entry.
func (entry *Entry) OPTIONS() *Entry {
	entry.handlers[options] = entry.handler
	return entry
}

// TRACE adds a TRACE HTTP method to the entry.
func (entry *Entry) TRACE() *Entry {
	entry.handlers[trace] = entry.handler
	return entry
}

// CONNECT adds a CONNECT HTTP method to the entry.
func (entry *Entry) CONNECT() *Entry {
	entry.handlers[connect] = entry.handler
	return entry
}

// All adds all HTTP method to the entry.
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
}

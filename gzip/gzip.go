package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"
	"bytes"
)

type GzipWriter struct {
	Writer *gzip.Writer
	w http.ResponseWriter
	gzip bool
	buf *bytes.Buffer
}
func NewGzipWriter(w http.ResponseWriter, r *http.Request)*GzipWriter  {
	g:=&GzipWriter{
		Writer:gzip.NewWriter(w),
		w:w,
	}
	g.ready(w,r)
	return g
}
func (g *GzipWriter)ready(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		g.gzip=false
		return
	}
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Del("Content-Length")
	g.gzip=true
}
func (g *GzipWriter) Write(b []byte) (int, error) {
	if g.gzip{
		if g.w.Header().Get("Content-Type")==""{
			g.w.Header().Set("Content-Type", http.DetectContentType(b))
		}
		return g.Writer.Write(b)
	}else {
		return g.w.Write(b)
	}
}
func (g *GzipWriter) Close() (error) {
	if g.gzip{
		return g.Writer.Close()
	}
	return nil
}

func WriteGzip(w http.ResponseWriter, r *http.Request, httpStatus int, b []byte) (err error) {
	w.WriteHeader(httpStatus)
	gz:=NewGzipWriter(w,r)
	gz.Write(b)
	return gz.Close()
}
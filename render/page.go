package render
import (
	"net/http"
	"strconv"
	"io/ioutil"
)
func RedirectUrl(w http.ResponseWriter, r *http.Request,url string) {
	http.Redirect(w, r, url, http.StatusFound)
}
func WriteBody(w http.ResponseWriter, r *http.Request,httpStatus int,Body []byte) bool {
	if _, ok := w.Header()["Content-Length"]; ok {
		w.Header().Set("Content-Length", strconv.Itoa(len(Body)))
	} else {
		w.Header().Add("Content-Length", strconv.Itoa(len(Body)))
	}
	if _, ok := w.Header()["Content-Type"]; ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
	}
	w.WriteHeader(httpStatus)
	w.Write(Body)
	return true
}
func WritePage(w http.ResponseWriter, r *http.Request,status int,f string) bool {
	var Body, err = LoadFile(f)
	if err!=nil{
		return false
	}
	if _, ok := w.Header()["Content-Length"]; ok {
		w.Header().Set("Content-Length", strconv.Itoa(len(Body)))
	} else {
		w.Header().Add("Content-Length", strconv.Itoa(len(Body)))
	}
	if _, ok := w.Header()["Content-Type"]; ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
	}
	w.WriteHeader(status)
	w.Write(Body)
	return true
}
func LoadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

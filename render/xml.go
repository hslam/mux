package render
import (
	"net/http"
	"encoding/xml"
)

func WriteXml(w http.ResponseWriter, r *http.Request, status int, obj interface{}) (err error) {
	var bytes []byte
	if r.FormValue("xml") != "" {
		bytes, err = xml.MarshalIndent(obj, "", "  ")
	} else {
		bytes, err = xml.Marshal(obj)
	}
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	return
}

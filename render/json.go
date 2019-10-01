package render
import (
	"net/http"
	"encoding/json"
)

func WriteJson(w http.ResponseWriter, r *http.Request, status int, obj interface{}) (err error) {
	var bytes []byte
	if r.FormValue("json") != "" ||r.FormValue("pretty") != ""{
		bytes, err = json.MarshalIndent(obj, "", "  ")
	} else {
		bytes, err = json.Marshal(obj)
	}
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	return
}

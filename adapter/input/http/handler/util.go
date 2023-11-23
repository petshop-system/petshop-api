package handler

import "net/http"

func responseReturn(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(body) != 0 {
		w.Write(body)
	}
}

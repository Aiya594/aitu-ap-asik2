package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if v != nil {
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			http.Error(w, `{"error":"failed to encode JSON"}`, http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

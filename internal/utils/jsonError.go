package utils

import (
	"encoding/json"
	"net/http"
)

func JsonError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(map[string]interface{}{
		"error": msg,
	})
	w.Write(resp)
}

package utils

import (
	"encoding/json"
	"net/http"
)

func CheckRequestMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, "Only method"+method+"allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		http.Error(w, "Invalid JSON paylod", http.StatusBadRequest)
		return false
	}
	return true
}

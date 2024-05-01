package common

import (
	"encoding/json"
	"net/http"
)

func GenerateResponse(data interface{}, err string) []byte {
	jsonBytes, jsonErr := json.Marshal(ResponseData{Data: data, Error: err})

	if jsonErr != nil {
		bytes := []byte(`{"error": "Failed to generate response"}`)
		return bytes
	}

	return jsonBytes
}

func RaiseBadRequest(w http.ResponseWriter, err string) {
	res := GenerateResponse(nil, err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}
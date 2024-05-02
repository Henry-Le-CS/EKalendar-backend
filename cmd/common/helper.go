package common

import (
	"encoding/json"
	"net/http"
	"time"
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

func ParseDate(date string, layout string) (time.Time, error) {
	// Issue: https://stackoverflow.com/questions/43456851/how-to-format-todays-date-in-go-as-dd-mm-yyyy
	if layout == "dd/MM/yyyy" {
		layout = "02/01/2006"
	} else if layout == "MM/dd/yyyy" {
		layout = "01/02/2006"
	}

    return time.Parse(layout, date)
}
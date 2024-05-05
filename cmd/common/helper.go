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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}

var COUNTRY_TZ = map[string]string{
	"Vietnam": "Asia/Ho_Chi_Minh",
}

func TimeIn(name string) (*time.Location, error){
    return time.LoadLocation(COUNTRY_TZ[name])
}


func ParseDate(date string, layout string) (time.Time, error) {
	// Issue: https://stackoverflow.com/questions/43456851/how-to-format-todays-date-in-go-as-dd-mm-yyyy
	if layout == "dd/MM/yyyy" {
		layout = "02/01/2006"
	} else if layout == "MM/dd/yyyy" {
		layout = "01/02/2006"
	}

    t, err := time.Parse(layout, date)

	if err != nil {
		return time.Time{}, err
	}

	var loc *time.Location

	if loc, err = TimeIn("Vietnam"); err != nil {
		return time.Time{}, err
	}

	return t.In(loc), nil
}
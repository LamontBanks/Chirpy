package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, msg string, statusCode int, errorToLog error) {
	if errorToLog != nil {
		log.Printf("%v", errorToLog)
	}

	type errorResponseJSON struct {
		Error string `json:"error"`
	}
	errorResp := errorResponseJSON{
		Error: msg,
	}

	sendJSONResponse(w, statusCode, errorResp)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, jsonStruct interface{}) {
	data, err := json.Marshal(jsonStruct)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendErrorJSONResponse(w http.ResponseWriter, msg string, statusCode int, errorToLog error) {
	if errorToLog != nil {
		log.Printf("%v", errorToLog)
	}

	errorResp := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	SendJSONResponse(w, statusCode, errorResp)
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, jsonStruct any) {
	data, err := json.Marshal(jsonStruct)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func sendResponse(w http.ResponseWriter, statusCode int, msgToLog string) {
	if msgToLog != "" {
		log.Printf("%v", msgToLog)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
}

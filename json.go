package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("responding with 5xx error: ", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	responseWithJson(w, code, errResponse{
		Error: msg,
	})
}

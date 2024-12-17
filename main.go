package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {

	http.HandleFunc("/text", text)

	log.Println("Server started :3000")
	http.ListenAndServe(":3000", nil)
}

func jsonResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
	})
}

func text(w http.ResponseWriter, r *http.Request) {
	content_type := r.Header.Get("Content-Type")

	if content_type != "application/json" {
		jsonResponse(w, http.StatusNotAcceptable, "JSON not available")
		return
	}

	jsonResponse(w, http.StatusOK, "Request Accepted")
}

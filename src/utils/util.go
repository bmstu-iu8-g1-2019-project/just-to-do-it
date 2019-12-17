package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string, error string) map[string]interface{} {
	return map[string]interface{} {"status" : status, "message" : message, "err" : error}
}

func Respond(w http.ResponseWriter, data map[string] interface{})  {
	w.Header().Add("Content-Type", "application/json")
	if data["err"] == "Unauthorized" {
		w.WriteHeader(http.StatusUnauthorized)
	} else if data["err"] == "Bad Request" {
		w.WriteHeader(http.StatusBadRequest)
	} else if data["err"] == "Internal Server Error" {
		w.WriteHeader(http.StatusInternalServerError)
	} else if data["err"] == "Forbidden"{
		w.WriteHeader(http.StatusForbidden)
	} else if data["err"] == "Not Found "{
		w.WriteHeader(http.StatusNotFound)
	} else if data["status"] == true {
		w.WriteHeader(http.StatusOK)
	}
	delete(data, "err")
	_ = json.NewEncoder(w).Encode(data)
}

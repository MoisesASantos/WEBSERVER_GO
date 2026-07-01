package config

import (
	"fmt"
	"encoding/json"
	"net/http"
)
func RespondWithError(w http.ResponseWriter, code int, msg string) error {

	type returnVals struct {
        Error string `json:"error"`
    }
    
	respBody := returnVals{
        Error: msg,
    }

    data, err := json.Marshal(respBody)
	
	if err != nil {
			fmt.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return err
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(data)
	return nil
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) error {

    data, err := json.Marshal(payload)
	
	if err != nil {
			fmt.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return err
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(data)
	return nil
}

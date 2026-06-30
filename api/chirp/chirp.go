package chirp


import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type returnVals struct {
        Error string `json:"error"`
    }
    
	respBody := returnVals{
        Error: msg,
    }

    dat, err := json.Marshal(respBody)
	
	if err != nil {
			fmt.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, value string) {
	
	type returnVals struct {
        Cleaned_body string `json:"cleaned_body"`
    }
    
	respBody := returnVals{
        Cleaned_body: value,
    }

    dat, err := json.Marshal(respBody)
	
	if err != nil {
			fmt.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

func ChirpRequestHandler(w http.ResponseWriter, r *http.Request) {
	
	type parameters struct {
        Body string `json:"body"`
    }

	mapWord := map[string]int{
		"kerfuffle":  1,
		"sharbert": 1,
		"fornax": 1,
	}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		fmt.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	if len(params.Body) <= 140 {
		splited := strings.Split(params.Body, " ")
		sliceResult := []string{}

		for _, value := range splited {
			_, ok := mapWord[strings.ToLower(value)]
			if ok {
				sliceResult = append(sliceResult, "****")
			} else {
				sliceResult = append(sliceResult, value)
			}
		}
		respondWithJSON(w, 200, strings.Join(sliceResult, " "))
	} else {
		respondWithError(w, 400, "Chirp is too long")
	}
}

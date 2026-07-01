package chirp


import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"github.com/MoisesASantos/WEBSERVER_GO/admin/config"
)

type returnVals struct {
        	Cleaned_body string `json:"cleaned_body"`
}

type parameters struct {
        Body string `json:"body"`
}


func ChirpRequestHandler(w http.ResponseWriter, r *http.Request) {
	
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
		return err
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
		
    
		respBody := returnVals{
        	Cleaned_body: strings.Join(sliceResult, " "),
    	}
		err = config.RespondWithJSON(w, 200, returnVals)
	} else {
		err =config.RespondWithError(w, 400, "Chirp is too long")
	}

	if err != nil {
		fmt.Printf("Error creating the response: %s", err)
		w.WriteHeader(500)
		return err
	}
	return nil;
}

package config


import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
	//"encoding/json"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"

)

type returnChirp struct {
    ID        uuid.UUID	`json:"id"`
	CreatedAt time.Time	`json:"created_id"`
	UpdatedAt time.Time	`json:"updated_id"`
	Body      string	`json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type parameters struct {
        Body	string `json:"body"`
        UserID	uuid.UUID `json:"user_id"`
}

func parsingString(entry string) string {

	//map the forbiden words
	mapWord := map[string]int{
		"kerfuffle":  1,
		"sharbert": 1,
		"fornax": 1,
	}

	splited := strings.Split(entry, " ")
	sliceResult := []string{}

	for _, value := range splited {
		_, ok := mapWord[strings.ToLower(value)]
		if ok {
			sliceResult = append(sliceResult, "****")
		} else {
			sliceResult = append(sliceResult, value)
		}
	}
	return strings.Join(sliceResult, " ")
}

func (cfg *ApiConfig) ChirpRequestHandler(w http.ResponseWriter, r *http.Request) {
	
	//Decode the json
	data, err := DecodeJSON[parameters](r)
	
	if len(data.Body) > 140 {
		err = RespondWithError(w, 400, "Chirp is too long")
		return 
	}
	//Check the jwt token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
    	w.WriteHeader(http.StatusUnauthorized)
    	fmt.Printf("error getting token: %s\n", err)
    	return
	}

	userID, err := auth.ValidateJWT(token, cfg.JwtKey)
	if err != nil {
	    w.WriteHeader(http.StatusUnauthorized)
	    fmt.Printf("invalid token: %s\n", err)
	    return
	}

	//parsing the string
	value := parsingString(data.Body)
	//create a chirp
	chirpBody := database.CreateChirpParams{Body: value, UserID: userID,}
	chirpResult, err := cfg.Db.CreateChirp(r.Context(), chirpBody)
	if err != nil {
		fmt.Printf("Error creating the response: %s", err)
		w.WriteHeader(500)
		return
	}
	
	//response body
	respBody := returnChirp{
		ID:        chirpResult.ID,
		CreatedAt: chirpResult.CreatedAt,
		UpdatedAt: chirpResult.UpdatedAt,
		Body:      chirpResult.Body,
		UserID:    chirpResult.UserID,
	}
	err = RespondWithJSON(w, 201, respBody)
	if err != nil {
		fmt.Printf("Error creating the response: %s", err)
		w.WriteHeader(500)
		return
	}
}

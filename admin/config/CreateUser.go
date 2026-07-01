package config

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"github.com/google/uuid"
)

type returnVals struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
}

type requestbody struct {
		Email string `json:"email"`
}

func (cfg *ApiConfig) UsersRequestHandler(w http.ResponseWriter, r *http.Request) {

	//Get the request body
	data, err := DecodeJSON[requestBody](r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}

	//create the user using a function created per sqlc
	userResult, err := cfg.Db.CreateUser(r.Context(), data.Email)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	//now create a response
	payload := returnVals{
		ID: 		userResult.ID,
		CreatedAt: 	userResult.CreatedAt,
		UpdatedAt: 	userResult.UpdatedAt,
		Email:		userResult.Email,
	}
	err = RespondWithJSON(w, 201, payload)
	if err != nil {
		return
	}
}

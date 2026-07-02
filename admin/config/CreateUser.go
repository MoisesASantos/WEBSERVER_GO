package config

import (
	"net/http"
	"time"
	"fmt"
	"github.com/google/uuid"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
)

type returnVals struct {
    ID uuid.UUID 		`json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string 		`json:"email"`
	Token string		`json:"token"`
}

type requestbody struct {
	Password 	string `json:"password"`
	Email 		string `json:"email"`
}

func (cfg *ApiConfig) UsersRequestHandler(w http.ResponseWriter, r *http.Request) {

	//Get the request body
	data, err := DecodeJSON[requestbody](r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}

	//create hash for password
	hash, err := auth.HashPassword(data.Password)
	//create the user using a function created per sqlc
	userparams := database.CreateUserParams{
		Email:          data.Email,
		HashedPassword: hash,
	}
	userResult, err := cfg.Db.CreateUser(r.Context(), userparams)
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

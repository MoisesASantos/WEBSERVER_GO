package config

import (
	"fmt"
	"net/http"
	"time"
	//"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
)

type requestbodyLogin struct {
	Password 			string `json:"password"`
	Email 				string `json:"email"`
	ExpiresInSeconds	int		`json:"expires_in_seconds"`
}

func (cfg *ApiConfig) LoginRequestHandler(w http.ResponseWriter, r *http.Request) {

	//Get the request body
	data, err := DecodeJSON[requestbodyLogin](r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}

	//get user data
	user, err := cfg.Db.GetUser(r.Context(), data.Email)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}

	//check the hash
	result, err := auth.CheckPasswordHash(data.Password, user.HashedPassword)
	if result == false {
		err = RespondWithError(w, 401, "Incorrect email or password")
		return 
	}
	
	//create the token
	expires := time.Hour // padrão

	if data.ExpiresInSeconds > 0 {
    	expires = time.Duration(data.ExpiresInSeconds) * time.Second
    	if expires > time.Hour {
        	expires = time.Hour
    	}
	}

	token, err := auth.MakeJWT(user.ID, cfg.JwtKey, expires)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}
	// return the response
	payload := returnVals{
		ID: 		user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
		Token:		token,
	}
	err = RespondWithJSON(w, 200, payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

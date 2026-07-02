package config

import (
	"fmt"
	"net/http"
	//"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
)

func (cfg *ApiConfig) LoginRequestHandler(w http.ResponseWriter, r *http.Request) {

	//Get the request body
	data, err := DecodeJSON[requestbody](r)
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

	// return the response
	payload := returnVals{
		ID: 		user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
	}
	err = RespondWithJSON(w, 200, payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

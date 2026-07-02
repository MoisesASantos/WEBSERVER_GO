package config

import (
	"fmt"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
)

type returnValsLogin struct {
    ID				uuid.UUID 	`json:"id"`
    CreatedAt 		time.Time	`json:"created_at"`
    UpdatedAt		time.Time	`json:"updated_at"`
    Email			string 		`json:"email"`
	Token			string		`json:"token"`
	RefreshToken	string		`json:"refresh_token"`
}

type requestbodyLogin struct {
	Password 			string `json:"password"`
	Email 				string `json:"email"`
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
	token, err := auth.MakeJWT(user.ID, cfg.JwtKey, expires)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}
	
	//create a refresh token
	refresh_token := auth.MakeRefreshToken()
	
	paramsRefreshToken := database.CreateRefreshTokenParams {
		Token:		refresh_token,
		ExpiresAt:	time.Now().Add(60 * 24 * time.Hour),
		UserID:		user.ID,
	}

	_, err = cfg.Db.CreateRefreshToken(r.Context(), paramsRefreshToken)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}
	
	// return the response
	payload := returnValsLogin{
		ID: 			user.ID,
		CreatedAt: 		user.CreatedAt,
		UpdatedAt: 		user.UpdatedAt,
		Email:			user.Email,
		Token:			token,
		RefreshToken:	refresh_token,
	}

	err = RespondWithJSON(w, 200, payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

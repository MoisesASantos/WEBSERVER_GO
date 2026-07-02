package config

import (
	"fmt"
	"net/http"
	//"time"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
)

type upRequestbody struct {
	Password 	string `json:"password"`
	Email 		string `json:"email"`
}

func (cfg *ApiConfig) UpdatedUserRequestHandler(w http.ResponseWriter, r *http.Request) {

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

	data, err := DecodeJSON[upRequestbody](r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}
	hash, err := auth.HashPassword(data.Password)

	params := database.UpdateUserParams {
		Email:          data.Email,
		HashedPassword: hash,
		ID:             userID,
	}
	err = cfg.Db.UpdateUser(r.Context(), params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}

	user, err := cfg.Db.GetUser(r.Context(), data.Email)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}

	payload := returnVals{
		ID: 		user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
	}
	err = RespondWithJSON(w, 200, payload)
	if err != nil {
		return
	}
}

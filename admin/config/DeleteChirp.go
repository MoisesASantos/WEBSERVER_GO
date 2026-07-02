package config

import (
	"fmt"
	"net/http"
	//"time"
	"github.com/google/uuid"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
	//"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
)


func (cfg *ApiConfig) DeleteChirpRequestHandler(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
    	w.WriteHeader(http.StatusUnauthorized)
    	fmt.Printf("error getting token: %s\n", err)
    	return
	}

	UserID, err := auth.ValidateJWT(token, cfg.JwtKey)
	if err != nil {
	    w.WriteHeader(403)
	    fmt.Printf("invalid token: %s\n", err)
	    return
	}

	id, err := uuid.Parse(chirpID)
	if err != nil {
    	http.Error(w, "invalid chirp id", http.StatusBadRequest)
    	return
	}

	chirpResult, err := cfg.Db.GetOneChirp(r.Context(), id)
	if err != nil {
		fmt.Printf("Error getting the chirp: %s", err)
		w.WriteHeader(404)
		return
	}

	if chirpResult.UserID != UserID {
		w.WriteHeader(403)
	    fmt.Printf("invalid token: %s\n", err)
	    return
	}

	err = cfg.Db.DeleteChirp(r.Context(), id)
	if err != nil {
	    fmt.Printf("Error delete: %s\n", err)
	    return
	}
	w.WriteHeader(204)
}

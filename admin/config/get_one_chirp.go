package config

import (
	"fmt"
	"net/http"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) GetOneChirpRequestHandler(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")

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

	respBody := returnChirp{
		ID:        chirpResult.ID,
		CreatedAt: chirpResult.CreatedAt,
		UpdatedAt: chirpResult.UpdatedAt,
		Body:      chirpResult.Body,
		UserID:    chirpResult.UserID,
	}
	err = RespondWithJSON(w, 200, respBody)
	if err != nil {
		fmt.Printf("Error creating the response: %s", err)
		w.WriteHeader(500)
		return
	}
}

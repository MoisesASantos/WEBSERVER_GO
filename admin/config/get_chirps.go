package config

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) ChirpsGetRequestHandler(w http.ResponseWriter, r *http.Request) {
	

	chirpsResult, err := cfg.Db.GetChirps(r.Context())
	if err != nil {
		fmt.Printf("Error creating the response: %s", err)
		w.WriteHeader(500)
		return
	}

	respBody := make([]returnChirp , 0, len(chirpsResult))

	for index, _ := range chirpsResult {
		respBody = append(respBody, returnChirp{
		ID:        chirpsResult[index].ID,
		CreatedAt: chirpsResult[index].CreatedAt,
		UpdatedAt: chirpsResult[index].UpdatedAt,
		Body:      chirpsResult[index].Body,
		UserID:    chirpsResult[index].UserID,
		})
	}
	err = RespondWithJSON(w, 200, respBody)
	if err != nil {
		fmt.Printf("Error creating the response: %s", err)
		w.WriteHeader(500)
		return
	}
}


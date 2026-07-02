package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
)

func (cfg *ApiConfig) RefreshRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Extrai o refresh token do header Authorization
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("error getting bearer token: %v\n", err)
		return
	}

	// Procura o refresh token na base de dados
	refreshToken, err := cfg.Db.GetRefreshToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("error retrieving refresh token: %v\n", err)
		return
	}

	// Verifica se foi revogado
	if refreshToken.RevokedAt.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verifica se expirou
	if refreshToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Cria um novo access token com duração de 1 hora
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.JwtKey, time.Hour)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("error creating access token: %v\n", err)
		return
	}

	// Resposta
	type response struct {
		Token string `json:"token"`
	}

	err = RespondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
	if err != nil {
		fmt.Printf("error writing response: %v\n", err)
	}
}

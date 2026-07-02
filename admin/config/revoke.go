package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MoisesASantos/WEBSERVER_GO/internal/auth"
)

func (cfg *ApiConfig) RevokeRequestHandler(w http.ResponseWriter, r *http.Request) {
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

	// Verifica se expirou
	if refreshToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verifica se foi revogado
	if refreshToken.RevokedAt.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} else {
		err := cfg.Db.RevokeRefreshToken(r.Context(), token)
		if err != nil {
			fmt.Printf("error revoking the refresh token: %v\n", err)
			return
		}
	}

	// Resposta
	w.WriteHeader(204)
}

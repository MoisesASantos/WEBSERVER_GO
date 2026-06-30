package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) writeHitRequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body := fmt.Sprintf(`
	<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
	    <p>Chirpy has been visited %d times!</p>
	  </body>
	</html>
	`, cfg.fileserverHits.Load())

	w.Write([]byte(body))
}

func (cfg *apiConfig) resetHitRequest(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}


func main() {

	apiconfig := apiConfig{}
	mux := http.NewServeMux()

	const filepathRoot = "."
	mux.Handle("/app/", apiconfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiconfig.writeHitRequest)
	mux.HandleFunc("POST /admin/reset", apiconfig.resetHitRequest)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error on server: %v", err)
		os.Exit(1)
	}
}

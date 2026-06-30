package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"encoding/json"
	"strings"
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


func respondWithError(w http.ResponseWriter, code int, msg string) {

	type returnVals struct {
        Error string `json:"error"`
    }
    
	respBody := returnVals{
        Error: msg,
    }

    dat, err := json.Marshal(respBody)
	
	if err != nil {
			fmt.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, value string) {
	
	type returnVals struct {
        Cleaned_body string `json:"cleaned_body"`
    }
    
	respBody := returnVals{
        Cleaned_body: value,
    }

    dat, err := json.Marshal(respBody)
	
	if err != nil {
			fmt.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

func chirpRequestHandler(w http.ResponseWriter, r *http.Request) {
	
	type parameters struct {
        Body string `json:"body"`
    }

	mapWord := map[string]int{
		"kerfuffle":  1,
		"sharbert": 1,
		"fornax": 1,
	}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		fmt.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	if len(params.Body) <= 140 {
		splited := strings.Split(params.Body, " ")
		sliceResult := []string{}

		for _, value := range splited {
			_, ok := mapWord[strings.ToLower(value)]
			if ok {
				sliceResult = append(sliceResult, "****")
			} else {
				sliceResult = append(sliceResult, value)
			}
		}
		respondWithJSON(w, 200, strings.Join(sliceResult, " "))
	} else {
		respondWithError(w, 400, "Chirp is too long")
	}
}

func main() {

	apiconfig := apiConfig{}
	mux := http.NewServeMux()

	const filepathRoot = "."
	mux.Handle("/app/", apiconfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiconfig.writeHitRequest)
	mux.HandleFunc("POST /admin/reset", apiconfig.resetHitRequest)
	mux.HandleFunc("POST /api/validate_chirp", chirpRequestHandler)

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

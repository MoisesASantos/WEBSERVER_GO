package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {


	mux := http.NewServeMux()

	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	defer server.Close()
	if err != nil {
		fmt.Printf("Error on server: %v", err)
		os.Exit(1)
	}
}

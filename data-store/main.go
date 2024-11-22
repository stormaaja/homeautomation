package main

import (
	"fmt"
	"log"
	"net/http"
	"stormaaja/go-ha/security/configvalidators"
	"stormaaja/go-ha/security/requestvalidators"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if !requestvalidators.IsApiTokenValid(r.Header) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "%s!", r.URL.Path[1:])
}

func main() {
	if err := configvalidators.IsConfigEnvironmentVariablesValid(); err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

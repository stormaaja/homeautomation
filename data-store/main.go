package main

import (
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/dataroutes"
	"stormaaja/go-ha/data-store/genericroutes"
	"stormaaja/go-ha/data-store/routes"
	"stormaaja/go-ha/data-store/store"
	"stormaaja/go-ha/security/configvalidators"
	"stormaaja/go-ha/security/requestvalidators"
	"strings"
)

var temperatureRoute = dataroutes.TemperatureRoute{Store: &store.MemoryStore{Data: make(map[string]interface{})}}
var healthcheckRoute = genericroutes.HealthcheckRoute{}

func GetRootPath(path string) (string, string) {
	splittedPath := strings.Split(path, "/")
	if len(splittedPath) > 2 {
		return splittedPath[1], splittedPath[2]
	} else if len(splittedPath) > 1 {
		return splittedPath[1], ""
	}
	return "", ""
}

func HandleRoute(route routes.RouteHandler, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		route.HandleGet(w, r)
	case "POST":
		route.HandlePost(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if !requestvalidators.IsApiTokenValid(r.Header) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("%s %s", r.Method, r.URL.Path) // TODO: Clean ids from logs

	rootPath, subPath := GetRootPath(r.URL.Path)

	switch rootPath {
	case "data":
		if subPath == "temperature" {
			HandleRoute(temperatureRoute, w, r)
		} else {
			http.Error(w, "Invalid path", http.StatusBadRequest)
		}
	case "healthcheck":
		HandleRoute(healthcheckRoute, w, r)
	default:
		http.Error(w, "Invalid path", http.StatusBadRequest)
	}
}

func main() {
	log.Println("Starting server...")
	if err := configvalidators.IsConfigEnvironmentVariablesValid(); err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	http.HandleFunc("/", handler)
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

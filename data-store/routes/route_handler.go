package routes

import "net/http"

type RouteHandler interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandlePost(w http.ResponseWriter, r *http.Request)
}

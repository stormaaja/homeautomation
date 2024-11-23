package genericroutes

import "net/http"

type HealthcheckRoute struct {
}

func (h HealthcheckRoute) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h HealthcheckRoute) HandlePost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

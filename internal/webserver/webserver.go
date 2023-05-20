package webserver

import (
	"net/http"
)

func RegisterRoutes(m *http.ServeMux) {
	m.Handle("/", rootHandler())
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func twitchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do twitch stuff
		w.WriteHeader(http.StatusOK)
	}
}

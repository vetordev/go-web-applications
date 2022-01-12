package server

import (
	"net/http"
)

func ViewController(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	var handler http.HandlerFunc

	switch method {
	case http.MethodGet:
		handler = makeHandler(ViewService)
	default:
		http.NotFound(w, r)
	}

	handler(w, r)
}

func EditController(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	var handler http.HandlerFunc

	switch method {
	case http.MethodGet:
		handler = makeHandler(EditService)
	default:
		http.NotFound(w, r)
	}

	handler(w, r)
}

func SaveController(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	var handler http.HandlerFunc

	switch method {
	case http.MethodPost:
		handler = makeHandler(SaveService)
	default:
		http.NotFound(w, r)
	}

	handler(w, r)
}

func makeHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := GetTitle(r.URL.Path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		service(w, r, title)
	}
}

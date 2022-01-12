package server

import (
	"log"
	"net/http"
)

const (
	PORT = ":8080"
)

func routes() {
	http.HandleFunc("/view/", ViewController)
	http.HandleFunc("/edit/", EditController)
	http.HandleFunc("/save/", SaveController)
}

func Serve() {
	routes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	TextFileExtension = ".txt"
	HtmlFileExtension = ".html"
)

func makePageFileName(title string) string {

	return title + TextFileExtension
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := makePageFileName(p.Title)

	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := makePageFileName(title)
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(writer http.ResponseWriter, name string, page *Page) {
	t, _ := template.ParseFiles(name + HtmlFileExtension)
	t.Execute(writer, page)
}

func viewPageHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/view/"):]
	page, _ := loadPage(title)

	renderTemplate(writer, "view", page)
}

func editPageHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/edit/"):]
	page, err := loadPage(title)

	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(writer, "edit", page)
}

func main() {
	http.HandleFunc("/view/", viewPageHandler)
	http.HandleFunc("/edit/", editPageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	templates = template.Must(template.ParseFiles("template/edit.html", "template/view.html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

const (
	TextFileExtension = ".txt"
	HtmlFileExtension = ".html"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := getFilename(p.Title)

	return os.WriteFile(filename, p.Body, 0600)
}

type Handler func(writer http.ResponseWriter, request *http.Request, title string)

func getFilename(title string) string {
	return title + TextFileExtension
}

func getTitle(path string) (string, error) {
	m := validPath.FindStringSubmatch(path)

	if m == nil {
		return "", errors.New("invalid page title")
	}

	return m[2], nil
}

func loadPage(title string) (*Page, error) {
	filename := getFilename(title)
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(writer http.ResponseWriter, name string, page *Page) {
	err := templates.ExecuteTemplate(writer, name+HtmlFileExtension, page)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func viewPageHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(writer, "view", page)
}

func editPageHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(writer, "edit", page)
}

func savePageHandler(writer http.ResponseWriter, request *http.Request, title string) {
	body := request.FormValue("body")

	page := &Page{Title: title, Body: []byte(body)}
	err := page.save()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func makeHandler(fn Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		title, err := getTitle(request.URL.Path)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		fn(writer, request, title)
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewPageHandler))
	http.HandleFunc("/edit/", makeHandler(editPageHandler))
	http.HandleFunc("/save/", makeHandler(savePageHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

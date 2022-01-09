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
	templates = template.Must(template.ParseFiles("edit.html", "view.html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

const (
	TextFileExtension = ".txt"
	HtmlFileExtension = ".html"
)

func getFileName(title string) string {
	return title + TextFileExtension
}

func getTitle(writer http.ResponseWriter, request *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(request.URL.Path)

	if m == nil {
		http.NotFound(writer, request)
		return "", errors.New("Invalid Page Title")
	}

	return m[2], nil
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := getFileName(p.Title)

	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := getFileName(title)
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

func viewPageHandler(writer http.ResponseWriter, request *http.Request) {
	title, err := getTitle(writer, request)
	if err != nil {
		return
	}

	page, err := loadPage(title)

	if err != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(writer, "view", page)
}

func editPageHandler(writer http.ResponseWriter, request *http.Request) {
	title, err := getTitle(writer, request)
	if err != nil {
		return
	}

	page, err := loadPage(title)

	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	title, err := getTitle(writer, request)
	if err != nil {
		return
	}

	body := request.FormValue("body")

	page := &Page{Title: title, Body: []byte(body)}
	err = page.save()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewPageHandler)
	http.HandleFunc("/edit/", editPageHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package server

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

type Service = func(w http.ResponseWriter, r *http.Request, title string)

var (
	templates = template.Must(template.ParseFiles("./server/template/edit.html", "./server/template/view.html"))
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

func getFilename(title string) string {
	return title + TextFileExtension
}

func GetTitle(path string) (string, error) {
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

func renderTemplate(w http.ResponseWriter, name string, page *Page) {
	err := templates.ExecuteTemplate(w, name+HtmlFileExtension, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewService(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", page)
}

func EditService(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(w, "edit", page)
}

func SaveService(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")

	page := &Page{Title: title, Body: []byte(body)}
	err := page.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

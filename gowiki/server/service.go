package server

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

type Service = func(w http.ResponseWriter, r *http.Request, title string)

const (
	TextFileExtension = ".txt"
	HtmlFileExtension = ".html"
	DataStore         = "\\server\\data\\"
	TemplateStore     = "\\server\\template\\"
)

var (
	projectPath, _ = os.Getwd()
	templates      = template.Must(
		template.ParseFiles(projectPath+TemplateStore+"edit.html", projectPath+TemplateStore+"view.html"),
	)
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filePath := getFilePath(p.Title)
	return os.WriteFile(filePath, p.Body, 0600)
}

func getFilePath(title string) string {
	return projectPath + DataStore + title + TextFileExtension
}

func GetTitle(path string) (string, error) {
	m := validPath.FindStringSubmatch(path)

	if m == nil {
		return "", errors.New("invalid page title")
	}

	return m[2], nil
}

func loadPage(title string) (*Page, error) {
	filePath := getFilePath(title)
	body, err := os.ReadFile(filePath)

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

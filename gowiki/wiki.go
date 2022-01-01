package main

import (
	"fmt"
	"net/http"
	"os"
)

const TextFileExtension = ".txt"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := makeFileName(p.Title)

	return os.WriteFile(filename, p.Body, 0600)
}

func main() {
	page := &Page{Title: "Test Page", Body: []byte("This is a sample page")}
	page.save()

	loadedPage, _ := loadPage("Test Page")
	fmt.Println(string(loadedPage.Body))
}

func viewPageHandler(response http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/view/"):]
	page, _ := loadPage(title)

	fmt.Fprintf(response, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)
}

func loadPage(title string) (*Page, error) {
	filename := makeFileName(title)
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func makeFileName(title string) string {

	return title + TextFileExtension
}

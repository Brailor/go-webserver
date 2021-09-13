package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)
var templates = template.Must(template.ParseFiles("./templ/index.html","./templ/edit.html", "./templ/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body []byte
}

func (page *Page) save() error {
	filename := page.Title + ".txt"

	return ioutil.WriteFile("./pages/" + filename, page.Body, 0600)
}

func loadPage (title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("./pages/" + filename)

	if err != nil {
		fmt.Println("error: ", err)

		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(writer http.ResponseWriter, req *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		http.Redirect(writer, req, "/edit/" + title, http.StatusFound)
		return 
	}

	renderTemplate(writer, "view.html", page)
}

func editHandler(writer http.ResponseWriter, req *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(writer, "edit.html", page)
}

func saveHandler(writer http.ResponseWriter, req *http.Request, title string) {
	body := req.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	err := page.save()

	if err != nil{
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, req, "/view/" + title, http.StatusFound)
}

func renderTemplate(writer http.ResponseWriter, temp string, page *Page) {
	err := templates.ExecuteTemplate(writer, temp, page)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func handler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Hi there, I love %s", req.URL.Path[1:])
}
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		m := validPath.FindStringSubmatch(req.URL.Path)
		if m == nil {
			http.NotFound(writer, req)
	
			return
		}

		fn(writer, req, m[2])
	}
 }
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		page, _ := loadPage("FrontPage")
		renderTemplate(writer, "index.html", page)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
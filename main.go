package main

import (
	"net/http"
	"html/template"
	"fmt"
	//"github.com/google/go-github/github"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct { }{ }
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, &data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Query string
	}{
		Query: r.URL.Query().Get("query"),
	}

	t, err := template.ParseFiles("templates/search.html")
	if err != nil {
		fmt.Print(err)
	}

	t.Execute(w, &data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", nil)
}

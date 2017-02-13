package main

import (
	"net/http"
	"html/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct { }{ }
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, &data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := struct { }{ }
	t, _ := template.ParseFiles("templates/search.html")
	t.Execute(w, &data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", nil)
}

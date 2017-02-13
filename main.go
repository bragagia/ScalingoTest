package main

import (
	"net/http"
	"html/template"
	"fmt"
	"github.com/google/go-github/github"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct { }{ }
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, &data)
}

func GetList(q string) {
	client := github.NewClient(nil)
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	orgs, _, err := client.Search.Repositories("created:>=2017-02-13T23:25:00+01:00", opt)

	if err != nil {
		fmt.Print(err)
	} else {
		for i, _ := range orgs.Repositories {
			fmt.Print("\n")
			fmt.Print(i)
		}
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Query string
	}{
		Query: r.URL.Query().Get("query"),
	}

	GetList(data.Query)

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

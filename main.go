package main

import (
	"net/http"
	"html/template"
	"fmt"
	"github.com/google/go-github/github"
)

type RepoInfos struct {
	FullName string
	Languages map[string]int
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func GetRepoInfos(r github.Repository, ghc *github.Client, c chan RepoInfos) {
	var ri RepoInfos
	ri.FullName = *r.FullName

	languages, _, err := ghc.Repositories.ListLanguages(*r.Owner.Login, *r.Name)
	if err != nil {
		//fmt.Print(err)
	} else {
		ri.Languages = languages
	}
	c <- ri
}

func GetList(q string) (*[100]RepoInfos) {
	client := github.NewClient(nil)
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	orgs, _, err := client.Search.Repositories("created:>=2017-02-13T23:25:00+01:00", opt)

	if err != nil {
		fmt.Print(err, "\n")
		return nil
	}

	c := make(chan RepoInfos)

	i := 0
	for _, r := range orgs.Repositories {
		go GetRepoInfos(r, client, c)
		i++
	}

	var res [100]RepoInfos
	for e := i; e > 0; e-- {
		ri := <-c
		res[e - 1] = ri
	} 

	return &res
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Query string
		Repos [100]RepoInfos
	}{
		Query: r.URL.Query().Get("query"),
	}

	data.Repos = *GetList(data.Query)

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

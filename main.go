package main

import (
	"net/http"
	"html/template"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type RepoInfos struct {
	FullName string
	Languages map[string]int
}

var ghClient *github.Client

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func GetRepoInfos(r github.Repository, c chan RepoInfos) {
	var ri RepoInfos
	ri.FullName = *r.FullName

	languages, _, err := ghClient.Repositories.ListLanguages(*r.Owner.Login, *r.Name)
	if err != nil {
		fmt.Print(err, "\n")
	} else {
		ri.Languages = languages
	}
	c <- ri
}

func GetList(q string) (*[100]RepoInfos) {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	orgs, _, err := ghClient.Search.Repositories("created:>=2017-02-13T23:25:00+01:00", opt)

	if err != nil {
		fmt.Print(err, "\n")
		return nil
	}

	c := make(chan RepoInfos)

	i := 0
	for _, r := range orgs.Repositories {
		go GetRepoInfos(r, c)
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

func ghAuth() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "537b73a46b51212bd8c394b9ec53504c585486cc"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	ghClient = github.NewClient(tc)
}

func main() {
	ghAuth()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", nil)
}

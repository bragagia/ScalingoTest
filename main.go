package main

import (
	"os"
	"time"
	"net/http"
	"html/template"
	"fmt"
	"strings"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type RepoInfos struct {
	FullName string
	Languages map[string]int
}

type LanguageInfos struct {
	Language string
	Repos map[string]int
	Total int
}

var ghClient *github.Client

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

func GetList() (*[100]RepoInfos) {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	current_time := time.Now().UTC()
	orgs, _, err := ghClient.Search.Repositories("created:>=" + current_time.Format("2006-01-02"), opt)

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

func GetLanguagesList(rl [100]RepoInfos) (*[]LanguageInfos) {
	var res []LanguageInfos

	for _, repo := range rl {
		for lang, langBytes := range repo.Languages {
			found := false
			for _, reslang := range res {
				if (reslang.Language == lang) {
					reslang.Repos[repo.FullName] = langBytes
					reslang.Total += langBytes
					found = true
				}
			}
			if !found {
				var li LanguageInfos
				li.Language = lang
				li.Repos = make(map[string]int)
				li.Repos[repo.FullName] = langBytes
				li.Total += langBytes
				res = append(res, li)
			}
		}
	}

	return &res
}

func FilterList(list []LanguageInfos, q string) ([]LanguageInfos) {
	var res []LanguageInfos

	f := func (li LanguageInfos) (bool) {
		if (strings.Contains(strings.ToLower(li.Language), strings.ToLower(q))) {
			return true
		} else {
			return false
		}
	}

	for _, x := range list {
		if f(x) {
			res = append(res, x)
		}
	}

	return res
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Query string
		Languages []LanguageInfos
	}{
		Query: r.URL.Query().Get("query"),
	}

	repos := *GetList()
	data.Languages = FilterList(*GetLanguagesList(repos), data.Query)

	t, err := template.ParseFiles("templates/search.html")
	if err != nil {
		fmt.Print(err)
	}

	t.Execute(w, &data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func ghAuth(t string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	ghClient = github.NewClient(tc)
}

func main() {
	if (len(os.Args) > 1) {
		ghAuth(os.Args[1])
	} else {
		fmt.Print("Running without login to API.\n")
		fmt.Print("run: ScalingoTest GITHUBAPIKEY\n")
		ghClient = github.NewClient(nil)
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", nil)
}

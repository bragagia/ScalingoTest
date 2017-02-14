# ScalingoTest
**This is a unmaintainted repo**

Technical test for Scalingo.

ScalingoTest allow you to search within 100 random GitHub repos of the current day.

It show the results sorted by language which are used in the repo.

## Installation

Run the following commands :
```
go get github.com/google/go-github/github
go get golang.org/x/oauth2
go get github.com/bragagia/ScalingoTest
go install ScalingoTest
```

### Usage

Run `ScalingoTest` inside of the git directory. Pass it your GitHub API Key as argument to unblock the rate limit.

Example :
```
~/goprojects/src/ScalingoTest> ScalingoTest 537xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx6cc
```

Access it through http://localhost:8080

## How it work

### main()
Handle program arguments and launch webserver

### ghAuth(t)
t: OAuth2 token to use

Enable OAuth2 Authentification.

### indexHandler(w, r)
w: Object to answer request

r: Request infos

Handle home page. Use ./templates/index.html as template.

### searchHandler(w, r)
w: Object to answer request

r: Request infos

Handle search results page. Use ./templates/search.html as template.

### getList()
return: Repo list

Use GitHub API to get 100 random repos of the current date.
Get languages informations asynchronously with getRepoInfos()

### getRepoInfos(r, c)
r: GitHub repository

c: Output chanel

Use GitHub API to get language use informations of the given repo.

### getLanguageList(rl)
rl: Repos list

return: Languages list

Convert the repo list containing a list of language used into a list of language with the list of repositories which use it.

### FilterList(list, q)
list: Language list

q: Word query

return: Language list filtered

Filter the language list by eliminating those which not contain the given word.

# License
Feel free to fork the repository and send pull request.
**This is a unmaintainted repo**

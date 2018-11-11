package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strings"
)

// RepoSuffix is the suffix to add at the end of the repository url.
var RepoSuffix = ""

var errIncompletePath = errors.New("incomplete path")

type parts struct {
	org, proj, repo string
}

func extractRepoPath(path string) (parts, error) {
	p := strings.SplitN(path[1:], "/", 4)
	if len(p) < 3 {
		return parts{}, errIncompletePath
	}
	return parts{
		org:  p[0],
		proj: p[1],
		repo: p[2],
	}, nil
}

var tmpl = template.Must(template.New("vanity").Parse(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
	<meta name="go-import" content="{{.Host}}/{{.Org}}/{{.Proj}}/{{.Repo}} git https://dev.azure.com/{{.Org}}/{{.Proj}}/_git/{{.Repo}}{{.Suffix}}">
	<meta http-equiv="refresh" content="0; url=https://dev.azure.com/{{.Org}}/{{.Proj}}/_git/{{.Repo}}{{.Suffix}}">	
    <title>Vanity URL for Go package hosted on Azure DevOps git repository</title>
  </head>
  <body>
    Please see <a href="https://dev.azure.com/{{.Org}}/{{.Proj}}/_git/{{.Repo}}{{.Suffix}}"> the package on Azure DevOps</a>.
  </body>
</html>
`))

func Index(w http.ResponseWriter, r *http.Request) {
	p, err := extractRepoPath(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours.
	err = tmpl.Execute(w, struct {
		Host   string
		Org    string
		Proj   string
		Repo   string
		Suffix string
	}{
		Host: r.Host,
		Org:  p.org,
		Proj: p.proj,
		Repo: p.repo,
		Suffix: RepoSuffix,
	})
	if err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

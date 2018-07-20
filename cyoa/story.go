package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Demo struct {
	Age  int
	Name string
}

func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, fmt.Errorf("Got error in d.Decode(): %s", err)
	}
	return story, nil
}

type handler struct {
	s Story
}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:] //remove '/' which is path[0]
	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something wrong in tpl.Execute()", http.StatusNotFound)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
	<meta chartset="utf-8">
	<title>Choose Your Own Adventure</title>
  </head>
  <body>
	  <h1>{{.Title}}</h1>
	  {{range .Paragraphs}}
	    <p>{{.}}</p>
	  {{end}}
	  <ul>
	  {{range .Options}}
		<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
	  {{end}}
	  </ul>
  </body>
</html>`

package main

import (
	"fmt"
	"net/http"

	"github.com/user/gophercises/urlshort"
)

func main() {
	mux := defaultMux()

	//get data
	pathToURLs, yaml := getData()

	//get instance of MapHandler
	mapHandler := urlshort.MapHandler(pathToURLs, mux)

	//get instance of YAMLHandler
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)

}

func defaultMux() *http.ServeMux {
	//NewServeMux allocates and returns a new ServeMux
	mux := http.NewServeMux()
	//HandleFunc registers the handler function for the given pattern
	mux.HandleFunc("/", hello)
	return mux
}

//handler that takes a ResponseWritter and a *http.Request
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Oops! The URL is wrong.")
}

func getData() (map[string]string, string) {
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	pathToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	return pathToUrls, yaml
}

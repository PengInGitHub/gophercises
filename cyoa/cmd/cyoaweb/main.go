package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gophercises/cyoa"
)

func main() {
	//get JSON file name from cmd
	fileName := flag.String("json", "gopher.json", "the json file of stories")
	port := flag.Int("port", 3000, "the port to start the CYOA web application")
	flag.Parse()

	//open JSON
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Errorf("Got error in os.Open(): %s", err)
	}

	//decode JSON
	story, err := cyoa.JSONStory(file)
	if err != nil {
		fmt.Errorf("Got error in cyoa.JSONStory(): %s", err)
	}

	//print JSON
	tpl := template.Must(template.New("").Parse("This is a temporary template"))
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))
	fmt.Printf("Start the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
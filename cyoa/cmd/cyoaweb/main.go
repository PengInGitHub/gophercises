package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/user/gophercises/cyoa"
)

func main() {
	//get JSON file name from cmd
	fileName := flag.String("json", "gopher.json", "the json file of stories")
	//get localhost port
	port := flag.Int("port", 3000, "the port to start the CYOA web application")

	flag.Parse()

	//open JSON file
	file, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	//decode JSON to Story
	story, err := cyoa.JSONStory(file)
	if err != nil {
		panic(err)
	}

	//create template
	tpl := template.Must(template.New("").Parse("This is a temporary template")) //temporary
	tpl = template.Must(template.New("").Parse(getTmplstring()))                 //template to use

	//WithTemplate assigns template instance to handler struct
	handlerOption := cyoa.WithTemplate(tpl)
	//create instance of handler with story and hadlerOption
	h := cyoa.NewHandler(story, handlerOption)
	fmt.Printf("Start the server on port: %d\n", *port)
	//ListenAndServe() starts an HTTP server with a given address and handler
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))

}

func getTmplstring() string {
	var defaultHandlerTmpl = `
	<!DOCTYPE html>
	<html>
	  <head>
		<meta charset="utf-8">
		<title>Choose Your Own Adventure</title>
	  </head>
	  <body>
		<section class="page">
		  <h1>{{.Title}}</h1>
		  {{range .Paragraphs}}
			<p>{{.}}</p>
		  {{end}}
		  {{if .Options}}
			<ul>
			{{range .Options}}
			  <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
			</ul>
		  {{else}}
			<h3>The End</h3>
		  {{end}}
		</section>
		<style>
		  body {
			font-family: helvetica, arial;
		  }
		  h1 {
			text-align:center;
			position:relative;
		  }
		  .page {
			width: 80%;
			max-width: 500px;
			margin: auto;
			margin-top: 40px;
			margin-bottom: 40px;
			padding: 80px;
			background: #FFFCF6;
			border: 1px solid #eee;
			box-shadow: 0 10px 6px -6px #777;
		  }
		  ul {
			border-top: 1px dotted #ccc;
			padding: 10px 0 0 0;
			-webkit-padding-start: 0;
		  }
		  li {
			padding-top: 10px;
		  }
		  a,
		  a:visited {
			text-decoration: none;
			color: #6295b5;
		  }
		  a:active,
		  a:hover {
			color: #7792a2;
		  }
		  p {
			text-indent: 1em;
		  }
		</style>
	  </body>
	</html>`

	return defaultHandlerTmpl
}

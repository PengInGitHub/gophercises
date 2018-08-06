package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

//Refine the Panic Recover MiddleWare created in the last exercise
//For example utilize the Chroma package to highlight syntax
//1.Create HTTP handler to render source files in browser
//2.Highlight syntax code via Chroma package
//3.Parse out stack trace
//4.Build links to source files
//5.Highlight lines
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", sourceCodeHandler)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	fmt.Println("start server on local host 3000")
	log.Fatal(http.ListenAndServe(":3000", devMW(mux)))
}

//parse source code to HTML
func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path") //a parameter used in url path
	lineStr := r.FormValue("line")
	lineNumber, err := strconv.Atoi(lineStr)
	if err != nil {
		lineNumber = -1
	}
	file, err := os.Open(path) // For read access.
	if err != nil {
		//it's fine to print out error directly in dev env
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//write the file into byte buffer then cast to string to feed quick.Highlight()
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var lines [][2]int
	if lineNumber > 0 {
		lines = append(lines, [2]int{lineNumber, lineNumber})
	}
	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, b.String())
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(), html.LineNumbersInTable(), html.HighlightLines(lines))
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<style>pre { font-size:1em; }</style>")
	formatter.Format(w, style, iterator)
	//format source text via quick.Highlight of Chroma package
	//quick.Highlight(w, b.String(), "go", "html", "monokai")

	//copy the file into the ResponseWriter
	//io.Copy(w, file)
}

func devMW(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, makeLinks(string(stack)))
			}
		}()

		//defer func() would be called if app.ServeHTTP panics
		//nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(w, r)
		//nw.flush()
	}
}

//just a func panics, no msg written to RW
func panicDemo(w http.ResponseWriter, r *http.Request) {
	// defer func() {
	// 	err := recover()
	// 	fmt.Fprint(w, err)
	// }() //() is used to to call this func
	funcThatPanics()
}

//panic AFTER writing some info to the ResponseWriter
//to deal with info already written into the RW but actually panic occurs
func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

//print out hellp msg
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func funcThatPanics() {
	panic("Oh no!")
}

func makeLinks(stack string) string {
	lines := strings.Split(stack, "\n")
	for lineIndex, line := range lines {
		if len(line) == 0 || line[0] != '\t' { //not start by tab
			continue
		}

		//parse out the link
		file := ""
		for i, ch := range line {
			if ch == ':' {
				file = line[1:i]
				break //quit this loop
			}
		}
		//line = 		/usr/local/go/src/runtime/panic.go:491 +0x283
		lineNumberStr := ""
		//this loop is not so efficient
		for i := len(file) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' { //not an int
				break
			}
			lineNumberStr = lineNumberStr + string(line[i])
		}

		//could use strings.Builder{} instead
		//var lineStr strings.Builder
		//do the loop
		//lineStr.WriteByte(line[i])

		//file =/usr/local/go/src/runtime/panic.go
		//build url
		//turn path into a url encoded value

		//encode href
		v := url.Values{}
		v.Set("path", file)
		v.Set("line", lineNumberStr)

		//href="https//..." so the quote in string should be canceled by \"
		aTagFirstPart := "\t<a href=\"/debug/?" + v.Encode() + "\">"

		linkContent := file

		lineNumberAndSuffix := line[len(file)+2+len(lineNumberStr):]
		lines[lineIndex] = aTagFirstPart + linkContent + ":" + lineNumberStr + "</a>" + lineNumberAndSuffix

	}
	return strings.Join(lines, "\n")
}

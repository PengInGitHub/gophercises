package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

//Utilize panic and recover to create a Panic Recover MiddleWare for a web server
//Capture error in app and render a usefull error page to the user
//In details, create a handler that wraps existing mux that is potentially to panic,
//and recovers from any panic and does following:

//1.log error and stack trace
//2.set status code to http.StatusInternalServerError (500)
//3.write "something went wrong" message when panic occurs
//4.do sth before panic starts
//5.let printed msg depends on env(dev/production)
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	fmt.Println("start server on local host 3000")
	log.Fatal(http.ListenAndServe(":3000", recoverMW(mux, true)))
}

//recover MiddleWare, a generic wrapper to handle error for the entire server
//HandlerFunc is a variable with type of func(w http.ResponseWriter, r *http.Request)
//in this way, this MW could return a func directly
func recoverMW(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				if !dev {
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, string(stack))
			}
		}()

		//defer func() would be called if app.ServeHTTP panics
		nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

type responseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(rw.writes), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
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

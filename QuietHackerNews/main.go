package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/user/gophercises/QuietHackerNews/hn"
)

func main() {
	//parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	//parse template
	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	//create handler
	http.HandleFunc("/", handler(numStories, tpl))

	//start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//stories, err := getTopStories(numStories)
		stories, err := getStories(numStories)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getStories(numStories int) ([]item, error) {
	var stories []item
	var client hn.Client
	IDs, err := client.GetTopItemsID()
	if err != nil {
		return nil, err
	}
	for _, ID := range IDs {
		hnItem, err := client.GetItems(ID)
		if err != nil {
			continue
		}
		item := parseHNItem(hnItem)
		if isStoryLink(item) {
			stories = append(stories, item)
			if len(stories) >= numStories {
				break
			}
		}
	}
	return stories, nil
}

func getTopStories(numStories int) ([]item, error) {
	var stories []item
	var client hn.Client
	IDs, err := client.GetTopItemsID()
	if err != nil {
		return nil, errors.New("Failed to get top items' IDs")
	}
	for _, id := range IDs {
		type result struct {
			item item
			err  error
		}
		resultCh := make(chan result)
		go func(id int) {
			hnItem, err := client.GetItems(id)
			if err != nil {
				resultCh <- result{err: err}
			}
			resultCh <- result{item: parseHNItem(hnItem)}
		}(id)

		res := <-resultCh
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
			if len(stories) >= numStories {
				break
			}
		}
	}
	return stories, nil
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

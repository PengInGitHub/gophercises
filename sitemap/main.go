package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/gophercises/link"
)

func main() {

	fileName := flag.String("url", "https://www.calhoun.io/", "the url which the sitemap is built for")

	maxDepth := flag.Int("depth", 3, "the maximum depth of links to traverse")
	flag.Parse()
	//links := getLinks(*fileName)
	links := bfs(*fileName, *maxDepth)
	for _, l := range links {
		fmt.Println(l)
	}
}

//empty helps to keep an unique set of seen urls WITHOUT creating additional data structure
//struct{} uses less memory
type empty struct{}

//get all urls used for sitemap in a Breadth First Search manner
func bfs(urlStr string, maxDepth int) []string {
	//return all the urls ever visted

	//maps' keys (url) are cashed so that it is faster than slice to look up if there is a value there
	seen := make(map[string]empty)

	//current queue, every key is the url needs to getLinks() with
	var q map[string]empty
	//next queue has a value inside
	nq := map[string]empty{
		urlStr: empty{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]empty)

		if len(q) == 0 {
			break //terminate the loop
		}

		//loop through map: k, v := range map{}, or k := range map{}
		//pull out each url string from the current queue
		for url := range q {

			//skip this round of loop if the url is already seen
			if _, ok := seen[url]; ok {
				continue
			}
			//not in seen so set the url to be seen
			seen[url] = empty{}
			//add obtained links to next queue
			for _, link := range getLinks(url) {
				nq[link] = empty{}
			}
		}
	}
	//length: 0, capacity: len(seen) - pre-allocate all memories
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func getLinks(urlStr string) []string {
	//request the website, get response
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	//build base url from the request url
	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	//get base URL
	base := baseURL.String()
	//compose full url: base+href
	hrefs := getHrefs(resp.Body, base)
	//leave urls have prefix of base only
	return filter(hrefs, withPrefix(base))
}

func getHrefs(r io.Reader, base string) []string {
	//parse links from the response body
	links, _ := link.Parse(r)
	var hrefs []string
	for _, l := range links {
		switch {
		//get full url
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		//get url directly
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	return hrefs
}

func filter(links []string, f func(string) bool) []string {
	var ret []string
	for _, link := range links {
		//base: https://gophercise.com
		if f(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

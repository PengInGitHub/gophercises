package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Parse() takes an HTML file and returns a slice of Link parsed from it
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	//Find <a> nodes in doc
	nodes := linkNodes(doc)

	//Build link for each node
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

type Link struct {
	Href string
	Text string
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = getText(n)
	return ret
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	// if n.Type == html.CommentNode {
	// 	return n.Data
	// }
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

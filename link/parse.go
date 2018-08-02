package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Parse takes an io.Reader, parse the HTML template in the string of the reader, returns a slice of Link
func Parse(r io.Reader) ([]Link, error) {
	//html.Parse(r) returns the parse tree for the HTML from the given reader
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	//get <a> nodes from parsed HTML nodes
	nodes := linkNodes(doc)

	//build link for each node
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

//compose Link from a html.Node
func buildLink(n *html.Node) Link {
	var link Link
	//loop through attributes of a node
	for _, attr := range n.Attr {
		//if attrbuite has key of 'href'
		if attr.Key == "href" {
			//add the Hyperlink to link object's Href field
			link.Href = attr.Val
			break //terminate a loop at this point
		}
	}
	//
	link.Text = getText(n)
	return link
}

func getText(n *html.Node) string {
	//obtain value of TextNode and return
	if n.Type == html.TextNode {
		return n.Data
	}
	// if n.Type == html.CommentNode {
	// 	return n.Data
	// }
	//return empty string if node is not ElementNode
	if n.Type != html.ElementNode {
		return ""
	}
	//loop through child nodes and call getText recursively
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	//return combined strings, seperate by space
	return strings.Join(strings.Fields(ret), " ")
}

//A Node consists of a NodeType and some Data and are part of a tree of Nodes
func linkNodes(n *html.Node) []*html.Node {
	//analysis the Type and Data of a html.Node
	//The <a> tag defines a Hyperlink, which is used to link one page to another
	//if an <a>, add into the slice and return
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	//if not an <a>, loop through all the child nodes, search, collect and return <a> nodes
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		//recursion
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

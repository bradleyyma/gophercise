package parse

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	linkNodes := linkNodes(doc)

	links := []Link{}

	for _, node := range linkNodes {
		link := Link{}
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				link.Href = attr.Val
			}
		}
		link.Text = getText(node)

		links = append(links, link)
	}

	for _, link := range links {

		fmt.Printf("Link: %+v\n", link)
	}
	return links, nil

}

func linkNodes(n *html.Node) []*html.Node {
	var nodes []*html.Node
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, linkNodes(c)...)
	}
	return nodes
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return strings.Join(strings.Fields(text), " ") // Normalize whitespace
}

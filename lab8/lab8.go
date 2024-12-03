package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Define the flag for the maximum number of comments
	max := flag.Int("max", 10, "Max number of comments to show")
	flag.Parse()

	// URL to crawl
	url := "https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html"

	// Fetch the page content
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch URL: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Parse the HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse HTML: %v\n", err)
		os.Exit(1)
	}

	// Extract comments
	comments := extractComments(doc)

	// Print the comments, limiting to the specified max, without extra newlines
	for i, comment := range comments {
		if i >= *max {
			break
		}

		fmt.Printf("%d. %s", i+1, comment)
	}
}

// extractComments traverses the HTML node tree and extracts text comments
func extractComments(n *html.Node) []string {
	var comments []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "push" {
					comment := extractPushContent(n)
					if comment != "" {
						comments = append(comments, comment)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return comments
}

// extractPushContent extracts the text content from a push comment node
func extractPushContent(n *html.Node) string {
	var name, text, time string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					switch attr.Val {
					case "f3 hl push-userid":
						name = getTextContent(n)
					case "f3 push-content":
						text = strings.TrimPrefix(getTextContent(n), ": ")
					case "push-ipdatetime":
						time = getTextContent(n)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	if name != "" && text != "" && time != "" {
		return fmt.Sprintf("名字：%s，留言: %s，時間：%s", name, text, time)
	}
	return ""
}

// getTextContent extracts the text content of an HTML node
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += getTextContent(c)
	}
	return result
}
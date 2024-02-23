package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getHref(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = getHref(links, c)
	}

	return links
}

func getHtmlPage(webPage string) (string, error) {

	resp, err := http.Get(webPage)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	var links []string
	for _, url := range os.Args[1:] {
		// get string html
		date, err := getHtmlPage(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		}

		// fmt.Println(date)
		doc, err := html.Parse(strings.NewReader(date))
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		}

		for _, link := range getHref(links, doc) {
			fmt.Println(link)
		}
	}
}

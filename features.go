package main

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Features branch
//
// Adding email scraping/phone number scrapping

func holder() {
	fmt.Println("Hello mars!")
}

func emailPhone(n *html.Node, emailRegex, phoneRegex *regexp.Regexp) {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if len(text) > 0 {
			emails := emailRegex.FindAllString(text, -1)
			phones := phoneRegex.FindAllString(text, -1)
			for _, e := range emails {
				fmt.Println("Found Email:", e)
			}
			for _, p := range phones {
				fmt.Println("Found phone:", p)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		emailPhone(c, emailRegex, phoneRegex)
	}
}

package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func (cfg *config) GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {

	// get the URL's from the HTML here

	// parse the URL data to break it down into nodes
	// Nodes are a type as per below:
	/*
			type Node struct {
			Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

			Type      NodeType
			DataAtom  atom.Atom
			Data      string
			Namespace string
			Attr      []Attribute
		}

	*/

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println("Error parsing baseURL string")
		return nil, err
	}

	links := []string{}

	// Below just scans the HTML, loads it into memory then outputs it as text
	htmmlReader := strings.NewReader(htmlBody)

	// Logic for node tree creation
	nodeTree, err := html.Parse(htmmlReader)
	if err != nil {
		fmt.Println("Error parsing HTML data to nodes")
		log.Fatal(err)
	}

	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	phoneRegex := regexp.MustCompile(`(?i)(?:\+61|61|\(\+61\))?[\s\-\.]*?(?:\(0?[2-478]\)[\s\-\.]*?\d{4}[\s\-\.]*?\d{4}|0?4\d{2,3}[\s\-\.]*?\d{3}[\s\-\.]*?\d{3})`)

	// Modify here for looking for other elements
	// Add phone number
	// add email
	// add external links
	// This needs to be put into a helper function
	//

	for n := range nodeTree.Descendants() {

		// Helper function for just email and phone
		cfg.emailPhone(n, emailRegex, phoneRegex)

		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					// Check if a.Val has a suffix here
					if strings.HasPrefix(a.Val, "http") {
						links = append(links, strings.TrimSpace(a.Val))
					} else {
						relativeURL, err := url.Parse(a.Val)
						if err != nil {
							fmt.Println("Error parsing relative URL string")
							return nil, err
						}
						finalURL := baseURL.ResolveReference(relativeURL)
						links = append(links, strings.TrimSpace(finalURL.String()))
						break

					}

					break
				}
			}
		}
		// This needs to be put into a helper function
		// Adding new logic here for grabbing numbers

	}

	return links, nil

}

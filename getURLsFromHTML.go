package main

import (
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {

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

	// Modify here for looking for other elements
	// Add phone number
	// add email
	// add external links
	// This needs to be put into a helper function
	//

	for n := range nodeTree.Descendants() {
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
		if n.Type == html.ElementNode && n.Data == "div" {

			//fmt.Printf("n.Data: %v\n", n.Data)
			for _, a := range n.Attr {
				// First check is for mobile numbers and usernames

				if a.Key == "class" && a.Val == "from_name" {
					if strings.Contains(n.FirstChild.Data, "+61") {
						//fmt.Printf("Mobiles: %v\n", n.FirstChild.Data)
						// Need to add a check for exists
						if slices.Contains(links, n.FirstChild.Data) {
							//fmt.Println("Slice contains number!")
							break
						} else {
							links = append(links, n.FirstChild.Data)
						}
					}
				}
			}
		}

	}

	return links, nil

}

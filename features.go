package main

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func (cfg *config) addToEmail(emailAdd string) {

	// Add thread safe mutex here
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, exists := cfg.email[emailAdd]; exists {
		cfg.email[emailAdd]++
		return
	} else {
		cfg.email[emailAdd] = 1
		return
	}
}

func (cfg *config) addToPhone(phoneNum string) {

	// Add thread safe mutex here
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, exists := cfg.phone[phoneNum]; exists {
		cfg.phone[phoneNum]++
		return
	} else {
		cfg.phone[phoneNum] = 1
		return
	}
}

// Features branch
//
// Adding email scraping/phone number scrapping

func (cfg *config) emailPhone(n *html.Node, emailRegex, phoneRegex *regexp.Regexp) {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if len(text) > 0 {
			emails := emailRegex.FindAllString(text, -1)
			phones := phoneRegex.FindAllString(text, -1)
			for _, e := range emails {
				//fmt.Println("Found Email:", e)
				cfg.addToEmail(e)
			}
			for _, p := range phones {
				//fmt.Println("Found phone:", p)
				cfg.addToPhone(p)
			}
		}
	}
	// Added check for html.ElementNodes also
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			// Check if attribute value contains emails/phones (like href="mailto:...")
			attrValue := strings.TrimSpace(attr.Val)
			emails := emailRegex.FindAllString(attrValue, -1)
			phones := phoneRegex.FindAllString(attrValue, -1)
			for _, e := range emails {

				cfg.addToEmail(e)
			}
			for _, p := range phones {
				cfg.addToPhone(p)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cfg.emailPhone(c, emailRegex, phoneRegex)
	}
}

func (cfg *config) printReportEmail(baseURL string) {
	fmt.Printf(`
=============================
  Report for Email addresses for %v
=============================
`, baseURL)

	for email, _ := range cfg.email {
		fmt.Printf("Email address: %v\n", email)
	}

}

func (cfg *config) printReportPhone(baseURL string) {
	fmt.Printf(`
=============================
  Report for Phone Numbers for %v
=============================
`, baseURL)

	for n, _ := range cfg.phone {
		fmt.Printf("Phone Number: %v\n", n)
	}

}

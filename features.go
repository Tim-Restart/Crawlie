package main

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func (cfg *config) addToEmail(emailAdd string) {

	// Add thread safe mutex here
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if slices.Contains(cfg.email, emailAdd) {
		cfg.email = append(cfg.email, emailAdd)
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

func (cfg *config) intToStringPages() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	for key, value := range cfg.pages {
		convertedValue := strconv.Itoa(value)
		cfg.pagesG[key] = convertedValue
	}
}

func (cfg *config) intToStringExternal() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	for key, value := range cfg.external {
		convertedValue := strconv.Itoa(value)
		cfg.externalG[key] = convertedValue
	}
}

func (cfg *config) intToStringEmail() {
	cfg.mu.Lock()
	defer cfg.mu.Lock()
	cfg.emailKeys = make([]string, 0, len(cfg.email))
	for k := range cfg.email {
		cfg.emailKeys = append(cfg.emailKeys, k)
	}
}

/*
func (cfg *config) intToStringEmail() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	for key, value := range cfg.email {
		convertedValue := strconv.Itoa(value)
		cfg.emailG[key] = convertedValue
	}
}

*/

func (cfg *config) intToStringPhone() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	for key, value := range cfg.phone {
		convertedValue := strconv.Itoa(value)
		cfg.phoneG[key] = convertedValue
	}
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

func GetHTML(rawURL string) (string, error) {

	res, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 399 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return "", err
	}
	if !strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
		err = fmt.Errorf("Header not text/html")
		return "", err
	}
	if err != nil {
		return "", err
	}
	return string(body), nil

}

func (cfg *config) crawlPage(rawCurrentURL string) {

	fmt.Printf("Crawling: %v\n", rawCurrentURL)
	//	cfg.concurrencyControl <- struct{}{}
	//	defer func() {
	//		<-cfg.concurrencyControl
	//		cfg.wg.Done()
	//	}()

	//	cfg.mu.Lock()
	//	if len(cfg.pages) >= cfg.maxPages {
	//		fmt.Printf("Max crawl pages reached: %v\n", cfg.maxPages)
	//		cfg.mu.Unlock()
	//		return
	//	}
	//	cfg.mu.Unlock()

	// Base and Current are the same for the first
	// Current is used to do the calls, base is used for a base case
	crawlURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("Error normalising URL")
		return
	}

	// Uses url.Parse to compare the baseURL and CurrentURL (after normalization)
	// to make sure we haven't left the page
	err = compareURL(cfg.baseURL, rawCurrentURL)
	if err != nil {
		cfg.addExternalPage(rawCurrentURL)
		fmt.Println(err)
		return
	}

	// Start the HTML collection and review
	if !cfg.addPageVisit(crawlURL) {
		return
	}

	html, err := GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error getting html for %v\n", rawCurrentURL)

	}

	links, err2 := GetURLsFromHTML(html, rawCurrentURL)
	if err2 != nil {
		fmt.Println("Error getting links from HTML")
		return
	}

	for _, newLink := range links {
		//	cfg.wg.Add(1)
		cfg.crawlPage(newLink) // removed go routine from here
	}

}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	//	cfg.mu.Lock()
	//	defer cfg.mu.Unlock()

	if _, exists := cfg.pages[normalizedURL]; exists {
		cfg.pages[normalizedURL]++
		return false
	}

	cfg.pages[normalizedURL] = 1
	return true
}

// Add external ages to a map in the cfg struct
func (cfg *config) addExternalPage(normalizedURL string) {
	//	cfg.mu.Lock()
	//	defer cfg.mu.Unlock()

	if _, exists := cfg.external[normalizedURL]; exists {
		cfg.external[normalizedURL]++
		return
	} else {
		cfg.external[normalizedURL] = 1
		return
	}

}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for %v internal links
=============================
`, baseURL)

	sortedPageMap := sorted(pages)

	for page, value := range sortedPageMap {
		fmt.Printf("Found %v internal links to %v\n", page, value)
	}

	// Sorting logic
	// Iterate through the list, looking for the highest number...
	// print that first, then continue on down...
	// data structure?
	// Append to new slice?
	// Maybe find a 1 then keep going up? Then print backwards...?
}

func printReportExternal(pages map[string]int, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for external Links for %v
=============================
`, baseURL)

	sortedPageMap := sorted(pages)

	for page, value := range sortedPageMap {
		fmt.Printf("Found %v external links to %v\n", page, value)
	}

	// Sorting logic
	// Iterate through the list, looking for the highest number...
	// print that first, then continue on down...
	// data structure?
	// Append to new slice?
	// Maybe find a 1 then keep going up? Then print backwards...?
}

type KeyValue struct {
	Key   string
	Value int
}

func sorted(pages map[string]int) []KeyValue {

	var sortedPages []KeyValue

	for k, v := range pages {
		sortedPages = append(sortedPages, KeyValue{Key: k, Value: v})
	}

	sort.Slice(sortedPages, func(i, j int) bool {
		return sortedPages[i].Value > sortedPages[j].Value
	})

	return sortedPages

}

package main

import (
	"fmt"

	//"log"
	"net/url"
	"os"

	//"strconv"
	//
	"sync"
)

type config struct {
	pages              map[string]int
	external           map[string]int
	email              map[string]int
	phone              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {

	//var maxConcurrency int
	//	var maxPagesSet int
	// Removed max page at this time as not anticipated to use

	var website string

	switch len(os.Args) {
	case 1:
		fmt.Println("no website provided")
		os.Exit(1)
	case 2:
		fmt.Printf("starting crawl of: %v\n", os.Args[1])
		website = os.Args[1]
	//	fmt.Println("Max pages set to default 10 and concurrency set to default 5")
	//	maxConcurrency = 5
	//	maxPagesSet = 10
	//case 3:
	//	fmt.Printf("starting crawl of: %v\n", os.Args[1])
	//	website = os.Args[1]
	//	maxConcurrency, _ = strconv.Atoi(os.Args[2])
	//	fmt.Printf("Max Concurrency set to : %v\n", maxConcurrency)
	//	fmt.Println("Max pages set to default 10")
	//	maxPagesSet = 10
	//case 4:
	//	fmt.Printf("starting crawl of: %v\n", os.Args[1])
	//	website = os.Args[1]
	//	maxConcurrency, _ = strconv.Atoi(os.Args[2])
	//	fmt.Printf("Max Concurrency set to : %v\n", maxConcurrency)
	//	maxPagesSet, _ = strconv.Atoi(os.Args[3])
	//	fmt.Printf("Max pages set to : %v\n", maxPagesSet)
	default:
		fmt.Println("Failed to set right parameters")
		return
	}

	//fmt.Println("exited switch Statment")
	baseURLParsed, err := stringToURL(website)
	if err != nil {
		fmt.Println("Failed to parse Base URL")
		return
	}

	// Pickup here with the Mu and channels stuff
	cfg := &config{
		pages:    make(map[string]int),
		external: make(map[string]int),
		email:    make(map[string]int),
		phone:    make(map[string]int),
		baseURL:  baseURLParsed,
		//	mu:      &sync.Mutex{},
		//	concurrencyControl: make(chan struct{}, maxConcurrency),
		//	wg: &sync.WaitGroup{},
		//maxPages:           maxPagesSet,
	}

	//cfg.wg.Add(1)
	cfg.crawlPage(website)
	//	cfg.wg.Wait()

	//	cfg.mu.Lock()

	printReport(cfg.pages, website)
	printReportExternal(cfg.external, website)
	cfg.printReportEmail(website)
	cfg.printReportPhone(website)

	//	cfg.mu.Unlock()
	return

}

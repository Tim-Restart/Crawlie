package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	var maxConcurrency int
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
		fmt.Println("Max pages set to default 10 and concurrency set to default 5")
		maxConcurrency = 5
	//	maxPagesSet = 10
	case 3:
		fmt.Printf("starting crawl of: %v\n", os.Args[1])
		website = os.Args[1]
		maxConcurrency, _ = strconv.Atoi(os.Args[2])
		fmt.Printf("Max Concurrency set to : %v\n", maxConcurrency)
		fmt.Println("Max pages set to default 10")
	//	maxPagesSet = 10
	case 4:
		fmt.Printf("starting crawl of: %v\n", os.Args[1])
		website = os.Args[1]
		maxConcurrency, _ = strconv.Atoi(os.Args[2])
		fmt.Printf("Max Concurrency set to : %v\n", maxConcurrency)
	//	maxPagesSet, _ = strconv.Atoi(os.Args[3])
	//	fmt.Printf("Max pages set to : %v\n", maxPagesSet)
	default:
		fmt.Println("Failed to set right parameters")
		return
	}
	log.Println(website)
}

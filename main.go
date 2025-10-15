package main

import (
	//"fmt"
	"net/url"
	"sync"

	//"time"
	//"github.com/briandowns/spinner"
	//"image/color"
	//"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Layout for the GUI - Currently setup for IPloc, to be changed to Crawlie specs
//
// inputDomain, inputConcurrency, inputDelay, crawlButton, phoneNumberLabel, emailLabel, internalLinksLabel, externalLinksLabel

func desktopLayout(inputDomain, inputConcurrency, inputDelay *widget.Entry, crawlButton *widget.Button, phoneNumberLabel, emailLabel, internalLinksLabel, externalLinksLabel *widget.Label) *fyne.Container {
	return container.NewGridWithRows(8, // 8 rows for results/input
		container.NewGridWithColumns(3, // 3 x 3 grid made
			layout.NewSpacer(), // First spacer on left, first column first row
			container.NewVBox( // Second Column in first row (actually two items)
				inputDomain,
				inputConcurrency,
				inputDelay,
			),
			crawlButton, // Third Column in first row
		),
		container.NewVBox(
			fixedLabels,
			rect,
			fixedContainer,
		),
	)
}

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
	delayRequest       int
}

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("Crawlie GUI")

	// Three inputs boxes which have defaults if not set
	// Domain - No default, what are you doing if you don't set this
	// Concurrency - sets how many go routines to spawn for searching concurrently
	// Max delay between requests - done in seconds

	inputDomain := widget.NewEntry()
	inputDomain.SetPlaceHolder("Enter Website to Crawl")
	inputDomain.Resize(fyne.NewSize(100, 20))

	inputConcurrency := widget.NewEntry()
	inputConcurrency.SetPlaceHolder("Enter how many threads")
	inputConcurrency.Resize(fyne.NewSize(50, 20))

	inputDelay := widget.NewEntry()
	inputDelay.SetPlaceHolder("Enter request delay")
	inputDelay.Resize(fyne.NewSize(50, 20))

	var website string

	crawlButton := widget.NewButton("Crawl", func() {
		table.Refresh()
		website = checkInput(input.Text)

	})

	var maxConcurrency int
	//	var maxPagesSet int
	// Removed max page at this time as not anticipated to use

	var website string
	delayReqs := 0 // Delay requests default to 0, set in case 3 if required, and done in seconds

	// CLI input code - left for reference whilst building GUI

	/*
			fmt.Printf("starting crawl of: %v\n", os.Args[1])
			website = os.Args[1]
			maxConcurrency, _ = strconv.Atoi(os.Args[2])
			delayReqs, _ = strconv.Atoi(os.Args[3])
			fmt.Printf("Max Concurrency set to : %v\nDelay Request set to: %v", maxConcurrency, delayReqs)
		//	maxPagesSet, _ = strconv.Atoi(os.Args[3])
		//	fmt.Printf("Max pages set to : %v\n", maxPagesSet)
	*/

	// Pickup here with the Mu and channels stuff
	cfg := &config{
		pages:              make(map[string]int),
		external:           make(map[string]int),
		email:              make(map[string]int),
		phone:              make(map[string]int),
		baseURL:            baseURLParsed,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		//maxPages:           maxPagesSet,
		delayRequest: delayReqs,
	}

	/* Not required for GUI Version

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Color("blue")
	s.Suffix = " Crawling... "

	*/

	cfg.wg.Add(1)
	s.Start()
	go cfg.crawlPage(website)
	cfg.wg.Wait()

	cfg.mu.Lock()
	s.Stop()
	printReport(cfg.pages, website)
	printReportExternal(cfg.external, website)
	cfg.printReportEmail(website)
	cfg.printReportPhone(website)

	cfg.mu.Unlock()

	// Layout for GUI defined here - likely can move to seperate package, but start from here

	phoneNumberLabel := widget.NewLabel("Phone Numbers")
	phoneNumberLabel.Resize(fyne.NewSize(80, 20))
	phoneNumberLabel.Move(fyne.NewPos(10, 0))

	emailLabel := widget.NewLabel("Emails")
	emailLabel.Resize(fyne.NewSize(80, 20))
	emailLabel.Move(fyne.NewPos(10, 0))

	internalLinksLabel := widget.NewLabel("Internal Links")
	internalLinksLabel.Resize(fyne.NewSize(80, 20))
	internalLinksLabel.Move(fyne.NewPos(10, 0))

	externalLinksLabel := widget.NewLabel("external Links")
	externalLinksLabel.Resize(fyne.NewSize(80, 20))
	externalLinksLabel.Move(fyne.NewPos(10, 0))

	myWindow.SetContent(desktopLayout(inputDomain, inputConcurrency, inputDelay, crawlButton, phoneNumberLabel, emailLabel, internalLinksLabel, externalLinksLabel))
	myWindow.ShowAndRun()

	// return - Likely not needed now as is gui

}

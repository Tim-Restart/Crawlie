package main

import (
	//"fmt"

	"net/url"
	"strconv"
	"sync"

	//"time"
	//"github.com/briandowns/spinner"
	//"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Layout for the GUI - Currently setup for IPloc, to be changed to Crawlie specs
//
// inputDomain, inputConcurrency, inputDelay, crawlButton, phoneNumberLabel, emailLabel, internalLinksLabel, externalLinksLabel

func desktopLayout(inputDomain, inputConcurrency, inputDelay *widget.Entry, crawlButton *widget.Button, phoneNumberLabel, emailLabel, internalLinksLabel, externalLinksLabel *widget.Label, phoneNumberTable, emailTable, externalTable, internalTable *widget.Table) *fyne.Container {
	return container.NewGridWithRows(9, // 8 rows for results/input
		container.NewGridWithColumns(3, // 3 x 3 grid made - This is row 1/9
			layout.NewSpacer(), // First spacer on left, first column first row
			container.NewVBox( // Second Column in first row (actually two items)
				inputDomain,
				inputConcurrency,
				inputDelay,
			),
			crawlButton, // Third Column in first row
		),
		phoneNumberLabel,   // This is row 2/9
		phoneNumberTable,   // This is row 3/9
		emailLabel,         // This is row 4/9
		emailTable,         // This is row 5/9
		internalLinksLabel, // This is row 6/9
		externalTable,      // This is row 7/9
		externalLinksLabel, // This is row 8/9
		internalTable,      // This is row 9/9
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

	log.Println("Started program")
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

	log.Println("Setup Inputs")

	cfg := &config{
		pages:    make(map[string]int),
		external: make(map[string]int),
		email:    make(map[string]int),
		phone:    make(map[string]int),
		//baseURL:            baseURLParsed,
		mu: &sync.Mutex{},
		//concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
		//maxPages:           maxPagesSet,
		//delayRequest: delayReqs,
	}
	log.Println("cfg initialized")

	var website string

	//var maxConcurrency int
	//	var maxPagesSet int
	// Removed max page at this time as not anticipated to use

	//var delayReqs int // Delay requests initialized, but default to 0, and done in seconds
	//var baseURLParsed *url.URL

	// Crawl button commences the crawl of the specified domain
	// If maxConcurrency is not set it defaults to 5
	// If delayReqs is not set it defaults to 0 - this is in seconds
	crawlButton := widget.NewButton("Crawl", func() {
		// table.Refresh() - No table made yet
		website = inputDomain.Text
		if inputConcurrency != nil {
			maxConcurrency, _ := strconv.Atoi(inputConcurrency.Text)
			cfg.concurrencyControl = make(chan struct{}, maxConcurrency)
		} else {
			cfg.concurrencyControl = make(chan struct{}, 5)
		}
		if inputDelay != nil {
			cfg.delayRequest, _ = strconv.Atoi(inputDelay.Text)
		} else {
			cfg.delayRequest = 0
		}
		cfg.baseURL, _ = stringToURL(website)
		/*
			if err != nil {
				fmt.Println("Failed to parse Base URL")
				os.Exit(1)
			}
		*/

		cfg.wg.Add(1)
		log.Println("This might print")
		go cfg.crawlPage(website)
		log.Println("Then this one")
		cfg.wg.Wait()
		log.Println("Doubt this will")

	})

	log.Println("Crawl Button Created")

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

	/* Not required for GUI Version

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Color("blue")
	s.Suffix = " Crawling... "

	*/

	//cfg.mu.Lock()

	//printReport(cfg.pages, website)
	//printReportExternal(cfg.external, website)
	//cfg.printReportEmail(website)
	//cfg.printReportPhone(website)

	//cfg.mu.Unlock()

	log.Println("Bet this doesn't print")

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

	var data = [][]string{[]string{"top left", "top right"},
		[]string{"bottom left", "bottom right"}}

	phoneNumberTable := widget.NewTable(
		func() (int, int) {
			return len(data), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("phone numbers")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col]) // Do a int conversion to string? Helper func perhaps?
		})

	emailTable := widget.NewTable(
		func() (int, int) {
			return len(cfg.email), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Email addresses")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col]) // Do a int conversion to string? Helper func perhaps?
		})

	internalTable := widget.NewTable(
		func() (int, int) {
			return len(cfg.pages), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Internal pages")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col]) // Do a int conversion to string? Helper func perhaps?
		})

	externalTable := widget.NewTable(
		func() (int, int) {
			return len(cfg.external), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("External Links")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col]) // Do a int conversion to string? Helper func perhaps?
		})

	myWindow.SetContent(desktopLayout(inputDomain, inputConcurrency, inputDelay, crawlButton, phoneNumberLabel, emailLabel, internalLinksLabel, externalLinksLabel, phoneNumberTable, emailTable, externalTable, internalTable))
	myWindow.ShowAndRun()

	// return - Likely not needed now as is gui

}

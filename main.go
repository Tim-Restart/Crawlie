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

func desktopLayout(inputDomain, inputConcurrency, inputDelay *widget.Entry, crawlButton *widget.Button, emailTable *widget.List) *fyne.Container {
	return container.NewGridWithRows(5, // 5 rows for results/input
		container.NewGridWithColumns(3, // 3 x 3 grid made - This is row 1/9
			layout.NewSpacer(), // First spacer on left, first column first row
			container.NewVBox( // Second Column in first row (actually two items)
				inputDomain,
				inputConcurrency,
				inputDelay,
			),
			crawlButton, // Third Column in first row
		),
		//phoneNumberTable, // This is row 2/5
		emailTable, // This is row 3/5
		//externalTable,    // This is row 4/5
		//internalTable,    // This is row 5/5
	)
}

type config struct {
	pages              []string
	external           []string
	email              []string
	phone              []string
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
		mu: &sync.Mutex{},
		//concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
		//maxPages:           maxPagesSet,
		//delayRequest: delayReqs,
		email: []string{},
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
		//log.Println("This might print")
		cfg.crawlPage(website)
		log.Println("Converting Maps")
		cfg.intToStringPages()
		cfg.intToStringExternal()
		cfg.intToStringPhone()
		cfg.intToStringEmail()
		cfg.wg.Wait()
		//log.Println("Doubt this will")

	})

	//log.Println("Crawl Button Created")

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

	//log.Println("Bet this doesn't print")

	// Layout for GUI defined here - likely can move to seperate package, but start from here

	//cfg.phoneG[string("445 445 556")] = string("1")
	cfg.email["test@gmail.com"] = 2
	//cfg.externalG["https://test.com"] = "5"
	//cfg.pagesG["http://dev/h.com"] = "1"

	/*

		phoneNumberTable := widget.NewTable(
			func() (int, int) {
				return len(cfg.phoneG), 2
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("phone numbers")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				labelPhone := o.(*widget.Label)
				labelPhone.SetText(cfg.phoneG[i.Row][i.Col])
			})

		phoneNumberTable.ShowHeaderRow = true
		phoneNumberTable.CreateHeader = func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		phoneNumberTable.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText("Number")
			case 1:
				label.SetText("Phone Number")
			}
		}

	*/

	emailTable := widget.NewList(
		func() int {
			return len(cfg.emailKeys)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""), // For key
				widget.NewLabel(""), // For value
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			keyLabel := container.Objects[0].(*widget.Label)
			valLabel := container.Objects[1].(*widget.Label)

			key := cfg.emailKeys[i]
			keyLabel.SetText(key)
			valInt, _ := strconv.Atoi(key)
			valLabel.SetText(cfg.emailKeys[valInt])
		},
	)
	/*
		internalTable := widget.NewTable(
			func() (int, int) {
				return len(cfg.pagesG), 2
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("Internal pages")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(cfg.pagesG[i.Row][i.Col])
			})

		externalTable := widget.NewTable(
			func() (int, int) {
				return len(cfg.externalG), 2
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("External Links")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(cfg.externalG[i.Row][i.Col])
			})

	*/

	myWindow.SetContent(desktopLayout(inputDomain, inputConcurrency, inputDelay, crawlButton, emailTable))
	myWindow.ShowAndRun()

	// return - Likely not needed now as is gui

}

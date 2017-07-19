package main

import (
	"flag"
	"fmt"
	"regexp"
	"time"

	"github.com/tebeka/selenium"
)

var idRegex *regexp.Regexp
var priceRegex *regexp.Regexp
var bedRegex *regexp.Regexp

// this is the container contains info of a listing
const infoContainerSelector = ".infoContainer_v72lrv"

// this is the container contains link to the listing, which has unique identifier
const linkContainerSelector = ".linkContainer_15ns6vh"

// available paramters
var searchTerms string
var numOfPages int

func init() {
	idRegex = regexp.MustCompile(`/[0-9]+`)
	priceRegex = regexp.MustCompile(`\$[0-9]+`)
	bedRegex = regexp.MustCompile(`\d beds?`)

	flag.StringVar(&searchTerms, "searchTerms", "Brooklyn, NY", "Location you wish to search for listings")
	flag.IntVar(&numOfPages, "pages", 1, "Number of pages you wish to search for")

	flag.Parse()
}

func main() {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "firefox"})
	if webDriver, err = selenium.NewRemote(caps, "http://localhost:4444/wd/hub"); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	defer webDriver.Quit()

	err = webDriver.Get(fmt.Sprintf("https://www.airbnb.com/s/%s/homes", searchTerms))
	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}

	firstURL, err := webDriver.CurrentURL()
	if err != nil {
		fmt.Printf("Failed to get CurrentURL: %s", err)
		return
	}

	var results []Result

	for i := 0; i < numOfPages; i++ {

		// first page is opened already, we will need to reload after first page
		if i > 0 {
			err = webDriver.Get(getOffsetPageURL(firstURL, i))
			if err != nil {
				fmt.Printf("Failed to load page: %s\n", err)
				return
			}

			// a hack to avoid findElements before the page is loaded, otherwise
			// no results because the page has not loaded completely
			<-time.After(5 * time.Second)
		}

		elems, err := webDriver.FindElements(selenium.ByCSSSelector, infoContainerSelector)
		if err != nil {
			fmt.Printf("Failed to find element: %s\n", err)
			return
		}

		pageResults := processElements(elems)
		results = append(results, pageResults...)
	}

	fmt.Println(fmt.Sprintf("Total Found: %d", len(results)))
	for _, result := range results {
		fmt.Println(result)
	}
}

func getOffsetPageURL(firstURLString string, offset int) string {
	return fmt.Sprintf("%s&section_offset=%d", firstURLString, offset)
}

func processElements(elements []selenium.WebElement) (results []Result) {
	for _, element := range elements {
		linkElement, err := element.FindElement(selenium.ByCSSSelector, linkContainerSelector)
		if err != nil {
			fmt.Printf("Failed to find link element: %s\n", err)
			continue
		}

		value, err := linkElement.GetAttribute("href")
		if err != nil {
			fmt.Printf("Failed to get attribute: %s\n", err)
			continue
		}

		var result Result
		result.ID = idRegex.FindString(value)

		if text, err := element.Text(); err == nil {
			result.Price = priceRegex.FindString(text)
			result.Beds = bedRegex.FindString(text)

			results = append(results, result)
		} else {
			fmt.Printf("Failed to get text of element: %s\n", err)
			continue
		}
	}

	return
}

// type
////
type Result struct {
	ID    string
	Price string
	Beds  string
}

func (r Result) String() string {
	return "ID: " + r.ID + " Price: " + r.Price + " - " + r.Beds
}

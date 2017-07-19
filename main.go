package main

import (
	"fmt"
	"regexp"
	"time"

	selenium "sourcegraph.com/sourcegraph/go-selenium"
)

var idRegex *regexp.Regexp
var priceRegex *regexp.Regexp
var bedRegex *regexp.Regexp

const infoContainerSelector = ".infoContainer_v72lrv"
const linkContainerSelector = ".linkContainer_15ns6vh"

// available paramters
const searchTerms = "jersey city"
const numOfPages = 2

func init() {
	idRegex = regexp.MustCompile(`/[0-9]+`)
	priceRegex = regexp.MustCompile(`\$[0-9]+`)
	bedRegex = regexp.MustCompile(`\d beds?`)
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

	err = webDriver.Get("https://www.airbnb.com/s/" + searchTerms + "/homes")
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
			<-time.After(3 * time.Second)
		}

		elems, err := findElements(webDriver)
		if err != nil {
			fmt.Printf("Failed to find element: %s\n", err)
			return
		}

		pageResults := processElements(elems)
		results = append(results, pageResults...)
	}

	fmt.Println(fmt.Sprintf("Total: %d", len(results)))
	for _, result := range results {
		fmt.Println(result)
	}
}

func findElements(webDriver selenium.WebDriver) (elements []selenium.WebElement, err error) {
	// ".infoContainer_v72lrv" is the container contains info of a listing
	elements, err = webDriver.FindElements(selenium.ByCSSSelector, infoContainerSelector)
	return
}

func getOffsetPageURL(firstURLString string, offset int) string {
	return firstURLString + "&section_offset=" + fmt.Sprintf("%d", offset)
}

func processElements(elements []selenium.WebElement) (results []Result) {
	for _, element := range elements {
		linkElement, err := element.FindElement(selenium.ByCSSSelector, linkContainerSelector)
		if err != nil {
			fmt.Printf("Failed to find link element: %s\n", err)
			return
		}

		value, err := linkElement.GetAttribute("href")
		if err != nil {
			fmt.Printf("Failed to get attribute: %s\n", err)
			return
		}

		var result Result
		result.ID = idRegex.FindString(value)

		if text, err := element.Text(); err == nil {
			result.Price = priceRegex.FindString(text)
			result.Beds = bedRegex.FindString(text)

			results = append(results, result)
		} else {
			fmt.Printf("Failed to get text of element: %s\n", err)
			return
		}
	}

	return
}

type Result struct {
	ID    string
	Price string
	Beds  string
}

func (r Result) String() string {
	return "ID: " + r.ID + " Price: " + r.Price + " - " + r.Beds
}

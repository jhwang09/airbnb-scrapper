# airbnb-scrapper
airbnb scrapper using golang

## Setup
1. Install `go`
2. create go workspace
3. go to the workspace and git pull this project, then `cd airbnb-scrapper`
4. ~~`go get sourcegraph.com/sourcegraph/go-selenium`~~ `go get github.com/tebeka/selenium`
5. Install Selenium, preferably using `brew`
6. add `geckodriver` to path so that selenium can pick up driver for firefox
7. `selenium-server -port 4444`
8. `go run main.go`

## Run Binary directly

run `./airbnb-scrapper -h` to get command line help menu

Sample:
```
./airbnb-scrapper -searchTerms="Boston" -pages=3
```

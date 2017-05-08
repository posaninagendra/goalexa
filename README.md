# goalexa
A concurrent go crawler.

## Dependencies

#### Languages
* [GO](https://golang.org/doc/install)

#### Packages
* [coop](https://github.com/rakyll/coop): `go get github.com/rakyll/coop`
* [boltdb](https://github.com/boltdb/bolt): `go get github.com/boltdb/bolt`
* [progress-bar](https://github.com/cheggaaa/pb): `go get github.com/cheggaaa/pb`
* [go-command-line-tool](https://github.com/codegangsta/cli): `go get github.com/codegangsta/cli`

#### Crawler
The crawler reads the alexa-1m website list and crawls the data and saves them in boltdb database. Using the cronjob we run the crawler daily and collect the data. Steps to run the crawler are given below.
1. Load the sites to crawl from the .csv.gz file:
   	```
   	cd crawler/go/src/goalexa/
	go run main.go load.go goalexa.go cache
   	```
2. Start crawling, by defualt the crawler uses 100 parallel jobs but you can specify using -j JOBS (upto 256):
   	```
  	go run main.go load.go goalexa.go start -j 100
   	```

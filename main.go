package main 

import (
	"log"
	"runtime"
	
	"github.com/codegangsta/cli"
)

var processes = runtime.NumCPU() * 2
var processBuffer = 128 * processes

const sitesNum int64 = 1e6

var app = cli.NewApp()

func init() {
	log.SetFlags(log.Lshortfile)

	app.Name = "crawl"
	app.Usage = "an utility for fetching and parsing Alexa Top-1M sites in parallel mode."
	app.Version = "1.0"
	app.Commands = []cli.Command{
		{
			Name:   "cache",
			Action: loadSites,
			Usage:  "caches the list of domains to crawl",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "db", Value: "sites.db", Usage: "a path to DB with cached domains"},
				cli.StringFlag{Name: "csv", Value: "top-1m.csv.gz", Usage: "csv with domains list in format `rank,domain`"},
			},
		},
		{
			Name:   "start",
			Action: crawlSites,
			Usage:  "start crawling the cached sites and check against some patterns",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "db", Value: "sites.db", Usage: "a path to DB with cached domains"},
				cli.StringFlag{Name: "out, o", Value: "crawled.db", Usage: "a file to output the results to"},
				cli.IntFlag{Name: "level, l", Value: 1, Usage: "how deeply the crawler should go (1-3)"},
				cli.IntFlag{Name: "jobs, j", Value: 1, Usage: "maximum parallel jobs allowed (1-64)"},
				cli.IntFlag{Name: "skip, s", Usage: "skips the defined number of top-positions"},
			},
		},
	}
}

func main() {
	app.RunAndExitOnError()
}
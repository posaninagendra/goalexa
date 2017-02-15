// Created by Nagendra Posani
// Date: Feb 15, 2017
// Crawler package crawls a given URI

package crawler

import (
	"fmt"
	"net/http"
	"bytes"
	"errors"
	"runtime"
	"strconv"
	"time"

	"github.com/rakyll/coop"
	"github.com/boltdb/bolt"
	"github.com/fern4lvarez/go-metainspector/metainspector"	
)


var processes = runtime.NumCPU() * 2
var processBuffer = 128 * processes

var BucketSites = []bytes("sites")

const connectionTimeout = time.Second * 10

type Config struct{
	Jobs int
	Level int
}

type Crawler struct{
	OnProgress func() int

	db *bolt.DB 
	config *Config
	uriChan chan string
}

type Result struct{
	URL         string
	Host 		string
	Title       string
	Description string
	Body 		string
	Language 	string
	Links		[]string
	Headers		map[string][]string
	LastModified [string]
}

func New(db *bolt.DB, config *Config) *Crawler{
	return &Crawler{db: db, config: config,}
}

func (cr *Crawler) Crawl(resultChan chan<-*Result, errorChan chan<-struct{}, skip int) error{
	cr.uriChan = make(chan string, processBuffer)
	skipB := strconv.AppendInt([]byte{}, int64(skip), 10)
	if err := cr.checkBucket(BucketSites); err != nil {
		return err
	}

	// fill the pool of URLs
	go func() {
		cr.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(BucketSites)
			b.ForEach(func(k, v []byte) error {
				if bytesLess(k, skipB) {
					return nil
				}
				cr.uriChan <- string(v)
				return nil
			})
			return nil
		})
		close(cr.uriChan)
	}()

	process := func() {
		for uri := range cr.uriChan {
			if cr.OnProgress != nil {
				cr.OnProgress()
			}
			result, err := cr.crawlURI(uri)
			if err != nil {
				errorChan <- struct{}{}
				continue
			}
			resultChan <- result
		}
	}

	// run processing jobs
	<-coop.Replicate(cr.config.Jobs, process)
	close(resultChan)
	close(errorChan)

return nil
}

func (cr *Crawler) crawlURI(uri string) (result *Result, err error) {
	mi, err := metainspector.New(uri, &metainspector.Config{
		Timeout: connectionTimeout,
	})
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(uri)
	if err != nil{
		return nil, err
	}
	resData, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println(err)
	}
	body := string(resData)
	headers := resp.Headers
	result = &Result{
		URL:         mi.Url(),
		Host: 		 mi.Host(),
		Title:       mi.Title(),
		Description: mi.Description(),
		Body:		 body,
		Language:    mi.Language(),
		Links: 		 mi.Links(),
		Headers: 	 headers,
		LastModified: headers["Last-Modified"]
	}
	
	return
}

func (cr *Crawler) checkBucket(name []byte) error {
	return cr.db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(name); b == nil {
			return errors.New("no such bucket: " + string(name))
		}
		return nil
	})
}

// bytesLess return true iff a < b.
func bytesLess(a, b []byte) bool {
	return bytes.Compare(a, b) < 0
}
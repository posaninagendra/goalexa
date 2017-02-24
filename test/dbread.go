package main 

import (
	"log"
	"github.com/boltdb/bolt"
	"encoding/json"
	"goalexa/crawler"
)

const openMode = 0644

func readFromDb(){
	log.Println("db read came")
	db := openDB("../crawled.db")
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("crawledSites"))
			b.ForEach(func(k, v []byte) error {
				var res crawler.Result
				err := json.Unmarshal(v, &res)
				if err != nil{
					log.Println("error",err)
				}
				log.Println(string(k), res)
				return nil
			})
			return nil
		})
}

func openDB(name string) *bolt.DB {
	db, err := bolt.Open(name, openMode, nil)
	if err != nil {
		log.Fatalln("error opening db:", err)
	}
	return db
}

func main() {
	readFromDb()
}
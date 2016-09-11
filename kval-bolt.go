package main

import (
   "fmt"
   "time"
   "github.com/boltdb/bolt"
)

func Connect() (kvalbolt, error) {
   fmt.Println("Create a DB and return a connection.")
   var kb kvalbolt
   db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 2 * time.Second})
   kb.db = db
   return kb, err
}

func Disconnect(db *bolt.DB) {
   db.Close()
}

func Query(query string) error {
   
   return nil
}

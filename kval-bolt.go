package main

import (
   "fmt"
   "time"
   "github.com/boltdb/bolt"
   "github.com/kval-access-language/KVAL-Parse"
)

func donothing() {
   fmt.Println("Create a DB and return a connection.")
}

func Connect(dbname string) (kvalbolt, error) {
   var kb kvalbolt
   db, err := bolt.Open(dbname, 0600, &bolt.Options{Timeout: 2 * time.Second})
   kb.db = db
   return kb, err
}

func Disconnect(kb kvalbolt) error {
   err := kb.db.Close()
   return err
}

func Query(kb kvalbolt, query string) (kvalresult, error) {
   var kr kvalresult
   var err error
   kq, err := kvalparse.Parse(query)
   if err != nil {
      return kr, err
   }
   kb.query = kq 
   kr, err = queryhandler(kb)
   if err != nil {
      return kr, err
   }  
   return kr, nil
}

func queryhandler(kb kvalbolt) (kvalresult, error) {
   var kr kvalresult
   fmt.Println("handle query.")
   return kr, nil
}

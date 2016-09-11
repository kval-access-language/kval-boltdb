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

/*type KQUERY struct { 
   function Token
   buckets []string  
   key string
   value string
   newname string
   regex bool
}*/

func queryhandler(kb kvalbolt) (kvalresult, error) {
   var kr kvalresult
   switch kb.query.Function {
   case kvalparse.INS:
      kr, err := insertHandler(kb.query)
      return kr, err
   case kvalparse.GET:
   case kvalparse.LIS:
   case kvalparse.DEL:
   case kvalparse.REN:
   default:
      fmt.Errorf("Function not implemented yet: %v", kb.query.Function)
   }
   return kr, nil
}

//INS
//More validation needed in the Parser
func insertHandler(kq kvalparse.KQUERY) (kvalresult, error) {
   var kr kvalresult
   if kq.Key == "" && kq.Value == "" {
      //we'll make new buckets
      kr, err := createboltbuckets(kq)
      return kr, err
   } else if kq.Key != "" && kq.Value != "" {
      //we'll create a key value
   } else if kq.Key != "" && kq.Value == "" {
      //create a nil key
   }
   return kr, nil
}


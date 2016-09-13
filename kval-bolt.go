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

func Disconnect(kb kvalbolt) {
   kb.db.Close()
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
      err := insHandler(kb)
      return kr, err
   case kvalparse.GET:
      if kb.query.Key == "" {
         //get all
         kr, err := getallHandler(kb)
         return kr, err
      } else {
         kr, err := getHandler(kb)
         return kr, err
      }
   case kvalparse.LIS:
   case kvalparse.DEL:
   case kvalparse.REN:
   default:
      fmt.Errorf("Function not implemented yet: %v", kb.query.Function)
   }
   return kr, nil
}

//INS
func insHandler(kb kvalbolt) error {
   //as long as there are buckets, we can create
   //anything we need. it all happens in a single
   //transaction, based on kval query...
   err := createboltentries(kb)
   if err != nil {
      return err
   }      
   return nil
}

//GET
func getHandler(kb kvalbolt) (kvalresult, error) {
   var kr kvalresult
   kr, err := viewboltentries(kb)
   if err != nil {
      return kr, err
   }
   return kr, nil
}

func getallHandler(kb kvalbolt) (kvalresult, error) {
   kr, err := getallfrombucket(kb)
   if err != nil {
      return kr, err
   }
   return kr, nil
}
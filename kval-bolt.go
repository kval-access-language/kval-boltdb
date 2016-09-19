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
      kr, err := lisHandler(kb)
      return kr, err
   case kvalparse.DEL:
      if kb.query.Key == "" {
         //we're deleting a bucket (and all contents)
         err := delbucketHandler(kb)
         return kr, err      
      } else if kb.query.Key == "_" {
         //we're making nil "" values for all keys
         //use case, we want the keys, we don't want the values
         err := delbucketkeysHandler(kb)
         return kr, err
      } else if kb.query.Key != "" && kb.query.Key != "_" && kb.query.Value != "_" {
         //we're deleting a key and its value
         err := delonekeyHandler(kb)
         return kr, err
      } else if kb.query.Value == "_" {
         //we're deleting a value and leaving the key
         err := nullifyvalHandler(kb)
         return kr, err                  
      } 
   case kvalparse.REN:
      if kb.query.Key == "" {
         renbucketHandler(kb)
      } else if kb.query.Key != "" {
         renkeyHandler(kb)
      }
   default:
      fmt.Errorf("Function not implemented yet: %v", kb.query.Function)
   }
   return kr, nil
}

//INS
func insHandler(kb kvalbolt) error {
   //as long as there are buckets, we can create anything we need. 
   //it all happens in a single transaction, based on kval query...
   err := createboltentries(kb)
   if err != nil {
      return err
   }      
   return nil
}

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

func delbucketHandler(kb kvalbolt) error {
   err := deletebucket(kb)
   if err != nil {
      return err
   }   
   return nil
}

func delbucketkeysHandler(kb kvalbolt) error {
   err := deletebucketkeys(kb)
   if err != nil {
      return err
   }
   return nil
}

func delonekeyHandler(kb kvalbolt) error {
   err := deletekey(kb)
   if err != nil {
      return err
   }
   return nil
}

func nullifyvalHandler(kb kvalbolt) error {
   err := nullifykeyvalue(kb)
   if err != nil {
      return err
   }
   return nil
}

func renbucketHandler(kb kvalbolt) error {
   err := renamebucket(kb)
   if err != nil {
      return err
   }
   return nil
}

func renkeyHandler(kb kvalbolt) error {
   err := renamekey(kb)
   if err != nil {
      return err 
   }
   return nil
}

func lisHandler(kb kvalbolt) (kvalresult, error) {
   kr, err := bucketkeyexists(kb)
   if err != nil {
      //Nil bucket returns an error we can use
      if kr.Exists == true {
         return kr, err
      }
      return kr, err
   }
   return kr, nil
}

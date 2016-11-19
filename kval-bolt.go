package main

import (
   "time"
   "github.com/boltdb/bolt"
   "github.com/pkg/errors"
   "github.com/kval-access-language/KVAL-Parse"
)

//Open a BoltDB with a given name to work with. 
//Our first most important function. Returns a KVAL Bolt structure
//with the details required for KBAL BoltDB to perform queries.
func Connect(dbname string) (kvalbolt, error) {
   var kb kvalbolt
   db, err := bolt.Open(dbname, 0600, &bolt.Options{Timeout: 2 * time.Second})
   kb.db = db
   return kb, err
}

//Disconnect from a BoltDB.
//Recommended that this is made a deferred function call where possible. 
func Disconnect(kb kvalbolt) {
   kb.db.Close()
}

//Query. Given a KVALBolt Structure, and a KVAL query string
//this function will do all of the work for you when interacting with
//BoltDB. Everything should become less programmatic making for cleaner code.
//The KVAL spec can be found here: https://github.com/kval-access-language/kval
func Query(kb kvalbolt, query string) (kvalresult, error) {
   var kr kvalresult
   var err error
   kq, err := kvalparse.Parse(query)
   if err != nil {
      return kr, errors.Wrapf(err, "%s: '%s'", err_parse, query)
   }
   kb.query = kq 
   kr, err = queryhandler(kb)
   if err != nil {
      return kr, err
   }  
   return kr, nil
}

//Abstracted away from Query() query handler is an unexported function that
//will route all queries as required by the application when given by the user.
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
      //function is parsed correctly but not recognised by binding
      return kr, errors.Wrapf(err_not_implemented, "%v", kb.query.Function)
   }
   return kr, nil
}

//INS (Insert Handler) handles INS capability of KVAL language
func insHandler(kb kvalbolt) error {
   //as long as there are buckets, we can create anything we need. 
   //it all happens in a single transaction, based on kval query...
   err := createboltentries(kb)
   if err != nil {
      return err
   }      
   return nil
}

//GET (Get Handler) handles GET capability of KVAL language
func getHandler(kb kvalbolt) (kvalresult, error) {
   var kr kvalresult
   kr, err := viewboltentries(kb)
   if err != nil {
      return kr, err
   }
   return kr, nil
}

//GET (Get Handler) handles GET (ALL) capability of KVAL language
func getallHandler(kb kvalbolt) (kvalresult, error) {
   kr, err := getallfrombucket(kb)
   if err != nil {
      return kr, err
   }
   return kr, nil
}

//DEL (Delete Handler) handles DEL bucket capability of KVAL language
func delbucketHandler(kb kvalbolt) error {
   err := deletebucket(kb)
   if err != nil {
      return err
   }   
   return nil
}

//DEL (Delete Handler) handles DEL all keys capability of KVAL language
func delbucketkeysHandler(kb kvalbolt) error {
   err := deletebucketkeys(kb)
   if err != nil {
      return err
   }
   return nil
}

//DEL (Delete Handler) handles DEL one key capability of KVAL language
func delonekeyHandler(kb kvalbolt) error {
   err := deletekey(kb)
   if err != nil {
      return err
   }
   return nil
}

//DEL (Delete Handler) Handles DEL (or in this case, NULL, capability of KVAL
func nullifyvalHandler(kb kvalbolt) error {
   err := nullifykeyvalue(kb)
   if err != nil {
      return err
   }
   return nil
}

//REN (Rename Handler) Handles rename bucket capability of KVAL
func renbucketHandler(kb kvalbolt) error {
   err := renamebucket(kb)
   if err != nil {
      return err
   }
   return nil
}

//REN (Rename Handler) Handles rename key capability of KVAL
func renkeyHandler(kb kvalbolt) error {
   err := renamekey(kb)
   if err != nil {
      return err 
   }
   return nil
}

//LIS (List Handler) Handles listing capability of KVAL (does (x) exist?) 
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

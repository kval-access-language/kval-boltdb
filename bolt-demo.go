package main

import (
   "os"
   "fmt"
) 

var nilresult kvalresult

func main() {

   kb, err := Connect("newdb.bolt")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error creating new bolt database: %v", err)
      os.Exit(1)
   }
   res, err := Query(kb, "INS triage bucket >> document bucket")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
   }
   
   if res != nilresult {
      fmt.Println(res)
   }
   err = Disconnect(kb)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error disconnecting database: %v", err)
   }
}

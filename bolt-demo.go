package main

import (
   "os"
   "fmt"
) 

func main() {

   kb, err := Connect("newdb.bolt")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error opening bolt database: %v", err)
      os.Exit(1)
   }
   defer Disconnect(kb)

   res, err := Query(kb, "INS triage bucket >> document bucket >> testbucket >>>> test :: value")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
   }
   
   if res != nilresult {
      fmt.Println(res)
   }
}

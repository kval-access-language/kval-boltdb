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

   var res kvalresult

   var testins = []string{
      "INS triage bucket >> document bucket >> testbucket >>>> test1 :: value1",
      "INS triage bucket >> document bucket >> testbucket >>>> test2 :: value2",
      "INS triage bucket >> document bucket >> testbucket >>>> test3 :: value3",
   }

   for _, value := range(testins) {
      _, err = Query(kb, value)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
      }
   }


   _, err = Query(kb, "INS triage bucket >> document bucket >> testbucket >>>> abc :: def")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
   }

   res, err = Query(kb, "GET triage bucket >> document bucket >> testbucket >>>> abc")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
   }   
   
   if res.Result != nil {
      fmt.Println("Result one:", res.Result)
   }

   var testget = []string{
      "GET triage bucket >> document bucket >> testbucket >>>> test1",
      "GET triage bucket >> document bucket >> testbucket >>>> test2",
      "GET triage bucket >> document bucket >> testbucket >>>> test3",
   }

   for _, value := range(testget) {
      res, err = Query(kb, value)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
      } else {
         fmt.Println("GET loop:", res.Result)
      }
   }

   res, err = Query(kb, "GET triage bucket >> document bucket >> testbucket")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error trying to get all.")
   }

   if res.Result != nil{
      fmt.Println("get all result:", res.Result)
   }

   res, err = Query(kb, "GET triage bucket >> document bucket")
   if err != nil {
      fmt.Fprintf(os.Stderr, "%v\n", err)
   } else if res.Result != nil{
      fmt.Println("get all result:", res.Result)
   }   
}
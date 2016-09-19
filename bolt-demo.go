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
      "INS triage bucket >> document bucket >> testbucket >>>> test4",   
      "INS triage bucket >> document bucket >> testbucket >> inline bucket",
      "INS triage bucket >> document bucket >> delbucket >>>> a1 :: b1",
      "INS triage bucket >> document bucket >> delbucket >>>> a2 :: b2",
      "INS triage bucket >> document bucket >> delbucket >>>> a3 :: b3",
      "INS triage bucket >> document bucket >> delbucket >>>> a4 :: b4",
      "INS triage bucket >> document bucket >> delbucket >> bucketbucket",    //test delete nested buckets in one go
      "INS triage bucket >> document bucket >> testbucket >> bucketbucket",    
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
      "GET triage bucket >> document bucket >> testbucket >>>> test4",      
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
   } 
   if res.Result != nil{
      fmt.Println("get all result:", res.Result)
   }   


   res, err = Query(kb, "GET triage bucket >> document bucket >> delbucket")
   if err != nil {
      fmt.Fprintf(os.Stderr, "%v\n", err)
   } 
   if res.Result != nil{
      fmt.Println("get all DEL result:", res.Result)
   }   


   var deltests = []string {
      //"DEL triage bucket >> document bucket >> testbucket",
      //"DEL triage bucket >> document bucket >> testbucket >>>> test2",
      //"DEL triage bucket >> document bucket >> testbucket >>>> bucketbucket",
      "DEL  triage bucket >> document bucket >> delbucket >>>> a2 :: _",
      //"DEL triage bucket >> document bucket >> delbucket >>>> _",     //del all keys from a bucket
   }

   for _, value := range(deltests) {
      _, err = Query(kb, value)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
      }
   }

   res, err = Query(kb, "GET triage bucket >> document bucket >> delbucket")
   if err != nil {
      fmt.Fprintf(os.Stderr, "%v\n", err)
   } 
   if res.Result != nil{
      fmt.Println("get all DEL result:", res.Result)
   }      

   var testrens = []string{
      "INS r1 >> r2 >> r3 >>>> rk1 :: vvvrv1",
      "INS r1 >> r2 >> r3 >>>> rk2 :: rv3",
      "INS r1 >> r2 >> r3 >> r4",
      "INS r1 >> r2 >> r3 >> r4 >>>> rk11 :: rv11",
      "INS r1 >> r2 >> r3 >> r4 >>>> rk21 :: rv21",
      "INS r1 >> r2 >> r3 >> r4 >>>> rk31 :: rv31",
      "INS r1 >> r2 >> r3 >>>> rk3 :: rv4",
      "REN r1 >> r2 >> r3 >>>> rk1 => rename1", 
      "REN r1 >> r2 => supersonic",  
   }

   for _, value := range(testrens) {
      _, err = Query(kb, value)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
      }
   }     

   res, err = Query(kb, "GET r1 >> r2")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Expected error result for r2: %v\n", err)
   } else {
      fmt.Println("REN result:", res.Result)
   }      

   res, err = Query(kb, "GET r1 >> supersonic >> r3")
   if err != nil {
      fmt.Fprintf(os.Stderr, "%v\n", err)
   } else {
      fmt.Println("REN result:", res.Result)
   } 

   var testlis = []string{
      "INS ins1 >> ins2 >> ins3 >>>> key1 :: value1",
      "LIS ins1 >> ins2",
      "lis ins1 >> ins2 >> ins3 >>>> key1",
      "lis ins1 >> ins2 >> ins3 >>>> keyunknown",
      "lis ins1 >> ins2 >> ins4",
   }

   for _, value := range(testlis) {
      kr, err := Query(kb, value)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
      } else {
         fmt.Println("Exist results:", kr.Exists)
      }
   }     
}

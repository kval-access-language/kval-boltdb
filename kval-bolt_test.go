package main

import (
   "os"   
   "log"
   "testing"
)

var (
   dbloc = "bolt-test-db/test-db.bolt"
   kb  kvalbolt
   err error
)

func init() {
   log.Println("Info: Initialize unit tests.")   
   setup()
}

func refreshdb() {
   clear()
   setup()
}

func setup() {
   kb, err = Connect(dbloc)
   if err != nil {
      log.Printf("Error opening bolt database: %v\n", err)
      os.Exit(1)
   }
}

func clear() {
   Disconnect(kb)
   f, err := os.Create(dbloc)
   if err != nil {
      log.Println("Error resetting db during teardown.")
   }
   f.Close()
}

func teardown() {
   log.Println("Info: Running tear-down.")
   clear()
}

//Populate a database with data to work with for testing
func doinserts() {
   //clear db when we need it afresh...
   refreshdb()  
   //baseline inserts...   
   for _, value := range(ins_tests) {
      _, err = Query(kb, value)
      if err != nil {
         log.Printf("Error querying db: %v\n", err)
      }
   }
}

//Test insert functions associated with KVAL capabilities
func testins(t *testing.T) {
   doinserts()

   // Utilise BoltDB Tree statistics.
   // KeyN  int // number of keys/value pairs
   // Depth int // number of levels in B+tree

   bs, _ := getbucketstats(kb, ins_getbuckets1)
   if bs.KeyN != ins_result1.keys && bs.Depth != ins_result1.depth {
      t.Errorf("Expected stats results for INS don't match")
   } 

   bs, _ = getbucketstats(kb, ins_getbuckets2)
   if bs.KeyN != ins_result2.keys && bs.Depth != ins_result2.depth {
      t.Errorf("Expected stats results for INS don't match")
   } 

   bs, _ = getbucketstats(kb, ins_getbuckets3)
   if bs.KeyN != ins_result3.keys && bs.Depth != ins_result3.depth {
      t.Errorf("Expected stats results for INS don't match")
   } 
}

//Test list functions associated with KVAL capabilities
func testlis(t *testing.T) {
   doinserts()
   for k, v := range(lis_results) {
      kq, err := Query(kb, k)
      if err != nil {
         log.Printf("Error querying db: %v\n", err)
      }
      if kq.Exists != v {
         t.Errorf("Expected %b got %b.\n", v, kq.Exists)
      } 
   }
}

func testdel(t *testing.T) {
   doinserts()
   for k, _ := range(del_results) {
      _, err := Query(kb, k)
      if err != nil {
         log.Printf("Error querying db: %v\n", err)
      }
   }
}

func TestQuery(t *testing.T) {
   defer teardown()
   //testins(t)   
   //testlis(t)
   testdel(t)
   //testget(t)
   //testren(t)
}

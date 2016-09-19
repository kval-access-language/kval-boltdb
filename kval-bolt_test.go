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
   kb, err = Connect(dbloc)
   if err != nil {
      log.Printf("Error opening bolt database: %v\n", err)
      os.Exit(1)
   }
}

func teardown() {
   log.Println("Info: Running tear-down.")
   Disconnect(kb)
   f, err := os.Create(dbloc)
   if err != nil {
      log.Println("Error resetting db during teardown.")
   }
   f.Close()
}

func TestQuery(t *testing.T) {
   defer teardown()

   for _, value := range(ins_tests) {
      _, err = Query(kb, value)
      if err != nil {
         log.Printf("Error querying db: %v\n", err)
      }
   }

   //KeyN || Depty
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

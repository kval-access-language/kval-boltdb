package main

import (
   "os"   
   "log"
   "fmt"   
   "reflect"
   "strings"
   "testing"
   "github.com/pkg/errors"
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

//---------------------------------------------------------------------------//

//Miscellaneous tests (tests that don't group nicely together)
func testnotimplementedfuncs(t *testing.T) {
   _, err = Query(kb, make_tea)
   if err == nil {
      log.Printf("Error expected from test but not returned.")
   } else {
      //TODO: rethink testing the error string... can github.com/pkg/errors help?
      teststr := fmt.Sprintf("%s", err)
      if !strings.Contains(teststr, err_parse) {
         log.Printf("Error querying db: %v\n", err)
      }
   }
}

//Test handling of unicode and big strings, e.g. for blogs...
func testbigstring(t *testing.T) {

   var unistrings = [...]string{bigstring_one, bigstring_two}
   var key = "str"   

   for i := range(unistrings) {
      _, err = Query(kb, "INS bigstring >>>> " + key + " :: " + unistrings[i])
      if err != nil {
         t.Errorf("Error returned when not expected while trying to store valid bigstring from BoltDB:", err)
      }

      res, err := Query(kb, "GET bigstring >>>> " + key)
      if err != nil {
         t.Errorf("Error returned when not expected while trying to retrieve valid bigstring from BoltDB:", err)
      }

      if res.Result[key] != unistrings[i] {
         t.Errorf("Unicode big strings: Error retrieving expected Unicode result back from BoltDB.")
      }
   }
}


//Test handling of Base64 strings, e.g. for blob encoding...
func testbase64(t *testing.T) {

   for _, value := range(ins_b64_values) {
      _, err = Query(kb, value)
      if err != nil {
         log.Printf("Error creating state for base64 unit tests: %v\n", err)
      }
   }

   for k, v := range(get_b64_results) {
      res, err := Query(kb, k)
      if err != nil {
         log.Printf("Error found when not expected in Base64 retrieve tests: %v\n", err)
      }
      if !reflect.DeepEqual(res.Result, v) {
         t.Errorf("Base64 retrieve failed for query: %v\n", k)
      }
   }
}

//---------------------------------------------------------------------------//

//Populate a database with data to work with for testing
func create_state_inserts() {
   //clear db when we need it afresh...
   refreshdb()  
   //baseline inserts...   
   for _, value := range(create_initial_state) {
      _, err = Query(kb, value)
      if err != nil {
         log.Printf("Error creating state for unit tests, exiting: %v\n", err)
         os.Exit(1)
      }
   }
}

//Test insert functions associated with KVAL capabilities
func testins(t *testing.T) {
   create_state_inserts()

   // Utilise BoltDB Tree statistics.
   // KeyN  int // number of keys/value pairs
   // Depth int // number of levels in B+tree

   for i := range(ins_checks_all) {
      bs, _ := getbucketstats(kb, ins_checks_all[i].buckets)
      if bs.KeyN != ins_checks_all[i].counts.keys && bs.Depth != ins_checks_all[i].counts.depth {
         t.Errorf("Expected stats results for INS don't match")
      } 
   }
}

//Test list functions associated with KVAL capabilities
func testlis(t *testing.T) {
   create_state_inserts()
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

//Test delete functions associated with KVAL capabilities
func testdel(t *testing.T) {

   //test results we expect to pass
   for k := range(good_del_results) {
      
      //create_state_inserts: slower but efficient for test writing... maintains
      //constant state throughout *any* delete we do...       
      create_state_inserts()

      //perform our queries on the Bolt DB...      
      _, err := Query(kb, good_del_results[k])
      if err != nil {
         t.Errorf("Invalid error for delete procedure (nil error expected): %v\n", err)
      }
   }

   //test results we expect to fail, and check fail result...
   for k, e := range(bad_del_results) {
      
      //create_state_inserts: slower but efficient for test writing... maintains
      //constant state throughout *any* delete we do...       
      create_state_inserts()

      switch k {
         case delnonekeytwo:
            //testing nul result where Bolt returns nil when trying to delete
            //a key that doesn't actually exist...
            bs, _ := getbucketstats(kb, bucket_nonekey)

            //KeyN  int // number of keys/value pairs
            //compare expected keys to remaining keys - should be identical
            expectedkeys := bs.KeyN

            _, err := Query(kb, k)
            if err != nil {
               if errors.Cause(err) != e {
                  t.Errorf("Invalid error for delete procedure (nil expected for none key): %v\n", err)
               }
            }      

            bs, _ = getbucketstats(kb, bucket_nonekey)
            remainingkeys := bs.KeyN
            if expectedkeys != remainingkeys {
               t.Errorf("Invalid error deleting nil key. Expected 'nil' return from BoltDB: %v\n", err)
            }

         default: 
            //perform our queries on the Bolt DB...      
            _, err := Query(kb, k)
            if err != nil {
               if errors.Cause(err) != e {
                  t.Errorf("Invalid error for DEL procedure (different error expected): %v\n", err)
               }
            }      
      }
   }

}

func testget(t *testing.T) {
   create_state_inserts()

   //test regular gets
   for k, v := range(get_sole_results) {
      res, err := Query(kb, k)
      if err != nil {
         t.Errorf("Invalid error for GET procedure (zero errors expected): %v\n", err)
      }
      if !reflect.DeepEqual(res.Result, v) {
         t.Errorf("Unexpected result value for GET: %s, expected: %s\n", res.Result, v)
      }
   }

   //test regex gets
   for k, v := range(get_regex_results) {
      res, err := Query(kb, k)
      if err != nil {
         t.Errorf("Unexpected error returned for GET regex: %v\n", err)
      }     
      if !reflect.DeepEqual(res.Result, v) {
         t.Errorf("Unexpected result value for GET: %s, expected: %s\n", res.Result, v)
      }
   }
}

func TestQuery(t *testing.T) {
   defer teardown()
   //testnotimplementedfuncs(t)
   //testbigstring(t)
   testbase64(t)
   //testins(t)   
   //testlis(t)
   //testdel(t)
   //testget(t)
   //testren(t)
}

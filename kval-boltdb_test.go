package kvalbolt

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	dbloc = "bolt-test-db/test-db.bolt"
	kb    Kvalboltdb
	err   error
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
	_, err = Query(kb, makeTea)
	if err == nil {
		log.Printf("Error expected from test but not returned.")
	} else {
		//TODO: rethink testing the error string... can github.com/pkg/errors help?
		teststr := fmt.Sprintf("%s", err)
		if !strings.Contains(teststr, errParse) {
			log.Printf("Error querying db: %v\n", err)
		}
	}
}

//Test handling of unicode and big strings, e.g. for blogs...
func testbigstring(t *testing.T) {
	var unistrings = [...]string{bigStringOne, bigStringTwo}
	var key = "str"
	for i := range unistrings {
		_, err = Query(kb, "INS bigstring >>>> "+key+" :: "+unistrings[i])
		if err != nil {
			t.Errorf("Error returned when not expected while trying to store valid bigstring from BoltDB:", err)
		}

		res, err := Query(kb, "GET bigstring >>>> "+key)
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

	for _, value := range insB64Values {
		_, err = Query(kb, value)
		if err != nil {
			log.Printf("Error creating state for base64 unit tests: %v\n", err)
		}
	}

	for k, v := range getB64Results {
		res, err := Query(kb, k)
		if err != nil {
			log.Printf("Error found when not expected in Base64 retrieve tests: %v\n", err)
		}
		if !reflect.DeepEqual(res.Result, v) {
			t.Errorf("Base64 retrieve failed for query: %v\n", k)
		}
	}
}

//Tests PutBlob with various scnarios, starting with a simple one
func testPutBlob(t *testing.T) {
	//Test complete cycle of behaviour that a user might go through...
	for encode, result := range simpleB64Results {
		err := StoreBlob(kb, "INS Blob Bucket >>>> Blob Key", "image/png", []byte(encode))
		if err != nil {
			t.Errorf("Putblob failed for query: %v\n", err)
		}
		res, err := Query(kb, "GET Blob Bucket >>>> Blob Key")
		if err != nil {
			t.Errorf("Retrieve failed for GET query: %v\n", err)
		}
		kvb, err := UnwrapBlob(res)
		if err != nil {
			t.Errorf("Unwrapblob failed: %v\n", err)
		}
		if kvb.Data != result {
			t.Errorf("Unwrap failed with incorrect result: %s expected: %s\n", kvb.Data, result)
		}
		data, err := GetBlobData(kvb)
		if err != nil {
			t.Errorf("Error with GetBlobData: %v.", err)
		} else {
			if string(data) != encode {
				t.Errorf("Base64 blob not decoded correctly: %v, expected: %v.", string(data), encode)
			}
		}
	}
}

//Test that stats.db is working as we hope, that is, test that each query
//returns with an appropriate BucketStats structure associated with it.
func teststatsdb(t *testing.T) {

	var statstest1 = map[string]int{"LIS one bucket": 0}
	var statstest2 = map[string]int{"INS one bucket >>>> only key": 1}
	var statstest3 = map[string]int{"GET one bucket >>>> only key": 1}
	var statstest4 = map[string]int{"INS one bucket >>>> another key": 2}
	var statstest5 = map[string]int{"REN one bucket >>>> another key => new name": 2}
	var statstest6 = map[string]int{"LIS one bucket >>>> new name": 2}
	var statstest7 = map[string]int{"DEL one bucket >>>> new name": 1}
	var statstest8 = map[string]int{"LIS one bucket": 1}
	var statstest9 = map[string]int{"REN one bucket => delete me": 1}
	var statstest10 = map[string]int{"LIS one bucket": 0}
	var statstest11 = map[string]int{"LIS delete me": 1}
	var statstest12 = map[string]int{"GET delete me": 1}
	var statstest13 = map[string]int{"DEL delete me": 0}
	var statstest14 = map[string]int{"LIS delete me": 0}

	var statstests = [...]map[string]int{statstest1, statstest2, statstest3,
		statstest4, statstest5, statstest6, statstest7,
		statstest8, statstest9, statstest10, statstest11,
		statstest12, statstest13, statstest14}

	//TODO: Have a think about making these tests work a bit more systematically
	for i := range statstests {
		for k, v := range statstests[i] {
			res, err := Query(kb, k)
			if err != nil {
				log.Println(k, err)
			}
			if res.Stats.KeyN != v {
				t.Errorf("Incorrect stats returned query: %s: %d, %d", k, res.Stats.KeyN, v)
			}
			// take this opportunity to do some additional tests
			if k == "LIS delete me" && v == 1 {
				if res.Exists != true {
					t.Error("Bucket should have been created as part of rename.")
				}
			} else if k == "LIS one bucket" && v == 0 {
				if res.Exists != false {
					t.Error("Bucket should have been deleted as part of rename at tx level.")
				}
			} else if k == "GET delete me" {
				if len(res.Result) != 1 {
					t.Error("Error spotted renaming bucket at the transaction level.")
				}
				if _, ok := res.Result["only key"]; !ok {
					t.Error("Error spotted renaming bucket at the transaction level.")
				}
			}
		}
	}
}

//---------------------------------------------------------------------------//

//Populate a database with data to work with for testing
func createStateInserts() {
	//clear db when we need it afresh...
	refreshdb()
	//baseline inserts...
	for _, value := range createInitialState {
		_, err := Query(kb, value)
		if err != nil {
			log.Printf("Error creating state for unit tests, exiting: %v\n", err)
			os.Exit(1)
		}
	}
}

//Test insert functions associated with KVAL capabilities
func testins(t *testing.T) {
	createStateInserts()

	// Utilise BoltDB Tree statistics.
	// KeyN  int // number of keys/value pairs
	// Depth int // number of levels in B+tree

	for i := range insChecksAll {
		bs, _ := getbucketstats(kb, insChecksAll[i].buckets)
		if bs.KeyN != insChecksAll[i].counts.keys && bs.Depth != insChecksAll[i].counts.depth {
			t.Errorf("Expected stats results for INS don't match")
		}
	}
}

//Test list functions associated with KVAL capabilities
func testlis(t *testing.T) {
	createStateInserts()
	for k, v := range lisResults {
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
	for k := range goodDelResults {

		//create_state_inserts: slower but efficient for test writing... maintains
		//constant state throughout *any* delete we do...
		createStateInserts()

		//perform our queries on the Bolt DB...
		_, err := Query(kb, goodDelResults[k])
		if err != nil {
			t.Errorf("Invalid error for delete procedure (nil error expected): %v\n", err)
		}

		switch goodDelResults[k] {
		case delkey:
			//"DEL bucket one >> bucket two >> bucket three >>>> test1" //delete key test1
			res, _ := Query(kb, "LIS bucket one >> bucket two >> bucket three >>>> test1")
			if res.Exists != false {
				t.Errorf("Delete failed for key 'test1', still exists.")
			}
		case nullvalue:
			//"DEL bucket one >> bucket two >> bucket three >>>> test3 :: _" //make value null without deleting key
			res, _ := Query(kb, "GET bucket one >> bucket two >> bucket three >>>> test3")
			if val, ok := res.Result["test3"]; ok {
				if val != "" {
					t.Errorf("Nullify key failed for key 'test3', found: %v", val)
				}
			} else {
				t.Errorf("Error querying key 'test3'. Key not found.")
			}
		case delkeys:
			//"DEL bucket one >> bucket two >> bucket three >>>> _" //del all keys from a bucket
			bs, err := getbucketstats(kb, []string{"bucket one", "bucket two", "bucket three"})
			if err != nil {
				t.Errorf("Error retrieving key count for bucket, %v", err)
			} else if bs.KeyN != 0 {
				t.Errorf("Key count following delete from bucket is incorrect, %d", bs.KeyN)
			}
		case delbucket:
			//"DEL bucket one >> bucket two" //delete bucket two
			res, _ := Query(kb, "LIS bucket one >> bucket two")
			if res.Exists != false {
				t.Errorf("Delete failed for bucket 'bucket two', still exists.")
			}
		}

	}

	//test results we expect to fail, and check fail result...
	for k, e := range badDelResults {

		//create_state_inserts: slower but efficient for test writing... maintains
		//constant state throughout *any* delete we do...
		createStateInserts()

		switch k {
		case delnonekeytwo:
			//testing nul result where Bolt returns nil when trying to delete
			//a key that doesn't actually exist...
			bs, _ := getbucketstats(kb, bucketNoneKey)

			//KeyN  int // number of keys/value pairs
			//compare expected keys to remaining keys - should be identical
			expectedkeys := bs.KeyN

			results, err := Query(kb, k)
			if err != nil {
				if errors.Cause(err) != e {
					t.Errorf("Invalid error for delete procedure (nil expected for none key): %v\n", err)
				}
			}
			//compare our expected keys statistic with the stats produced from
			//running the query. Also allows us to make use of statsdb again.
			if expectedkeys != results.Stats.KeyN {
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
	createStateInserts()

	//test regular gets
	for k, v := range getSoleResults {
		res, err := Query(kb, k)
		if err != nil {
			t.Errorf("Invalid error for GET procedure (zero errors expected): %v %s\n", err, k)
		}
		if !reflect.DeepEqual(res.Result, v) {
			t.Errorf("Unexpected result value for GET: %s, expected: %s\n", res.Result, v)
		}
	}

	//test regex gets
	for k, v := range getRegexResults {
		res, err := Query(kb, k)
		if err != nil {
			t.Errorf("Unexpected error returned for GET regex: %v\n", err)
		}
		if !reflect.DeepEqual(res.Result, v) {
			t.Errorf("Unexpected result value for GET: %s, expected: %s\n", res.Result, v)
		}
	}
}

// just test getroot
func testgetroot(t *testing.T) {
	createStateInserts()
	//test regular gets
	for k, v := range getRootResults {
		res, err := Query(kb, k)
		if err != nil {
			t.Errorf("Invalid error for GET procedure (zero errors expected): %v %s\n", err, k)
		}
		if !reflect.DeepEqual(res.Result, v) {
			t.Errorf("Unexpected result value for GET: %s, expected: %s\n", res.Result, v)
		}
	}
}

func renamestate(t *testing.T) {
	//setup new state for rename functions
	for i := range renameState {
		_, err := Query(kb, renameState[i])
		if err != nil {
			t.Errorf("Unexpected error returned setting up rename state: %v\n", err)
		}
	}
}

func testren(t *testing.T) {

	//run tests...
	for k, v := range renameTests {
		var oldstats Kvalresult

		//setup state...
		renamestate(t)

		//grab a count for bucket rename tests
		switch k {
		case r2:
			//query how many keys are in our bucket before renaming...
			//allows us to doubly check LIS still works and statdb works...
			oldstats, _ = Query(kb, renOldList)
		}

		_, err := Query(kb, v)
		if err != nil {
			t.Errorf("Error with rename function: %v\n", err)
		}

		//check list results for renames
		switch k {
		case r1:
			for i := range renLis1 {
				res, _ := Query(kb, renLis1[i])
				switch i {
				case 0:
					if res.Exists != false {
						t.Errorf("Rename key failed, bucket or key still exists.")
					}
				case 1:
					if res.Exists != true {
						t.Errorf("Rename key failed, bucket or key doesn't exist.")
					}
				}
			}
		case r2:
			for i := range renLis2 {
				res, _ := Query(kb, renLis2[i])
				switch i {
				case 0:
					if res.Exists != false {
						t.Errorf("Rename bucket failed, bucket or key still exists.")
					}
				case 1:
					if res.Exists != true {
						t.Errorf("Rename bucket failed, bucket or key doesn't exist.")
					}
				}
			}
			//now query the count for our newly renamed bucket, and compare to the old count...
			//allows us to doubly check LIS still works and statdb works...
			newstats, _ := Query(kb, renNewList)
			if newstats.Stats.KeyN != oldstats.Stats.KeyN {
				t.Errorf("Bucket count following rename doesn't match: %d, old: %d", newstats.Stats.KeyN, oldstats.Stats.KeyN)
			}
			if newstats.Stats.Depth != oldstats.Stats.Depth {
				t.Errorf("Bucket count following rename doesn't match: %d, old: %d", newstats.Stats.Depth, oldstats.Stats.Depth)
			}
		}
	}
}

func TestQuery(t *testing.T) {
	defer teardown()
	testnotimplementedfuncs(t)
	testbigstring(t)
	testbase64(t)
	testPutBlob(t)
	teststatsdb(t)
	testins(t)
	testlis(t)
	testdel(t)
	testget(t)
	testren(t)
	testgetroot(t)
}

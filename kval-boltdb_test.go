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
	kb    kvalbolt
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
	_, err = Query(kb, make_tea)
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
	var unistrings = [...]string{bigstring_one, bigstring_two}
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

	for _, value := range ins_b64_values {
		_, err = Query(kb, value)
		if err != nil {
			log.Printf("Error creating state for base64 unit tests: %v\n", err)
		}
	}

	for k, v := range get_b64_results {
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

	for encode, result := range simple_b64_results {
		err := Putblob(kb, "INS Blob Bucket >>>> Blob Key", "image/png", []byte(encode))
		if err != nil {
			t.Errorf("Putblob failed for query: %v\n", err)
		}
		res, err := Query(kb, "GET Blob Bucket >>>> Blob Key")
		if err != nil {
			t.Errorf("Retrieve failed for GET query: %v\n", err)
		}
		kvb, err := Unwrapblob(res)
		if err != nil {
			t.Errorf("Unwrapblob failed: %v\n", err)
		}
		if kvb.Data != result {
			t.Errorf("Unwrap failed with incorrect result: %s expected: %s\n", kvb.Data, result)
		}
	}
}

//---------------------------------------------------------------------------//

//Populate a database with data to work with for testing
func create_state_inserts() {
	//clear db when we need it afresh...
	refreshdb()
	//baseline inserts...
	for _, value := range create_initial_state {
		_, err := Query(kb, value)
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

	for i := range ins_checks_all {
		bs, _ := getbucketstats(kb, ins_checks_all[i].buckets)
		if bs.KeyN != ins_checks_all[i].counts.keys && bs.Depth != ins_checks_all[i].counts.depth {
			t.Errorf("Expected stats results for INS don't match")
		}
	}
}

//Test list functions associated with KVAL capabilities
func testlis(t *testing.T) {
	create_state_inserts()
	for k, v := range lis_results {
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
	for k := range good_del_results {

		//create_state_inserts: slower but efficient for test writing... maintains
		//constant state throughout *any* delete we do...
		create_state_inserts()

		//perform our queries on the Bolt DB...
		_, err := Query(kb, good_del_results[k])
		if err != nil {
			t.Errorf("Invalid error for delete procedure (nil error expected): %v\n", err)
		}

		switch good_del_results[k] {
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
	for k, e := range bad_del_results {

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
	for k, v := range get_sole_results {
		res, err := Query(kb, k)
		if err != nil {
			t.Errorf("Invalid error for GET procedure (zero errors expected): %v\n", err)
		}
		if !reflect.DeepEqual(res.Result, v) {
			t.Errorf("Unexpected result value for GET: %s, expected: %s\n", res.Result, v)
		}
	}

	//test regex gets
	for k, v := range get_regex_results {
		res, err := Query(kb, k)
		if err != nil {
			t.Errorf("Unexpected error returned for GET regex: %v\n", err)
		}
		if !reflect.DeepEqual(res.Result, v) {
			t.Errorf("Unexpected result value for GET: %s, expected: %s\n", res.Result, v)
		}
	}
}

func renamestate(t *testing.T) {
	//setup new state for rename functions
	for i := range rename_state {
		_, err := Query(kb, rename_state[i])
		if err != nil {
			t.Errorf("Unexpected error returned setting up rename state: %v\n", err)
		}
	}
}

func testren(t *testing.T) {

	//run tests...
	for k, v := range rename_tests {
		var oldcount, olddepth int

		//setup state...
		renamestate(t)

		//grab a count for bucket rename tests
		switch k {
		case r2:
			//query how many keys are in our bucket before renaming...
			bs, _ := getbucketstats(kb, ren_slice_old)
			oldcount = bs.KeyN
			olddepth = bs.Depth
		}

		_, err := Query(kb, v)
		if err != nil {
			t.Errorf("Error with rename function: %v\n", err)
		}

		//check list results for renames
		switch k {
		case r1:
			for i := range ren_lis1 {
				res, _ := Query(kb, ren_lis1[i])
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
			for i := range ren_lis2 {
				res, _ := Query(kb, ren_lis2[i])
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
			newcount, _ := getbucketstats(kb, ren_slice_new)
			if newcount.KeyN != oldcount {
				t.Errorf("Bucket count following rename doesn't match: %d, old: %d", newcount.KeyN, oldcount)
			}
			if newcount.Depth != olddepth {
				t.Errorf("Bucket count following rename doesn't match: %d, old: %d", newcount.Depth, olddepth)
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
	testins(t)
	testlis(t)
	testdel(t)
	testget(t)
	testren(t)
}

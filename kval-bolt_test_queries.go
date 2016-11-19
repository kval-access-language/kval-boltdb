package main

import "github.com/boltdb/bolt"

//test invalid/non-implemented capabilities

var make_tea = "TEA bucket one >> bucket two >>>> cup :: saucer"

//---------------------------------------------------------------------------//

//test insert procedures
var create_initial_state = []string{
   "INS bucket one >> bucket two >> bucket three >>>> test1 :: value1",
   "INS bucket one >> bucket two >> bucket three >>>> test2 :: value2",
   "INS bucket one >> bucket two >> bucket three >>>> test3 :: value3",
   "INS bucket one >> bucket two >>>> test4 :: value4",
   "INS bucket one >> bucket two >>>> test5 :: value5",
   "INS bucket one >>>> test6 :: value6", 
}

var ins_getbuckets1 = []string{"bucket one", "bucket two", "bucket three"}
var ins_getbuckets2 = []string{"bucket one", "bucket two"}
var ins_getbuckets3 = []string{"bucket one"}

var ins_result1 = insresult{3, 1}
var ins_result2 = insresult{6, 2}
var ins_result3 = insresult{8, 3}

type ins_check struct {
   buckets []string 
   counts insresult 
}

type insresult struct {
   keys int
   depth int
}

var i1 = ins_check{ins_getbuckets1, ins_result1}
var i2 = ins_check{ins_getbuckets2, ins_result2}
var i3 = ins_check{ins_getbuckets3, ins_result3}

var ins_checks_all = [...]ins_check{i1, i2, i3}

//---------------------------------------------------------------------------//

//test delete procedures 

//good delete procedures - nil error expected... 
var delkey = "DEL bucket one >> bucket two >> bucket three >>>> test1"           //delete key test1
var nullvalue = "DEL bucket one >> bucket two >> bucket three >>>> test3 :: _"   //make value null without deleting key
var delkeys = "DEL bucket one >> bucket two >> bucket three >>>> _"              //del all keys from a bucket
var delbucket = "DEL bucket one >> bucket two"                                   //delete bucket two        

var good_del_results = [...]string{delkey, nullvalue, delkeys, delbucket}

//bad delete procedures - error of certain types are expected... 
var delnonekey = "DEL zero bucket >>>> nonkey"
var delnonebucket = "DEL zero bucket"
var delnonebuckettwo = "DEL bucket one >> zero bucket two"

var delnonekeytwo = "DEL bucket one >>>> nonkey"      //silent fail of non-existent key is all BoltDB does
var bucket_nonekey = []string{"bucket one"}

var bad_del_results = map[string]error{
   delnonekey: err_nil_bucket,
   delnonekeytwo: nil,     //we only get a silent fail, test may have little value, but it's here...
   delnonebucket: bolt.ErrBucketNotFound,
   delnonebuckettwo: bolt.ErrBucketNotFound,
}

//---------------------------------------------------------------------------//

//test get procedures
var get_test1 = "GET bucket one >> bucket two >> bucket three >>>> test1"
var get_test2 = "GET bucket one >> bucket two >> bucket three >>>> test2"
var get_bucket_three = "GET bucket one >> bucket two >> bucket three"
var get_bucket_one = "GET bucket one"

var get_sole_results = map[string]map[string]string {
   get_test1: map[string]string{"test1": "value1"},
   get_test2: map[string]string{"test2": "value2"},
   get_bucket_three: map[string]string{"test1": "value1", "test2": "value2", "test3": "value3"},
   get_bucket_one: map[string]string{"bucket two": NESTEDBUCKET, "test6": "value6"},
}

//---------------------------------------------------------------------------//

//test rename procedures
var ren_tests = []string{
   "INS ren1 >> ren2 >> ren3 >>>> r1 :: v1",
   "INS ren1 >> ren2 >> ren3 >>>> r2 :: v2",
   "INS ren1 >> ren2 >> ren3 >> ren4",
   "INS ren1 >> ren2 >> ren3 >> ren4 >>>> r3 :: v3",
   "INS ren1 >> ren2 >> ren3 >> ren4 >>>> r4 :: v4",
   "INS ren1 >> ren2 >> ren3 >> ren4 >>>> r5 :: v5",
   "INS ren1 >> ren2 >> ren3 >>>> r6 :: v6",  
}

var ren_key = "REN ren1 >> ren2 >> ren3 >>>> r7 => renkey" 
var ren_bucket = "REN ren1 >> ren2 => renbuckets"

var ren_results = map[string]bool {
   
}

//---------------------------------------------------------------------------//

//example kvalresults
//a: kvalresult{map[string]string{"test1": "value1", "test2": "value2", "test3": "value3"}, false},
//b: kvalresult{map[string]string{"bucket two": NESTEDBUCKET, "test6": "value6"}, false},

//test list procedures
var lis_bucket_two = "LIS bucket one >> bucket two"
var lis_test1 = "LIS bucket one >> bucket two >> bucket three >>>> test1"
var lis_unknown_key = "LIS bucket one >> bucket two >> bucket three >>>> nokey"
var lis_unknown_bucket = "LIS ins1 >> ins2 >> no-bucket"

var lis_results = map[string]bool{
   lis_bucket_two: true,
   lis_test1: true,
   lis_unknown_key: false,
   lis_unknown_bucket: false,
}

//---------------------------------------------------------------------------//

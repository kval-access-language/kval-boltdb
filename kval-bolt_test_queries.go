package main

import "github.com/boltdb/bolt"

//test insert procedures
var ins_tests = []string{
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

type insresult struct {
   keys int
   depth int
}

var ins_result1 = insresult{3, 1}
var ins_result2 = insresult{6, 2}
var ins_result3 = insresult{8, 3}

//test get procedures
var get_test1 = "GET bucket one >> bucket two >> bucket three >>>> test1"
var get_test2 = "GET bucket one >> bucket two >> bucket three >>>> test2"
var get_bucket_three = "GET bucket one >> bucket two >> bucket three"
var get_bucket_one = "GET INS bucket one"

var get_results = map[string]bool {

}

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

//test delete procedures (reinstate all state before each test 'doinserts()')
//good procedures - nil error expected... 
var delkey = "DEL bucket one >> bucket two >> bucket three >>>> test1"           //delete key test1
var nullvalue = "DEL bucket one >> bucket two >> bucket three >>>> test3 :: _"   //make value null without deleting key
var delkeys = "DEL bucket one >> bucket two >> bucket three >>>> _"              //del all keys from a bucket
var delbucket = "DEL bucket one >> bucket two"                                   //delete bucket two        

//results that shouldn't cause an error...
var good_del_results = [...]string{delkey, nullvalue, delkeys, delbucket}

//bad procedures - error of certain types are expected... 
var delnonekey = "DEL zero bucket >>>> nonkey"
var delnonekeytwo = "DEL bucket one >>>> nonkey"      //silent fail of non-existent key is all BoltDB does
var delnonebucket = "DEL zero bucket"
var delnonebuckettwo = "DEL bucket one >> zero bucket two"

var bad_del_results = map[string]error{
   delnonekey: err_nil_bucket,
   delnonekeytwo: nil,     //we only get a silent fail, test may have little value, but it's here...
   delnonebucket: bolt.ErrBucketNotFound,
   delnonebuckettwo: bolt.ErrBucketNotFound,
}

//test list procedures
//run testins again
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


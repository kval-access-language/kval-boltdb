package kvalbolt

import "github.com/boltdb/bolt"

//test invalid/non-implemented capabilities

var make_tea = "TEA bucket one >> bucket two >>>> cup :: saucer"

//---------------------------------------------------------------------------//

//test insert procedures
var create_initial_state = []string{
	"ins bucket one >> bucket two >> bucket three >>>> test1 :: value1",
	"INS bucket one >> bucket two >> bucket three >>>> test2 :: value2",
	"INS bucket one >> bucket two >> bucket three >>>> test3 :: value3",
	"INS bucket one >> bucket two >>>> test4 :: value4",
	"INS bucket one >> bucket two >>>> test5 :: value5",
	"INS bucket one >>>> test6 :: value6",
	"INS code bucket >>>> code example :: GET bucket one >> bucket two >>>> key1 :: key2",
	"INS regex bucket >>>> regex example one :: middle regex string middle",
	"INS regex bucket >>>> regex example two :: regex string beginning beginning",
	"INS regex bucket >>>> regex example three :: end end regex string",
	"INS regex bucket >>>> regex example four :: regex shouldn't match",
	"INS regex bucket >>>> regex example five :: regex string",
	"INS regex bucket >> regex bucket two >>>> regex example six :: nil bucket test regex string",
}

var ins_getbuckets1 = []string{"bucket one", "bucket two", "bucket three"}
var ins_getbuckets2 = []string{"bucket one", "bucket two"}
var ins_getbuckets3 = []string{"bucket one"}

// Utilise BoltDB Tree statistics.
// KeyN  int // number of keys/value pairs
// Depth int // number of levels in B+tree

var ins_result1 = insresult{3, 1}
var ins_result2 = insresult{6, 2}
var ins_result3 = insresult{8, 3}

type ins_check struct {
	buckets []string
	counts  insresult
}

type insresult struct {
	keys  int
	depth int
}

var i1 = ins_check{ins_getbuckets1, ins_result1}
var i2 = ins_check{ins_getbuckets2, ins_result2}
var i3 = ins_check{ins_getbuckets3, ins_result3}

var ins_checks_all = [...]ins_check{i1, i2, i3}

//---------------------------------------------------------------------------//

//test delete procedures

//good delete procedures - nil error expected...
var delkey = "DEL bucket one >> bucket two >> bucket three >>>> test1"         //delete key test1
var nullvalue = "DEL bucket one >> bucket two >> bucket three >>>> test3 :: _" //make value null without deleting key
var delkeys = "del bucket one >> bucket two >> bucket three >>>> _"            //del all keys from a bucket
var delbucket = "DEL bucket one >> bucket two"                                 //delete bucket two

var good_del_results = [...]string{delkey, nullvalue, delkeys, delbucket}

//bad delete procedures - error of certain types are expected...
var delnonekey = "DEL zero bucket >>>> nonkey"
var delnonebucket = "DEL zero bucket"
var delnonebuckettwo = "DEL bucket one >> zero bucket two"

var delnonekeytwo = "DEL bucket one >>>> nonkey" //silent fail of non-existent key is all BoltDB does
var bucket_nonekey = []string{"bucket one"}

var bad_del_results = map[string]error{
	delnonekey:       err_nil_bucket,
	delnonekeytwo:    nil, //we only get a silent fail, test may have little value, but it's here...
	delnonebucket:    bolt.ErrBucketNotFound,
	delnonebuckettwo: bolt.ErrBucketNotFound,
}

//---------------------------------------------------------------------------//

//test get procedures
var get_test1 = "GET bucket one >> bucket two >> bucket three >>>> test1"
var get_test2 = "GET bucket one >> bucket two >> bucket three >>>> test2"
var get_bucket_three = "GET bucket one >> bucket two >> bucket three"
var get_bucket_one = "GET bucket one"
var get_code_bucket = "GET code bucket >>>> code example"

var get_sole_results = map[string]map[string]string{
	get_test1:        {"test1": "value1"},
	get_test2:        {"test2": "value2"},
	get_bucket_three: {"test1": "value1", "test2": "value2", "test3": "value3"},
	get_bucket_one:   {"bucket two": NESTEDBUCKET, "test6": "value6"},
	get_code_bucket:  {"code example": "GET bucket one >> bucket two >>>> key1 :: key2"},
}

//---------------------------------------------------------------------------//

//test get regex procedures
//GET Prime Bucket >> Secondary Bucket >> Tertiary Bucket >>>> {PAT}
//GET Prime Bucket >> Secondary Bucket >> Tertiary Bucket >>>> _ :: Value
//GET Prime Bucket >> Secondary Bucket >> Tertiary Bucket >>>> _ :: {PAT}
var get_regex_test1 = "GET bucket one >> bucket two >> bucket three >>>> {^test\\d+$}"
var get_regex_test2 = "GET bucket one >> bucket two >> bucket three >>>> _ :: value3"
var get_regex_test3 = "GET regex bucket >>>> _ :: {regex string}"

var regex_res1 = map[string]string{"test1": "value1", "test2": "value2", "test3": "value3"}
var regex_res2 = map[string]string{"test3": "value3"}
var regex_res3 = map[string]string{"regex example one": "middle regex string middle", "regex example two": "regex string beginning beginning", "regex example three": "end end regex string", "regex example five": "regex string"}

var get_regex_results = map[string]map[string]string{
	//var get_sole_results = map[string]map[string]string {
	get_regex_test1: regex_res1,
	get_regex_test2: regex_res2,
	get_regex_test3: regex_res3,
}

//---------------------------------------------------------------------------//

//example Kvalresults
//a: Kvalresult{map[string]string{"test1": "value1", "test2": "value2", "test3": "value3"}, false},
//b: Kvalresult{map[string]string{"bucket two": NESTEDBUCKET, "test6": "value6"}, false},

//test list procedures
var lis_bucket_two = "LIS bucket one >> bucket two"
var lis_test1 = "LIS bucket one >> bucket two >> bucket three >>>> test1"
var lis_unknown_key = "LIS bucket one >> bucket two >> bucket three >>>> nokey"
var lis_unknown_bucket = "LIS ins1 >> ins2 >> no-bucket"

var lis_results = map[string]bool{
	lis_bucket_two:     true,
	lis_test1:          true,
	lis_unknown_key:    false,
	lis_unknown_bucket: false,
}

//---------------------------------------------------------------------------//

//test rename procedures
var rename_state = []string{
	"INS ren1 >> ren2 >> ren3 >>>> r1 :: v1",         //2
	"INS ren1 >> ren2 >> ren3 >>>> r2 :: v2",         //3
	"INS ren1 >> ren2 >> ren3 >> ren4",               //4
	"INS ren1 >> ren2 >> ren3 >> ren4 >>>> r3 :: v3", //5
	"INS ren1 >> ren2 >> ren3 >> ren4 >>>> r4 :: v4", //6
	"INS ren1 >> ren2 >> ren3 >> ren4 >>>> r5 :: v5", //7
	"INS ren1 >> ren2 >> ren3 >>>> r6 :: v6",         //8
	"INS ren1 >> ren2 >>>> r6 :: v6",                 //9
	"INS ren1 >> ren2 >>>> r7 :: v6",                 //10
	"INS ren1 >> ren2 >>>> r8 :: v6",                 //11
	"INS ren1 >> ren2 >>>> r1 :: v1",                 //12
	"INS ren1 >> renamekey >>>> key :: value",        //key to rename...
}

var r1 = "ren_key"
var r2 = "ren_bucket"

//though few, this should prove our capability adequately...
var rename_tests = map[string]string{
	r1: "REN ren1 >> renamekey >>>> key => newkey", //rename key
	r2: "REN ren1 >> ren2 => rnew",                 //rename bucket
}

//FALSE :: TRUE if rename has worked, we'll see true for second value
var ren_lis1 = [2]string{"LIS ren1 >> renamekey >>>> key", "LIS ren1 >> renamekey >>>> newkey"}
var ren_lis2 = [2]string{"LIS ren1 >> ren2", "LIS ren1 >> rnew"}

//grab stats dynamically as well
var ren_slice_old = []string{"ren1", "ren2"}
var ren_slice_new = []string{"ren1", "rnew"}

//---------------------------------------------------------------------------//

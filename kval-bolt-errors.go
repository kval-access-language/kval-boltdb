package kvalbolt

import "github.com/pkg/errors"

var err_nil_bucket = errors.New("Cannot GOTO bucket, bucket not found")
var err_empty_bucket_slice = errors.New("Cannot GOTO bucket, empty buckets slice provided")
var err_not_implemented = errors.New("KVAL Function not implemented")
var err_no_kv_in_bucket = errors.New("No Keys: There are no key::value pairs in this bucket")

var err_blob_key = errors.New("No Key: attempting to add blob but key value is empty or '_'")
var err_blob_val = errors.New("Value added: attempting to add blob but have specified value")
var err_blob_ins = errors.New("INS Only: Can only use INS to PUT blob")
var err_blob_len = errors.New("Blob data supplied is not a blob, or is wrapped incorrectly")
var err_blob_map_len = errors.New("Maplen not equal to one. May be multiple value context, or zero")

var err_stat = errors.New("Cannot stat database. Stats struct will remain empty.")

//Other non-error error strings...
var err_parse = "Query parse failed"


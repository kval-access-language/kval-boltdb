package main

import "github.com/pkg/errors"

var err_nil_bucket = errors.New("Cannot GOTO bucket, bucket not found")
var err_empty_bucket_slice = errors.New("Cannot GOTO bucket, empty buckets slice provided")
var err_not_implemented = errors.New("KVAL Function not implemented")

//Other non-error error strings...
var err_parse = "Query parse failed"


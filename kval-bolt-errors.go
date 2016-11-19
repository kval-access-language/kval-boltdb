package main

import "github.com/pkg/errors"

var err_nil_bucket = errors.New("Cannot GOTO bucket, bucket not found.")
var err_empty_bucket_slice = errors.New("Cannot GOTO bucket, empty buckets slice provided.")


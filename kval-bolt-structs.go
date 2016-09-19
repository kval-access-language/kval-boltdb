package main

import (
   "github.com/boltdb/bolt"   
   "github.com/kval-access-language/KVAL-Parse"
)

const NESTEDBUCKET = "NestedBucket"

type kvalbolt struct {
   db       *bolt.DB
   fname    string
   query    kvalparse.KQUERY
}

type kvalresult struct {
   Result   map[string]string
   Exists   bool                 //IF LIS QEURY, however, idiomatic enough?
}

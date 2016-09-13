package main

import (
   "github.com/boltdb/bolt"   
   "github.com/kval-access-language/KVAL-Parse"
)

//nil result for nil comparisons with kvalresult struct
var nilresult kvalresult

type kvalbolt struct {
   db       *bolt.DB
   fname    string
   query    kvalparse.KQUERY
}

type kvalresult struct {
   String      string
}

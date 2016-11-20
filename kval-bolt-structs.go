package main

import (
   "github.com/boltdb/bolt"   
   "github.com/kval-access-language/kval-parse"
)

const NESTEDBUCKET = "NestedBucket"
const DATA = "data"
const BASE64 = "base64"

type kvalbolt struct {
   db       *bolt.DB
   fname    string
   query    kvalparse.KQUERY
}

type kvalresult struct {
   Result   map[string]string
   Exists   bool                 //IF LIS QEURY, however, idiomatic enough?
}

type kvalblob struct {
   query       string
   datatype    string
   mimetype    string
   encoding    string
   data        string
}

func initkvalblob(query string, mime string, data string) kvalblob {
   return kvalblob{query, DATA, mime, BASE64, data}
}

func queryfromkvb(kvb kvalblob) string {
   query := kvb.query + " :: " + kvb.datatype + ":" + kvb.mimetype + ":" + kvb.encoding + ":" + kvb.data
   return query
}

package main

import (
   "strings"
   "github.com/boltdb/bolt"   
   "github.com/kval-access-language/kval-parse"
)

const NESTEDBUCKET = "NestedBucket"
const DATA = "data"
const BASE64 = "base64"
const BLOBLEN = 4          //"data:<mimetype>:<encoding type>:<data>"

type kvalbolt struct {
   db       *bolt.DB
   fname    string
   query    kvalparse.KQUERY
}

type kvalresult struct {
   Result   map[string]string
   Exists   bool   //If LIS query then we just need a flag to say if a value is there...                
}

type kvalblob struct {
   Query       string
   Datatype    string
   Mimetype    string
   Encoding    string
   Data        string
}

func initkvalblob(query string, mimetype string, data string) kvalblob {
   return kvalblob{query, DATA, mimetype, BASE64, data}
}

func queryfromkvb(kvb kvalblob) string {
   query := kvb.Query + " :: " + kvb.Datatype + ":" + kvb.Mimetype + ":" + kvb.Encoding + ":" + kvb.Data
   return query
}

func blobfromkvalresult(kv kvalresult) (kvalblob, error) {
   var kvb kvalblob
   for k, v := range(kv.Result) {
      kvb.Query = k
      reslice := strings.Split(v, ":")
      if len(reslice) != 4 {
         return kvb, err_blob_len
      }
      kvb.Datatype = reslice[0]
      kvb.Mimetype = reslice[1]
      kvb.Encoding = reslice[2]
      kvb.Data     = reslice[3]
   }   
   return kvb, nil
}

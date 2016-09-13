package main

import (
   "github.com/boltdb/bolt"   
   "github.com/kval-access-language/KVAL-Parse"
)

type kvalbolt struct {
   db       *bolt.DB
   fname    string
   query    kvalparse.KQUERY
}

type kvalresult struct {
   Result   map[string]string
}

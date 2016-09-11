package main

import (
   "github.com/boltdb/bolt"   
)

type kvalbolt struct {
   db *bolt.DB
}

package main

import (
   "fmt"
   "github.com/boltdb/bolt"   
)

func createboltbuckets(kb kvalbolt) (kvalresult, error) {

   var kr kvalresult
   var kq = kb.query

   fmt.Println("creating buckets")

   err := kb.db.Update(func(tx *bolt.Tx) error {
      var bucket *bolt.Bucket 
      var err error
      for index, bucketname := range(kq.Buckets) {
         if index == 0 {
            bucket, err = tx.CreateBucketIfNotExists([]byte(bucketname))   
            if err != nil {
               return err
            }    
         } else{
            bucket, err = bucket.CreateBucketIfNotExists([]byte(bucketname))
            if err != nil {
               return err
            }        
         }

      }
      return nil
   })

   if err != nil {
      return kr, err
   }

   return kr, nil
}
package main

import (
   "fmt"
   "github.com/boltdb/bolt"
)

func createboltentries(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {
      var bucket *bolt.Bucket    //we only ever need the 'last' bucket in memory
      var err error
      //create buckets
      for index, bucketname := range kq.Buckets {
         if index == 0 {
            bucket, err = tx.CreateBucketIfNotExists([]byte(bucketname))
            if err != nil {
               return err
            }
         } else {
            bucket, err = bucket.CreateBucketIfNotExists([]byte(bucketname))
            if err != nil {
               return err
            }
         }
      }
      //func (b *Bucket) Put(key []byte, value []byte) error
      //create key::values
      if kq.Key != "" {         
         if kq.Value != "" {
            //write value...
           err = bucket.Put([]byte(kq.Key), []byte(kq.Value))
         } else {
            //write blank value if allowed... (UC: User may want to know unknown)
            err = bucket.Put([]byte(kq.Key), []byte(""))
         }
         if err != nil {
            return err
         }
      }

      //commit transaction
      return nil
   })
   if err != nil {
      return err
   }
   return nil
}

func viewboltentries(kb kvalbolt) (kvalresult, error) {
   var kr kvalresult
   var kq = kb.query
   err := kb.db.View(func(tx *bolt.Tx) error {
      var bucket *bolt.Bucket
      for index, bucketname := range kq.Buckets {
         if index == 0 {
            bucket = tx.Bucket([]byte(bucketname)) 
            if bucket == nil {
               return fmt.Errorf("Nil Bucket: Bucket does not exist", "\n")
            }
         } else {
            bucket = bucket.Bucket([]byte(bucketname))
            if bucket == nil {
               return fmt.Errorf("Nil Bucket: Bucket does not exist", "\n")
            }
         }
      }   
      if bucket != nil {
         val := bucket.Get([]byte(kq.Key))
         kr.Res = string(val)
      } 
      return nil
   })
   if err != nil {
      return kr, err
   }
   return kr, nil
} 

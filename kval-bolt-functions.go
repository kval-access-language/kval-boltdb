package main

import (
   "fmt"
   "github.com/boltdb/bolt"
)

func initkvalresult() (kvalresult) {
   kr := kvalresult{
      Result: map[string]string{},
   }
   return kr
}

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
   return err
}

func viewboltentries(kb kvalbolt) (kvalresult, error) {
   var kr = initkvalresult()
   var kq = kb.query
   err := kb.db.View(func(tx *bolt.Tx) error {
      bucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      if bucket != nil {
         val := bucket.Get([]byte(kq.Key))
         kr.Result[kq.Key] = string(val)
      }
      //commit transaction 
      return nil
   })
   return kr, err
} 

func getallfrombucket(kb kvalbolt) (kvalresult, error) {
   var kq = kb.query
   var kr = initkvalresult()
   err := kb.db.View(func(tx *bolt.Tx) error {
      bucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      if bucket != nil {
         bs := bucket.Stats()
         if bs.KeyN > 0 {
            cursor := bucket.Cursor()
            k,v := cursor.First()
            for k != nil {
               if v == nil {
                  kr.Result[string(k)] = "Is a nested Bucket."
               } else {
                  kr.Result[string(k)] = string(v)
               }

               k, v = cursor.Next()
            }
         } else {
            return fmt.Errorf("No Keys: There are no key :: value pairs in this bucket.")
         }
      }
      //commit transaction
      return nil
   })
   return kr, err
}

func deletebucket(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {
      var searchindex = len(kq.Buckets)-1      
      bucket, err := gotobucket(tx, kq.Buckets[:searchindex])
      if err != nil {
         return err
      }
      err = bucket.DeleteBucket([]byte(kq.Buckets[len(kq.Buckets)-1]))
      if err != nil {
         return err
      }
      return nil
   })
   return err
}

func gotobucket(tx *bolt.Tx, bucketslice []string) (*bolt.Bucket, error) {
   var bucket *bolt.Bucket
   for index, bucketname := range bucketslice {
      if index == 0 {
         bucket = tx.Bucket([]byte(bucketname)) 
         if bucket == nil {
            return bucket, fmt.Errorf("Nil Bucket: Bucket does not exist", "\n")
         }
      } else {
         bucket = bucket.Bucket([]byte(bucketname))
         if bucket == nil {
            return bucket, fmt.Errorf("Nil Bucket: Bucket does not exist", "\n")
         }
      }
   }   
   return bucket, nil
}
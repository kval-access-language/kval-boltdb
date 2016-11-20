package kvalbolt

import (
   "fmt"
   "regexp"
   "github.com/pkg/errors"
   "github.com/boltdb/bolt"
)

func initkvalresult() (kvalresult) {
   kr := kvalresult{
      Result: map[string]string{},
   }
   return kr
}

//Create bucket and key/value entries in BoltDB from a kval structure
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

//Retrieve an entry from a BoltDB from a kval structure
func getboltentry(kb kvalbolt) (kvalresult, error) {
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

//Retrieve an entry from a BoltDB using regular expression
func getboltkeyregex(kb kvalbolt) (kvalresult, error) {
   var kq = kb.query
   var kr = initkvalresult()
   re, err := regexp.Compile(kq.Value)
   if err != nil {
      return kr, err
   }
   err = kb.db.View(func(tx *bolt.Tx) error {
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
               if re.MatchString(string(k)) {
                  if v == nil {
                     kr.Result[string(k)] = NESTEDBUCKET
                  } else {
                     kr.Result[string(k)] = string(v)
                  }
               }
               k, v = cursor.Next()
            }
         } else {
            return err_no_kv_in_bucket
         }
      }
      //commit transaction
      return nil
   })
   return kr, nil
} 

//Retrieve an entry from a BoltDB using regular expression
func getboltvalueregex(kb kvalbolt) (kvalresult, error) {
   var kq = kb.query
   var kr = initkvalresult()
   re, err := regexp.Compile(kq.Value)
   if err != nil {
      return kr, err
   }
   err = kb.db.View(func(tx *bolt.Tx) error {
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
               if v != nil {
                  //nil means nested bucket: can't work with nested buckets for this search
                  if re.MatchString(string(v)) {
                     kr.Result[string(k)] = string(v)
                  }
               }
               k, v = cursor.Next()
            }
         } else {
            return err_no_kv_in_bucket
         }
      }
      //commit transaction
      return nil
   })
   return kr, nil
} 

//Retrieve all values from a single bucket per KVAL syntax
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
                  kr.Result[string(k)] = NESTEDBUCKET
               } else {
                  kr.Result[string(k)] = string(v)
               }
               k, v = cursor.Next()
            }
         } else {
            return err_no_kv_in_bucket
         }
      }
      //commit transaction
      return nil
   })
   return kr, err
}

//Delete a single bucket from a BoltDB from a KVAL structure
func deletebucket(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {
      //as we're deleting a bucket we need a pointer to 
      //bucket level we're deleting minus one, that is
      //the container of the bucket we're deleting
      var delname = kq.Buckets[len(kq.Buckets)-1]
      var searchindex = len(kq.Buckets)-1
      if searchindex == 0 {
         //reset to one? this is the bucket we're deleting
         searchindex = 1
         delname = kq.Buckets[0]
         err := tx.DeleteBucket([]byte(delname))
         if err != nil {
            return errors.Wrapf(err, "Bucket name: '%s'", delname)
         }
      } else {
         bucketname := kq.Buckets[:searchindex]
         bucket, err := gotobucket(tx, bucketname)
         if err != nil {
            return err
         }
         err = bucket.DeleteBucket([]byte(delname))
         if err != nil {
            return errors.Wrapf(err, "Bucket name: '%s'", delname)
         }
      }
      return nil
   })
   return err
}

//Delete all the keys in a BoltDB bucket leaving Bucket in tact
func deletebucketkeys(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {     
      bucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      cursor := bucket.Cursor()
      k,_ := cursor.First()
      for k != nil {
         err := bucket.Delete(k)
         if err != nil {
            if err == bolt.ErrIncompatibleValue {
               //likely we're trying to delete a nested bucket
               err = bucket.DeleteBucket(k)
               if err != nil {
                  return err
               }
            } else {
               return err
            }
         }
         k, _ = cursor.Next()
      }     
      return err 
   })
   return err   
}

//Delete a key and its corresponding value from a BoltDB
func deletekey(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {   
      bucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      err = bucket.Delete([]byte(kb.query.Key))
      if err != nil {
         if err == bolt.ErrIncompatibleValue {
            //likely we're trying to delete a nested bucket
            err = bucket.DeleteBucket([]byte(kb.query.Key))
            if err != nil {
               return err
            }
         } else {
            return err
         }
      }
      return err
   }) 
   return err  
}

//Turn a value for a given key to NULL based on KVAL capabilities
func nullifykeyvalue(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {   
      bucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      err = bucket.Put([]byte(kq.Key), []byte(""))
      if err != nil {
         return err
      }
      return err
   })
   return err
}

//Rename a bucket (full OR empty) in a BoltDB
func renamebucket(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {   
      //the bucket containing the one we're renaming
      var searchindex = len(kq.Buckets)-1      
      containerbucket, err := gotobucket(tx, kq.Buckets[:searchindex])
      if err != nil {
         return err
      }
      //the bucket we're renaming      
      oldbucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      //gotta create the new bucket here...
      newbucket, err := containerbucket.CreateBucketIfNotExists([]byte(kq.Newname))
      if err != nil {
         return err
      }      
      err = copybuckets(oldbucket, newbucket)
      if err != nil {
         return err
      }
      //delete the origial bucket
      oldname := []byte(kq.Buckets[len(kq.Buckets)-1:][0])
      err = containerbucket.DeleteBucket(oldname)
      if err != nil {
         return err
      }
      //complete the transaction
      return nil
   })   
   return err
}

//Helper function for rename to copy a bucket to a newly named bucket
func copybuckets(from, to *bolt.Bucket) error {
   bs := from.Stats()
   if bs.KeyN > 0 {
      cursor := from.Cursor()
      k,v := cursor.First()
      for k != nil {
         if v == nil {
            //nested bucket 
            to_nested, err := to.CreateBucketIfNotExists(k)
            if err != nil {
               return err
            }
            from_nested := from.Bucket(k)
            copybuckets(from_nested, to_nested)
         } else {
            to.Put(k,v)
         }
         k, v = cursor.Next()
      }
   } else {
      return err_no_kv_in_bucket
   }
   return nil
}

//Rename a key in a BoltDB based on described KVAL capabilities such as rename
func renamekey(kb kvalbolt) error {
   var kq = kb.query
   err := kb.db.Update(func(tx *bolt.Tx) error {   
      //the bucket containing the key we're renaming
      bucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      v := bucket.Get([]byte(kq.Key))
      if v == nil {
         return fmt.Errorf("Nil Value: Key doesn't exist or points to a nested bucket.")
      }
      err = bucket.Put([]byte(kq.Newname), v)
      if err != nil {
         return err
      }
      err = bucket.Delete([]byte(kq.Key))
      if err != nil {
         return err
      }
      return nil
   })

   return err
}

//Check to see if a key exists in a BoltDB bucket, per KVAL LIS capabilities
func bucketkeyexists(kb kvalbolt) (kvalresult, error) {
   var kq = kb.query
   var kr = initkvalresult()
   err := kb.db.Update(func(tx *bolt.Tx) error {   
      //the bucket containing the key we're renaming
      bucket, err := gotobucket(tx, kq.Buckets)
      if err != nil {
         return err
      }
      if kq.Key != "" {
         k := bucket.Get([]byte(kq.Key))         
         if k == nil {
            return fmt.Errorf("Key '%s' does not exist.", kq.Key)
         }
      }
      return nil
   })
   if err == nil {
      kr.Exists = true
   }
   return kr, nil   
}

//Retrieve a bucket pointer to work with from the BoltDB
func gotobucket(tx *bolt.Tx, bucketslice []string) (*bolt.Bucket, error) {
   var bucket *bolt.Bucket
   if len(bucketslice) > 0 {
      for index, bucketname := range bucketslice {
         if index == 0 {   //need a bucket from our transaction pointer first
            bucket = tx.Bucket([]byte(bucketname)) 
            if bucket == nil {   //only ever get nil if our root bucket doesn't exist
               return bucket, err_nil_bucket
            }
            if len(bucketslice) == 1 && bucket != nil {
               //return early, we've got out bucket
               return bucket, nil
            }
         } else {   //nested buckets, only returning if nil...
            bucket = bucket.Bucket([]byte(bucketname))
            if bucket == nil {
               return bucket, err_nil_bucket
            }
         }
      }   
   } else {
      //gold plating at this point, easily handled elsewhere...
      return bucket, err_empty_bucket_slice
   }
   return bucket, nil
}


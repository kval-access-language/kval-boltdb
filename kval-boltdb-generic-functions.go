package kvalbolt

import "github.com/boltdb/bolt"

// Constant values for statdb to be able to externalize bucket stats
// for anyone using kvalresults structs without having to go back to bolt ptr
const (
	opcodeNormal       int = iota // opcodeNormal refers to LIS, GET, INS, and certain REN/DEL functions
	opcodeDelBucket               // delete functions that delete a bucket thus can't be handled as easily
	opcodeRenameBucket            // rename functions that rename a bucket thus can't be handled as easily
)

//get boltdb bucket stats
func getbucketstats(kb Kvalboltdb, buckets []string) (bolt.BucketStats, error) {
	var bs bolt.BucketStats
	err := kb.DB.View(func(tx *bolt.Tx) error {
		bucket, err := gotobucket(tx, buckets)
		if err != nil {
			return err
		}
		bs = bucket.Stats()
		return err
	})
	return bs, err
}

func statdb() {

	//switch(

}

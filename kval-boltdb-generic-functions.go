package kvalbolt

import "github.com/boltdb/bolt"

//get boltdb bucket stats
func getbucketstats(kb kvalbolt, buckets []string) (bolt.BucketStats, error) {
	var bs bolt.BucketStats
	err := kb.db.View(func(tx *bolt.Tx) error {
		bucket, err := gotobucket(tx, buckets)
		if err != nil {
			return err
		}
		bs = bucket.Stats()
		return err
	})
	return bs, err
}

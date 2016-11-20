# KVAL-BoltDB

[![Build Status](https://travis-ci.org/kval-access-language/kval-boltdb.svg?branch=master)](https://travis-ci.org/kval-access-language/kval-boltdb)
[![GoDoc](https://godoc.org/github.com/kval-access-language/kval-boltdb?status.svg)](https://godoc.org/github.com/kval-access-language/kval-boltdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/kval-access-language/kval-boltdb)](https://goreportcard.com/report/github.com/kval-access-language/kval-boltdb)

BoltDB bindings for KVAL
 
###Key Value Access Language

I have created a modest specification for a key value access langauge. 
It allows for input and access of values to a key value store such as Golang's
[BoltDB](https://github.com/boltdb/). 

The language specification: https://github.com/kval-access-language/KVAL 

###Features 

* Single function entry-point:
    * res, err := Query(INS B1 >> B2 >> B3 >>>> KEY :: VAL) &nbsp; &nbsp; //(will create three buckets, plus k/v in one-go)
    * res, err := Query(GET B1 >> B2 >> B3 >>>> KEY) &nbsp; &nbsp; &nbsp; //(will retrieve that entry in one-go)
* Start using BoltDB immediately without writing partial wrappers for your code
* KVAL-Parse enables handling of Base64 binary BLOBS
* Regular Expression based searching for key names and values
* [KVAL Language](https://github.com/kval-access-language/KVAL) specifies easy bulk or singleton DEL and RENAME capabilities
* Language specification at highest abstraction, so other bindings for other DBs are hoped for (hint: [NOMS!](https://github.com/attic-labs/noms)) 

###Usage

Use is simple. There is one function which accepts a string formatted to KVAL's
specification:

    res, err = Query(kb, "GET Bucket One >> Bucket Two >>>> Requested Key")
    if err != nil {
       fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
    } else {
       //Access our (result structure)[https://github.com/kval-access-language/kval-boltdb/blob/master/kval-bolt-structs.go#L16]: res.Result (a map[string]string)
    } 

For write operations we simply check for the existence of an error, else the
operation passed as expected: 

    res, err = Query(kb, "INS Bucket One >> Bucket Two >>>> Insert Key :: Insert Value")
    if err != nil {
       fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
    }

###License

**[GPL Version 3](http://choosealicense.com/licenses/gpl-3.0/)**: https://github.com/kval-access-language/KVAL-BoltDB/blob/master/LICENSE

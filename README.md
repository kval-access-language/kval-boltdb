# KVAL-BoltDB

BoltDB bindings for KVAL

###Key Value Access Language

I have created a modest specification for a key value access langauge. 
It allows for input and access of values to a key value store such as Golang's
[BoltDB](https://github.com/boltdb/). 

The language specification: https://github.com/kval-access-language/KVAL 

###Usage

Use is simple. There is one function which accepts a string formatted to KVAL's
specification:

    res, err = Query(kb, "GET Bucket One >> Bucket Two >>>> Requested Key")
    if err != nil {
       fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
    } else {
       //Access our result structure: res.Result (a map[string]string)
    } 

For write operations we simply check for the existence of an error, else the
operation passed as expected: 

    res, err = Query(kb, "INS Bucket One >> Bucket Two >>>> Insert Key :: Insert Value")
    if err != nil {
       fmt.Fprintf(os.Stderr, "Error querying db: %v", err)
    }

###License

**[http://choosealicense.com/licenses/gpl-3.0/](GPL Version 3)**: https://github.com/kval-access-language/KVAL-BoltDB/blob/master/LICENSE

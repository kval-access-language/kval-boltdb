package kvalbolt

import (
	b64 "encoding/base64"
	"github.com/boltdb/bolt"
	"github.com/kval-access-language/kval-parse"
	"github.com/kval-access-language/kval-scanner"
	"github.com/pkg/errors"
	"time"
)

//Open a BoltDB with a given name to work with.
//Our first most important function. Returns a KVAL Bolt structure
//with the details required for KBAL BoltDB to perform queries.
func Connect(dbname string) (Kvalbolt, error) {
	var kb kvalbolt
	db, err := bolt.Open(dbname, 0600, &bolt.Options{Timeout: 2 * time.Second})
	kb.db = db
	return kb, err
}

//Disconnect from a BoltDB.
//Recommended that this is made a deferred function call where possible.
func Disconnect(kb kvalbolt) {
	kb.db.Close()
}

//Retrieve a pointer to BoltDB at any time for working with it manually.
func GetBolt(kb kvalbolt) *bolt.DB {
	return kb.db
}

//Query. Given a KVALBolt Structure, and a KVAL query string
//this function will do all of the work for you when interacting with
//BoltDB. Everything should become less programmatic making for cleaner code.
//The KVAL spec can be found here: https://github.com/kval-access-language/kval
func Query(kb kvalbolt, query string) (Kvalresult, error) {
	var kr Kvalresult
	var err error
	kq, err := kvalparse.Parse(query)
	if err != nil {
		return kr, errors.Wrapf(err, "%s: '%s'", errParse, query)
	}
	kb.query = kq
	kr, err = queryhandler(kb)
	if err != nil {
		return kr, err
	}
	return kr, nil
}

//Wrap a blob of data, KVAL-Bolt/KVAL proposes a standard encoding for
//this data inside Key-Value databases, that goes like this:
//data:mimetype;base64;{base64 data}. Use Unwrap to get the datastream back
//location should be specified in the form of a query, e.g. INS bucket >>>> key
func StoreBlob(kb kvalbolt, loc string, mime string, data []byte) error {

	//Check location query parses correctly...
	kq, err := kvalparse.Parse(loc)
	if err != nil {
		return errors.Wrapf(err, "%s: '%s'", errParse, loc)
	}

	//Validate for certain features...
	if kq.Function != kvalscanner.INS {
		return errBlobIns
	} else if kq.Key == "" || kq.Key == "_" {
		return errBlobKey
	} else if kq.Value != "" {
		return errBlobVal
	}

	//Encode our data as base64
	encoded := b64.StdEncoding.EncodeToString([]byte(data))

	//Convert to known datatype and retrieve a standardised value from it
	kvb := initKvalblob(loc, mime, encoded)
	query := queryfromkvb(kvb)

	//Check our new query including base64 string validates okay
	kq, err = kvalparse.Parse(query)
	if err != nil {
		return errors.Wrapf(err, "%s: '%s'", errParse, query)
	}

	//Finally... do the rest of the work with one of our other Exported functions
	_, err = Query(kb, query)
	return err
}

//If you retrieve a blob via GET, unwrap it here to see what
//you asked for...
func UnwrapBlob(kv Kvalresult) (Kvalblob, error) {
	var kvb Kvalblob
	if len(kv.Result) != 1 {
		return kvb, errBlobMapLen
	}
	kvb, err := blobfromKvalresult(kv)
	return kvb, err
}

//Abstracted away from Query() query handler is an unexported function that
//will route all queries as required by the application when given by the user.
func queryhandler(kb kvalbolt) (Kvalresult, error) {
	var kr Kvalresult
	switch kb.query.Function {
	case kvalscanner.INS:
		err := insHandler(kb)
		return kr, err
	case kvalscanner.GET:
		if kb.query.Key == "" {
			//get all
			kr, err := getallHandler(kb)
			return kr, err
		} else if kb.query.Regex {
			kr, err := getregexHandler(kb)
			return kr, err
		} else {
			kr, err := getHandler(kb)
			return kr, err
		}
	case kvalscanner.LIS:
		kr, err := lisHandler(kb)
		return kr, err
	case kvalscanner.DEL:
		if kb.query.Key == "" {
			//we're deleting a bucket (and all contents)
			err := delbucketHandler(kb)
			return kr, err
		} else if kb.query.Key == "_" {
			//we're making nil "" values for all keys
			//use case, we want the keys, we don't want the values
			err := delbucketkeysHandler(kb)
			return kr, err
		} else if kb.query.Key != "" && kb.query.Key != "_" && kb.query.Value != "_" {
			//we're deleting a key and its value
			err := delonekeyHandler(kb)
			return kr, err
		} else if kb.query.Value == "_" {
			//we're deleting a value and leaving the key
			err := nullifyvalHandler(kb)
			return kr, err
		}
	case kvalscanner.REN:
		if kb.query.Key == "" {
			renbucketHandler(kb)
		} else if kb.query.Key != "" {
			renkeyHandler(kb)
		}
	default:
		//function is parsed correctly but not recognised by binding
		return kr, errors.Wrapf(errNotImplemented, "%v", kb.query.Function)
	}
	return kr, nil
}

//INS (Insert Handler) handles INS capability of KVAL language
func insHandler(kb kvalbolt) error {
	//as long as there are buckets, we can create anything we need.
	//it all happens in a single transaction, based on kval query...
	err := createboltentries(kb)
	if err != nil {
		return err
	}
	return nil
}

//GET (Get Handler) handles GET capability of KVAL language
func getHandler(kb kvalbolt) (Kvalresult, error) {
	if kb.query.Key == "_" {
		//turn our value into a regular expression for better search
		kb.query.Value = "^" + kb.query.Value + "$"
		return getregexHandler(kb)
	}
	var kr Kvalresult
	kr, err := getboltentry(kb)
	if err != nil {
		return kr, err
	}
	return kr, nil
}

//GET (Get Handler) handles GET (ALL) capability of KVAL language
func getallHandler(kb kvalbolt) (Kvalresult, error) {
	kr, err := getallfrombucket(kb)
	if err != nil {
		return kr, err
	}
	return kr, nil
}

//GET (Get Handler) handles GET (REGEX) capability of KVAL language
func getregexHandler(kb kvalbolt) (Kvalresult, error) {
	var kr Kvalresult
	var err error
	if kb.query.Value == "" {
		kr, err = getboltkeyregex(kb)
		if err != nil {
			return kr, err
		}
	} else if kb.query.Key != "" && kb.query.Value != "" {
		kr, err = getboltvalueregex(kb)
		if err != nil {
			return kr, err
		}
	}
	return kr, nil
}

//DEL (Delete Handler) handles DEL bucket capability of KVAL language
func delbucketHandler(kb kvalbolt) error {
	err := deletebucket(kb)
	if err != nil {
		return err
	}
	return nil
}

//DEL (Delete Handler) handles DEL all keys capability of KVAL language
func delbucketkeysHandler(kb kvalbolt) error {
	err := deletebucketkeys(kb)
	if err != nil {
		return err
	}
	return nil
}

//DEL (Delete Handler) handles DEL one key capability of KVAL language
func delonekeyHandler(kb kvalbolt) error {
	err := deletekey(kb)
	if err != nil {
		return err
	}
	return nil
}

//DEL (Delete Handler) Handles DEL (or in this case, NULL, capability of KVAL
func nullifyvalHandler(kb kvalbolt) error {
	err := nullifykeyvalue(kb)
	if err != nil {
		return err
	}
	return nil
}

//REN (Rename Handler) Handles rename bucket capability of KVAL
func renbucketHandler(kb kvalbolt) error {
	err := renamebucket(kb)
	if err != nil {
		return err
	}
	return nil
}

//REN (Rename Handler) Handles rename key capability of KVAL
func renkeyHandler(kb kvalbolt) error {
	err := renamekey(kb)
	if err != nil {
		return err
	}
	return nil
}

//LIS (List Handler) Handles listing capability of KVAL (does (x) exist?)
func lisHandler(kb kvalbolt) (Kvalresult, error) {
	kr, err := bucketkeyexists(kb)
	if err != nil {
		//Nil bucket returns an error we can use
		if kr.Exists == true {
			return kr, err
		}
		return kr, err
	}
	return kr, nil
}

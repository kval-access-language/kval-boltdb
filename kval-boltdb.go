package kvalbolt

import (
	b64 "encoding/base64"
	"github.com/boltdb/bolt"
	"github.com/kval-access-language/kval-parse"
	"github.com/kval-access-language/kval-scanner"
	"github.com/pkg/errors"
	"time"
)

// Connect should first be used to open a connection to a BoltDB with a given
// name. Returns a KVAL Bolt structure with the details required for future
// KVAL BoltDB operations.
func Connect(dbname string) (Kvalbolt, error) {
	var kb Kvalbolt
	db, err := bolt.Open(dbname, 0600, &bolt.Options{Timeout: 2 * time.Second})
	kb.db = db
	return kb, err
}

// Disconnect lets us disconnect from a BoltDB. It is recommended that this
// function is deffered where possible.
func Disconnect(kb Kvalbolt) {
	kb.db.Close()
}

// GetBolt retrieves a pointer to BoltDB at any time for working with it manually.
func GetBolt(kb Kvalbolt) *bolt.DB {
	return kb.db
}

// Query is our primary function once we've opened a BoltDB connection.
// Given a KVALBolt Structure, and a KVAL query string
// this function will do all of the work for you when interacting with
// BoltDB. Everything should become less programmatic making for cleaner code.
// The KVAL spec can be found here: https://github.com/kval-access-language/kval.
func Query(kb Kvalbolt, query string) (Kvalresult, error) {
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

// StoreBlob is used to wrap a blob of data.
// KVAL-Bolt/KVAL proposes a standard encoding for this data inside Key-Value
// databases, that goes like this: data:mimetype;base64;{base64 data}. Use
// Unwrap to get the datastream back and futher GetBlobdata as a shortcut to
// decode it from Base64.
// Location for StoreBlob should be specified in the form of a query:
// e.g. INS bucket >>>> key
func StoreBlob(kb Kvalbolt, loc string, mime string, data []byte) error {

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

// UnwrapBlob can be used to unwrap a blob you are retrieving via GET
func UnwrapBlob(kv Kvalresult) (Kvalblob, error) {
	var kvb Kvalblob
	if len(kv.Result) != 1 {
		return kvb, errBlobMapLen
	}
	kvb, err := blobfromKvalresult(kv)
	return kvb, err
}

// GetBlobData decodes the Base64 data stored in a Kvalblob object
func GetBlobData(kvb Kvalblob) ([]byte, error) {
	data, err := b64.StdEncoding.DecodeString(kvb.Data)
	return data, err
}

//Abstracted away from Query() query handler is an unexported function that
//will route all queries as required by the application when given by the user.
func queryhandler(kb Kvalbolt) (Kvalresult, error) {
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
func insHandler(kb Kvalbolt) error {
	//as long as there are buckets, we can create anything we need.
	//it all happens in a single transaction, based on kval query...
	err := createboltentries(kb)
	if err != nil {
		return err
	}
	return nil
}

//GET (Get Handler) handles GET capability of KVAL language
func getHandler(kb Kvalbolt) (Kvalresult, error) {
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
func getallHandler(kb Kvalbolt) (Kvalresult, error) {
	kr, err := getallfrombucket(kb)
	if err != nil {
		return kr, err
	}
	return kr, nil
}

//GET (Get Handler) handles GET (REGEX) capability of KVAL language
func getregexHandler(kb Kvalbolt) (Kvalresult, error) {
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
func delbucketHandler(kb Kvalbolt) error {
	err := deletebucket(kb)
	if err != nil {
		return err
	}
	return nil
}

//DEL (Delete Handler) handles DEL all keys capability of KVAL language
func delbucketkeysHandler(kb Kvalbolt) error {
	err := deletebucketkeys(kb)
	if err != nil {
		return err
	}
	return nil
}

//DEL (Delete Handler) handles DEL one key capability of KVAL language
func delonekeyHandler(kb Kvalbolt) error {
	err := deletekey(kb)
	if err != nil {
		return err
	}
	return nil
}

//DEL (Delete Handler) Handles DEL (or in this case, NULL, capability of KVAL
func nullifyvalHandler(kb Kvalbolt) error {
	err := nullifykeyvalue(kb)
	if err != nil {
		return err
	}
	return nil
}

//REN (Rename Handler) Handles rename bucket capability of KVAL
func renbucketHandler(kb Kvalbolt) error {
	err := renamebucket(kb)
	if err != nil {
		return err
	}
	return nil
}

//REN (Rename Handler) Handles rename key capability of KVAL
func renkeyHandler(kb Kvalbolt) error {
	err := renamekey(kb)
	if err != nil {
		return err
	}
	return nil
}

//LIS (List Handler) Handles listing capability of KVAL (does (x) exist?)
func lisHandler(kb Kvalbolt) (Kvalresult, error) {
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

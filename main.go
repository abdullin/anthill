package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"binary"

	"github.com/abdullin/events"
	"github.com/abdullin/fdb-go/fdb/subspace"

	"github.com/bmatsuo/lmdb-go/lmdb"
)

func handler(w http.ResponseWriter, r *http.Request) {
}

// Context is our test context
type Context struct {
	Lmdb  *lmdb.Env
	Dbi   lmdb.DBI
	Space subspace.Subspace
}

// Increment counter

var counter = []byte("counter")
var encoding := binary.LittleEndian


func (c Context) Increment(txn *lmdb.Txn, pos []byte) (uint32, err) {
	res, err := txn.Get(c.Dbi, pos)
	if err == nil {
		encoding.
		
		
	}
	if err != nil {
		if lmdb.IsNotFound(err) {
			if r, err := txn.PutReserve(c.Dbi, pos, 8, 0); err != nil {
				return 0, err
			} else {
				return 0, nil
			}
		}
		return 0, err
	}
	

}

func (c Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]

	switch path {
	case "put":
		id := events.NewSequentialUUID()

		full := c.Space.Sub(id).Bytes()

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(err)
		}

		err = c.Lmdb.Update(func(txn *lmdb.Txn) (err error) {
			txn.RawRead = true

			b, err := txn.Get(c.Dbi, []byte("counter"))

			err = txn.Put(c.Dbi, full, dump, 0)

			return err

		})
		if err != nil {
			panic(err)
		}
		w.Write([]byte("OK"))
		return
	}

	fmt.Fprintf(w, "Hi there, I love %s!", path)

}

func main() {

	env, err := lmdb.NewEnv()

	if err != nil {
		panic(err)
	}
	defer env.Close()

	// configure and open the environment.  most configuration must be done
	// before opening the environment.
	err = env.SetMaxDBs(1)
	if err != nil {
		log.Fatalf("Failed to configure env: %s", err)
	}
	err = env.SetMapSize(1 << 30)
	if err != nil {
		log.Fatalf("Failed to set map size to %d", 1<<30)
	}

	os.MkdirAll("db", os.ModePerm)
	err = env.Open("db", 0, 0644)
	if err != nil {
		log.Fatalf("Failed to open db")
	}

	// open a database that can be used as long as the enviroment is mapped.
	var dbi lmdb.DBI
	err = env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err = txn.CreateDBI("agg")
		return err
	})
	if err != nil {
		log.Fatalf("failed to open database")
	}

	ctx := &Context{
		Lmdb:  env,
		Dbi:   dbi,
		Space: subspace.Sub("bench"),
	}

	http.Handle("/", ctx)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"text/tabwriter"
	"time"

	capnp "zombiezen.com/go/capnproto2"

	"github.com/abdullin/lex-go/subspace"
	"github.com/abdullin/lex-go/tuple"
	"github.com/pborman/uuid"

	"github.com/bmatsuo/lmdb-go/lmdb"
)

var read uint64
var saved uint64

var (
	maxCountOption uint64
	maxMapSizeMb   int64
	envCount       int

	batchProducts int
	reuseDB       bool
	writeMap      bool
	asyncWrites   bool
)

type handle struct {
	env  *lmdb.Env
	name string
}

func newEnv(name string) (h handle) {

	var err error
	var env *lmdb.Env

	if !reuseDB {
		fmt.Println("Deleting the DB")
		if err := os.RemoveAll(name); err != nil {
			log.Fatalf("Failed to cleanup db folder: %s", err)
		}
	}

	env, err = lmdb.NewEnv()

	if err != nil {
		log.Fatalf("Failed to create env: %s", err)
	}

	// configure and open the environment.  most configuration must be done
	// before opening the environment.

	err = env.SetMaxDBs(5)
	if err != nil {
		log.Fatalf("Failed to configure env: %s", err)
	}

	fmt.Println("Setting map size to", maxMapSizeMb, "MB")

	err = env.SetMapSize(maxMapSizeMb * 1024 * 1024)
	if err != nil {
		log.Fatalf("Failed to set map size to %d", maxMapSizeMb)
	}

	var envFlags uint

	if asyncWrites {
		envFlags |= lmdb.NoSync
		fmt.Println("  env: NoSync (let OS flush pages to disk whenever it wants)")
	}

	if err := env.SetFlags(envFlags); err != nil {
		log.Fatalf("Failed to set flags %s", err)
	}

	os.MkdirAll(name, os.ModePerm)
	err = env.Open(name, 0, 0644)
	if err != nil {
		log.Fatalf("Failed to open db")
	}

	// open a database that can be used as long as the enviroment is mapped.
	err = env.Update(func(txn *lmdb.Txn) (err error) {
		_, err = txn.CreateDBI("agg")
		return err
	})
	if err != nil {
		log.Fatalf("failed to open database")
	}

	return handle{env, name}
}

func main() {

	flag.BoolVar(&asyncWrites, "async", false, "Enables no flush mode (makes LMDB ACI instead of ACID)")
	flag.BoolVar(&reuseDB, "reuse-db", true, "Keeps the database file")
	flag.BoolVar(&writeMap, "write-map", false, "Use writeable memory")
	flag.Uint64Var(&maxCountOption, "max", 1000000, "Max number of records to write")
	flag.Int64Var(&maxMapSizeMb, "mb", 1024, "Max map size")
	flag.IntVar(&batchProducts, "batch", 1, "Products batching")
	flag.IntVar(&envCount, "env-count", 1, "Environment count")

	flag.Parse()

	var txFlags uint
	if writeMap {
		txFlags |= lmdb.WriteMap
		fmt.Println("   tx: WriteMap (use writeable map pages)")
	}

	log.Printf("Using %d environments", envCount)

	var handles = make([]handle, envCount, envCount)

	for e := 0; e < envCount; e++ {
		name := "db" + strconv.Itoa(e)

		h := newEnv(name)
		defer h.env.Close()
		handles[e] = h
	}

	for e := 0; e < envCount; e++ {

		h := handles[e]

		go func() {
			benchWrites(h, txFlags)
			benchLookups(h)
		}()
	}

	ticker := time.NewTicker(1 * time.Second)

	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 12, 0, padding, ' ', tabwriter.AlignRight|tabwriter.TabIndent)
	fmt.Fprintln(w, "write tx/s", "\t", "read tx/s", "\t", "total", "\t", "Size MB", "\t")
	w.Flush()

	for {
		savedStart := saved
		readStart := read
		select {
		case <-ticker.C:

			fi, e := os.Stat("db/data.mdb")
			if e != nil {
				panic(e)
			}
			// get the size
			size := fi.Size() / 1024 / 1024

			fmt.Fprintln(w, (saved - savedStart), "\t", read-readStart, "\t", saved, "\t", size, "\t")

			w.Flush()
		}
	}

}

// BenchLookups looks up a random sku, then loads the associated
// product and verifies that its SKU is the one we expected
func benchLookups(h handle) {

	var dbi lmdb.DBI
	var err error
	var env = h.env
	fmt.Println("Product sku lookup benchmark")

	err = env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err = txn.OpenDBI("agg", 0)
		return err
	})
	if err != nil {
		log.Fatalf("failed to open database")
	}

	txn, err := env.BeginTxn(nil, lmdb.Readonly)

	defer txn.Abort()

	if err != nil {
		panic(err)
	}

	for {
		txn.Reset()
		txn.Renew()

		err = handleRead(txn, dbi)

		if err != nil {
			panic(err)
		}
	}

}

func handleRead(txn *lmdb.Txn, dbi lmdb.DBI) (err error) {

	curr := atomic.AddUint64(&read, 1)

	id := curr % saved

	num := strconv.Itoa(int(id))

	sku := "sku" + num

	skuIndexKey := skuIndex.Pack(tuple.Tuple{sku})

	txn.RawRead = true
	var data []byte
	data, err = txn.Get(dbi, skuIndexKey)

	if err != nil {
		return err
	}

	productID := data

	productKey := prodTable.Pack(tuple.Tuple{productID})

	data, err = txn.Get(dbi, productKey)
	if err != nil {
		return
	}

	reader := bytes.NewReader(data)

	msg, err := capnp.NewPackedDecoder(reader).Decode()
	if err != nil {
		return
	}
	product, err := ReadRootProduct(msg)
	if err != nil {
		return
	}

	realSku, err := product.Sku()
	if err != nil {
		return
	}
	if strings.Compare(realSku, sku) != 0 {
		panic("Expected and actual SKU don't match")
	}

	return err

}

func benchWrites(h handle, txFlags uint) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var dbi lmdb.DBI
	env := h.env

	err = h.env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err = txn.OpenDBI("agg", 0)
		return err
	})
	if err != nil {
		log.Fatalf("failed to open database")
	}

	// pin this routine to a single thread. This allows us to use
	// locked version of LMDB txn update

	fmt.Println("Product append benchmark with", maxCountOption, "records")

	var i uint64

	var iterations = maxCountOption / uint64(batchProducts)

	for i = 0; i < iterations; i++ {
		err = env.RunTxn(txFlags, func(txn *lmdb.Txn) (err error) {

			for j := 0; j < batchProducts; j++ {

				setProduct(txn, dbi, saved)
				saved++
			}

			setCounter(txn, dbi, saved-1)

			return err
		})

		if err != nil {
			log.Fatalf("failed to open database")
		}
	}

	if err := env.Sync(true); err != nil {
		log.Fatalf("Failed to fsync %s", err)
	}

}

var checkKey = []byte("counter")

var order = binary.LittleEndian

var codeIndex = subspace.Sub("code")
var skuIndex = subspace.Sub("sku")
var prodTable = subspace.Sub("prod")

func setProduct(txn *lmdb.Txn, dbi lmdb.DBI, id uint64) (err error) {
	// Make a brand new empty message.  A Message allocates Cap'n Proto structs.

	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return err
	}

	prod, err := NewRootProduct(seg)
	if err != nil {
		return err
	}

	classification, err := prod.NewClassification()
	if err != nil {
		return err
	}

	classification.SetId(id)
	classification.SetName("Bing")

	num := strconv.Itoa(int(id))

	productID := uuid.NewRandom()

	code := "code" + productID.String()
	sku := "sku" + num

	prod.SetDescription("description")
	prod.SetCode(code)
	prod.SetSku(sku)
	prod.SetId(id)
	prod.SetClassification(classification)

	var buffer bytes.Buffer
	var writer = bufio.NewWriter(&buffer)
	var encoder = capnp.NewPackedEncoder(writer)
	err = encoder.Encode(msg)
	if err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}

	var keyBuffer = []byte(productID)

	codeIndexKey := codeIndex.Pack(tuple.Tuple{code})
	skuIndexKey := skuIndex.Pack(tuple.Tuple{sku})
	prodValueKey := prodTable.Pack(tuple.Tuple{keyBuffer})

	if err = txn.Put(dbi, codeIndexKey, keyBuffer, 0); err != nil {
		return err
	}

	if err = txn.Put(dbi, skuIndexKey, keyBuffer, 0); err != nil {
		return
	}

	return txn.Put(dbi, prodValueKey, buffer.Bytes(), 0)
}

func setCounter(txn *lmdb.Txn, dbi lmdb.DBI, counter uint64) (err error) {
	var buf []byte
	buf, err = txn.PutReserve(dbi, checkKey, 8, 0)

	if err != nil {
		return err
	}

	order.PutUint64(buf, counter)
	return nil
}

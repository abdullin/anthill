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

var (
	maxCountOption uint64
	maxMapSizeMb   int64
	envCount       int
	readCount      int

	batchProducts int
	reuseDB       bool
	writeMap      bool
	asyncWrites   bool

	parallelRW bool
)

type handle struct {
	env   *lmdb.Env
	name  string
	saved uint64
	read  uint64
}

func newEnv(name string) (h *handle) {

	var err error
	var env *lmdb.Env

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

	sizeMbs := maxMapSizeMb / int64(envCount)

	fmt.Println(name, "setting map size to", sizeMbs, "MB * ", envCount, "(size split between writers)")

	err = env.SetMapSize(sizeMbs * 1024 * 1024)
	if err != nil {
		log.Fatalf("Failed to set map size to %d", sizeMbs)
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

	var dbi lmdb.DBI
	err = env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err = txn.CreateDBI("agg")
		return err
	})
	if err != nil {
		log.Fatalf("failed to create database")
	}

	env.CloseDBI(dbi)

	return &handle{env, name, 0, 0}
}

func dbName(e int) string {
	return "db/" + strconv.Itoa(e)
}

func main() {

	flag.BoolVar(&asyncWrites, "async", false, "Enables no flush mode (makes LMDB ACI instead of ACID)")
	flag.BoolVar(&reuseDB, "reuse-db", false, "Keeps the database file")
	flag.BoolVar(&writeMap, "write-map", false, "Use writeable memory")

	flag.BoolVar(&parallelRW, "parallel-rw", false, "Reads and writes happen in parallel")
	flag.Uint64Var(&maxCountOption, "max", 1000000, "Max number of records to write across all envs")
	flag.Int64Var(&maxMapSizeMb, "mb", 1024, "Max map size for all envs")
	flag.IntVar(&batchProducts, "batch", 1, "Products batching")
	flag.IntVar(&envCount, "env-count", 1, "Environement and write thread count")
	flag.IntVar(&readCount, "read-count", 1, "Read threads PER ENV")

	flag.Parse()

	var txFlags uint
	if writeMap {
		txFlags |= lmdb.WriteMap
		fmt.Println("   tx: WriteMap (use writeable map pages)")
	}

	fmt.Println("Writers and environments:", envCount)

	if !reuseDB {
		fmt.Println("Deleting the DB")
		if err := os.RemoveAll("db"); err != nil {
			log.Fatalf("Failed to cleanup db folder: %s", err)
		}
	} else {
		fmt.Println("Keeping the db")
	}

	var handles = make([]*handle, envCount, envCount)

	for e := 0; e < envCount; e++ {
		name := dbName(e)

		h := newEnv(name)
		defer h.env.Close()
		handles[e] = h
	}

	for e := 0; e < envCount; e++ {

		h := handles[e]

		if parallelRW {
			go benchWrites(h, txFlags)

			for r := 0; r < readCount; r++ {
				go benchLookups(h)
			}

		} else {
			go func() {
				benchWrites(h, txFlags)
				for r := 0; r < readCount; r++ {
					go benchLookups(h)
				}
			}()
		}
	}

	ticker := time.NewTicker(1 * time.Second)

	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 12, 0, padding, ' ', tabwriter.AlignRight|tabwriter.TabIndent)
	fmt.Fprintln(w, "write tx/s", "\t", "read tx/s", "\t", "total", "\t", "Size MB", "\t")
	w.Flush()

	for {

		var savedStart, readStart uint64

		for e := 0; e < envCount; e++ {
			savedStart += handles[e].saved
			readStart += handles[e].read
		}

		select {
		case <-ticker.C:

			var size int64
			var savedCurrent, readCurrent uint64

			for e := 0; e < envCount; e++ {

				fi, err := os.Stat(dbName(e) + "/data.mdb")
				if err != nil {
					panic(err)
				}
				// get the size
				size += fi.Size() / 1024 / 1024
				savedCurrent += handles[e].saved
				readCurrent += handles[e].read
			}

			fmt.Fprintln(w, (savedCurrent - savedStart), "\t", readCurrent-readStart, "\t", savedCurrent, "\t", size, "\t")

			w.Flush()
		}
	}

}

// BenchLookups looks up a random sku, then loads the associated
// product and verifies that its SKU is the one we expected
func benchLookups(h *handle) {

	var dbi lmdb.DBI
	var err error
	var env = h.env

	err = env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err = txn.OpenDBI("agg", 0)
		return err
	})

	defer env.CloseDBI(dbi)
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

		err = handleRead(txn, dbi, h)

		if err != nil {
			panic(err)
		}
	}

}

func handleRead(txn *lmdb.Txn, dbi lmdb.DBI, h *handle) (err error) {

	saved := atomic.LoadUint64(&h.saved)
	if saved == 0 {
		time.Sleep(time.Millisecond)
		return nil
	}

	curr := atomic.AddUint64(&h.read, 1)

	id := curr % saved

	num := strconv.Itoa(int(id))

	sku := "sku" + num

	skuIndexKey := skuIndex.Pack(tuple.Tuple{sku})

	txn.RawRead = true
	var data []byte
	data, err = txn.Get(dbi, skuIndexKey)

	if err != nil {
		return fmt.Errorf("Failed to find product %s: %s", num, err)
	}

	productID := data

	productKey := prodTable.Pack(tuple.Tuple{productID})

	data, err = txn.Get(dbi, productKey)
	if err != nil {
		return fmt.Errorf("Failed to load product %s: %s", num, err)
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

func benchWrites(h *handle, txFlags uint) {

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

	var i uint64

	var iterations = maxCountOption / uint64(batchProducts) / uint64(envCount)

	var savedProducts uint64

	for i = 0; i < iterations; i++ {
		err = env.RunTxn(txFlags, func(txn *lmdb.Txn) (err error) {

			for j := 0; j < batchProducts; j++ {

				setProduct(txn, dbi, savedProducts)
				savedProducts++
			}

			setCounter(txn, dbi, savedProducts-1)

			return err
		})
		atomic.StoreUint64(&h.saved, savedProducts)

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

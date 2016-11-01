package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	capnp "zombiezen.com/go/capnproto2"

	_ "github.com/abdullin/events"
	"github.com/abdullin/fdb-go/fdb/subspace"
	"github.com/abdullin/fdb-go/fdb/tuple"

	"github.com/bmatsuo/lmdb-go/lmdb"
)

func main() {

	env, err := lmdb.NewEnv()

	if err != nil {
		log.Fatalf("Failed to create env: %s", err)
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

	var counter uint64

	ticker := time.NewTicker(1 * time.Second)
	go func() {

		for {
			start := counter
			select {
			case <-ticker.C:
				fmt.Println("1 sec ", counter-start, " total ", counter)
				// do stuff
			}
		}
	}()

	for {

		err = env.Update(func(txn *lmdb.Txn) (err error) {
			setProduct(txn, dbi, counter)
			setCounter(txn, dbi, counter)

			//dbi, err = txn.CreateDBI("agg")
			return err
		})
		if err != nil {
			log.Fatalf("failed to open database")
		}
		counter++

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

	classification, err := NewClassification(seg)
	if err != nil {
		return err
	}

	classification.SetId(id)
	classification.SetName("Bing")

	num := strconv.Itoa(int(id))

	code := "code" + num
	sku := "sku" + num

	prod.SetDescription("description")
	prod.SetCode(code)
	prod.SetSku(sku)
	prod.SetId(id)
	prod.SetClassification(classification)

	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	enc := capnp.NewPackedEncoder(buffer)
	err = enc.Encode(msg)
	if err != nil {
		return err
	}
	if err = buffer.Flush(); err != nil {
		return err
	}

	// key
	buf := make([]byte, 8, 8)
	order.PutUint64(buf, uint64(id))

	codeIndexKey := codeIndex.Pack(tuple.Tuple{code})
	skuIndexKey := skuIndex.Pack(tuple.Tuple{sku})
	prodValueKey := prodTable.Pack(tuple.Tuple{id})

	if err = txn.Put(dbi, codeIndexKey, buf, 0); err != nil {
		return err
	}

	if err = txn.Put(dbi, skuIndexKey, buf, 0); err != nil {
		return
	}

	return txn.Put(dbi, prodValueKey, b.Bytes(), 0)
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

package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

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

	opts := &testOptions{}

	opts.createProb = 0.3
	opts.deleteProb = 0.1
	opts.changeSkuProb = 0.1
	opts.readProb = 0.5
	opts.tenants = 3
	opts.partitions = 1
	opts.iterations = 10

	firehose := make(chan request, 1000)
	go generate(opts, firehose)

	for r := range firehose {
		fmt.Println(reflect.TypeOf(r).String(), r)
	}

	// flag.BoolVar(&asyncWrites, "async", false, "Enables no flush mode (makes LMDB ACI instead of ACID)")
	// flag.BoolVar(&reuseDB, "reuse-db", false, "Keeps the database file")
	// flag.BoolVar(&writeMap, "write-map", false, "Use writeable memory")

	// flag.BoolVar(&parallelRW, "parallel-rw", false, "Reads and writes happen in parallel")
	// flag.Uint64Var(&maxCountOption, "max", 1000000, "Max number of records to write across all envs")
	// flag.Int64Var(&maxMapSizeMb, "mb", 1024, "Max map size for all envs")
	// flag.IntVar(&batchProducts, "batch", 1, "Products batching")
	// flag.IntVar(&envCount, "env-count", 1, "Environement and write thread count")
	// flag.IntVar(&readCount, "read-count", 1, "Read threads PER ENV")

	// flag.Parse()

}

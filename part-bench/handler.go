package main

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"github.com/abdullin/lex-go/subspace"
	"github.com/abdullin/lex-go/tuple"
	"github.com/bmatsuo/lmdb-go/lmdb"
	"github.com/pborman/uuid"
	capnp "zombiezen.com/go/capnproto2"
)

type partition struct {
	env  *lmdb.Env
	dbis map[TenantID]lmdb.DBI
}

type table byte

const (
	codeSpace table = iota
	skuSpace
	prodSpace
)

var codeIndex = subspace.Sub(byte(codeSpace))
var skuIndex = subspace.Sub(byte(skuSpace))
var prodTable = subspace.Sub(byte(prodSpace))

func (p *partition) getDBI(tenant TenantID, txn *lmdb.Txn) (dbi lmdb.DBI, err error) {
	if dbi, ok := p.dbis[tenant]; ok {
		return dbi, nil
	}

	n := fmt.Sprintf("%d", tenant)

	if dbi, err = txn.OpenDBI(n, lmdb.Create); err != nil {
		err = fmt.Errorf("Failed to create db for tenant %d: %s", tenant, err)
		return
	}

	p.dbis[tenant] = dbi
	return dbi, nil

}

func (p *partition) processWriteBatch(rs []request) {
	err := p.env.RunTxn(0, func(txn *lmdb.Txn) (err error) {

		for r := range rs{
			switch r.(type) {
			case createProduct:
				
			}
		}

			for j := 0; j < batchProducts; j++ {

				setProduct(txn, dbi, savedProducts)
				savedProducts++
			}

			setCounter(txn, dbi, savedProducts-1)

			return err
		})

}
func (p *partition) handleWrites(c chan request) {

	// batch size
	const capacity = 100
	batch := make([]request, 0, capacity)

	for {
		select {
		case x := <-c:
			batch := append(batch, x)
			if len(batch) == cap(batch) {
				p.processWriteBatch(batch)
				batch = batch[:0]
			}
		default:

			// nothing in the queue
			if len(batch) > 0 {
				p.processBatch(batch)
				// clear the batch
				batch = batch[:0]
			} else {
				time.Sleep(time.Millisecond)
			}
		}
	}
}






func (p *partition) handleCreate(txn *lmdb.Txn, dbi lmdb.DBI, e createProduct) (err error) {

	// options: DBI per tenant or all tenants together
	// options: nested transaction with abort or careful checks

	

	prodValueKey := prodTable.Pack(tuple.Tuple{uint64(e.id)})	
	codeIndexKey := codeIndex.Pack(tuple.Tuple{string(e.code)})
	skuIndexKey := skuIndex.Pack(tuple.Tuple{string(e.sku)})

	{
		
	}

	if _, err := txn.Get(dbi, codeIndexKey); err == lmdb.NotFound

	
	if dbi, err := p.getDBI(e.tenant, txn); err != nil {
		return err
	}

	

	// tenant is encoded in the dbi
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

	classification.SetId(uint64(e.id))
	classification.SetName("Bing")

	productID := uuid.NewRandom()

	prod.SetDescription("description")
	prod.SetCode(string(e.code))
	prod.SetSku(string(e.sku))
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


	if err = txn.Put(dbi, codeIndexKey, keyBuffer, 0); err != nil {
		return err
	}

	if err = txn.Put(dbi, skuIndexKey, keyBuffer, 0); err != nil {
		return
	}

	return txn.Put(dbi, prodValueKey, buffer.Bytes(), 0)
}

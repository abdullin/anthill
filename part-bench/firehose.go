package main

import (
	"fmt"
	"strconv"

	"github.com/pborman/uuid"
)

// TenantID - positive tenant id
type TenantID uint32

// ProductID - product id
type ProductID uint64

// ProductSku - random product id
type ProductSku string

// ProductCode - fixed product id
type ProductCode string

type request interface {
	info() TenantID
}

type createProduct struct {
	tenant TenantID
	id     ProductID
	sku    ProductSku
	code   ProductCode
}

func (c createProduct) info() TenantID {
	return c.tenant
}

type changeProductSku struct {
	tenant TenantID
	id     ProductID
	sku    ProductSku
}

func (c changeProductSku) info() TenantID {
	return c.tenant
}

type deleteProduct struct {
	tenant TenantID
	id     ProductID
}

func (c deleteProduct) info() TenantID {
	return c.tenant
}

type queryProductByCode struct {
	tenant TenantID
	code   ProductCode
}

func (c queryProductByCode) info() TenantID {
	return c.tenant
}

type testOptions struct {
	tenants       int
	partitions    int
	createProb    float32
	deleteProb    float32
	readProb      float32
	changeSkuProb float32
	iterations    int
}

// generate
func generate(opts *testOptions, c chan request) {

	max := opts.iterations

	testBase := 10

	createRatio := int(opts.createProb * float32(testBase))
	deleteRatio := int(opts.deleteProb * float32(testBase))
	changeSkuRatio := int(opts.changeSkuProb * float32(testBase))

	fmt.Println("Given", testBase, "throws, we will:")

	fmt.Println("   Create    ", createRatio)
	fmt.Println("   Delete    ", deleteRatio)
	fmt.Println("   Change Sku", changeSkuRatio)

	tenants := opts.tenants
	counters := make([]int, tenants, tenants)

	seq := 0
	for i := 0; i < max; i++ {
		seq++

		if (seq % testBase) < createRatio {

			tenant := i % tenants

			counter := counters[tenant] + 1
			counters[tenant] = counter

			id := uuid.NewUUID()
			sku := "sku" + id.String()
			code := "code" + strconv.Itoa(counter)

			c <- &createProduct{
				TenantID(tenant),
				ProductID(counter),
				ProductSku(sku),
				ProductCode(code),
			}
		}

		seq++
		if (seq % testBase) < changeSkuRatio {

			tenant := i % tenants
			counter := counters[tenant]
			id := i % counter

			skuBase := uuid.NewUUID()
			sku := "sku" + skuBase.String()

			c <- &changeProductSku{
				TenantID(tenant),
				ProductID(id),
				ProductSku(sku),
			}
		}
		seq++
		// delete by id
		if (seq % testBase) < deleteRatio {
			tenant := i % tenants
			counter := counters[tenant]
			id := i % counter

			c <- &deleteProduct{
				TenantID(tenant),
				ProductID(id),
			}
		}

	}

}

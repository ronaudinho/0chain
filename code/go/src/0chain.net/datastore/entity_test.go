package datastore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"0chain.net/common"
)

/*Company - a test data type */
type Company struct {
	CollectionIDField
	//ID     Key    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name,omitempty"`
}

func (c *Company) GetEntityName() string {
	return "company"
}

/*
func (c *Company) SetKey(key Key) {
	c.ID = key
}

func (c *Company) GetKey() Key {
	return c.ID
}

func (c *Company) ComputeProperties() {

}

func (c *Company) Validate(ctx context.Context) error {
	return nil
} */

func (c *Company) Read(ctx context.Context, id Key) error {
	return Read(ctx, id, c)
}

func (c *Company) Write(ctx context.Context) error {
	return Write(ctx, c)
}

func (c *Company) Delete(ctx context.Context) error {
	return Delete(ctx, c)
}

var companyEntityCollection = &EntityCollection{CollectionName: "collection.company", CollectionSize: 10000, CollectionDuration: time.Hour}

/*TransactionProvider - entity provider for client object */
func CompanyProvider() interface{} {
	c := &Company{}
	c.CollectionIDField.EntityCollection = companyEntityCollection
	return c
}

func TestEntityWriteRead(t *testing.T) {
	fmt.Printf("time : %v\n", time.Now().UnixNano()/int64(time.Millisecond))
	common.SetupRootContext(context.Background())
	ctx := WithConnection(common.GetRootContext())
	zeroChain := CompanyProvider().(*Company)
	zeroChain2 := CompanyProvider().(*Company)
	keys := []Key{ToKey([]byte("0chain.net")), ToKey("0chain.io")}
	entities := []Entity{zeroChain, zeroChain2}
	fmt.Printf("keys : %v\n", keys)
	err := MultiRead(ctx, keys, entities)
	if err != nil {
		fmt.Printf("error reading : %v\n", err)
	} else {
		fmt.Printf("e1 %v\n", entities[0])
		fmt.Printf("e2 %v\n", entities[1])
	}
	zeroChain.Domain = "0chain.net"
	zeroChain.Name = "0chain"
	zeroChain.ID = ToKey(zeroChain.Domain)
	zeroChain.CollectionIDField.EntityCollection = companyEntityCollection
	err = InsertIfNE(ctx, zeroChain)
	if err != nil {
		fmt.Printf("error ifne: %v\n", err)
	}
	zeroChain2.Domain = "0chain.io"
	err = Read(ctx, ToKey(zeroChain2.Domain), zeroChain2)
	if err != nil {
		fmt.Printf("error reading: %v\n", err)
	} else {
		fmt.Printf("zc2 = %v\n", zeroChain2)
	}
	zeroChain2.InitCollectionScore()
	zeroChain2.SetCollectionScore(zeroChain2.GetCollectionScore() + 10)
	MultiWrite(ctx, []Entity{zeroChain, zeroChain2})

	fmt.Printf("iterating\n")
	IterateCollection(ctx, zeroChain.GetCollectionName(), PrintIterator, CompanyProvider)

}

/*
func TestEntityCollectionTrimmer(t *testing.T) {
	CollectionTrimmer("collection.company", 100, time.Second)
} */

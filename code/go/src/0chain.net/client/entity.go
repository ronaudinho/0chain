package client

import (
	"context"
	"time"

	"0chain.net/common"
	"0chain.net/datastore"
	"0chain.net/encryption"
)

/*Client - data structure that holds the client data */
type Client struct {
	datastore.IDField
	datastore.CreationDateField
	PublicKey string `json:"public_key"`
}

/*GetEntityName - implementing the interface */
func (c *Client) GetEntityName() string {
	return "client"
}

/*Validate - implementing the interface */
func (c *Client) Validate(ctx context.Context) error {
	if datastore.IsEmpty(c.ID) {
		return common.InvalidRequest("client id is required")
	}
	if !datastore.IsEqual(c.ID, datastore.ToKey(encryption.Hash(c.PublicKey))) {
		return common.InvalidRequest("client id is not a SHA3-256 hash of the public key")
	}
	return nil
}

/*Read - datastore read */
func (c *Client) Read(ctx context.Context, key datastore.Key) error {
	return datastore.Read(ctx, key, c)
}

/*Write - datastore read */
func (c *Client) Write(ctx context.Context) error {
	return datastore.Write(ctx, c)
}

/*Delete - datastore read */
func (c *Client) Delete(ctx context.Context) error {
	return datastore.Delete(ctx, c)
}

/*Verify - given a signature and hash verify it with client's public key */
func (c *Client) Verify(signature string, hash string) (bool, error) {
	return encryption.Verify(c.PublicKey, signature, hash)
}

/*Provider - entity provider for client object */
func Provider() interface{} {
	c := &Client{}
	c.InitializeCreationDate()
	return c
}

/*SetupEntity - setup the entity */
func SetupEntity() {
	datastore.RegisterEntityProvider("client", Provider)

	var collectionOptions = datastore.CollectionOptions{
		EntityBufferSize: 1024,
		MaxHoldupTime:    500 * time.Millisecond,
		NumChunkCreators: 1,
		ChunkSize:        64,
		ChunkBufferSize:  16,
		NumChunkStorers:  2,
	}
	ClientEntityChannel = datastore.SetupWorkers(common.GetRootContext(), &collectionOptions)
}

var ClientEntityChannel chan datastore.Entity

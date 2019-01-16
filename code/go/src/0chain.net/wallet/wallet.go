package wallet

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"

	"0chain.net/client"
	"0chain.net/datastore"
	"0chain.net/encryption"
	"0chain.net/transaction"
)

/*Wallet - a struct representing the client's wallet */
type Wallet struct {
	SignatureScheme encryption.SignatureScheme
	PublicKeyBytes  []byte
	ClientID        string
	Balance         int64
}

/*Initialize - initialize a wallet with public/private keys */
func (w *Wallet) Initialize() error {
	sigScheme := encryption.NewED25519Scheme()
	err := sigScheme.GenerateKeys() // generate a new pair of public/private key combination
	if err != nil {
		return err
	}
	return w.SetSignatureScheme(sigScheme)
}

/*SetSignatureScheme - sets the keys for the wallet */
func (w *Wallet) SetSignatureScheme(signatureScheme encryption.SignatureScheme) error {
	w.SignatureScheme = signatureScheme
	publicKeyBytes, err := hex.DecodeString(signatureScheme.GetPublicKey())
	if err != nil {
		return err
	}
	w.ClientID = encryption.Hash(publicKeyBytes) // ClientID is the hexadecimal encoding of the hash of the Public key bytes
	return nil
}

/*Register - register a wallet using the server side api */
func (w *Wallet) Register(ctx context.Context) error {
	c := clientMetadataProvider.Instance().(*client.Client)
	c.PublicKey = w.SignatureScheme.GetPublicKey()
	c.ID = w.ClientID
	_, err := client.PutClient(ctx, c)
	return err
}

var transactionMetadataProvider datastore.EntityMetadata
var clientMetadataProvider datastore.EntityMetadata

/*SetupWallet - setup the wallet package */
func SetupWallet() {
	transactionMetadataProvider = datastore.GetEntityMetadata("txn")
	clientMetadataProvider = datastore.GetEntityMetadata("client")
}

/*CreateRandomSendTransaction - create a transaction */
func (w *Wallet) CreateRandomSendTransaction(toClient string) *transaction.Transaction {
	value := rand.Int63n(100) * 1000000000
	if value == 0 {
		value = 100000000
	}
	msg := fmt.Sprintf("0chain zerochain zipcode Europe rightthing Oriental California honest accurate India network %v %v", rand.Int63(), value)
	return w.CreateSendTransaction(toClient, value, msg)
}

/*CreateSendTransaction - create a send transaction */
func (w *Wallet) CreateSendTransaction(toClient string, value int64, msg string) *transaction.Transaction {
	txn := transactionMetadataProvider.Instance().(*transaction.Transaction)
	txn.ClientID = w.ClientID
	txn.ToClientID = toClient
	txn.Value = value
	txn.TransactionData = msg
	txn.Sign(w.SignatureScheme)
	return txn
}

/*CreateRandomDataTransaction - creat a random data transaction */
func (w *Wallet) CreateRandomDataTransaction() *transaction.Transaction {
	msg := fmt.Sprintf("storing some random data - 1234567890 abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ %v", rand.Int63())
	return w.CreateDataTransaction(msg)
}

/*CreateDataTransaction - create a data transaction */
func (w *Wallet) CreateDataTransaction(msg string) *transaction.Transaction {
	txn := transactionMetadataProvider.Instance().(*transaction.Transaction)
	txn.ClientID = w.ClientID
	txn.TransactionData = msg
	txn.TransactionType = transaction.TxnTypeData
	txn.Sign(w.SignatureScheme)
	return txn
}

package extension

import (
	"context"

	"github.com/corverroos/quorum"
	"github.com/corverroos/quorum/common"
	"github.com/corverroos/quorum/core/types"
	"github.com/corverroos/quorum/ethclient"
)

type Client interface {
	SubscribeToLogs(query ethereum.FilterQuery) (<-chan types.Log, ethereum.Subscription, error)
	NextNonce(from common.Address) (uint64, error)
	TransactionByHash(hash common.Hash) (*types.Transaction, error)
	TransactionInBlock(blockHash common.Hash, txIndex uint) (*types.Transaction, error)
	Close()
}

type InProcessClient struct {
	client *ethclient.Client
}

func NewInProcessClient(client *ethclient.Client) *InProcessClient {
	return &InProcessClient{
		client: client,
	}
}

func (client *InProcessClient) SubscribeToLogs(query ethereum.FilterQuery) (<-chan types.Log, ethereum.Subscription, error) {
	retrievedLogsChan := make(chan types.Log)
	sub, err := client.client.SubscribeFilterLogs(context.Background(), query, retrievedLogsChan)
	return retrievedLogsChan, sub, err
}

func (client *InProcessClient) NextNonce(from common.Address) (uint64, error) {
	return client.client.PendingNonceAt(context.Background(), from)
}

func (client *InProcessClient) TransactionByHash(hash common.Hash) (*types.Transaction, error) {
	tx, _, err := client.client.TransactionByHash(context.Background(), hash)
	return tx, err
}

func (client *InProcessClient) TransactionInBlock(blockHash common.Hash, txIndex uint) (*types.Transaction, error) {
	tx, err := client.client.TransactionInBlock(context.Background(), blockHash, txIndex)
	if err != nil {
		return nil, err
	}

	// Fetch the underlying private tx if we got a Privacy Marker Transaction
	if tx.IsPrivacyMarker() {
		return client.client.GetPrivateTransaction(context.Background(), tx.Hash())
	}

	return tx, nil
}

func (client *InProcessClient) Close() {
	client.client.Close()
}

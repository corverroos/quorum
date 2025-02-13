package extension

import (
	"github.com/corverroos/quorum"
	"github.com/corverroos/quorum/core/types"
	"github.com/corverroos/quorum/ethclient"
	"github.com/corverroos/quorum/log"
	"github.com/corverroos/quorum/node"
	"github.com/corverroos/quorum/private"
)

type subscriptionHandler struct {
	facade  ManagementContractFacade
	client  Client
	service *PrivacyService
}

func NewSubscriptionHandler(node *node.Node, psi types.PrivateStateIdentifier, ptm private.PrivateTransactionManager, service *PrivacyService) *subscriptionHandler {
	rpcClient, err := node.AttachWithPSI(psi)
	if err != nil {
		panic("extension: could not connect to ethereum client rpc")
	}

	client := ethclient.NewClientWithPTM(rpcClient, ptm)

	return &subscriptionHandler{
		facade:  NewManagementContractFacade(client),
		client:  NewInProcessClient(client),
		service: service,
	}
}

func (handler *subscriptionHandler) createSub(query ethereum.FilterQuery, logHandlerCb func(types.Log)) error {
	incomingLogs, subscription, err := handler.client.SubscribeToLogs(query)

	if err != nil {
		return err
	}

	go func() {
		stopChan, stopSubscription := handler.service.subscribeStopEvent()
		defer stopSubscription.Unsubscribe()

		for {
			select {
			case err := <-subscription.Err():
				log.Error("Contract extension watcher subscription error", "error", err)
				break
			case foundLog := <-incomingLogs:
				logHandlerCb(foundLog)
			case <-stopChan:
				return
			}
		}
	}()

	return nil
}

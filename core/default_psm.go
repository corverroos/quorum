package core

import (
	"context"

	"github.com/corverroos/quorum/common"
	"github.com/corverroos/quorum/core/mps"
	"github.com/corverroos/quorum/core/rawdb"
	"github.com/corverroos/quorum/core/state"
	"github.com/corverroos/quorum/core/types"
	"github.com/corverroos/quorum/ethdb"
	"github.com/corverroos/quorum/rpc"
	"github.com/corverroos/quorum/trie"
)

type DefaultPrivateStateManager struct {
	// Low level persistent database to store final content in
	db        ethdb.Database
	repoCache state.Database
}

func newDefaultPrivateStateManager(db ethdb.Database, config *trie.Config) *DefaultPrivateStateManager {
	return &DefaultPrivateStateManager{
		db:        db,
		repoCache: state.NewDatabaseWithConfig(db, config),
	}
}

func (d *DefaultPrivateStateManager) StateRepository(blockHash common.Hash) (mps.PrivateStateRepository, error) {
	return mps.NewDefaultPrivateStateRepository(d.db, d.repoCache, blockHash)
}

func (d *DefaultPrivateStateManager) ResolveForManagedParty(_ string) (*mps.PrivateStateMetadata, error) {
	return mps.DefaultPrivateStateMetadata, nil
}

func (d *DefaultPrivateStateManager) ResolveForUserContext(ctx context.Context) (*mps.PrivateStateMetadata, error) {
	psi, ok := rpc.PrivateStateIdentifierFromContext(ctx)
	if !ok {
		psi = types.DefaultPrivateStateIdentifier
	}
	return &mps.PrivateStateMetadata{ID: psi, Type: mps.Resident}, nil
}

func (d *DefaultPrivateStateManager) PSIs() []types.PrivateStateIdentifier {
	return []types.PrivateStateIdentifier{
		types.DefaultPrivateStateIdentifier,
	}
}

func (d *DefaultPrivateStateManager) NotIncludeAny(_ *mps.PrivateStateMetadata, _ ...string) bool {
	// with default implementation, all managedParties are members of the psm
	return false
}

func (d *DefaultPrivateStateManager) CheckAt(root common.Hash) error {
	_, err := state.New(rawdb.GetPrivateStateRoot(d.db, root), d.repoCache, nil)
	return err
}

func (d *DefaultPrivateStateManager) TrieDB() *trie.Database {
	return d.repoCache.TrieDB()
}

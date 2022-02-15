package extensionContracts

import (
	"github.com/corverroos/quorum/core/state"
)

type AccountWithMetadata struct {
	State state.DumpAccount `json:"state"`
}

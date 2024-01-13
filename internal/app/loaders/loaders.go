package loaders

import (
	"github.com/graph-gophers/dataloader"
	"gqlgen-starter/internal/ent"
	"time"
)

type Loaders struct {
	UserLoader *dataloader.Loader
}

// NewLoaders instantiates data loaders
func NewLoaders(entClient *ent.Client) *Loaders {
	// define the data loader
	ur := &userReader{entClient}
	return &Loaders{
		UserLoader: dataloader.NewBatchedLoader(ur.GetUsersBatchFn, dataloader.WithWait(time.Millisecond)),
	}
}

// handleError creates slice of dataloader.Result with the same error instance in each
func handleError(numItems int, err error) []*dataloader.Result {
	result := make([]*dataloader.Result, numItems)
	for i := 0; i < numItems; i++ {
		result[i] = &dataloader.Result{Error: err}
	}
	return result
}

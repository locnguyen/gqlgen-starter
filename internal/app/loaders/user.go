package loaders

import (
	"context"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/ent/user"
	"gqlgen-starter/internal/oops"
	"gqlgen-starter/internal/utils"
	"net/http"
)

type userReader struct {
	EntClient *ent.Client
}

func (r *userReader) GetUsersBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	id64s := make([]int64, 0, len(keys))
	keyOrder := make(map[int64]int, len(keys))

	for idx, k := range keys {
		id64, err := utils.ID64(k.String())
		if err != nil {
			return handleError(len(keys), err)
		}
		keyOrder[id64] = idx
		id64s = append(id64s, id64)
	}

	entC := ent.FromContext(ctx)
	users, err := entC.User.Query().
		Where(user.IDIn(id64s...)).
		All(ctx)

	if err != nil {
		return handleError(len(keys), err)
	}

	results := make([]*dataloader.Result, len(keys))
	for _, u := range users {
		ix, ok := keyOrder[u.ID]

		if ok {
			results[ix] = &dataloader.Result{Data: u, Error: nil}
			delete(keyOrder, u.ID)
		}
	}

	for userID, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: &oops.CodedError{
			HumanMessage: fmt.Sprintf("user %d not found", userID),
			Context:      fmt.Sprintf("user %d not found while dataloading", userID),
			HttpStatus:   http.StatusNotFound,
			Err:          nil,
		}}
	}
	return results
}

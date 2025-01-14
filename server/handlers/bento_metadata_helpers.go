package handlers

import (
	"context"
	"konbini/server/db"
	"sync"
)

type metadataChannelResult struct {
	Err  error
	Data []string
	From string
}

func getBentoIngredientIDs(ctx context.Context, q *db.Queries, bentoID string, wg *sync.WaitGroup, results chan<- metadataChannelResult) {
	defer wg.Done()
	ids, err := q.GetBentoIngredientIDsInBento(ctx, bentoID)
	if err != nil {
		results <- metadataChannelResult{
			Err:  err,
			From: "get_bento_ingredients",
		}
	}
	results <- metadataChannelResult{Err: nil, Data: ids, From: "get_bento_ingredients"}
}

func getUsersWithAccess(ctx context.Context, q *db.Queries, bentoID string, wg *sync.WaitGroup, results chan<- metadataChannelResult) {
	defer wg.Done()
	ids, err := q.GetUserIDsWithBentoAccess(ctx, bentoID)
	if err != nil {
		results <- metadataChannelResult{
			Err:  err,
			From: "get_bento_access",
		}
	}
	results <- metadataChannelResult{Err: nil, Data: ids, From: "get_bento_access"}
}

func getGroupsWithAccess(ctx context.Context, q *db.Queries, bentoID string, wg *sync.WaitGroup, results chan<- metadataChannelResult) {
	defer wg.Done()
	ids, err := q.GetGroupIDsWithBentoAccess(ctx, bentoID)
	if err != nil {
		results <- metadataChannelResult{
			Err:  err,
			From: "get_group_access",
		}
	}
	results <- metadataChannelResult{Err: nil, Data: ids, From: "get_group_access"}
}

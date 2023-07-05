package cosmos

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

var ErrNoMoreItems = errors.New("no more items")

func IsNoMoreItems(err error) bool {
	return errors.Is(err, ErrNoMoreItems)
}

func Query(client *azcosmos.ContainerClient, partitionKey azcosmos.PartitionKey, query string) (next func(ctx context.Context, v any) (err error)) {
	pager := client.NewQueryItemsPager(query, partitionKey, nil)
	var curr [][]byte

	return func(ctx context.Context, v any) (err error) {
		if len(curr) == 0 {
			if !pager.More() {
				return ErrNoMoreItems
			}

			resp, err := pager.NextPage(ctx)
			if err != nil {
				return err
			}
			if len(resp.Items) == 0 {
				return ErrNoMoreItems
			}

			curr = resp.Items
		}

		err, curr = json.Unmarshal(curr[0], v), curr[1:]
		return
	}
}

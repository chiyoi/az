package cosmos

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/chiyoi/az"
)

func KeyExist(ctx context.Context, client *azcosmos.ContainerClient, partitionKey azcosmos.PartitionKey, itemID string) (exist bool, err error) {
	if _, err = client.ReadItem(ctx, partitionKey, itemID, nil); err != nil {
		if az.IsNotFound(err) {
			return false, nil
		}
		return
	} else {
		return true, nil
	}
}

func CreateItem(ctx context.Context, client *azcosmos.ContainerClient, partitionKey azcosmos.PartitionKey, item any) (err error) {
	data, err := json.Marshal(item)
	if err != nil {
		return
	}

	_, err = client.CreateItem(ctx, partitionKey, data, nil)
	return
}

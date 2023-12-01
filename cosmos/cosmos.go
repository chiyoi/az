package cosmos

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/chiyoi/az"
)

func NewClient(endpoint string, o *azcosmos.ClientOptions) func(cred azcore.TokenCredential) (*azcosmos.Client, error) {
	return func(cred azcore.TokenCredential) (*azcosmos.Client, error) {
		return azcosmos.NewClient(endpoint, cred, o)
	}
}

func KeyExist(ctx context.Context, client *azcosmos.ContainerClient, partitionKey azcosmos.PartitionKey, itemID string) (exist bool, err error) {
	if itemID == "" {
		return false, nil
	}

	_, err = client.ReadItem(ctx, partitionKey, itemID, nil)
	if err != nil {
		if az.IsNotFound(err) {
			return false, nil
		}
		return
	}
	return true, nil
}

func CreateItem(ctx context.Context, client *azcosmos.ContainerClient, partitionKey azcosmos.PartitionKey, item any) (err error) {
	data, err := json.Marshal(item)
	if err != nil {
		return
	}

	_, err = client.CreateItem(ctx, partitionKey, data, nil)
	return
}

type Schema struct {
	PartitionKeyPath string
	TimeToLiveOn     bool
	TimeToLive       int32
}

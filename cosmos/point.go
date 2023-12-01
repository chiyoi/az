package cosmos

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func PointKeyExist(ctx context.Context, key string) func(client *azcosmos.ContainerClient) (exist bool, err error) {
	return func(client *azcosmos.ContainerClient) (exist bool, err error) {
		return KeyExist(ctx, client, azcosmos.NewPartitionKeyString(key), key)
	}
}

func PointCreate(ctx context.Context, key string, item any) func(client *azcosmos.ContainerClient) (err error) {
	return func(client *azcosmos.ContainerClient) (err error) {
		data, err := json.Marshal(item)
		if err != nil {
			return
		}
		_, err = client.CreateItem(ctx, azcosmos.NewPartitionKeyString(key), data, nil)
		return
	}
}

func PointDelete(ctx context.Context, key string) func(client *azcosmos.ContainerClient) (err error) {
	return func(client *azcosmos.ContainerClient) (err error) {
		_, err = client.DeleteItem(ctx, azcosmos.NewPartitionKeyString(key), key, nil)
		return
	}
}

func PointRead(ctx context.Context, key string, item any) func(client *azcosmos.ContainerClient) (err error) {
	return func(client *azcosmos.ContainerClient) (err error) {
		itemResponse, err := client.ReadItem(ctx, azcosmos.NewPartitionKeyString(key), key, nil)
		if err != nil {
			return
		}
		return json.Unmarshal(itemResponse.Value, item)
	}
}

func PointPatch(ctx context.Context, key string, patch azcosmos.PatchOperations) func(client *azcosmos.ContainerClient) (err error) {
	return func(client *azcosmos.ContainerClient) (err error) {
		_, err = client.PatchItem(ctx, azcosmos.NewPartitionKeyString(key), key, patch, nil)
		return
	}
}

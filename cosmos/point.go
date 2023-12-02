package cosmos

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func PointExist(ctx context.Context, key string) func(client *azcosmos.ContainerClient) (exist bool, err error) {
	return func(client *azcosmos.ContainerClient) (exist bool, err error) {
		return Exist(ctx, client, azcosmos.NewPartitionKeyString(key), key)
	}
}

func PointCreate(ctx context.Context, key string, value any) func(client *azcosmos.ContainerClient) (err error) {
	return func(client *azcosmos.ContainerClient) (err error) {
		var item = struct {
			ID    string `json:"id"`
			Value any    `json:"value"`
		}{key, value}
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

func PointRead(ctx context.Context, key string, value any) func(client *azcosmos.ContainerClient) (err error) {
	return func(client *azcosmos.ContainerClient) (err error) {
		resp, err := client.ReadItem(ctx, azcosmos.NewPartitionKeyString(key), key, nil)
		if err != nil {
			return
		}
		var item struct {
			Value json.RawMessage `json:"value"`
		}
		err = json.Unmarshal(resp.Value, &item)
		if err != nil {
			return
		}
		err = json.Unmarshal(item.Value, value)
		return
	}
}

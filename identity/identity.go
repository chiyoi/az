package identity

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func DefaultCredential() (credential azcore.TokenCredential, err error) {
	if credential, err = azidentity.NewDefaultAzureCredential(nil); err != nil {
		return
	}
	return RetryCredential(credential), err
}

func RetryCredential(credential azcore.TokenCredential) azcore.TokenCredential {
	return retryCredential{credential}
}

type retryCredential struct {
	azcore.TokenCredential
}

func (c retryCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (token azcore.AccessToken, err error) {
	delay := []int{0, 2, 6, 14, 30}
	for _, d := range delay {
		if token, err = c.TokenCredential.GetToken(ctx, options); err == nil {
			return
		}

		time.Sleep(time.Duration(d) * time.Second)
	}
	return
}

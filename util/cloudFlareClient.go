package util

import (
	"github.com/cloudflare/cloudflare-go"
)

type CloudFlareClient struct {
	Client    *cloudflare.API
	AccountId string
}

func NewCloudFlareClient(accountId string, apiKey string, email string) (*CloudFlareClient, error) {
	api, err := cloudflare.New(apiKey, email)
	if err != nil {
		return nil, err
	}
	return &CloudFlareClient{Client: api, AccountId: accountId}, nil
}

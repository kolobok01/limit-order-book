package data

import (
	"github.com/coinbase-samples/advanced-trade-sdk-go/client"
	"github.com/coinbase-samples/advanced-trade-sdk-go/credentials"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/kolobublik/limit-order-book/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewFIXClient,
	NewOrderBookRepo,
)

// Data .
type Data struct {
	CoinbaseClient client.RestClient
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	creds := &credentials.Credentials{
		AccessKey:     c.Coinbase.ApiKey,
		PrivatePemKey: c.Coinbase.ApiSecret,
	}

	httpClient, err := client.DefaultHttpClient()
	if err != nil {
		return nil, nil, err
	}

	coinbaseClient := client.NewRestClient(creds, httpClient)

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{CoinbaseClient: coinbaseClient}, cleanup, nil
}

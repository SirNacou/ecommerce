package clients

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	"github.com/SirNacou/ecommerce/services/catalog-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/catalog-service/gen/v1/catalogv1connect")

// CatalogClient is an internal client for calling the catalog service to
// resolve authoritative product prices.
type CatalogClient struct {
	client catalogv1connect.CatalogServiceClient
}

func NewCatalogClient(baseURL string) *CatalogClient {
	httpClient := &http.Client{}
	return &CatalogClient{
		client: catalogv1connect.NewCatalogServiceClient(httpClient, baseURL),
	}
}

// ResolvePrices returns the authoritative products for the given ids, keyed by
// product id. Products that are not found are simply omitted.
func (c *CatalogClient) ResolvePrices(ctx context.Context, productIDs []string) (map[string]int64, error) {
	if len(productIDs) == 0 {
		return map[string]int64{}, nil
	}

	resp, err := c.client.GetProductsByIds(ctx, connect.NewRequest(&catalogv1.GetProductsByIdsRequest{
		Ids: productIDs,
	}))
	if err != nil {
		return nil, fmt.Errorf("catalog get products by ids: %w", err)
	}

	prices := make(map[string]int64, len(resp.Msg.GetProducts()))
	for _, p := range resp.Msg.GetProducts() {
		prices[p.GetId()] = p.GetPriceCents()
	}

	return prices, nil
}
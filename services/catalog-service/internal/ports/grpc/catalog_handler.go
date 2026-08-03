package grpc

import (
	"context"

	"connectrpc.com/connect"

	catalogv1 "github.com/SirNacou/ecommerce/services/catalog-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/catalog-service/gen/v1/catalogv1connect"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/app"
)

type CatalogHandler struct {
	catalogv1connect.UnimplementedCatalogServiceHandler
	createProductUC *app.CreateProductUseCase
}

func NewCatalogHandler(createProductUC *app.CreateProductUseCase) *CatalogHandler {
	return &CatalogHandler{
		createProductUC: createProductUC,
	}
}

func (h *CatalogHandler) CreateProduct(
	ctx context.Context,
	req *connect.Request[catalogv1.CreateProductRequest],
) (*connect.Response[catalogv1.CreateProductResponse], error) {

	product, err := h.createProductUC.Execute(ctx, app.CreateProductInput{
		CategoryID:    req.Msg.GetCategoryId(),
		Name:          req.Msg.GetName(),
		Description:   req.Msg.GetDescription(),
		PriceCents:    req.Msg.GetPriceCents(),
		StockQuantity: req.Msg.GetStockQuantity(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&catalogv1.CreateProductResponse{
		Product: &catalogv1.Product{
			Id:            product.ID().String(),
			CategoryId:    product.CategoryID().String(),
			Name:          product.Name(),
			Description:   product.Description(),
			PriceCents:    product.PriceCents(),
			StockQuantity: product.StockQuantity(),
			CreatedAt:     product.CreatedAt().Format("2006-01-02T15:04:05Z"),
		},
	}), nil
}

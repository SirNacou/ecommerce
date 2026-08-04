package grpc

import (
	"context"
	"errors"
	"strconv"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	catalogv1 "github.com/SirNacou/ecommerce/services/catalog-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/catalog-service/gen/v1/catalogv1connect"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/app"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type CatalogHandler struct {
	catalogv1connect.UnimplementedCatalogServiceHandler
	createProductCmd  *app.CreateProductCommandHandler
	getProductQry     *app.GetProductQueryHandler
	listProductsQry   *app.ListProductsQueryHandler
	createCategoryCmd *app.CreateCategoryCommandHandler
	listCategoriesQry *app.ListCategoriesQueryHandler
}

func NewCatalogHandler(
	createProductCmd *app.CreateProductCommandHandler,
	getProductQry *app.GetProductQueryHandler,
	listProductsQry *app.ListProductsQueryHandler,
	createCategoryCmd *app.CreateCategoryCommandHandler,
	listCategoriesQry *app.ListCategoriesQueryHandler,
) *CatalogHandler {
	return &CatalogHandler{
		createProductCmd:  createProductCmd,
		getProductQry:     getProductQry,
		listProductsQry:   listProductsQry,
		createCategoryCmd: createCategoryCmd,
		listCategoriesQry: listCategoriesQry,
	}
}

// -----------------------------------------------------------------------------
// Category Methods
// -----------------------------------------------------------------------------

func (h *CatalogHandler) CreateCategory(
	ctx context.Context,
	req *connect.Request[catalogv1.CreateCategoryRequest],
) (*connect.Response[catalogv1.CreateCategoryResponse], error) {
	category, err := h.createCategoryCmd.Handle(ctx, app.CreateCategoryCommand{
		Name: req.Msg.GetName(),
		Slug: req.Msg.GetSlug(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidName) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&catalogv1.CreateCategoryResponse{
		Category: toProtoCategory(category),
	}), nil
}

func (h *CatalogHandler) ListCategories(
	ctx context.Context,
	req *connect.Request[catalogv1.ListCategoriesRequest],
) (*connect.Response[catalogv1.ListCategoriesResponse], error) {
	categories, err := h.listCategoriesQry.Handle(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoCategories := make([]*catalogv1.Category, 0, len(categories))
	for _, cat := range categories {
		protoCategories = append(protoCategories, toProtoCategory(cat))
	}

	return connect.NewResponse(&catalogv1.ListCategoriesResponse{
		Categories: protoCategories,
	}), nil
}

// -----------------------------------------------------------------------------
// Product Methods
// -----------------------------------------------------------------------------

func (h *CatalogHandler) CreateProduct(
	ctx context.Context,
	req *connect.Request[catalogv1.CreateProductRequest],
) (*connect.Response[catalogv1.CreateProductResponse], error) {
	product, err := h.createProductCmd.Handle(ctx, app.CreateProductCommand{
		CategoryID:    req.Msg.GetCategoryId(),
		Name:          req.Msg.GetName(),
		Description:   req.Msg.GetDescription(),
		PriceCents:    req.Msg.GetPriceCents(),
		StockQuantity: req.Msg.GetStockQuantity(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidName) ||
			errors.Is(err, domain.ErrInvalidPrice) ||
			errors.Is(err, domain.ErrCategoryNotFound) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&catalogv1.CreateProductResponse{
		Product: toProtoProduct(product),
	}), nil
}

func (h *CatalogHandler) GetProduct(
	ctx context.Context,
	req *connect.Request[catalogv1.GetProductRequest],
) (*connect.Response[catalogv1.GetProductResponse], error) {
	product, err := h.getProductQry.Handle(ctx, req.Msg.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&catalogv1.GetProductResponse{
		Product: toProtoProduct(product),
	}), nil
}

func (h *CatalogHandler) ListProducts(
	ctx context.Context,
	req *connect.Request[catalogv1.ListProductsRequest],
) (*connect.Response[catalogv1.ListProductsResponse], error) {
	pageSize := req.Msg.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}

	var offset int32 = 0
	if token := req.Msg.GetPageToken(); token != "" {
		if parsed, err := strconv.Atoi(token); err == nil && parsed >= 0 {
			offset = int32(parsed)
		}
	}

	var categoryId *uuid.UUID
	if id, err := uuid.Parse(req.Msg.GetCategoryId()); err == nil {
		categoryId = &id
	}

	query := app.ListProductsQuery{
		PageSize:   pageSize,
		Offset:     offset,
		CategoryID: categoryId,
	}

	products, err := h.listProductsQry.Handle(ctx, query)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoProducts := make([]*catalogv1.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, toProtoProduct(p))
	}

	// Simple offset-based cursor calculation for next_page_token
	var nextPageToken string
	if int32(len(products)) == pageSize {
		nextPageToken = strconv.Itoa(int(offset + pageSize))
	}

	return connect.NewResponse(&catalogv1.ListProductsResponse{
		Products:      protoProducts,
		NextPageToken: nextPageToken,
	}), nil
}

// -----------------------------------------------------------------------------
// Proto Mapper Helpers
// -----------------------------------------------------------------------------

func toProtoProduct(p *domain.Product) *catalogv1.Product {
	return &catalogv1.Product{
		Id:            p.ID.String(),
		CategoryId:    p.CategoryID.String(),
		Name:          p.Name,
		Description:   p.Description,
		PriceCents:    p.PriceCents,
		StockQuantity: p.StockQuantity,
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
}

func toProtoCategory(c *domain.Category) *catalogv1.Category {
	return &catalogv1.Category{
		Id:   c.ID.String(),
		Name: c.Name,
		Slug: c.Slug,
	}
}

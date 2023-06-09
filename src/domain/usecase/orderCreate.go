package usecase

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/products-orders-service/src/domain/action/command"
	"github.com/unq-arq2-ecommerce-team/products-orders-service/src/domain/action/query"
	"github.com/unq-arq2-ecommerce-team/products-orders-service/src/domain/model"
	"github.com/unq-arq2-ecommerce-team/products-orders-service/src/domain/model/exception"
	"time"
)

type CreateOrder struct {
	baseLogger           model.Logger
	createOrderCmd       command.CreateOrder
	findProductByIdQuery query.FindProductById
}

func NewCreateOrder(baseLogger model.Logger, createOrderCmd command.CreateOrder, findProductByIdQuery query.FindProductById) *CreateOrder {
	return &CreateOrder{
		baseLogger:           baseLogger.WithFields(model.LoggerFields{"useCase": "CreateOrder"}),
		createOrderCmd:       createOrderCmd,
		findProductByIdQuery: findProductByIdQuery,
	}
}

func (u CreateOrder) Do(ctx context.Context, customerId, productId int64, deliveryDate time.Time, deliveryAddress model.Address) (int64, error) {
	log := u.baseLogger.WithRequestId(ctx).WithFields(model.LoggerFields{"customerId": customerId, "productId": productId, "deliveryDate": deliveryDate, "deliveryAddress": deliveryAddress})
	product, err := u.findProductByIdQuery.Do(ctx, productId)
	if err != nil {
		log.WithFields(model.LoggerFields{"error": err}).Errorf("error when find product")
		return 0, err
	}
	if !product.ReduceStock() {
		log.Infof("product with stock %v is not available", product.Stock)
		return 0, exception.ProductWithNoStock{Id: productId}
	}
	order := model.NewOrder(customerId, product, deliveryDate, deliveryAddress)
	orderId, err := u.createOrderCmd.Do(ctx, order)
	if err != nil {
		log.WithFields(model.LoggerFields{"error": err, "order": order}).Errorf("error when create order")
		return 0, err
	}
	log.Infof("successful order created with id %v", orderId)
	return orderId, nil
}

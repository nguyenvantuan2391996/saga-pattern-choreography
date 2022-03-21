package api

import (
	"context"
	"github.com/itimofeev/go-saga"
	"github.com/jinzhu/gorm"
	"saga-pattern-choreography/order-service/model"
	"saga-pattern-choreography/order-service/payment-client"
	"saga-pattern-choreography/order-service/stock-client"
	"saga-pattern-choreography/proto/order"
)

type orderServiceServer struct {
	db             *gorm.DB
	paymentService *payment_client.PaymentClient
	stockService   *stock_client.StockClient
}

func NewOrderServiceServer(db *gorm.DB, paymentClient *payment_client.PaymentClient, stockClient *stock_client.StockClient) order.OrderServiceServer {
	return &orderServiceServer{
		db:             db,
		paymentService: paymentClient,
		stockService:   stockClient,
	}
}

func (s *orderServiceServer) CreateOrder(ctx context.Context, request *order.CreateRequest) (*order.CreateResponse, error) {
	paymentDB, err := s.paymentService.GetPayment(ctx, request.Id, request.UserId)
	if err != nil {
		return nil, err
	}

	stockDB, err := s.stockService.GetStock(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	sagaPattern := saga.NewSaga("order-product")
	err = sagaPattern.AddStep(&saga.Step{
		Name: "create order",
		Func: func(context.Context) error {
			err := s.db.Create(&model.Order{
				ID:     request.Id,
				UserID: request.UserId,
				Status: request.Status,
			}).Error
			if err != nil {
				return err
			}
			return nil
		},
		CompensateFunc: func(context.Context) error {
			_, err = s.UpdateOrder(ctx, &order.CreateRequest{
				Id:     request.Id,
				UserId: request.UserId,
				Status: "pending",
			})
			if err != nil {
				return err
			}
			return nil
		},
		Options: nil,
	})
	if err != nil {
		return nil, err
	}

	err = sagaPattern.AddStep(&saga.Step{
		Name: "update payment",
		Func: func(context.Context) error {
			_, err = s.paymentService.MinusPayment(ctx, request.Id, request.UserId, 500)
			if err != nil {
				return err
			}
			return nil
		},
		CompensateFunc: func(context.Context) error {
			_, err = s.paymentService.UpdatePayment(ctx, request.Id, request.UserId, paymentDB.Money)
			if err != nil {
				return err
			}
			return nil
		},
		Options: nil,
	})
	if err != nil {
		return nil, err
	}

	err = sagaPattern.AddStep(&saga.Step{
		Name: "update stock",
		Func: func(context.Context) error {
			_, err = s.stockService.MinusStock(ctx, request.Id, 8)
			if err != nil {
				return err
			}
			return nil
		},
		CompensateFunc: func(context.Context) error {
			_, err = s.stockService.UpdateStock(ctx, request.Id, stockDB.Stock)
			if err != nil {
				return err
			}
			return nil
		},
		Options: nil,
	})
	if err != nil {
		return nil, err
	}

	err = sagaPattern.AddStep(&saga.Step{
		Name: "update order success",
		Func: func(context.Context) error {
			_, err = s.UpdateOrder(ctx, &order.CreateRequest{
				Id:     request.Id,
				UserId: request.UserId,
				Status: "success",
			})
			if err != nil {
				return err
			}
			return nil
		},
		CompensateFunc: func(context.Context) error {
			_, err = s.UpdateOrder(ctx, &order.CreateRequest{
				Id:     request.Id,
				UserId: request.UserId,
				Status: "pending",
			})
			if err != nil {
				return err
			}
			return nil
		},
		Options: nil,
	})
	if err != nil {
		return nil, err
	}

	store := saga.New()
	c := saga.NewCoordinator(ctx, ctx, sagaPattern, store)
	err = c.Play().ExecutionError
	if err != nil {
		return nil, err
	}

	return &order.CreateResponse{
		Message: "success",
	}, nil
}

func (s *orderServiceServer) UpdateOrder(ctx context.Context, request *order.CreateRequest) (*order.UpdateResponse, error) {
	err := s.db.Model(&model.Order{}).Where("id = ? AND user_id = ?", request.Id, request.UserId).Update("status", request.Status).Error

	if err != nil {
		return &order.UpdateResponse{
			Message: "Failed",
		}, err
	}

	return &order.UpdateResponse{
		Message: request.Status,
	}, nil
}

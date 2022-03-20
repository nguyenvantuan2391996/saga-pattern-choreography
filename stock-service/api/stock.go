package api

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	stock2 "saga-pattern-choreography/proto/stock"
	"saga-pattern-choreography/stock-service/model"
)

type stockServiceServer struct {
	db *gorm.DB
}

func NewStockServiceServer(db *gorm.DB) stock2.StockServiceServer {
	return &stockServiceServer{db: db}
}

func (s *stockServiceServer) UpdateStock(ctx context.Context, request *stock2.StockRequest) (*stock2.StockResponse, error) {
	stockDB := &model.Stock{}
	if err := s.db.Where("id = ?", request.Id).Find(&stockDB).Error; err != nil {
		return nil, err
	}
	if stockDB.Stock-request.Stock < 0 {
		return nil, errors.New("stock < 0")
	} else {
		if err := s.db.Save(&model.Stock{
			ID:    request.Id,
			Stock: stockDB.Stock - request.Stock,
		}).Error; err != nil {
			return nil, err
		}
	}
	return &stock2.StockResponse{
		Message: "success",
	}, nil
}

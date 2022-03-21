package api

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"saga-pattern-choreography/proto/stock"
	"saga-pattern-choreography/stock-service/model"
)

type stockServiceServer struct {
	db *gorm.DB
}

func NewStockServiceServer(db *gorm.DB) stock.StockServiceServer {
	return &stockServiceServer{db: db}
}

func (s *stockServiceServer) GetStock(ctx context.Context, request *stock.StockRequest) (*stock.Stock, error) {
	stockDB := &model.Stock{}
	if err := s.db.Where("id = ?", request.Id).Find(&stockDB).Error; err != nil {
		return nil, err
	}
	return &stock.Stock{
		Id:    stockDB.ID,
		Stock: stockDB.Stock,
	}, nil
}

func (s *stockServiceServer) MinusStock(ctx context.Context, request *stock.StockRequest) (*stock.StockResponse, error) {
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
	return &stock.StockResponse{
		Message: "success",
	}, nil
}

func (s *stockServiceServer) UpdateStock(ctx context.Context, request *stock.StockRequest) (*stock.StockResponse, error) {
	err := s.db.Model(&model.Stock{}).Where("id = ?", request.Id).Update("stock", request.Stock).Error
	if err != nil {
		return nil, err
	}
	return &stock.StockResponse{
		Message: "success",
	}, nil
}

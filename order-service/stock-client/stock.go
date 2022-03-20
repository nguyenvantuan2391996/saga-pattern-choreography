package stock_client

import (
	"context"
	"google.golang.org/grpc"
	"saga-pattern-choreography/proto/stock"
)

type StockClient struct {
	stockClient stock.StockServiceClient
}

func NewStockClient(stockURL string) *StockClient {
	conn, err := grpc.Dial(stockURL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	stockClient := stock.NewStockServiceClient(conn)
	return &StockClient{stockClient: stockClient}
}

func (s *StockClient) UpdateStock(ctx context.Context, id, stockBuy int32) (*stock.StockResponse, error) {
	res, err := s.stockClient.UpdateStock(ctx, &stock.StockRequest{
		Id:    id,
		Stock: stockBuy,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

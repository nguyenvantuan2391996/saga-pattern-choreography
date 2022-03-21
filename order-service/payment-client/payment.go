package payment_client

import (
	"context"
	"google.golang.org/grpc"
	"saga-pattern-choreography/proto/payment"
)

type PaymentClient struct {
	paymentClient payment.PaymentServiceClient
}

func NewPaymentClient(paymentURL string) *PaymentClient {
	conn, err := grpc.Dial(paymentURL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	paymentClient := payment.NewPaymentServiceClient(conn)
	return &PaymentClient{paymentClient: paymentClient}
}

func (p *PaymentClient) GetPayment(ctx context.Context, id, userID int32) (*payment.Payment, error) {
	res, err := p.paymentClient.GetPayment(ctx, &payment.PaymentRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *PaymentClient) MinusPayment(ctx context.Context, id, userID, money int32) (*payment.PaymentResponse, error) {
	res, err := p.paymentClient.MinusPayment(ctx, &payment.PaymentRequest{
		Id:     id,
		UserId: userID,
		Money:  money,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *PaymentClient) UpdatePayment(ctx context.Context, id, userID, money int32) (*payment.PaymentResponse, error) {
	res, err := p.paymentClient.UpdatePayment(ctx, &payment.PaymentRequest{
		Id:     id,
		UserId: userID,
		Money:  money,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

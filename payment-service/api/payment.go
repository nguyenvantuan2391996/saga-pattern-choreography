package api

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"saga-pattern-choreography/payment-service/model"
	"saga-pattern-choreography/proto/payment"
)

type paymentServiceServer struct {
	db *gorm.DB
}

func NewPaymentServiceServer(db *gorm.DB) payment.PaymentServiceServer {
	return &paymentServiceServer{db: db}
}

func (s *paymentServiceServer) GetPayment(ctx context.Context, request *payment.PaymentRequest) (*payment.Payment, error) {
	paymentDB := &model.Payment{}
	if err := s.db.Where("id = ? AND user_id = ?", request.Id, request.UserId).Find(&paymentDB).Error; err != nil {
		return nil, err
	}
	return &payment.Payment{
		Id:     paymentDB.ID,
		UserId: paymentDB.UserID,
		Money:  paymentDB.Balance,
	}, nil
}

func (s *paymentServiceServer) MinusPayment(ctx context.Context, request *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	paymentDB := &model.Payment{}
	if err := s.db.Where("id = ? AND user_id = ?", request.Id, request.UserId).Find(&paymentDB).Error; err != nil {
		return nil, err
	}
	if paymentDB.Balance-request.Money < 0 {
		return nil, errors.New("balance < money")
	} else {
		if err := s.db.Save(&model.Payment{
			ID:      request.Id,
			UserID:  request.UserId,
			Balance: paymentDB.Balance - request.Money,
		}).Error; err != nil {
			return nil, err
		}
	}
	return &payment.PaymentResponse{
		Message: "success",
	}, nil
}

func (s *paymentServiceServer) UpdatePayment(ctx context.Context, request *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	err := s.db.Model(&model.Payment{}).Where("id = ? AND user_id = ?", request.Id, request.UserId).Update("balance", request.Money).Error
	if err != nil {
		return nil, err
	}
	return &payment.PaymentResponse{
		Message: "success",
	}, nil
}

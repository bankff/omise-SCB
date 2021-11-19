package scb

import (
	"context"
	"omisescb/domain/dto"
	"omisescb/domain/model"
	omiseClient "omisescb/external/omise"
	"omisescb/repository"
	"strconv"
	"time"

	"github.com/omise/omise-go"
)

type SCBPayment interface {
	CreatePayment(context.Context, SCBPaymentRequest) (SCBPaymentReponse, error)
	GetPaymentCallBack(context.Context, string) (SCBPaymentCallBackReponse, error)
	GetPayment(context.Context, SCBGetPaymentRequest) (SCBGetPaymentResponse, error)
}

type scbPayment struct {
	omise omiseClient.OmiseEx
	db    repository.OrderTransaction
}

func NewSCBPayment(o omiseClient.OmiseEx,
	db repository.OrderTransaction) SCBPayment {
	return scbPayment{
		omise: o,
		db:    db,
	}
}

func (o scbPayment) CreatePayment(ctx context.Context, s SCBPaymentRequest) (SCBPaymentReponse, error) {
	orderID := time.Now().UnixNano()
	source, err := o.omise.CreateSoure(s.Amount, dto.TH, dto.InternetBankingSCB)
	if err != nil {
		return SCBPaymentReponse{}, err
	}

	charge, err := o.omise.CreateCharge(source.Amount, source.Currency, source.ID, orderID)
	if err != nil {
		return SCBPaymentReponse{}, err
	}

	err = o.db.CreateOrder(ctx, model.OrderTransaction{
		OrderID:       orderID,
		ChargeID:      charge.ID,
		Amount:        charge.Amount,
		Currency:      charge.Currency,
		PaymentStatus: string(charge.Status),
		SoureID:       charge.Source.ID,
		SoureType:     charge.Source.Type,
		Create_at:     &charge.Created,
	})
	if err != nil {
		return SCBPaymentReponse{}, err
	}

	return SCBPaymentReponse{
		Charge: *charge,
	}, nil
}

func (o scbPayment) GetPaymentCallBack(ctx context.Context, ID string) (SCBPaymentCallBackReponse, error) {
	orderID, _ := strconv.ParseInt(ID, 10, 64)
	order, err := o.db.GetByOrderID(ctx, orderID)
	if err != nil {
		return SCBPaymentCallBackReponse{}, err
	}

	charge, err := o.omise.GetCharge(order.ChargeID)
	if err != nil {
		return SCBPaymentCallBackReponse{}, err
	}
	if charge.Status == omise.ChargeSuccessful {
		order.PaidAt = tmp(time.Now())
	}

	if err := o.db.Update(ctx, model.OrderTransaction{
		ID:            order.ID,
		PaymentStatus: string(charge.Status),
		PaidAt:        order.PaidAt,
	}); err != nil {
		return SCBPaymentCallBackReponse{}, err
	}

	return SCBPaymentCallBackReponse{
		OrderID:       ID,
		PaymentStatus: string(charge.Status),
	}, nil
}

func (o scbPayment) GetPayment(ctx context.Context, req SCBGetPaymentRequest) (SCBGetPaymentResponse, error) {
	orders, err := o.db.Get(ctx, req.OrderID, req.PaymentStatus)
	if err != nil {
		return SCBGetPaymentResponse{}, err
	}
	return SCBGetPaymentResponse{
		Status: "success",
		Data:   orders,
	}, nil
}

func tmp(t time.Time) *time.Time {
	return &t
}

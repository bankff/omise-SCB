package scb

import (
	"context"
	"omisescb/domain/model"

	"github.com/go-kit/kit/endpoint"
	"github.com/omise/omise-go"
)

type SCBPaymentRequest struct {
	Amount int64 `json:"amount"`
}
type SCBPaymentReponse struct {
	omise.Charge
}

type SCBPaymentCallBackReponse struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
}

type SCBGetPaymentRequest struct {
	OrderID       int64
	PaymentStatus string
}
type SCBGetPaymentResponse struct {
	Status string                   `json:"status"`
	Data   []model.OrderTransaction `json:"Data"`
}

func makeCreateSCBPaymentEndPoint(s SCBPayment) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SCBPaymentRequest)
		response, err = s.CreatePayment(ctx, req)
		if err != nil {
			return err, nil
		}
		return response, nil
	}
}

func makeGetSCBPaymentCallBackEndPoint(s SCBPayment) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		response, err = s.GetPaymentCallBack(ctx, req)
		if err != nil {
			return err, nil
		}
		return response, nil
	}
}

func makeGetSCBPaymentEndPoint(s SCBPayment) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SCBGetPaymentRequest)
		response, err = s.GetPayment(ctx, req)
		if err != nil {
			return err, nil
		}
		return response, nil
	}
}

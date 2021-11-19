package scb

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeCreateSCBPaymentHandler(s SCBPayment) http.Handler {
	return kithttp.NewServer(
		makeCreateSCBPaymentEndPoint(s),
		decodeCreateSCBPaymentRequest,
		encodeResponse(),
	)
}

func MakeGetCallBackSCBPaymentHandler(s SCBPayment) http.Handler {
	return kithttp.NewServer(
		makeGetSCBPaymentCallBackEndPoint(s),
		decodeGetCallbackSCBPaymentRequest,
		encodeResponse(),
	)
}

func MakeGetSCBPaymentHandler(s SCBPayment) http.Handler {
	return kithttp.NewServer(
		makeGetSCBPaymentEndPoint(s),
		decodeGetSCBPaymentRequest,
		encodeResponse(),
	)
}

func decodeCreateSCBPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body SCBPaymentRequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}
	return body, nil
}

func decodeGetCallbackSCBPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	sourceID := mux.Vars(r)["source_id"]
	return sourceID, nil
}

func decodeGetSCBPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SCBGetPaymentRequest
	id := r.URL.Query().Get("orderID")
	if id != "" {
		orderID, err := strconv.ParseInt(r.URL.Query()["orderID"][0], 10, 64)
		if err != nil {
			return req, err
		}
		req.OrderID = orderID
	}
	req.PaymentStatus = r.URL.Query().Get("status")
	return req, nil
}

func encodeResponse() func(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(error); ok {
			return e
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		return json.NewEncoder(w).Encode(response)
	}
}

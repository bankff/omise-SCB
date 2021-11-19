package model

import "time"

type OrderTransaction struct {
	ID            int
	OrderID       int64      `json:"order_id"`
	ChargeID      string     `json:"charge_id"`
	Amount        int64      `json:"amount"`
	Currency      string     `json:"currency"`
	PaymentStatus string     `json:"payment_status"`
	SoureID       string     `json:"soure_id"`
	SoureType     string     `json:"soure_type"`
	PaidAt        *time.Time `json:"paid_at"`
	Create_at     *time.Time `json:"create_at"`
}

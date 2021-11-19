package repository

import (
	"context"
	"database/sql"
	"fmt"
	"omisescb/domain/model"
)

type OrderTransaction interface {
	CreateOrder(ctx context.Context, order model.OrderTransaction) error
	Update(ctx context.Context, order model.OrderTransaction) error
	GetByOrderID(ctx context.Context, orderID int64) (model.OrderTransaction, error)
	Get(ctx context.Context, orderID int64, status string) ([]model.OrderTransaction, error)
}
type orderTransaction struct {
	db *sql.DB
}

func NewOrderTransaction(db *sql.DB) OrderTransaction {
	return orderTransaction{
		db: db,
	}
}

func (o orderTransaction) CreateOrder(ctx context.Context, order model.OrderTransaction) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	stmt, _ := tx.PrepareContext(ctx, "insert into orders (order_id,charge_id,amount,currency,payment_status,soure_id,soure_type,create_at) values (?,?,?,?,?,?,?,?)")
	_, err = stmt.Exec(order.OrderID, order.ChargeID, order.Amount, order.Currency, order.PaymentStatus, order.SoureID, order.SoureType, order.Create_at)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (o orderTransaction) Update(ctx context.Context, order model.OrderTransaction) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	stmt, _ := tx.PrepareContext(ctx, "update orders set payment_status=?,paid_at=? where id=?")
	_, err = stmt.Exec(order.PaymentStatus, order.PaidAt, order.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (o orderTransaction) GetByOrderID(ctx context.Context, orderID int64) (model.OrderTransaction, error) {
	var order model.OrderTransaction
	rows, err := o.db.QueryContext(ctx, "SELECT * From orders where order_id=?", orderID)
	if err != nil {
		return order, err
	}
	for rows.Next() {
		err = rows.Scan(&order.ID, &order.OrderID, &order.ChargeID, &order.Amount, &order.Currency, &order.PaymentStatus, &order.SoureID, &order.SoureType, &order.PaidAt, &order.Create_at)
		if err != nil {
			return order, err

		}
	}
	defer rows.Close()

	return order, nil
}

func (o orderTransaction) Get(ctx context.Context, orderID int64, status string) ([]model.OrderTransaction, error) {
	var (
		orders []model.OrderTransaction
		order  model.OrderTransaction
		query  string
	)
	query = fmt.Sprintf("SELECT * From orders")

	if orderID != 0 && status != "" {
		query = fmt.Sprintf("%s where order_id = %d AND payment_status = '%v'", query, orderID, status)
	} else if orderID != 0 {
		query = fmt.Sprintf("%s where order_id = %d", query, orderID)
	} else if status != "" {
		query = fmt.Sprintf("%s where payment_status = '%v'", query, status)
	}
	rows, err := o.db.QueryContext(ctx, query)
	if err != nil {
		return orders, err
	}
	for rows.Next() {
		err = rows.Scan(&order.ID, &order.OrderID, &order.ChargeID, &order.Amount, &order.Currency, &order.PaymentStatus, &order.SoureID, &order.SoureType, &order.PaidAt, &order.Create_at)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	defer rows.Close()

	return orders, nil
}

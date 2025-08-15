package postgres

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
)

type OrderRepository struct {
	db *Database
}

func NewOrderRepository(db *Database) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (o *OrderRepository) GetLastNOrders(ctx context.Context, n int) (orders []domain.Order, err error) {
	err = o.db.db.SelectContext(ctx, &orders, `SELECT * FROM orders ORDER BY date_created DESC LIMIT $1`, n)
	if err != nil {
		return
	}

	if len(orders) == 0 {
		return
	}

	orderIDs := make([]string, len(orders))
	orderMap := make(map[string]*domain.Order, len(orders))
	for i := range orders {
		orderIDs[i] = orders[i].OrderUID
		orderMap[orders[i].OrderUID] = &orders[i]
	}

	var payments []domain.Payment
	query, args, err := sqlx.In(`SELECT * FROM payment WHERE order_id in (?);`, orderIDs)
	if err != nil {
		return nil, err
	}
	query = o.db.db.Rebind(query)
	err = o.db.db.SelectContext(ctx, &payments, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	for _, p := range payments {
		if order, ok := orderMap[p.OrderId]; ok {
			order.Payment = p
		}
	}

	var deliveries []domain.Delivery
	query, args, err = sqlx.In(`SELECT * FROM delivery WHERE order_id in (?);`, orderIDs)
	if err != nil {
		return nil, err
	}
	query = o.db.db.Rebind(query)
	err = o.db.db.SelectContext(ctx, &deliveries, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	for _, d := range deliveries {
		if order, ok := orderMap[d.OrderId]; ok {
			order.Delivery = d
		}
	}

	var items []domain.Item
	query, args, err = sqlx.In(`SELECT * FROM items WHERE order_id in (?);`, orderIDs)
	if err != nil {
		return nil, err
	}
	query = o.db.db.Rebind(query)
	err = o.db.db.SelectContext(ctx, &items, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	for _, item := range items {
		if order, ok := orderMap[item.OrderId]; ok {
			order.Items = append(order.Items, item)
		}
	}

	return
}

func (o *OrderRepository) GetByID(ctx context.Context, uid string) (order domain.Order, err error) {
	err = o.db.db.GetContext(ctx, &order, `SELECT * FROM orders WHERE id = $1`, uid)
	if err != nil {
		return
	}
	err = o.db.db.GetContext(ctx, &order.Payment, `SELECT * FROM payment WHERE order_id = $1`, uid)
	if err != nil {
		return
	}
	err = o.db.db.GetContext(ctx, &order.Delivery, `SELECT * FROM delivery WHERE order_id = $1`, uid)
	if err != nil {
		return
	}
	err = o.db.db.SelectContext(ctx, &order.Items, `SELECT * FROM items WHERE order_id = $1`, uid)
	if err != nil {
		return
	}
	return
}

func (o *OrderRepository) Create(ctx context.Context, order domain.Order) (err error) {
	tx, err := o.db.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				return
			}
		} else {
			err = tx.Commit()
			if err != nil {
				return
			}
		}
	}()

	query, args, err := goqu.Insert("orders").Rows(order).Returning(goqu.C("id")).ToSQL()
	if err != nil {
		return err
	}
	var orderId string
	err = tx.GetContext(ctx, &orderId, query, args...)
	if err != nil {
		return err
	}
	order.Delivery.OrderId = orderId
	order.Payment.OrderId = orderId
	for i, item := range order.Items {
		item.OrderId = orderId
		order.Items[i] = item
	}
	if err = o.createDelivery(ctx, tx, order.Delivery); err != nil {
		return err
	}
	if err = o.createPayment(ctx, tx, order.Payment); err != nil {
		return err
	}
	if err = o.createItems(ctx, tx, order.Items); err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) createItems(ctx context.Context, tx *sqlx.Tx, items []domain.Item) (err error) {
	rows := make([]any, len(items))
	for i, item := range items {
		rows[i] = item
	}
	query, args, err := goqu.Insert("items").Rows(rows...).ToSQL()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) createPayment(ctx context.Context, tx *sqlx.Tx, payment domain.Payment) (err error) {
	query, args, err := goqu.Insert("payment").Rows(payment).ToSQL()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) createDelivery(ctx context.Context, tx *sqlx.Tx, delivery domain.Delivery) (err error) {
	query, args, err := goqu.Insert("delivery").Rows(delivery).ToSQL()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

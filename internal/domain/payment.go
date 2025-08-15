package domain

type Payment struct {
	Id           int     `json:"-" db:"id" goqu:"skipinsert"`
	OrderId      string  `json:"-" db:"order_id"`
	Transaction  string  `json:"transaction" db:"transaction"`
	RequestID    string  `json:"request_id" db:"request_id"`
	Currency     string  `json:"currency" db:"currency"`
	Provider     string  `json:"provider" db:"provider"`
	Amount       float64 `json:"amount" db:"amount"`
	PaymentDt    int     `json:"payment_dt" db:"payment_dt"`
	Bank         string  `json:"bank" db:"bank"`
	DeliveryCost float64 `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total" db:"goods_total"`
	CustomFee    float64 `json:"custom_fee" db:"custom_fee"`
}

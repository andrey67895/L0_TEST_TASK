package domain

type Item struct {
	Id          int     `json:"-" db:"id" goqu:"skipinsert"`
	OrderId     string  `json:"-" db:"order_id"`
	ChrtID      int     `json:"chrt_id" db:"chrt_id"`
	TrackNumber string  `json:"track_number" db:"track_number"`
	Price       float64 `json:"price" db:"price"`
	Rid         string  `json:"rid" db:"rid"`
	Name        string  `json:"name" db:"name"`
	Sale        float64 `json:"sale" db:"sale"`
	Size        string  `json:"size" db:"size"`
	TotalPrice  float64 `json:"total_price" db:"total_price"`
	NmID        int     `json:"nm_id" db:"nm_id"`
	Brand       string  `json:"brand" db:"brand"`
	Status      int     `json:"status" db:"status"`
}

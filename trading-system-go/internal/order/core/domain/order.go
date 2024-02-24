package order


type Order struct {
	OrderId     int         `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"-"`
	UserId      int32       `json:"user_id"`
	OrderType   OrderType   `json:"order_type"`
	Type        string      `json:"type"`
	Price       float64     `json:"price"`
	Volume      float64     `json:"volume"`
	BuyingPair  string      `json:"buying_pair"`
	SellingPair string      `json:"selling_pair"`
	OrderStatus OrderStatus `json:"-"`
	CreatedAt   int64       `gorm:"autoCreateTime" json:"-"`
}

type OrderType int

const (
	_ OrderType = iota
	MarketOrderBuy
	LimitOrderBuy
	MarketOrderSell
	LimitOrderSell
)

type OrderStatus string

const (
	PartialyFilled OrderStatus = "PartialyFilled"
	Filled         OrderStatus = "Filled"
	NotFilled      OrderStatus = "NotFilled"
);
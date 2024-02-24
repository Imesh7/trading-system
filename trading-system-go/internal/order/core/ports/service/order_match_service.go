package ports

type OrderMatchService interface {
	MatchOrder(orderId int)
}
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	order "trading-system-go/internal/order/core/domain"
	ports "trading-system-go/internal/order/core/ports/repo"
	order_book "trading-system-go/internal/order_book/core/ports"
)

type orderMatchService struct {
	repo                  ports.OrderRepository
	orderBookService      order_book.OrderBookService
	orderBookKafkaService order_book.OrderBookKafka
}

func NewOrderMatchService(repo ports.OrderRepository, orderBookService order_book.OrderBookService, orderBookKafkaService order_book.OrderBookKafka) *orderMatchService {
	return &orderMatchService{
		repo:                  repo,
		orderBookService:      orderBookService,
		orderBookKafkaService: orderBookKafkaService,
	}
}

func (orderMatchService *orderMatchService) MatchOrder(orderId int) {
	//var order order.Order
	//first get the unfilled or partfilled order with opposite
	//order type of this particular order
	fmt.Println("order id.........",orderId)

	orderData, err := orderMatchService.repo.FindOrder(orderId)
	if err != nil {
		log.Panic(err)
		return
	}
	//database.DB.DataBase.Where("order_id = ?", orderId).First(&order)
	//try to match with the price.
	//if the price match , remove from the orderbook & update the order status filled
	ctx := context.Background()
	buyOrderkey := fmt.Sprintf("buy-orders")
	sellOrderkey := fmt.Sprintf("sell-orders")
	//Buy order
	if orderData.Type == "bid" {
		//first try to match with the redis after that get that order & validate it &
		//remove from the order & create a sell order against,

		sellOrderList := orderMatchService.orderBookService.GetOrderBook(ctx, sellOrderkey)
		fmt.Println("sell-order-list", sellOrderList)
		if len(sellOrderList) == 0 {
			//add a new buy order to the orderbook
			err := orderMatchService.orderBookService.AddToOrderBook(ctx, buyOrderkey, orderData.OrderId, serializeOrder(orderData))
			if err != nil {
				log.Fatal(err)
			}
			orderMatchService.orderBookKafkaService.CreateOrderBookProducer("order-book", orderData)

		} else {
			//sort the order with the price low to high
			sort.SliceStable(sellOrderList, func(i, j int) bool {
				//if price is equal sort with order createdAt(first fill the oldest order)
				if sellOrderList[i].Price == sellOrderList[j].Price {
					return sellOrderList[i].CreatedAt < sellOrderList[j].CreatedAt
				}
				return sellOrderList[i].Price < sellOrderList[j].Price
			})

			for _, v := range sellOrderList {
				//if bid price is equal or higher than the sell order in the order book
				if orderData.Price >= v.Price {
					if orderData.Volume <= v.Volume {
						//set the buy order status as fully filled
						//update the postgres order as filled

						//set the redis & postgres sell order to or if equal fully filled / partially filled
						balanceVolume := v.Volume - orderData.Volume
						if balanceVolume == 0 {
							//update the orderbook in fully filled
							v.OrderStatus = order.Filled
							//remove item from the order book
							errs := orderMatchService.orderBookService.RemoveFromOrderBook(ctx, sellOrderkey, v.OrderId)
							if errs != nil {
								fmt.Println(errs)
							}
						} else {
							//update the orderbook in partially filled
							v.OrderStatus = order.PartialyFilled
							v.Volume = v.Volume - orderData.Volume
							orderMatchService.orderBookService.UpdateOrderBook(ctx, sellOrderkey, v.OrderId, serializeOrder(&v))
						}
						break
					} else {
						//buy order should be partially filled
						//sell fully filled
						v.OrderStatus = order.Filled
						//update the orderbook in partially filled
						/* err := UpdateOrderBook(ctx, sellOrderkey, v.OrderId, serializeOrder(v))
						if err != nil {
							fmt.Println(err)
							} */
						orderData.Volume = orderData.Volume - v.Volume
						errs := orderMatchService.orderBookService.RemoveFromOrderBook(ctx, sellOrderkey, v.OrderId)
						if errs != nil {
							fmt.Println(errs)
						}
					}
				} else {
					//add a new buy order to the orderbook
					orderMatchService.orderBookService.AddToOrderBook(ctx, buyOrderkey, orderData.OrderId, serializeOrder(orderData))
					orderMatchService.orderBookKafkaService.CreateOrderBookProducer("order-book", orderData)
					break
				}
			}
		}
	} else if orderData.Type == "ask" {
		//first try to match with the redis after that get that order & validate it &
		//remove from the order & create a sell order against,
		buyOrderList := orderMatchService.orderBookService.GetOrderBook(ctx, buyOrderkey)
		if len(buyOrderList) == 0 {
			orderMatchService.orderBookService.AddToOrderBook(ctx, sellOrderkey, orderData.OrderId, serializeOrder(orderData))
			orderMatchService.orderBookKafkaService.CreateOrderBookProducer("order-book", orderData)
		} else {
			//sort the order with the price low to high
			sort.SliceStable(buyOrderList, func(i, j int) bool {
				//if price is equal sort with order createdAt(first fill the oldest order)
				if buyOrderList[i].Price == buyOrderList[j].Price {
					return buyOrderList[i].CreatedAt < buyOrderList[j].CreatedAt
				}
				return buyOrderList[i].Price > buyOrderList[j].Price
			})

			for _, v := range buyOrderList {

				//if bid price is equal or higher than the sell order in the order book
				if orderData.Price <= v.Price {
					if orderData.Volume <= v.Volume {
						//set the buy order status as fully filled
						//update the postgres order as filled

						//set the redis & postgres sell order to or if equal fully filled / partially filled
						balanceVolume := v.Volume - orderData.Volume
						if balanceVolume == 0 {
							//update the orderbook in fully filled
							v.OrderStatus = order.Filled
							/* err := UpdateOrderBook(ctx, buyOrderkey, v.OrderId, serializeOrder(v))
							if err != nil {
								fmt.Println(err)
							} */
							//remove item from the order book
							errs := orderMatchService.orderBookService.RemoveFromOrderBook(ctx, buyOrderkey, v.OrderId)
							if errs != nil {
								fmt.Println(errs)
							}
						} else {
							//update the orderbook in partially filled
							v.OrderStatus = order.PartialyFilled
							v.Volume = v.Volume - orderData.Volume
							orderMatchService.orderBookService.UpdateOrderBook(ctx, buyOrderkey, v.OrderId, serializeOrder(&v))
						}
						break
					} else {
						//buy order should be partially filled
						//sell fully filled
						v.OrderStatus = order.Filled
						//update the orderbook in partially filled
						/* err := UpdateOrderBook(ctx, buyOrderkey, v.OrderId, serializeOrder(v))
						if err != nil {
							fmt.Println(err)
						} */
						orderData.Volume = orderData.Volume - v.Volume
						errs := orderMatchService.orderBookService.RemoveFromOrderBook(ctx, buyOrderkey, v.OrderId)
						if errs != nil {
							fmt.Println(errs)
						}

						//recurrsion
					}
				} else {
					//add a new buy order to the orderbook
					orderMatchService.orderBookService.AddToOrderBook(ctx, sellOrderkey, orderData.OrderId, serializeOrder(orderData))
					orderMatchService.orderBookKafkaService.CreateOrderBookProducer("order-book", orderData)
					break
				}
			}
		}
	}
}

func serializeOrder(order *order.Order) string {
	data, _ := json.Marshal(order)
	return string(data)
}

package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"trading-system-go/internal/order/core/domain"

	"github.com/redis/go-redis/v9"
)

type orderBookRepository struct {
	client *redis.Client
}

func NewOrderBookepository(client *redis.Client) *orderBookRepository {
	return &orderBookRepository{
		client: client,
	}
}


// get latest orderbook data
func (repository orderBookRepository) GetOrderBook(ctx context.Context, orderType string) []order.Order {
	var orderList []order.Order
	result := repository.client.Keys(ctx, orderType)
	for _, v := range result.Val() {
		var order order.Order
		data := repository.client.HGetAll(ctx, v)
		for _, v := range data.Val() {
			json.Unmarshal([]byte(v), &order)
			orderList = append(orderList, order)
		}
	}
	return orderList
}

// add order(create) to orderbook
func (repository orderBookRepository) AddToOrderBook(ctx context.Context, orderKey string, orderId int, value string) error {
	err := repository.client.HSet(ctx, orderKey, orderId, value).Err()
	if err != nil {
		//return err
		panic(err)
	}

	return nil
}

// if order cancel or fully filled
func (repository orderBookRepository) RemoveFromOrderBook(ctx context.Context, orderType string, orderId int) error {
	fieldsToDelete := strconv.Itoa(orderId)
	err := repository.client.HDel(ctx, orderType, fieldsToDelete).Err()
	if err != nil {
		//return err
		fmt.Printf("Error deleting fields from hash %s: %v\n", fieldsToDelete, err)
		panic(err)
	}
	return nil
}

// if order partially filled
func (repository orderBookRepository) UpdateOrderBook(ctx context.Context, orderType string, orderId int, updatedOrder string) error {
	err := repository.client.HSet(ctx, orderType, orderId, updatedOrder).Err()
	if err != nil {
		return err
	}
	return nil
}

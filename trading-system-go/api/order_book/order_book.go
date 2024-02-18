package order_book

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"trading-system-go/database"
	"trading-system-go/internal/data/order"
)

// get latest orderbook data
func GetOrderBook(ctx context.Context, orderType string) []order.Order {
	var orderList []order.Order
	result := database.RedisDB.Client.Keys(ctx, orderType)
	for _, v := range result.Val() {
		var order order.Order
		data := database.RedisDB.Client.HGetAll(ctx, v)
		for _, v := range data.Val() {
			json.Unmarshal([]byte(v), &order)
			orderList = append(orderList, order)
		}
	}
	return orderList
}

// add order(create) to orderbook
func AddToOrderBook(ctx context.Context, orderKey string, orderId int, value string) error {
	err := database.RedisDB.Client.HSet(ctx, orderKey, orderId, value).Err()
	if err != nil {
		//return err
		panic(err)
	}

	return nil
}

// if order cancel or fully filled
func RemoveFromOrderBook(ctx context.Context, orderType string, orderId int) error {
	fieldsToDelete := strconv.Itoa(orderId)
	err := database.RedisDB.Client.HDel(ctx, orderType, fieldsToDelete).Err()
	if err != nil {
		//return err
		fmt.Printf("Error deleting fields from hash %s: %v\n", fieldsToDelete, err)
		panic(err)
	}
	return nil
}

// if order partially filled
func UpdateOrderBook(ctx context.Context, orderType string, orderId int, updatedOrder string) error {
	err := database.RedisDB.Client.HSet(ctx, orderType, orderId, updatedOrder).Err()
	if err != nil {
		return err
	}
	return nil
}

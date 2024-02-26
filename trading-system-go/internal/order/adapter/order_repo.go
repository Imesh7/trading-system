package adapter

import (
	"fmt"
	order "trading-system-go/internal/order/core/domain"

	"gorm.io/gorm"
)

type orderRepository struct {
	dbconn *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{
		dbconn: db,
	}
}

func (OrderRepositoryImpl *orderRepository) CreateOrder(orderData *order.Order) (*order.Order, error) {
	//first create the order in DB
	err := OrderRepositoryImpl.dbconn.Create(&orderData).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return orderData, nil
}

func (OrderRepositoryImpl *orderRepository) GetOrders() (*[]order.Order, error) {
	var orderList []order.Order

	err := OrderRepositoryImpl.dbconn.Find(&orderList).Error
	if err != nil {
		return nil, err
	}
	return &orderList, nil
}

func (OrderRepositoryImpl *orderRepository) FindOrder(orderId int) (*order.Order, error) {
	var order order.Order
	err := OrderRepositoryImpl.dbconn.Where("order_id = ?", orderId).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

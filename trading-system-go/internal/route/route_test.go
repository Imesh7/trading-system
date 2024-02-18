package routes_test

import (
	"bytes"
	//"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"trading-system-go/database"
	routes "trading-system-go/internal/route"

	"github.com/go-playground/assert/v2"
	"github.com/gofiber/fiber/v2"
)

func TestCreateOrderRoute(t *testing.T) {
	app := fiber.New()
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "1234")
	t.Setenv("DB_NAME", "tradingsystem")
	t.Setenv("DB_HOST", "127.0.0.1")

	// Define the route and use CreateOrder as the handler
	app.Post("/create-order", routes.CreateOrderRoute)

	/* newOrder:=order.Order{
		UserId: 1000,
		Type: "ask",
		OrderType: order.LimitOrderSell,
		Price: 10,
		Volume: 10000,
		BuyingPair: "btc",
		SellingPair: "usdt",
	} */

	orderPayload := []byte(`{
		"user_id": 1000000,
		"order_type": 4,
		"type": "ask",
		"price": 100,
		"volume": 150,
		"buying_pair": "usd",
		"selling_pair": "btc"
	}`)

	//orderData,err:=json.Marshal(newOrder)
	// Create a request with the desired method, path, and payload
	req := httptest.NewRequest(http.MethodPost, "/create-order", bytes.NewBuffer(orderPayload))
	req.Header.Set("Content-Type", "application/json")
	// Create a ResponseRecorder to capture the response
	//rr := httptest.NewRecorder()

	database.ConnectDatabase()
	// Call the Fiber app's handler with the request
	res, err := app.Test(req)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
}




# Trading system using go + next.js

Trading system with go + postgres + kafka + redis + websocket + next.js

https://github.com/Imesh7/trading-system/assets/50042375/9980ef08-aa57-4bcb-8841-d9038234b3b8


Simulate Order process in a trading system.
 - Create a order
 - Maintain orderbook
 - Match the order with orderbook

 Tech stack used
 --------------------
 - Backend server - Go
 - Kafka
 - Redis (orderbook)
 - Postgres (db)

## Running the Application

To run the application, use the following command.  
Make sure you already installed [Docker](https://www.docker.com/).


```bash
docker-compose up
```



### Create Order Api Enpoint
```bash
http://localhost:8000/create-order
```

Request Body
```json
//create a Ask(sell) order
{
    "user_id":1000000,
    "order_type":4,
    "type":"ask",
    "price":100,
    "volume":150,
    "buying_pair":"usd",
    "selling_pair":"btc"
}

//create a Bid(buy) order
{
    "user_id":1000000,
    "order_type":4,
    "type":"bid",
    "price":100,
    "volume":150,
    "buying_pair":"usd",
    "selling_pair":"btc"
}
```

Get orders API
```
http://localhost:8000/get-orders
```


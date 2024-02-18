"use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import { FixedSizeList as List } from "react-window";
import useWindowDimensions from "./hook/use_window_dimetion";
import AutoSizer from "react-virtualized-auto-sizer";

enum OrderType {
  bid = "bid",
  ask = "ask",
}
interface Order {
  order_type: number;
  type: OrderType;
  price: number;
  volume: number;
  buying_pair: string;
  selling_pair: string;
}

export default function Home() {
  const [orderBookBids, setOrderBookBids] = useState<Order[]>([]);
  const [orderBookAsks, setOrderBookAsks] = useState<Order[]>([]);

  useEffect(() => {
    console.log("Printed on the server-side");
    const socket = new WebSocket(
      "ws://localhost:8000/order-book-update?pair=btc"
    );
    socket.onopen = () => {
      console.log("Connected to the WebSocket server");
    };

    // WebSocket onmessage event listener
    socket.onmessage = (event) => {
      console.log("Received message from server:", event.data);
      const order: Order = JSON.parse(event.data);
      console.log("decodeed data :", orderBookBids);
      if (order.type === "bid") {
        if (orderBookBids.length === 0) {
          setOrderBookBids((old) => [...old, order]);
        }
        setOrderBookBids((old) => [...old, order]);
        console.log("decodeed data :", order);
      } else {
        setOrderBookAsks((old) => [...old, order]);
      }
    };

    // WebSocket onclose event listener
    socket.onclose = (event) => {
      console.log(
        "Disconnected from the WebSocket server. Code:",
        event.code,
        "Reason:",
        event.reason
      );
    };

    // WebSocket onerror event listener
    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
      // Handle the error, e.g., display an error message to the user
    };

    return () => {
      socket.close();
    };
  }, []);

  function Component() {
    return (
      <div
        className="grid grid-cols-2 gap-4 w-full max-w-2xl mx-auto py-4 px-4 md:px2"
        style={{ height: "90vh", width: "100%" }}
      >
        <div className="space-y-3">
          <h2 className="text-2xl font-bold tracking-tight">Order Book</h2>
          {/* <div className="grid gap-4 bg-gray-100 dark:bg-gray-800 rounded-lg p-3"> */}
          <OrderBookCardTile />
          <AutoSizer>
            {({ height, width }) => (
              <List
                className="grid gap-4 bg-gray-100 dark:bg-gray-800 rounded-lg p-3"
                itemSize={orderBookAsks.length}
                itemCount={orderBookAsks.length}
                height={height / 2}
                width={width}
              >
                {({ index, style }) => (
                  <AskOrderTile
                    data={orderBookAsks[index]}
                    index={index}
                    style={style}
                  />
                )}
              </List>
            )}
          </AutoSizer>
          <AutoSizer>
            {({ height, width }) => (
              <List
                className="grid gap-4 bg-gray-100 dark:bg-gray-800 rounded-lg p-3"
                itemSize={orderBookBids.length}
                itemCount={orderBookBids.length}
                height={height / 2}
                width={width}
              >
                {({ index, style }) => (
                  <BidOrderTile
                    data={orderBookBids[index]}
                    index={index}
                    style={style}
                  />
                )}
              </List>
            )}
          </AutoSizer>
        </div>
      </div>
      /* </div> */
    );
  }

  const OrderBookCardTile: React.FC = ({}) => {
    return (
      <div className="grid grid-cols-3 gap-4 text-sm font-medium text-gray-500 dark:text-gray-400">
        <span>Price</span>
        <span>Quantity</span>
        <span>Total</span>
      </div>
    );
  };

  const AskOrderTile: React.FC<{
    data: Order;
    index: number;
    style: React.CSSProperties;
  }> = ({ data, index, style }) => {
    return (
      data && (
        <div
          key={index}
          className="grid grid-cols-3 gap-4 text-sm font-medium text-green-500 dark:text-green-400"
        >
          <span>{data.price}</span>
          <span>{data.volume}</span>
          <span>{data.price * data.volume}</span>
        </div>
      )
    );
  };

  const BidOrderTile: React.FC<{
    data: Order;
    index: number;
    style: React.CSSProperties;
  }> = ({ data, index, style }) => {
    return (
      data && (
        <div
          key={index}
          className="grid grid-cols-3 gap-4 text-sm font-medium text-red-500 dark:text-red-400"
        >
          <span>{data.price}</span>
          <span>{data.volume}</span>
          <span>{data.price * data.volume}</span>
        </div>
      )
    );
  };

  const itemCount = 1000;
  const itemSize = 50; // Height of each item in pixels

  // Render function for list items
  const Row = ({
    index,
    style,
  }: {
    index: number;
    style: React.CSSProperties;
  }) => {
    return (
      <div
        style={{
          ...style,
          backgroundColor: index % 2 === 0 ? "#f0f0f0" : "#ffffff",
        }}
      >
        Item {index}
      </div>
    );
  };

  const MyComponent = () => {
    return (
      <div style={{ height: "400px", width: "100%" }}>
        <List
          height={300} // Total height of the list
          itemCount={itemCount} // Total number of items in the list
          itemSize={itemSize} // Height of each item
          width={"100%"} // Width of the list
        >
          {Row}
        </List>
      </div>
    );
  };

  return (
    <main className="flex">
      <Component />
    </main>
  );
}

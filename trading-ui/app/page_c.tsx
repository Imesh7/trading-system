/* "use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import { FixedSizeList as List } from "react-window";
import useWindowDimensions from "./hook/use_window_dimetion";
import AutoSizer from "react-virtualized-auto-sizer";

interface Order {
  order_type: number;
  type: string;
  price: number;
  volume: number;
  buying_pair: string;
  selling_pair: string;
}

const orderList: Order[] = [
  {
    order_type: 4,
    type: "bid",
    price: 90,
    volume: 150,
    buying_pair: "usd",
    selling_pair: "btc",
  },
];
export default function Home() {
  const [orderBookBids, setOrderBookBids] = useState<Order[]>(orderList);

  useEffect(() => {
    console.log("Printed on the server-side");
  }, []);

  function Component() {
    return (
      <div style={{ height: "90vh", width: "100%" }}>
        
        <OrderBookCardTile />
        <AutoSizer>
          {({ height, width }) => (
            <List
              className="Bid price"
              itemSize={orderBookBids.length}
              itemCount={orderBookBids.length}
              height={height}
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

  const BidOrderTile: React.FC<{
    data: Order;
    index: number;
    style: React.CSSProperties;
  }> = ({ data, index, style }) => {
    if (!data) {
      return null; // or display a placeholder
    }
    return (
      <div
        key={index}
        className="grid grid-cols-3 gap-4 text-sm font-medium text-red-500 dark:text-red-400"
      >
        <span>{data.price}</span>
        <span>{data.volume}</span>
        <span>{data.price * data.volume}</span>
      </div>
    );
  };

  const itemSize = 10; // Height of each item in pixels

  return (
    <main>
      <Component />
    </main>
  );
}
 */

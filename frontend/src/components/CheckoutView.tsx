import { useMutation, useQuery } from "@connectrpc/connect-query"
import { useState } from "react"
import { getCart } from "../gen/v1/cart-CartService_connectquery.ts"
import { createOrder } from "../gen/v1/order-OrderService_connectquery.ts"

function formatPrice(cents: bigint): string {
  return (Number(cents) / 100).toLocaleString("en-US", { style: "currency", currency: "USD" })
}

export function CheckoutView({ onPaid }: { onPaid: (orderId: string, totalCents: bigint) => void }) {
  const [error, setError] = useState("")
  const cartQ = useQuery(getCart)
  const orderMut = useMutation(createOrder, {
    onSuccess: (data) => {
      const order = data.order
      if (!order) return
      onPaid(order.id, order.totalCents)
    },
    onError: (err) => setError(err.message),
  })

  const cart = cartQ.data?.cart

  if (cartQ.isPending) return <p>Loading cart…</p>
  if (cartQ.isError) return <p className="error">Failed to load cart: {cartQ.error.message}</p>
  if (!cart || cart.items.length === 0) return <p>Your cart is empty — nothing to check out.</p>

  const handleOrder = () => {
    setError("")
    orderMut.mutate({
      items: cart.items.map((item) => ({ productId: item.productId, quantity: item.quantity })),
    })
  }

  return (
    <div className="checkout">
      <h2>Checkout</h2>
      <table className="table">
        <thead>
          <tr>
            <th>Product</th>
            <th>Qty</th>
            <th>Price</th>
          </tr>
        </thead>
        <tbody>
          {cart.items.map((item) => (
            <tr key={item.id}>
              <td>{item.productId}</td>
              <td>{item.quantity}</td>
              <td>{formatPrice(item.priceCents)}</td>
            </tr>
          ))}
        </tbody>
        <tfoot>
          <tr>
            <td colSpan={2}>Total</td>
            <td>{formatPrice(cart.totalCents)}</td>
          </tr>
        </tfoot>
      </table>
      {error && <p className="error">{error}</p>}
      <button className="btn" disabled={orderMut.isPending} onClick={handleOrder}>
        {orderMut.isPending ? "Placing order…" : "Place order"}
      </button>
    </div>
  )
}

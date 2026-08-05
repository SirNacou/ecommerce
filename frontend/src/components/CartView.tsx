import { useMutation, useQuery } from "@connectrpc/connect-query"
import {
  getCart,
  removeItem,
  updateItemQuantity,
} from "../gen/v1/cart-CartService_connectquery.ts"

function formatPrice(cents: bigint): string {
  return (Number(cents) / 100).toLocaleString("en-US", { style: "currency", currency: "USD" })
}

export function CartView({ onCheckout }: { onCheckout: () => void }) {
  const cartQ = useQuery(getCart)
  const updateMut = useMutation(updateItemQuantity, {
    onSuccess: () => cartQ.refetch(),
  })
  const removeMut = useMutation(removeItem, {
    onSuccess: () => cartQ.refetch(),
  })

  const cart = cartQ.data?.cart

  if (cartQ.isPending) return <p>Loading cart…</p>
  if (cartQ.isError) return <p className="error">Failed to load cart: {cartQ.error.message}</p>
  if (!cart || cart.items.length === 0) return <p>Your cart is empty.</p>

  return (
    <div className="cart">
      <h2>Cart</h2>
      <table className="table">
        <thead>
          <tr>
            <th>Product</th>
            <th>Qty</th>
            <th>Price</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {cart.items.map((item) => (
            <tr key={item.id}>
              <td>{item.productId}</td>
              <td>
                <button
                  className="qty"
                  onClick={() => updateMut.mutate({ productId: item.productId, quantity: item.quantity - 1 })}
                >
                  −
                </button>
                <span className="qty-value">{item.quantity}</span>
                <button
                  className="qty"
                  onClick={() => updateMut.mutate({ productId: item.productId, quantity: item.quantity + 1 })}
                >
                  +
                </button>
              </td>
              <td>{formatPrice(item.priceCents)}</td>
              <td>
                <button className="btn ghost" onClick={() => removeMut.mutate({ productId: item.productId })}>
                  Remove
                </button>
              </td>
            </tr>
          ))}
        </tbody>
        <tfoot>
          <tr>
            <td colSpan={2}>Total</td>
            <td>{formatPrice(cart.totalCents)}</td>
            <td></td>
          </tr>
        </tfoot>
      </table>
      <button className="btn" onClick={onCheckout}>Checkout</button>
    </div>
  )
}

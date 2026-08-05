import { useMutation, useQuery } from "@connectrpc/connect-query"
import { cancelOrder, listOrders } from "../gen/v1/order-OrderService_connectquery.ts"

function formatPrice(cents: bigint): string {
  return (Number(cents) / 100).toLocaleString("en-US", { style: "currency", currency: "USD" })
}

function formatDate(seconds: bigint | undefined): string {
  if (seconds === undefined) return ""
  return new Date(Number(seconds) * 1000).toLocaleString()
}

function formatStatus(status: string): string {
  switch (status) {
    case "PAID":
      return "Completed"
    case "PENDING":
      return "Pending"
    case "CANCELLED":
      return "Cancelled"
    default:
      return status
  }
}

export function OrdersView() {
  const ordersQ = useQuery(listOrders, { pageSize: 100, pageToken: "" }, { refetchInterval: 5000 })
  const cancelMut = useMutation(cancelOrder, {
    onSuccess: () => ordersQ.refetch(),
  })

  if (ordersQ.isPending) return <p>Loading orders…</p>
  if (ordersQ.isError) return <p className="error">Failed to load orders: {ordersQ.error.message}</p>

  const orders = ordersQ.data?.orders ?? []
  if (orders.length === 0) return <p>No orders yet.</p>

  return (
    <div className="orders">
      <h2>Orders</h2>
      <table className="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Status</th>
            <th>Total</th>
            <th>Created</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td>{order.id}</td>
              <td>{formatStatus(order.status)}</td>
              <td>{formatPrice(order.totalCents)}</td>
              <td>{order.createdAt ? formatDate(order.createdAt.seconds) : ""}</td>
              <td>
                {order.status !== "CANCELLED" && order.status !== "PAID" && (
                  <button className="btn ghost" onClick={() => cancelMut.mutate({ id: order.id, reason: "cancelled by user" })}>
                    Cancel
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

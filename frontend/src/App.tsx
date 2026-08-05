import { useEffect, useState } from "react"
import { queryClient } from "./api/query.ts"
import { getToken, logout, onAuthChanged } from "./api/client.ts"
import { AuthView } from "./components/AuthView.tsx"
import { ShopView } from "./components/ShopView.tsx"
import { CartView } from "./components/CartView.tsx"
import { CheckoutView } from "./components/CheckoutView.tsx"
import { PaymentView } from "./components/PaymentView.tsx"
import { OrdersView } from "./components/OrdersView.tsx"
import { NotificationsView } from "./components/NotificationsView.tsx"
import "./App.css"

type View = "shop" | "cart" | "checkout" | "payment" | "orders" | "notifications"

export function App() {
  const [loggedIn, setLoggedIn] = useState(() => getToken() !== null)
  const [view, setView] = useState<View>("shop")
  const [payOrder, setPayOrder] = useState<{ id: string; totalCents: bigint } | null>(null)

  useEffect(() => {
    return onAuthChanged(() => {
      queryClient.clear()
      setLoggedIn(false)
      setView("shop")
    })
  }, [])

  const handleAuthed = () => {
    setLoggedIn(true)
    setView("shop")
  }

  const handleLogout = () => {
    logout()
    queryClient.clear()
    setLoggedIn(false)
  }

  if (!loggedIn) {
    return <AuthView onAuthed={handleAuthed} />
  }

  return (
    <div className="app">
      <nav className="nav">
        <span className="nav-brand">Store</span>
        <div className="nav-links">
          <button className="nav-link" onClick={() => setView("shop")}>Shop</button>
          <button className="nav-link" onClick={() => setView("cart")}>Cart</button>
          <button className="nav-link" onClick={() => setView("orders")}>Orders</button>
          <button className="nav-link" onClick={() => setView("notifications")}>Notifications</button>
        </div>
        <button className="nav-logout" onClick={handleLogout}>Log out</button>
      </nav>
      <main className="content">
        {view === "shop" && <ShopView onCheckout={() => setView("cart")} />}
        {view === "cart" && <CartView onCheckout={() => setView("checkout")} />}
        {view === "checkout" && (
          <CheckoutView
            onPaid={(id, totalCents) => {
              setPayOrder({ id, totalCents })
              setView("payment")
            }}
          />
        )}
        {view === "payment" && payOrder && (
          <PaymentView orderId={payOrder.id} totalCents={payOrder.totalCents} onDone={() => setView("orders")} />
        )}
        {view === "orders" && <OrdersView />}
        {view === "notifications" && <NotificationsView />}
      </main>
    </div>
  )
}

export default App

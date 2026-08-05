import { useMutation } from "@connectrpc/connect-query"
import { useState } from "react"
import { clearCart } from "../gen/v1/cart-CartService_connectquery.ts"
import { processPayment } from "../gen/v1/payment-PaymentService_connectquery.ts"

function formatPrice(cents: bigint): string {
  return (Number(cents) / 100).toLocaleString("en-US", { style: "currency", currency: "USD" })
}

const METHODS = ["card", "paypal", "bank_transfer", "wallet"]

export function PaymentView({
  orderId,
  totalCents,
  onDone,
}: {
  orderId: string
  totalCents: bigint
  onDone: () => void
}) {
  const [method, setMethod] = useState("card")
  const [cardName, setCardName] = useState("")
  const [cardNumber, setCardNumber] = useState("")
  const [cardExpiry, setCardExpiry] = useState("")
  const [cardCvc, setCardCvc] = useState("")
  const [error, setError] = useState("")
  const [paid, setPaid] = useState<{ transactionId: string; status: string } | null>(null)

  const payMut = useMutation(processPayment, {
    onSuccess: (data) => {
      setPaid({
        transactionId: data.payment?.transactionId ?? "",
        status: data.payment?.status ?? "",
      })
    },
    onError: (err) => setError(err.message),
  })
  const clearMut = useMutation(clearCart)

  const handlePay = () => {
    setError("")
    payMut.mutate({
      orderId,
      amountCents: totalCents,
      currency: "USD",
      paymentMethod: method,
    })
  }

  const handleFinish = async () => {
    await clearMut.mutateAsync({})
    onDone()
  }

  return (
    <div className="payment">
      <h2>Payment</h2>

      {!paid ? (
        <>
          <p className="summary">
            Order <code>{orderId}</code>
            <br />
            Total: <strong>{formatPrice(totalCents)}</strong>
          </p>

          <div className="methods">
            {METHODS.map((m) => (
              <label key={m} className="method">
                <input type="radio" name="method" value={m} checked={method === m} onChange={() => setMethod(m)} />
                {m.replace("_", " ")}
              </label>
            ))}
          </div>

          <div className="form">
            {method === "card" && (
              <>
                <input type="text" placeholder="Cardholder name" value={cardName} onChange={(e) => setCardName(e.target.value)} />
                <input type="text" placeholder="Card number" value={cardNumber} onChange={(e) => setCardNumber(e.target.value)} />
                <div className="card-row">
                  <input type="text" placeholder="MM/YY" value={cardExpiry} onChange={(e) => setCardExpiry(e.target.value)} />
                  <input type="text" placeholder="CVC" value={cardCvc} onChange={(e) => setCardCvc(e.target.value)} />
                </div>
              </>
            )}
            {method !== "card" && <p className="muted">You will be redirected to {method.replace("_", " ")} to complete the payment. (Simulated)</p>}
          </div>

          {error && <p className="error">{error}</p>}
          <button className="btn" disabled={payMut.isPending} onClick={handlePay}>
            {payMut.isPending ? "Processing…" : `Pay ${formatPrice(totalCents)}`}
          </button>
        </>
      ) : (
        <div className="card success">
          <h3>Payment successful</h3>
          <p><strong>Status:</strong> {paid.status}</p>
          <p><strong>Transaction ID:</strong> {paid.transactionId}</p>
          <button className="btn" onClick={handleFinish}>View orders</button>
        </div>
      )}
    </div>
  )
}

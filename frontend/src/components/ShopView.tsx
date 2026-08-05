import { useMutation, useQuery } from "@connectrpc/connect-query"
import { useState } from "react"
import { addItem } from "../gen/v1/cart-CartService_connectquery.ts"
import { listCategories, listProducts } from "../gen/v1/catalog-CatalogService_connectquery.ts"

function formatPrice(cents: bigint): string {
  return (Number(cents) / 100).toLocaleString("en-US", { style: "currency", currency: "USD" })
}

export function ShopView({ onCheckout }: { onCheckout: () => void }) {
  const [categoryId, setCategoryId] = useState("")
  const [added, setAdded] = useState<string | null>(null)

  const categoriesQ = useQuery(listCategories)
  const productsQ = useQuery(listProducts, { pageSize: 100, pageToken: "", categoryId })

  const addMut = useMutation(addItem, {
    onSuccess: () => {
      setAdded("Item added to cart")
      setTimeout(() => setAdded(null), 2000)
    },
    onError: (err) => {
      setAdded(`Failed: ${err.message}`)
      setTimeout(() => setAdded(null), 4000)
    },
  })

  return (
    <div className="shop">
      <div className="shop-header">
        <h2>Products</h2>
        <button className="btn" onClick={onCheckout}>View cart</button>
      </div>
      {added && <p className="notice">{added}</p>}

      <div className="categories">
        <button className={categoryId === "" ? "chip active" : "chip"} onClick={() => setCategoryId("")}>
          All
        </button>
        {categoriesQ.data?.categories.map((c) => (
          <button
            key={c.id}
            className={categoryId === c.id ? "chip active" : "chip"}
            onClick={() => setCategoryId(c.id)}
          >
            {c.name}
          </button>
        ))}
      </div>

      {productsQ.isPending && <p>Loading products…</p>}
      {productsQ.isError && <p className="error">Failed to load products: {productsQ.error.message}</p>}

      <div className="grid">
        {productsQ.data?.products.map((p) => (
          <div key={p.id} className="product">
            <h3>{p.name}</h3>
            <p className="price">{formatPrice(p.priceCents)}</p>
            <p className="desc">{p.description}</p>
            <button
              className="btn"
              disabled={addMut.isPending}
              onClick={() => addMut.mutate({ productId: p.id, quantity: 1 })}
            >
              Add to cart
            </button>
          </div>
        ))}
      </div>
      {productsQ.data && productsQ.data.products.length === 0 && <p>No products in this category.</p>}
    </div>
  )
}

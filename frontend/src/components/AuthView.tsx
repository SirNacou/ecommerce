import { useMutation } from "@connectrpc/connect-query"
import { useState } from "react"
import { login, register } from "../gen/v1/user-UserService_connectquery.ts"
import { setToken } from "../api/client.ts"

export function AuthView({ onAuthed }: { onAuthed: () => void }) {
  const [mode, setMode] = useState<"login" | "register">("login")
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")

  const loginMut = useMutation(login, {
    onSuccess: (data) => {
      setToken(data.accessToken, data.refreshToken)
      onAuthed()
    },
    onError: (err) => setError(err.message),
  })

  const registerMut = useMutation(register, {
    onSuccess: async (data) => {
      await loginMut.mutateAsync({ email: data.email, password })
    },
    onError: (err) => setError(err.message),
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    if (mode === "login") {
      loginMut.mutate({ email, password })
    } else {
      registerMut.mutate({ email, password, name })
    }
  }

  return (
    <div className="auth">
      <h1>Store</h1>
      <div className="card">
        <div className="tabs">
          <button className={mode === "login" ? "tab active" : "tab"} onClick={() => setMode("login")}>Log in</button>
          <button className={mode === "register" ? "tab active" : "tab"} onClick={() => setMode("register")}>Register</button>
        </div>
        <form onSubmit={handleSubmit} className="form">
          {mode === "register" && (
            <input
              type="text"
              placeholder="Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          )}
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {error && <p className="error">{error}</p>}
          <button type="submit" disabled={loginMut.isPending || registerMut.isPending}>
            {mode === "login" ? "Log in" : "Register"}
          </button>
        </form>
      </div>
    </div>
  )
}

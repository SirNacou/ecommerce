import { createConnectTransport } from "@connectrpc/connect-web"
import { Code, ConnectError, type Interceptor } from "@connectrpc/connect"
import { isTokenExpired } from "./session.ts"

const TOKEN_KEY = "access_token"
const REFRESH_KEY = "refresh_token"

type AuthListener = () => void
const authListeners = new Set<AuthListener>()

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_KEY)
}

export function setToken(accessToken: string | null, refreshToken?: string | null) {
  if (accessToken) {
    localStorage.setItem(TOKEN_KEY, accessToken)
    if (refreshToken !== undefined) {
      if (refreshToken) localStorage.setItem(REFRESH_KEY, refreshToken)
      else localStorage.removeItem(REFRESH_KEY)
    }
  } else {
    localStorage.removeItem(TOKEN_KEY)
    if (refreshToken !== undefined && !refreshToken) localStorage.removeItem(REFRESH_KEY)
  }
}

function notifyAuthChanged() {
  authListeners.forEach((listener) => listener())
}

export function onAuthChanged(listener: AuthListener): () => void {
  authListeners.add(listener)
  return () => authListeners.delete(listener)
}

export function logout() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(REFRESH_KEY)
  notifyAuthChanged()
}

const authInterceptor: Interceptor = (next) => async (req) => {
  const token = getToken()
  if (token && isTokenExpired(token)) {
    logout()
    throw new ConnectError("Session expired, please log in again", Code.Unauthenticated)
  }
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`)
  }
  try {
    return await next(req)
  } catch (err) {
    if (err instanceof ConnectError && err.code === Code.Unauthenticated && getToken()) {
      logout()
    }
    throw err
  }
}

export const transport = createConnectTransport({
  baseUrl: "http://localhost:9080",
  interceptors: [authInterceptor],
})

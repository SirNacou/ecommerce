export interface TokenClaims {
  user_id?: string
  email?: string
  token_type?: string
  exp?: number
}

function decodeBase64Url(input: string): string {
  const base64 = input.replace(/-/g, "+").replace(/_/g, "/")
  const padded = base64.padEnd(Math.ceil(base64.length / 4) * 4, "=")
  return decodeURIComponent(
    atob(padded)
      .split("")
      .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
      .join(""),
  )
}

export function decodeToken(token: string): TokenClaims | null {
  try {
    const parts = token.split(".")
    if (parts.length !== 3) return null
    return JSON.parse(decodeBase64Url(parts[1])) as TokenClaims
  } catch {
    return null
  }
}

export function isTokenExpired(token: string, now: number = Date.now()): boolean {
  const claims = decodeToken(token)
  if (!claims || typeof claims.exp !== "number") return false
  return claims.exp * 1000 <= now
}

export function tokenExpiresInMs(token: string): number | null {
  const claims = decodeToken(token)
  if (!claims || typeof claims.exp !== "number") return null
  return claims.exp * 1000 - Date.now()
}

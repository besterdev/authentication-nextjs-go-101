const ACCESS_TOKEN_KEY = "access_token"
const REFRESH_TOKEN_KEY = "refresh_token"
const TOKEN_EXPIRES_AT_KEY = "token_expires_at"

export type StoredTokens = {
  accessToken: string
  refreshToken: string
  expiresAt: number
}

const isBrowser = () => typeof window !== "undefined"

let tokenCache: StoredTokens | null | undefined

const readTokensFromStorage = (): StoredTokens | null => {
  if (!isBrowser()) return null

  const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY)
  const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
  const expiresAtRaw = localStorage.getItem(TOKEN_EXPIRES_AT_KEY)

  if (!accessToken || !refreshToken || !expiresAtRaw) {
    return null
  }

  const expiresAt = Number(expiresAtRaw)
  if (Number.isNaN(expiresAt)) {
    return null
  }

  return { accessToken, refreshToken, expiresAt }
}

export const getStoredTokens = (): StoredTokens | null => {
  if (!isBrowser()) return null

  if (tokenCache !== undefined) {
    return tokenCache
  }

  tokenCache = readTokensFromStorage()
  return tokenCache
}

export const saveTokens = (tokens: {
  access_token: string
  refresh_token: string
  expires_in: number
}) => {
  if (!isBrowser()) return

  const expiresAt = Date.now() + tokens.expires_in * 1000
  localStorage.setItem(ACCESS_TOKEN_KEY, tokens.access_token)
  localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token)
  localStorage.setItem(TOKEN_EXPIRES_AT_KEY, String(expiresAt))

  tokenCache = {
    accessToken: tokens.access_token,
    refreshToken: tokens.refresh_token,
    expiresAt,
  }
}

export const clearTokens = () => {
  if (!isBrowser()) return

  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  localStorage.removeItem(TOKEN_EXPIRES_AT_KEY)
  tokenCache = null
}

export const getAccessToken = (): string | null =>
  getStoredTokens()?.accessToken ?? null

export const getRefreshToken = (): string | null =>
  getStoredTokens()?.refreshToken ?? null

export const getTokenExpiresAt = (): number | null =>
  getStoredTokens()?.expiresAt ?? null

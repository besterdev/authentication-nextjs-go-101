import {
  clearTokens,
  getRefreshToken,
  getTokenExpiresAt,
  saveTokens,
} from "@/lib/auth-storage";
import { http } from "@/lib/axios-client";
import type { AuthTokens } from "@/lib/types";

const REFRESH_BUFFER_MS = 60_000;

let refreshPromise: Promise<boolean> | null = null;
let refreshTimer: ReturnType<typeof setTimeout> | null = null;

const fetchRefresh = async (refreshToken: string): Promise<AuthTokens> => {
  const { data } = await http.post<AuthTokens>("/auth/refresh", {
    refresh_token: refreshToken,
  });
  return data;
};

export const clearRefreshTimer = () => {
  if (refreshTimer) {
    clearTimeout(refreshTimer);
    refreshTimer = null;
  }
};

export const performTokenRefresh = async (): Promise<boolean> => {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = (async () => {
    const refreshToken = getRefreshToken();
    if (!refreshToken) {
      clearTokens();
      return false;
    }

    try {
      const tokens = await fetchRefresh(refreshToken);
      saveTokens(tokens);
      scheduleProactiveRefresh();
      return true;
    } catch {
      clearTokens();
      return false;
    } finally {
      refreshPromise = null;
    }
  })();

  return refreshPromise;
};

export const scheduleProactiveRefresh = () => {
  clearRefreshTimer();

  const expiresAt = getTokenExpiresAt();
  if (!expiresAt) return;

  const delay = Math.max(expiresAt - Date.now() - REFRESH_BUFFER_MS, 0);

  refreshTimer = setTimeout(() => {
    void performTokenRefresh();
  }, delay);
};

import { create } from "zustand";
import {
  getMe,
  loginAndStore,
  logout as apiLogout,
  register as apiRegister,
} from "@/lib/api";
import { clearTokens, getStoredTokens } from "@/lib/auth-storage";
import {
  clearRefreshTimer,
  performTokenRefresh,
  scheduleProactiveRefresh,
} from "@/lib/token-refresh";
import { ApiError, type User } from "@/lib/types";

type AuthState = {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
  hydrate: () => Promise<void>;
  cancelHydrate: () => void;
};

let hydrateId = 0;

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  isLoading: true,

  refreshUser: async () => {
    const profile = await getMe();
    set({ user: profile });
  },

  hydrate: async () => {
    const id = ++hydrateId;
    const tokens = getStoredTokens();

    if (!tokens) {
      if (id === hydrateId) {
        set({ isLoading: false });
      }
      return;
    }

    scheduleProactiveRefresh();

    try {
      const profile = await getMe();
      if (id === hydrateId) {
        set({ user: profile });
      }
    } catch (error) {
      if (error instanceof ApiError && error.code === "AUTH_INVALID") {
        const refreshed = await performTokenRefresh();
        if (refreshed) {
          try {
            const profile = await getMe();
            if (id === hydrateId) {
              set({ user: profile });
            }
            return;
          } catch {
            // fall through to clear
          }
        }
      }
      clearTokens();
      clearRefreshTimer();
      if (id === hydrateId) {
        set({ user: null });
      }
    } finally {
      if (id === hydrateId) {
        set({ isLoading: false });
      }
    }
  },

  cancelHydrate: () => {
    hydrateId += 1;
    clearRefreshTimer();
  },

  login: async (email, password) => {
    await loginAndStore({ email, password });
    scheduleProactiveRefresh();
    await get().refreshUser();
  },

  register: async (email, password) => {
    await apiRegister({ email, password });
  },

  logout: async () => {
    try {
      await apiLogout();
    } catch {
      // Clear local session even if the API call fails.
    } finally {
      clearTokens();
      clearRefreshTimer();
      set({ user: null });
    }
  },
}));

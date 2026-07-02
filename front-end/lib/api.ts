import { isAxiosError, type Method } from "axios";
import {
  clearTokens,
  getAccessToken,
  saveTokens,
} from "@/lib/auth-storage";
import { http } from "@/lib/axios-client";
import { performTokenRefresh } from "@/lib/token-refresh";
import {
  ApiError,
  type ApiErrorCode,
  type AuthTokens,
  type LoginRequest,
  type RegisterRequest,
  type User,
} from "@/lib/types";

type RequestOptions = {
  method?: Method;
  body?: unknown;
  auth?: boolean;
  retry?: boolean;
};

const parseError = (error: unknown): ApiError => {
  if (error instanceof ApiError) {
    return error;
  }

  if (isAxiosError(error)) {
    const data = error.response?.data as
      | { error?: string; code?: ApiErrorCode }
      | undefined;

    if (data?.error || data?.code) {
      return new ApiError(
        data.error ?? "Request failed",
        data.code ?? "INTERNAL_ERROR",
      );
    }
  }

  return new ApiError("Request failed", "INTERNAL_ERROR");
};

const request = async <T>(
  path: string,
  { method = "GET", body, auth = false, retry = true }: RequestOptions = {},
): Promise<T> => {
  const headers: Record<string, string> = {};

  if (auth) {
    const accessToken = getAccessToken();
    if (!accessToken) {
      throw new ApiError("missing authorization header", "AUTH_MISSING");
    }
    headers.Authorization = `Bearer ${accessToken}`;
  }

  try {
    const response = await http.request<T>({
      url: path,
      method,
      data: body,
      headers,
    });

    if (response.status === 204) {
      return undefined as T;
    }

    return response.data;
  } catch (error) {
    if (
      isAxiosError(error) &&
      error.response?.status === 401 &&
      auth &&
      retry
    ) {
      const apiError = parseError(error);

      if (apiError.code === "AUTH_INVALID" || apiError.code === "AUTH_MISSING") {
        const refreshed = await performTokenRefresh();
        if (refreshed) {
          return request<T>(path, { method, body, auth, retry: false });
        }
        clearTokens();
      }

      throw apiError;
    }

    throw parseError(error);
  }
};

export const register = (data: RegisterRequest): Promise<User> =>
  request<User>("/auth/register", { method: "POST", body: data });

export const login = (data: LoginRequest): Promise<AuthTokens> =>
  request<AuthTokens>("/auth/login", { method: "POST", body: data });

export const refreshTokens = (refreshToken: string): Promise<AuthTokens> =>
  request<AuthTokens>("/auth/refresh", {
    method: "POST",
    body: { refresh_token: refreshToken },
  });

export const getMe = (): Promise<User> =>
  request<User>("/auth/me", { auth: true });

export const logout = (): Promise<void> =>
  request<void>("/auth/logout", { method: "POST", auth: true });

export const loginAndStore = async (data: LoginRequest): Promise<AuthTokens> => {
  const tokens = await login(data);
  saveTokens(tokens);
  return tokens;
};

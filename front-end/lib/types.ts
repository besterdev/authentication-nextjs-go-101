export type User = {
  id: string;
  email: string;
  created_at: string;
};

export type AuthTokens = {
  access_token: string;
  refresh_token: string;
  expires_in: number;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type RegisterRequest = {
  email: string;
  password: string;
};

export type RefreshRequest = {
  refresh_token: string;
};

export type ApiErrorCode =
  | "INTERNAL_ERROR"
  | "AUTH_INVALID"
  | "AUTH_MISSING"
  | "REFRESH_INVALID"
  | "EMAIL_EXISTS";

export class ApiError extends Error {
  code: ApiErrorCode;

  constructor(message: string, code: ApiErrorCode) {
    super(message);
    this.name = "ApiError";
    this.code = code;
  }
}

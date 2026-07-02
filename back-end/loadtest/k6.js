import http from "k6/http"
import { check, sleep } from "k6"
import { Rate, Trend } from "k6/metrics"

// --- Metrics -----------------------------------------------------------------

const metrics = {
  errors: new Rate("errors"),
  healthDuration: new Trend("health_duration"),
  meDuration: new Trend("me_duration"),
  loginDuration: new Trend("login_duration"),
}

// --- Config ------------------------------------------------------------------

const env = {
  baseUrl: __ENV.BASE_URL || "http://localhost:8080",
  email: __ENV.TEST_EMAIL || "loadtest@example.com",
  password: __ENV.TEST_PASSWORD || "password123",
  testDuration: __ENV.TEST_DURATION || "30s",
  authRateLimitPerMin: Number(__ENV.AUTH_RATE_LIMIT_PER_MIN || 10),
  healthRate: Number(__ENV.HEALTH_RATE || 100),
  meRate: Number(__ENV.ME_RATE || 50),
  loginRatePerMin: __ENV.LOGIN_RATE_PER_MIN,
}

const SETUP_AUTH_REQUESTS = 2

const authBudgetPerMin = Math.max(
  1,
  env.authRateLimitPerMin - SETUP_AUTH_REQUESTS,
)

const loginRate = Number(
  env.loginRatePerMin || Math.max(1, authBudgetPerMin),
)

const rates = { me: env.meRate, login: loginRate }

// --- Helpers -----------------------------------------------------------------

const jsonHeaders = { "Content-Type": "application/json" }

const postJson = (url, body) =>
  http.post(url, JSON.stringify(body), { headers: jsonHeaders })

const isOk = (res) => res.status === 200

const track = (trend, res, label) => {
  trend.add(res.timings.duration)
  const passed = check(res, { [label]: isOk })
  metrics.errors.add(!passed)
  return passed
}

const arrivalScenario = (name, { rate, timeUnit, exec, startTime }) => {
  const scenario = {
    executor: "constant-arrival-rate",
    rate,
    timeUnit,
    duration: env.testDuration,
    preAllocatedVUs: name === "health" ? 20 : 2,
    maxVUs: name === "health" ? 100 : 5,
    exec,
  }
  if (startTime) scenario.startTime = startTime
  return scenario
}

// --- k6 options --------------------------------------------------------------

export const options = {
  scenarios: {
    health: arrivalScenario("health", {
      rate: env.healthRate,
      timeUnit: "1s",
      exec: "healthCheck",
    }),
    me: arrivalScenario("me", {
      rate: rates.me,
      timeUnit: "1s",
      exec: "meFlow",
      startTime: "2s",
    }),
    login: arrivalScenario("login", {
      rate: rates.login,
      timeUnit: "1m",
      exec: "loginFlow",
      startTime: "5s",
    }),
  },
  thresholds: {
    errors: ["rate<0.01"],
    health_duration: ["p(99)<50"],
    me_duration: ["p(99)<50"],
    login_duration: ["p(99)<500"],
  },
}

// --- Scenarios ---------------------------------------------------------------
// Rate limit applies only to register/login/refresh (10 req/min per IP).
// /auth/me and /auth/logout are not rate-limited.

export const setup = () => {
  const registerRes = postJson(`${env.baseUrl}/auth/register`, {
    email: env.email,
    password: env.password,
  })

  const registerOk = registerRes.status === 201 || registerRes.status === 409
  if (!registerOk) {
    throw new Error(`setup register failed: ${registerRes.status} ${registerRes.body}`)
  }

  const loginRes = postJson(`${env.baseUrl}/auth/login`, {
    email: env.email,
    password: env.password,
  })

  if (!check(loginRes, { "setup login ok": isOk })) {
    throw new Error(`setup login failed: ${loginRes.status} ${loginRes.body}`)
  }

  return { accessToken: loginRes.json().access_token }
}

export const healthCheck = () => {
  track(metrics.healthDuration, http.get(`${env.baseUrl}/health`), "health 200")
  sleep(0.01)
}

export const meFlow = ({ accessToken }) => {
  if (!accessToken) {
    metrics.errors.add(true)
    return
  }

  const res = http.get(`${env.baseUrl}/auth/me`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  })

  track(metrics.meDuration, res, "me 200")
  sleep(1)
}

export const loginFlow = () => {
  const res = postJson(`${env.baseUrl}/auth/login`, {
    email: env.email,
    password: env.password,
  })

  track(metrics.loginDuration, res, "login 200")
  sleep(1)
}

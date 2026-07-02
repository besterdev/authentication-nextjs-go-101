# Authentication Frontend

Next.js frontend for the Go JWT authentication API. Provides login, registration, token management, and a protected dashboard.

## Prerequisites

- Node.js 20+
- Running backend at `http://localhost:8080` (see `../back-end/README.md`)

## Setup

1. Install dependencies:

```bash
npm install
```

2. Configure environment:

```bash
cp .env.local.example .env.local
```

3. Start the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Features

- **Login / Register** — forms with validation and API error handling
- **JWT storage** — access and refresh tokens in `localStorage`
- **Auto refresh** — proactive refresh before expiry and reactive refresh on 401
- **Protected dashboard** — client-side route guard with profile from `GET /auth/me`
- **Logout** — revokes refresh tokens via API and clears local session

## Routes

| Route | Description |
| ----- | ----------- |
| `/` | Redirects to `/dashboard` or `/login` |
| `/login` | Sign in |
| `/register` | Create account (redirects to login on success) |
| `/dashboard` | Protected profile dashboard |

## Environment Variables

| Variable | Description | Default |
| -------- | ----------- | ------- |
| `NEXT_PUBLIC_API_URL` | Backend API base URL | `http://localhost:8080` |

Ensure the backend `CORS_ORIGIN` matches `http://localhost:3000`.

## Scripts

```bash
npm run dev      # Start dev server
npm run build    # Production build
npm run start    # Start production server
npm run lint     # Run ESLint
```

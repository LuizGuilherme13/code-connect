# Code Connect

Full-stack monorepo: React 19 + TypeScript + Vite 6 frontend, Go 1.22 backend.

## Structure

```
code-connect/
├── backend/        # Go HTTP server (multi-package)
│   ├── main.go     # Entry point, routes, server
│   ├── handlers/   # HTTP handlers (Register, Login, GetProfile)
│   ├── middleware/  # JWT auth middleware
│   ├── models/     # Structs and DTOs
│   ├── store/      # PostgreSQL storage (database/sql + lib/pq)
│   └── docs/       # Generated Swagger docs (swag init)
├── frontend/       # React + Vite SPA
│   └── src/
│       ├── components/
│       │   ├── atoms/        # Input, Button, Checkbox, Link
│       │   ├── molecules/    # FormField, SocialButton
│       │   ├── organisms/    # Banner, LoginForm, RegisterForm
│       │   ├── templates/    # AuthLayout
│       │   └── pages/        # LoginPage, RegisterPage
│       └── __tests__/        # Accessibility tests
├── dev.sh          # Runs both concurrently
└── package.json    # Root convenience scripts
```

## Commands

| Command | What it does |
|---|---|
| `npm run dev` | Runs `dev.sh` — starts backend (`:8080`) + frontend (`:5173`) in parallel |
| `npm run dev:frontend` | Frontend only via Vite |
| `npm run dev:backend` | Backend only via `go run .` |
| `npm run build` | `tsc && vite build` (frontend) |
| `npm run build:backend` | `go build -o ../server .` in `backend/` |
| `npm run start:backend` | `go run .` in `backend/` |
| `npm run test` | Run frontend accessibility tests (vitest) |
| `npm run test:watch` | Run frontend tests in watch mode |
| `npm run test:backend` | Run backend unit tests (`go test ./...`) |
| `npm run test:backend:watch` | Run backend tests (no caching) |

## Dev server quirks

- Vite proxies `/api/*` → `http://localhost:8080` — frontend dev expects the backend running.
- `dev.sh` spawns both processes and kills them on Ctrl+C.
- Backend port overridable via `PORT` env var (default `8080`).
- Backend uses `JWT_SECRET` env var for JWT signing (default: `default-dev-secret`).
- Both `tsconfig.json` and `tsconfig.node.json` exist; the former covers `src/`, the latter covers `vite.config.ts`.

## PostgreSQL

- A `docker-compose.yml` at the project root runs PostgreSQL 16 on port `5432`.
- Start the database before running the backend: `docker compose up -d`
- Connection configured via env vars (defaults match docker-compose):

| Env var | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `codeconnect` | Database user |
| `DB_PASSWORD` | `codeconnect` | Database password |
| `DB_NAME` | `codeconnect` | Database name |

- Backend uses `database/sql` with `github.com/lib/pq` driver, no ORM.
- Schema auto-created on startup via `CREATE TABLE IF NOT EXISTS`.

## Code conventions

- Backend: multi-package Go project, standard `net/http` with `ServeMux` route patterns (`"GET /api/..."`), no framework.
- ServeMux method patterns required (`"GET /api/..."`, not `"/api/..."`).
- Backend follows **REST** conventions: resource-based URL paths (`/api/users`, `/api/posts/:id`), proper HTTP methods (GET/POST/PUT/DELETE), meaningful status codes.
- Backend Swagger documentation via **swaggo** annotations — run `swag init` in `backend/` after adding/changing annotations.
- Frontend: React 19 with `react-jsx` transform, strict TypeScript (`noUnusedLocals`, `noUnusedParameters` on), plain CSS (no CSS-in-JS or Tailwind).
- Frontend components follow **Atomic Design**: `atoms/`, `molecules/`, `organisms/`, `templates/`, `pages/` under `frontend/src/components/`. New components go in the correct layer — atoms have no dependencies on molecules or above.
- Frontend components must have a **component test** co-located or under `__tests__/`, using a framework like Vitest + Testing Library (add config if not yet present).
- **All commits follow Conventional Commits**: `feat:`, `fix:`, `chore:`, `refactor:`, `docs:`, `test:`, etc. Scope optional: `feat(api): ...`.

## Routing

- React Router v7 (`react-router-dom`) for client-side routing.
- Routes defined in `App.tsx`, wrapped in `BrowserRouter` in `main.tsx`.
- Current routes: `/` (login), `/cadastro` (register).
- Use `<Link to="...">` from atoms/Link for navigation (wraps react-router-dom Link).

## API Endpoints

| Method | Route | Auth | Description |
|---|---|---|---|
| `POST` | `/api/register` | No | Register user (name, email, password) |
| `POST` | `/api/login` | No | Login, returns JWT |
| `GET` | `/api/users/me` | Bearer JWT | Get authenticated user data |
| `GET` | `/swagger/*` | No | Swagger UI docs |

## Testing

### Frontend (Vitest)

- **Framework**: Vitest + Testing Library + axe-core
- **Config**: `vite.config.ts` (test section), `src/test-setup.ts`
- **Location**: `src/__tests__/` for accessibility tests
- **Commands**: `npm run test`, `npm run test:watch`

#### Accessibility Testing (WCAG 2 AA)

- Use `axe-core` to detect violations automatically.
- Test files: `login-accessibility.test.tsx`, `register-accessibility.test.tsx`
- Tests cover: heading hierarchy, form labels, button/link accessibility, color contrast.
- Run `npm run test` before commits to catch regressions.

### Backend (Go testing)

- **Framework**: Go standard `testing` package + `net/http/httptest`
- **Location**: `*_test.go` files co-located in each package (`store/`, `handlers/`, `middleware/`)
- **Commands**: `npm run test:backend`, `cd backend && go test ./... -v`
- **Patterns**: Table-driven tests, subtests with `t.Run`, `httptest.NewRecorder` for handler tests
- **Requires PostgreSQL**: tests connect to the database using the same env vars as the server. Run `docker compose up -d` first.

## Design tokens

```css
:root {
  --bg-primary: #0F1923;
  --bg-card: #1A2633;
  --bg-input: #2A3A4A;
  --accent: #6CFFA0;
  --accent-hover: #5CE690;
  --text-primary: #FFFFFF;
  --text-secondary: #8899AA;
  --border-radius: 12px;
  --border-radius-sm: 8px;
}
```

## Assets

Public images in `frontend/public/`:
- `banner-login.png`, `banner-cadastro.png` — auth page banners
- `github.png`, `gmail.png` — social login icons

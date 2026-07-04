# Code Connect

Full-stack monorepo: React 19 + TypeScript + Vite 6 frontend, Go 1.22 backend.

## Structure

```
code-connect/
├── backend/        # Go HTTP server (single main.go)
├── frontend/       # React + Vite SPA
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

No linter or typecheck scripts exist yet. Plan to add them if working on tooling.

## Dev server quirks

- Vite proxies `/api/*` → `http://localhost:8080` — frontend dev expects the backend running.
- `dev.sh` spawns both processes and kills them on Ctrl+C.
- Backend port overridable via `PORT` env var (default `8080`).
- Both `tsconfig.json` and `tsconfig.node.json` exist; the former covers `src/`, the latter covers `vite.config.ts`.

## Code conventions

- Backend: single `main.go`, standard `net/http` with `ServeMux` route patterns (`"GET /api/..."`), no framework.
- ServeMux method patterns required (`"GET /api/..."`, not `"/api/..."`).
- Backend follows **REST** conventions: resource-based URL paths (`/api/users`, `/api/posts/:id`), proper HTTP methods (GET/POST/PUT/DELETE), meaningful status codes.
- Frontend: React 19 with `react-jsx` transform, strict TypeScript (`noUnusedLocals`, `noUnusedParameters` on), plain CSS (no CSS-in-JS or Tailwind).
- Frontend components follow **Atomic Design**: `atoms/`, `molecules/`, `organisms/`, `templates/`, `pages/` under `frontend/src/components/`. New components go in the correct layer — atoms have no dependencies on molecules or above.
- Frontend components must have a **component test** co-located or under `__tests__/`, using a framework like Vitest + Testing Library (add config if not yet present).
- **All commits follow Conventional Commits**: `feat:`, `fix:`, `chore:`, `refactor:`, `docs:`, `test:`, etc. Scope optional: `feat(api): ...`.

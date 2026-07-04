#!/bin/bash
set -e

cleanup() {
  echo ""
  echo "Stopping all processes..."
  kill $backend_pid $frontend_pid 2>/dev/null
  wait $backend_pid $frontend_pid 2>/dev/null
  echo "Done."
  exit 0
}
trap cleanup SIGINT SIGTERM

echo "Starting Code Connect..."
echo "  Frontend: http://localhost:5173"
echo "  Backend:  http://localhost:8080"
echo ""

(cd backend && go run .) &
backend_pid=$!

(cd frontend && npx vite) &
frontend_pid=$!

wait

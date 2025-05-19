@echo off
echo Starting Go backend...
start cmd /k "cd backend && go run main.go"

echo Starting Vite frontend...
start cmd /k "cd frontend && npm run dev"

echo Both servers started.

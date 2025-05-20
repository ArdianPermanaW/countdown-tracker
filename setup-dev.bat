@echo off
echo Installing frontend dependencies...
cd frontend
call npm install

echo Installing backend dependencies...
cd ../backend
call go mod tidy

echo Starting Go backend...
start cmd /k "cd backend && go run main.go"

echo Starting Vite frontend...
start cmd /k "cd frontend && npm run dev"

echo All done!

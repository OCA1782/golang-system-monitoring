@echo off
echo ============================================
echo  SysMonitor - System Resource Monitor
echo ============================================

:: Start backend
echo [1/2] Starting Go backend on :8080...
cd /d "%~dp0backend"
start "SysMonitor Backend" cmd /k "go run . 2>&1"
timeout /t 3 /nobreak >nul

:: Start frontend
echo [2/2] Starting Vue frontend on :3000...
cd /d "%~dp0frontend"
start "SysMonitor Frontend" cmd /k "npm run dev 2>&1"
timeout /t 4 /nobreak >nul

:: Open browser
echo Opening browser...
start http://localhost:3000

echo.
echo Backend:  http://localhost:8080
echo Frontend: http://localhost:3000
echo.
echo Both windows opened. Close them to stop the servers.

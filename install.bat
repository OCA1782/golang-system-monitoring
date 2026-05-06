@echo off
echo ============================================
echo  SysMonitor - First Time Setup
echo ============================================

echo [1/3] Installing Go dependencies...
cd /d "%~dp0backend"
go mod tidy
if errorlevel 1 (
    echo ERROR: go mod tidy failed. Make sure Go is installed.
    pause
    exit /b 1
)

echo [2/3] Installing Node.js dependencies...
cd /d "%~dp0frontend"
npm install
if errorlevel 1 (
    echo ERROR: npm install failed. Make sure Node.js is installed.
    pause
    exit /b 1
)

echo [3/3] Done!
echo.
echo Run start.bat to launch the application.
pause

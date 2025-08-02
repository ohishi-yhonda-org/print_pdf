@echo off
echo Managing PDF Generator API Service...

echo.
echo Service Status:
sc query "PDFGeneratorAPI"

echo.
echo Available commands:
echo 1. Start service
echo 2. Stop service
echo 3. Restart service
echo 4. Check service status
echo 5. View recent event logs
echo 6. Exit

echo.
set /p choice="Enter your choice (1-6): "

if "%choice%"=="1" (
    echo Starting service...
    sc start "PDFGeneratorAPI"
) else if "%choice%"=="2" (
    echo Stopping service...
    sc stop "PDFGeneratorAPI"
) else if "%choice%"=="3" (
    echo Restarting service...
    sc stop "PDFGeneratorAPI"
    timeout /t 5 >nul
    sc start "PDFGeneratorAPI"
) else if "%choice%"=="4" (
    echo Service status:
    sc query "PDFGeneratorAPI"
) else if "%choice%"=="5" (
    echo Recent event logs for PDF Generator API Service:
    wevtutil qe Application /c:10 /rd:true /f:text /q:"*[System[Provider[@Name='PDF Generator API Service']]]"
) else if "%choice%"=="6" (
    exit /b 0
) else (
    echo Invalid choice!
)

echo.
pause
goto :eof

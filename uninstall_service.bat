@echo off
echo Uninstalling PDF Generator API Service...

REM 管理者権限で実行されているかチェック
net session >nul 2>&1
if %errorLevel% == 0 (
    echo Running with administrator privileges
) else (
    echo This script must be run as administrator!
    echo Right-click and select "Run as administrator"
    pause
    exit /b 1
)

REM サービスの状態をチェック
sc query "PDFGeneratorAPI" >nul 2>&1
if %errorLevel% == 0 (
    echo Service found. Proceeding with uninstallation...
    
    REM サービスを停止
    echo Stopping service...
    sc stop "PDFGeneratorAPI"
    
    REM 少し待機
    timeout /t 5 >nul
    
    REM サービスを削除
    echo Removing service...
    sc delete "PDFGeneratorAPI"
    
    if %errorLevel% == 0 (
        echo Service uninstalled successfully!
    ) else (
        echo Failed to uninstall service!
    )
) else (
    echo Service not found or already uninstalled.
)

echo.
pause

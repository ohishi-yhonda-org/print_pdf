@echo off
echo Installing PDF Generator API Service...

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

REM サービスを停止（存在する場合）
sc query "PDFGeneratorAPI" >nul 2>&1
if %errorLevel% == 0 (
    echo Stopping existing service...
    sc stop "PDFGeneratorAPI"
    timeout /t 5 >nul
)

REM 既存のサービスを削除（存在する場合）
sc query "PDFGeneratorAPI" >nul 2>&1
if %errorLevel% == 0 (
    echo Removing existing service...
    sc delete "PDFGeneratorAPI"
    timeout /t 2 >nul
)

REM 実行ファイルのパスを取得
set SERVICE_PATH=%~dp0print_pdf.exe

REM サービスを作成
echo Creating service...
sc create "PDFGeneratorAPI" binPath= "%SERVICE_PATH%" DisplayName= "PDF Generator API Service" start= auto

if %errorLevel% == 0 (
    echo Service created successfully!
    
    REM サービスの説明を設定
    sc description "PDFGeneratorAPI" "HTTP API service for generating PDF documents"
    
    REM サービスを開始
    echo Starting service...
    sc start "PDFGeneratorAPI"
    
    if %errorLevel% == 0 (
        echo Service started successfully!
        echo Service is now running on http://localhost:8081
        echo.
        echo To check service status: sc query "PDFGeneratorAPI"
        echo To stop service: sc stop "PDFGeneratorAPI"
        echo To uninstall service: run uninstall_service.bat
    ) else (
        echo Failed to start service!
        echo Check Windows Event Log for details.
    )
) else (
    echo Failed to create service!
)

echo.
pause

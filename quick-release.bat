@echo off
:: Quick Release Batch Script
:: Usage: quick-release.bat [version]
:: Example: quick-release.bat v1.0.15

echo.
echo 🚀 Quick Release Script Starting...
echo.

:: Check if version parameter provided
set "VERSION=%1"

:: Step 1: Commit current changes
echo 📝 Committing current changes...
git add .
if "%VERSION%"=="" (
    set /p "COMMIT_MSG=Enter commit message (or press Enter for default): "
    if "!COMMIT_MSG!"=="" (
        for /f "tokens=1-3 delims=/: " %%a in ('echo %date% %time%') do set "TIMESTAMP=%%a-%%b-%%c_%%d"
        set "COMMIT_MSG=feat: release update !TIMESTAMP!"
    )
) else (
    set "COMMIT_MSG=feat: release %VERSION%"
)

git commit -m "%COMMIT_MSG%"
if %errorlevel% equ 0 (
    echo ✅ Changes committed
) else (
    echo ⚠️  No changes to commit or commit failed
)

:: Step 2: Push changes
echo.
echo 📤 Pushing changes to main...
git push origin main

:: Step 3: Trigger release
echo.
if "%VERSION%"=="" (
    echo 🔄 Triggering automatic release...
    echo GitHub Actions will auto-increment version and create release
) else (
    echo 🔄 Triggering manual release with version: %VERSION%
    git tag %VERSION%
    git push origin %VERSION%
    echo ✅ Tag %VERSION% created and pushed
)

:: Step 4: Open GitHub Actions page
echo.
echo 🌐 Opening GitHub Actions page...
start https://github.com/ohishi-yhonda-org/print_pdf/actions

echo.
echo 🎉 Release process initiated!
echo 📊 Check the Actions tab to monitor progress
echo 📦 Release will be available in 2-5 minutes
echo.
echo Next steps:
echo   1. Wait for CI to complete
echo   2. Check GitHub Releases page
echo   3. Download and test the new release
echo.
pause

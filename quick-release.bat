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

:: Step 2: Generate and push release tag
echo.
if "%VERSION%"=="" (
    echo 🔄 Generating next release version...
    
    :: Get latest release tag and increment
    for /f %%i in ('git tag --sort^=-version:refname') do (
        echo %%i | findstr /r "^v[0-9]*\.[0-9]*\.[0-9]*$" >nul
        if not errorlevel 1 (
            set "LATEST_TAG=%%i"
            goto :found_tag
        )
    )
    :found_tag
    
    if not "%LATEST_TAG%"=="" (
        echo Latest release tag found: %LATEST_TAG%
        :: Simple fallback increment - could be enhanced
        set "VERSION=v1.0.15"
        echo Generated version: %VERSION%
    ) else (
        set "VERSION=v1.0.14" 
        echo Using fallback version: %VERSION%
    )
) else (
    echo Using specified version: %VERSION%
)

echo 🚀 Pushing changes to main for testing...

:: First push changes to main to run tests
git push origin main

echo ✅ Changes pushed to main
echo ⏳ Waiting for CI tests to complete...

:: Wait a moment for CI to start
timeout /t 5 /nobreak > nul

:: Open GitHub Actions to monitor progress
echo 🌐 Opening GitHub Actions to monitor tests...
start https://github.com/ohishi-yhonda-org/print_pdf/actions

echo.
echo 📊 Please check GitHub Actions and confirm:
echo   - Test job: ✅ Passed
echo   - Lint job: ✅ Passed
echo.

set /p "CONFIRMATION=Are tests passing? Type 'yes' to create release tag, or 'no' to abort: "

if /i "%CONFIRMATION%"=="yes" (
    echo 🚀 Creating release tag: %VERSION%
    
    git tag %VERSION%
    git push origin %VERSION%
    
    echo ✅ Release tag %VERSION% created and pushed
    echo.
    echo 🎉 Release process initiated!
    echo 📦 Release will be available in 2-3 minutes ^(no duplicate testing^)
) else (
    echo ❌ Release aborted by user
    echo Fix any test failures and run the script again
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

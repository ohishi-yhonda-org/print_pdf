@echo off
REM PDF Generator Service Management Script

set SERVICE_NAME=PDFGeneratorService
set EXECUTABLE_PATH=%~dp0print_pdf.exe
set SERVICE_DESCRIPTION=PDF Generator HTTP API Service

echo PDF Generator Service Management
echo =================================
echo.
echo 1. Install Service
echo 2. Start Service  
echo 3. Stop Service
echo 4. Remove Service
echo 5. Check Status
echo 6. View Service Logs
echo.
set /p choice=Choose an option (1-6): 

if "%choice%"=="1" goto install
if "%choice%"=="2" goto start
if "%choice%"=="3" goto stop
if "%choice%"=="4" goto remove
if "%choice%"=="5" goto status
if "%choice%"=="6" goto logs
goto end

:install
echo Installing PDF Generator Service...
sc create %SERVICE_NAME% binPath= "%EXECUTABLE_PATH%" start= auto DisplayName= "PDF Generator API Service" 
sc description %SERVICE_NAME% "%SERVICE_DESCRIPTION%"
echo Service installed successfully!
goto end

:start
echo Starting PDF Generator Service...
sc start %SERVICE_NAME%
echo Service started!
goto end

:stop
echo Stopping PDF Generator Service...
sc stop %SERVICE_NAME%
echo Service stopped!
goto end

:remove
echo Removing PDF Generator Service...
sc stop %SERVICE_NAME%
sc delete %SERVICE_NAME%
echo Service removed!
goto end

:status
echo Checking service status...
sc query %SERVICE_NAME%
goto end

:logs
echo Opening Event Viewer for Application logs...
eventvwr.msc /c:Application
goto end

:end
echo.
pause

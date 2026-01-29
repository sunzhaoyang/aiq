@echo off
setlocal enabledelayedexpansion

REM AIQ Installation Script for Windows
REM Automatically detects latest version and installs aiq.exe binary

set "REPO=sunzhaoyang/aiq"
set "GITHUB_API=https://api.github.com/repos/%REPO%/releases/latest"
set "GITHUB_RELEASES=https://github.com/%REPO%/releases/download"

REM Installation directory
set "INSTALL_DIR=%USERPROFILE%\.aiq\bin"
set "BINARY_NAME=aiq.exe"

echo Detecting latest version...
REM Try to get latest version using PowerShell
for /f "delims=" %%i in ('powershell -Command "(Invoke-RestMethod -Uri '%GITHUB_API%').tag_name" 2^>nul') do set "LATEST_VERSION=%%i"

if "%LATEST_VERSION%"=="" (
    echo Warning: Failed to fetch latest version from GitHub API. Using v0.0.1 as fallback.
    set "LATEST_VERSION=v0.0.1"
)

echo Latest version: %LATEST_VERSION%

REM Detect architecture (assume amd64 for Windows)
set "ARCH=amd64"
set "PLATFORM=windows-%ARCH%"
echo Detected platform: %PLATFORM%

REM Construct download URL
set "GITHUB_URL=%GITHUB_RELEASES%/%LATEST_VERSION%/aiq-%PLATFORM%.exe"

REM Create installation directory
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
echo Install directory: %INSTALL_DIR%

REM Download binary
echo Downloading binary...
echo Download URL: %GITHUB_URL%
set "BINARY_PATH=%INSTALL_DIR%\%BINARY_NAME%"

REM Download from GitHub using PowerShell (no timeout - user can Ctrl+C if too slow)
powershell -Command "$ProgressPreference = 'Continue'; Invoke-WebRequest -Uri '%GITHUB_URL%' -OutFile '%BINARY_PATH%.tmp'" 2>&1
if %errorlevel% equ 0 (
    echo Downloaded successfully
) else (
    echo Error: Failed to download binary
    echo Please check your network or download manually from:
    echo   %GITHUB_URL%
    exit /b 1
)

REM Move temp file to final location
move /y "%BINARY_PATH%.tmp" "%BINARY_PATH%" >nul 2>&1

REM Verify installation
echo Verifying installation...
if exist "%BINARY_PATH%" (
    echo Installation successful!
    echo.
    
    REM Check if PATH contains INSTALL_DIR
    echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
    if %errorlevel% equ 0 (
        echo PATH already contains %INSTALL_DIR%
    ) else (
        echo To use 'aiq' command, add it to your PATH:
        echo   setx PATH "%%PATH%%;%INSTALL_DIR%"
        echo.
        echo Note: PATH changes will take effect in new terminal windows.
        echo Please close and reopen your terminal after running the setx command.
    )
) else (
    echo Warning: Installation completed but verification failed.
    echo Please check if %BINARY_PATH% exists.
    exit /b 1
)

endlocal

@echo off
setlocal enabledelayedexpansion

echo.
echo ==========================================
echo      GitHub Release Creator
echo ==========================================
echo.

rem Check if we're in a git repository
git status >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: This is not a git repository!
    pause
    exit /b 1
)

rem Check if we have uncommitted changes
git diff-index --quiet HEAD --
if %errorlevel% neq 0 (
    echo WARNING: You have uncommitted changes!
    echo Please commit your changes before creating a release.
    echo.
    git status --short
    echo.
    set /p "continue=Continue anyway? (y/N): "
    if /i not "!continue!"=="y" (
        echo Cancelled.
        pause
        exit /b 1
    )
)

rem Get current branch
for /f "tokens=*" %%i in ('git branch --show-current') do set CURRENT_BRANCH=%%i
echo Current branch: %CURRENT_BRANCH%

rem Check if we're on main branch
if not "%CURRENT_BRANCH%"=="main" (
    echo WARNING: You're not on the main branch!
    set /p "continue=Continue anyway? (y/N): "
    if /i not "!continue!"=="y" (
        echo Cancelled.
        pause
        exit /b 1
    )
)

rem Get the latest tag to suggest next version
echo.
echo Getting latest tags...
git fetch --tags >nul 2>&1

rem Get latest version tag
for /f "tokens=*" %%i in ('git tag -l "v*" --sort=-version:refname 2^>nul') do (
    set LATEST_TAG=%%i
    goto :found_tag
)
set LATEST_TAG=v0.0.0
:found_tag

echo Latest tag: %LATEST_TAG%
echo.

rem Suggest next versions
if "%LATEST_TAG%"=="v0.0.0" (
    echo Suggested versions:
    echo   1^) v1.0.0 ^(first release^)
    echo   2^) v0.1.0 ^(initial version^)
) else (
    rem Extract version numbers
    set VERSION_STR=%LATEST_TAG:~1%
    for /f "tokens=1,2,3 delims=." %%a in ("%VERSION_STR%") do (
        set /a MAJOR=%%a
        set /a MINOR=%%b
        set /a PATCH=%%c
    )
    
    set /a NEXT_PATCH=!PATCH!+1
    set /a NEXT_MINOR=!MINOR!+1
    set /a NEXT_MAJOR=!MAJOR!+1
    
    echo Suggested versions:
    echo   1^) v!MAJOR!.!MINOR!.!NEXT_PATCH! ^(patch - bug fixes^)
    echo   2^) v!MAJOR!.!NEXT_MINOR!.0 ^(minor - new features^)
    echo   3^) v!NEXT_MAJOR!.0.0 ^(major - breaking changes^)
)

echo.
set /p "NEW_VERSION=Enter new version (e.g., v1.0.0): "

rem Validate version format
echo %NEW_VERSION% | findstr /r "^v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*" >nul
if %errorlevel% neq 0 (
    echo ERROR: Invalid version format! Use format like v1.0.0
    pause
    exit /b 1
)

rem Check if tag already exists
git tag -l "%NEW_VERSION%" | findstr /x "%NEW_VERSION%" >nul
if %errorlevel% equ 0 (
    echo ERROR: Tag %NEW_VERSION% already exists!
    pause
    exit /b 1
)

echo.
echo Creating release %NEW_VERSION%...
echo.

rem Get release notes
echo Enter release notes (or just press Enter for default):
set /p "RELEASE_NOTES=>"
if not defined RELEASE_NOTES set "RELEASE_NOTES=Release %NEW_VERSION%"

echo.
echo ==========================================
echo Summary:
echo   Version: %NEW_VERSION%
echo   Branch: %CURRENT_BRANCH%
echo   Notes: !RELEASE_NOTES!
echo ==========================================
echo.

set /p "confirm=Create this release? (Y/n): "
if /i "!confirm!"=="n" (
    echo Cancelled.
    pause
    exit /b 0
)

echo.
echo Step 1/4: Pulling latest changes...
git pull origin %CURRENT_BRANCH%
if %errorlevel% neq 0 (
    echo ERROR: Failed to pull latest changes!
    pause
    exit /b 1
)

echo Step 2/4: Creating tag...
git tag -a "%NEW_VERSION%" -m "!RELEASE_NOTES!"
if %errorlevel% neq 0 (
    echo ERROR: Failed to create tag!
    pause
    exit /b 1
)

echo Step 3/4: Pushing tag to GitHub...
git push origin "%NEW_VERSION%"
if %errorlevel% neq 0 (
    echo ERROR: Failed to push tag!
    echo Removing local tag...
    git tag -d "%NEW_VERSION%"
    pause
    exit /b 1
)

echo Step 4/4: Done!
echo.
echo ==========================================
echo SUCCESS!
echo.
echo Release %NEW_VERSION% has been created!
echo.
echo What happens next:
echo 1. GitHub Actions will start building binaries
echo 2. Release will appear in GitHub Releases section
echo 3. Binaries will be attached automatically
echo.
echo You can monitor the progress at:
echo https://github.com/keshon/screen-tester/actions
echo ==========================================
echo.

pause
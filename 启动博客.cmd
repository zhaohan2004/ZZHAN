@echo off
setlocal
"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -NonInteractive -ExecutionPolicy Bypass -File "%~dp0start-local.ps1"
set "BLOG_EXIT_CODE=%ERRORLEVEL%"
echo.
pause
exit /b %BLOG_EXIT_CODE%

@echo off
setlocal
set TPS_INSTALL_SUBPROCESS=1
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1" %*

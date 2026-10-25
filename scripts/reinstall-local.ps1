# Rebuild tps from this repo into %LOCALAPPDATA%\Programs\tps (same path as install.ps1).
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$go = Get-Command go -ErrorAction SilentlyContinue
if (-not $go) {
    $candidate = Join-Path ${env:ProgramFiles} 'Go\bin\go.exe'
    if (Test-Path -LiteralPath $candidate) { $go = Get-Command $candidate }
}
if (-not $go) {
    Write-Host 'go not found. Install Go or open a terminal where go is on PATH.' -ForegroundColor Red
    exit 1
}
$dest = Join-Path $env:LOCALAPPDATA 'Programs\tps\tps.exe'
New-Item -ItemType Directory -Force -Path (Split-Path -Parent $dest) | Out-Null
Write-Host "Building $dest ..."
& $go.Source build -o $dest (Join-Path $root 'cmd\tps')
& $dest version
Write-Host 'Done. Restart tps in this terminal (or open a new one) before detaching jobs.' -ForegroundColor Green

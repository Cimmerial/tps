#requires -Version 5.1
$ErrorActionPreference = "Stop"

$InstallDir = if ($env:TPS_INSTALL_DIR) { $env:TPS_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\tps" }

Get-Process -Name tps -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 200

if (Test-Path $InstallDir) {
    Remove-Item -Recurse -Force -LiteralPath $InstallDir
    Write-Host "removed $InstallDir"
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not [string]::IsNullOrEmpty($userPath)) {
    $parts = @($userPath -split ';' | Where-Object { $_ -and $_ -ne $InstallDir })
    $newUserPath = ($parts -join ';')
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
}

$profile = $PROFILE
if (Test-Path -LiteralPath $profile) {
    $lines = @(Get-Content -LiteralPath $profile)
    $filtered = @($lines | Where-Object {
        $_ -notlike "*$InstallDir*" -and $_ -notmatch 'Programs\\tps'
    })
    if ($filtered.Count -ne $lines.Count) {
        $substantive = @($filtered | Where-Object { $_.Trim() -ne '' -and $_.Trim() -notmatch '^\s*#' })
        if ($substantive.Count -eq 0) {
            Remove-Item -LiteralPath $profile -Force
        } else {
            Set-Content -LiteralPath $profile -Value $filtered
        }
    }
}

$machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$env:Path = if ([string]::IsNullOrEmpty($machinePath)) { $userPath } else { "$machinePath;$userPath" }

Write-Host "tps uninstalled (config in %AppData%\tps\config.yaml was kept)"

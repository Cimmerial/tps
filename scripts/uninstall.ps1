#requires -Version 5.1
$ErrorActionPreference = "Stop"

$InstallDir = if ($env:TPS_INSTALL_DIR) { $env:TPS_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\tps" }

Get-Process -Name tps -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 200

if (Test-Path $InstallDir) {
    Remove-Item -LiteralPath $InstallDir -Recurse -Force
    Write-Host "removed $InstallDir"
} else {
    Write-Host "no install dir at $InstallDir"
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not [string]::IsNullOrEmpty($userPath)) {
    $parts = $userPath -split ';' | Where-Object {
        $_ -and ($_ -ne $InstallDir) -and ($_.TrimEnd('\') -ne $InstallDir.TrimEnd('\'))
    }
    $newUserPath = ($parts | Select-Object -Unique) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
    Write-Host "removed $InstallDir from user PATH"
}

$profile = $PROFILE
if (Test-Path $profile) {
    $original = @(Get-Content -LiteralPath $profile)
    $filtered = @($original | Where-Object { $_ -notmatch [regex]::Escape($InstallDir) })
    if ($filtered.Count -ne $original.Count) {
        Set-Content -LiteralPath $profile -Value $filtered -Encoding utf8
        Write-Host "removed PATH line from $profile"
    }
}

$machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$env:Path = if ([string]::IsNullOrEmpty($machinePath)) { $userPath } else { "$machinePath;$userPath" }

Write-Host "tps uninstalled (this window PATH updated; restart other terminals or your editor if tps still resolves there)"

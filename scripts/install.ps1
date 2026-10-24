#requires -Version 5.1
$ErrorActionPreference = "Stop"

$Repo = "Cimmerial/tps"
$InstallDir = if ($env:TPS_INSTALL_DIR) { $env:TPS_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\tps" }
$Version = if ($env:TPS_VERSION) { $env:TPS_VERSION } else { "latest" }

function Install-GoModule {
    param([string]$Ver = "latest")
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        throw "go not found; install Go or publish a GitHub release"
    }
    $modVer = $Ver
    if ($modVer -eq "latest") {
        $rel = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $modVer = $rel.tag_name
    }
    if ($modVer -notmatch "^v") {
        $modVer = "v$modVer"
    }
    if ((Test-Path go.mod) -and (Test-Path cmd\tps)) {
        Write-Host "installing from local clone..."
        go install ./cmd/tps
    } else {
        Write-Host "installing from module $modVer..."
        go install "github.com/$Repo/cmd/tps@$modVer"
    }
    $goBin = Join-Path (go env GOPATH) "bin"
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    Install-Binary (Join-Path $goBin "tps.exe")
}

function Stop-TpsIfRunning {
    Get-Process -Name tps -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
    Start-Sleep -Milliseconds 200
}

function Install-Binary {
    param([string]$SourceExe)
    if (-not (Test-Path $SourceExe)) {
        throw "missing binary: $SourceExe"
    }
    Stop-TpsIfRunning
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $dest = Join-Path $InstallDir "tps.exe"
    $staging = "$dest.new"
    Copy-Item -LiteralPath $SourceExe -Destination $staging -Force
    try {
        Move-Item -LiteralPath $staging -Destination $dest -Force
    } catch {
        throw @"
cannot replace $dest (is tps running?). Close tps and PowerShell tabs using it, then run:

  Move-Item -LiteralPath '$staging' -Destination '$dest' -Force

Or: taskkill /IM tps.exe /F
then re-run: irm https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.ps1 | iex
"@
    }
}

function Install-Release {
    param([string]$Ver)
    if ($Ver -eq "latest") {
        $rel = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $Ver = $rel.tag_name
    }
    $tag = $Ver.TrimStart("v")
    $asset = "tps_${tag}_windows_amd64.zip"
    $base = "https://github.com/$Repo/releases/download/v$tag"
    $tmp = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_.FullName }
    Write-Host "downloading $asset..."
    $zipPath = Join-Path $tmp.FullName "archive.zip"
    Invoke-WebRequest -Uri "$base/$asset" -OutFile $zipPath -UseBasicParsing
    Expand-Archive -Path $zipPath -DestinationPath $tmp.FullName -Force
    $src = Join-Path $tmp.FullName "tps.exe"
    if (-not (Test-Path $src)) {
        throw "archive did not contain tps.exe"
    }
    Install-Binary $src
    Remove-Item -Recurse -Force $tmp.FullName
    return $tag
}

function Remove-LegacyProfilePath {
    $profile = $PROFILE
    if (-not (Test-Path -LiteralPath $profile)) {
        return
    }
    $lines = @(Get-Content -LiteralPath $profile -ErrorAction SilentlyContinue)
    $filtered = @($lines | Where-Object {
        $_ -notlike "*$InstallDir*" -and $_ -notmatch 'Programs\\tps'
    })
    if ($filtered.Count -eq $lines.Count) {
        return
    }
    $substantive = @($filtered | Where-Object { $_.Trim() -ne '' -and $_.Trim() -notmatch '^\s*#' })
    if ($substantive.Count -eq 0) {
        Remove-Item -LiteralPath $profile -Force
        Write-Host "removed legacy PowerShell profile hook (user PATH is enough)"
    } else {
        Set-Content -LiteralPath $profile -Value $filtered
        Write-Host "removed legacy PowerShell profile hook (user PATH is enough)"
    }
}

function Refresh-SessionPath {
    $machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $env:Path = if ([string]::IsNullOrEmpty($machinePath)) { $userPath } else { "$machinePath;$userPath" }
}

function Show-ParentSessionPathHint {
    Write-Host ""
    Write-Host "install ran in a subprocess; refresh PATH in this window, then run tps version:"
    Write-Host '$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")' -ForegroundColor Cyan
}

function Ensure-Path {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        if ([string]::IsNullOrEmpty($userPath)) {
            [Environment]::SetEnvironmentVariable("Path", $InstallDir, "User")
        } else {
            [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        }
        Write-Host "added $InstallDir to user PATH"
    }
    Remove-LegacyProfilePath
    Refresh-SessionPath
}

$resolvedTag = $null
try {
    $resolvedTag = Install-Release -Ver $Version
    Write-Host "installed to $InstallDir\tps.exe"
} catch {
    Write-Host "release install failed: $($_.Exception.Message)" -ForegroundColor Yellow
    if (-not $resolvedTag) {
        try {
            $rel = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
            $resolvedTag = $rel.tag_name.TrimStart("v")
        } catch {
            $resolvedTag = "latest"
        }
    }
    Write-Host "trying go install..." -ForegroundColor Yellow
    Install-GoModule -Ver $resolvedTag
    Write-Host "installed to $InstallDir\tps.exe"
}

Ensure-Path
if ($env:TPS_INSTALL_SUBPROCESS -eq '1') {
    Show-ParentSessionPathHint
} else {
    Write-Host "run: tps version"
}

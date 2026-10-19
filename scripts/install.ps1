#requires -Version 5.1
$ErrorActionPreference = "Stop"

$Repo = "Cimmerial/tps"
$InstallDir = if ($env:TPS_INSTALL_DIR) { $env:TPS_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\tps" }
$Version = if ($env:TPS_VERSION) { $env:TPS_VERSION } else { "latest" }

function Install-GoModule {
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        throw "go not found; install Go or publish a GitHub release"
    }
    if ((Test-Path go.mod) -and (Test-Path cmd\tps)) {
        Write-Host "installing from local clone..."
        go install ./cmd/tps
    } else {
        Write-Host "installing from module..."
        go install "github.com/$Repo/cmd/tps@latest"
    }
    $goBin = Join-Path (go env GOPATH) "bin"
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    Copy-Item -Force (Join-Path $goBin "tps.exe") (Join-Path $InstallDir "tps.exe")
}

function Install-Release {
    param([string]$Ver)
    if ($Ver -eq "latest") {
        $rel = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $Ver = $rel.tag_name
    }
    $Ver = $Ver.TrimStart("v")
    $asset = "tps_${Ver}_windows_amd64.zip"
    $base = "https://github.com/$Repo/releases/download/v$Ver"
    $tmp = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_.FullName }
    Write-Host "downloading $asset..."
    Invoke-WebRequest -Uri "$base/$asset" -OutFile (Join-Path $tmp.FullName "archive.zip")
    Expand-Archive -Path (Join-Path $tmp.FullName "archive.zip") -DestinationPath $tmp.FullName -Force
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    Copy-Item -Force (Join-Path $tmp.FullName "tps.exe") (Join-Path $InstallDir "tps.exe")
    Remove-Item -Recurse -Force $tmp.FullName
}

function Ensure-Path {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        Write-Host "added $InstallDir to user PATH (open a new terminal)"
    }
    $profile = $PROFILE
    if (-not (Test-Path $profile)) {
        New-Item -ItemType File -Path $profile -Force | Out-Null
    }
    $line = '$env:Path += ";' + $InstallDir + '"'
    if (-not (Select-String -Path $profile -Pattern $InstallDir -Quiet -ErrorAction SilentlyContinue)) {
        Add-Content -Path $profile -Value $line
    }
}

try {
    Install-Release -Ver $Version
    Write-Host "installed to $InstallDir\tps.exe"
} catch {
    Write-Host "release install failed; trying go install..." -ForegroundColor Yellow
    Install-GoModule
}

Ensure-Path
Write-Host "run: tps version"

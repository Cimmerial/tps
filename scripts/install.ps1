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

function Install-Binary {
    param([string]$SourceExe)
    if (-not (Test-Path $SourceExe)) {
        throw "missing binary: $SourceExe"
    }
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
    if (-not (Select-String -Path $profile -Pattern $InstallDir -SimpleMatch -Quiet -ErrorAction SilentlyContinue)) {
        Add-Content -Path $profile -Value $line
    }
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
Write-Host "run: tps version"

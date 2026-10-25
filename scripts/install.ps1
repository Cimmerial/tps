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

function Broadcast-EnvironmentChange {
    try {
        if (-not ([System.Management.Automation.PSTypeName]'TpsNativeMethods').Type) {
            Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public static class TpsNativeMethods {
    [DllImport("user32.dll", CharSet = CharSet.Auto, SetLastError = true)]
    public static extern IntPtr SendMessageTimeout(IntPtr hWnd, int Msg, UIntPtr wParam, string lParam, int fuFlags, int uTimeout, out IntPtr lpdwResult);
}
'@ -ErrorAction Stop
        }
        $unused = [IntPtr]::Zero
        [TpsNativeMethods]::SendMessageTimeout([IntPtr]0xffff, 0x1a, [UIntPtr]::Zero, "Environment", 2, 5000, [ref]$unused) | Out-Null
    } catch {
        # non-fatal
    }
}

function Show-ParentSessionPathHint {
    Write-Host ""
    Write-Host "install ran in a subprocess; refresh PATH in this window, then run tps version:"
    Write-Host '$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")' -ForegroundColor Cyan
}

function Find-GitExe {
    $cmd = Get-Command git -ErrorAction SilentlyContinue
    if ($cmd) {
        return $cmd.Source
    }
    $candidates = @(
        (Join-Path ${env:ProgramFiles} "Git\cmd\git.exe"),
        (Join-Path ${env:ProgramFiles(x86)} "Git\cmd\git.exe")
    )
    foreach ($p in $candidates) {
        if ($p -and (Test-Path -LiteralPath $p)) {
            return $p
        }
    }
    $desktop = Get-ChildItem -Path (Join-Path $env:LOCALAPPDATA "GitHubDesktop\app-*\resources\app\git\cmd\git.exe") -ErrorAction SilentlyContinue |
        Sort-Object FullName -Descending |
        Select-Object -First 1
    if ($desktop) {
        return $desktop.FullName
    }
    return $null
}

function Ensure-GitOnUserPath {
    $gitExe = Find-GitExe
    if (-not $gitExe) {
        return
    }
    $gitDir = Split-Path -Parent $gitExe
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -like "*$gitDir*") {
        return
    }
    if ([string]::IsNullOrEmpty($userPath)) {
        [Environment]::SetEnvironmentVariable("Path", $gitDir, "User")
    } else {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$gitDir", "User")
    }
    Write-Host "added $gitDir to user PATH (git)"
    Broadcast-EnvironmentChange
    Refresh-SessionPath
}

function Test-InstallInteractive {
    try {
        return [Environment]::UserInteractive
    } catch {
        return $false
    }
}

function Test-SkipGitPrompt {
    return $env:TPS_SKIP_GIT_PROMPT -eq '1' -or $env:TPS_SKIP_TOOLCHAIN_PROMPT -eq '1'
}

function Test-SkipGoPrompt {
    return $env:TPS_SKIP_GO_PROMPT -eq '1' -or $env:TPS_SKIP_TOOLCHAIN_PROMPT -eq '1'
}

function Add-DirToUserPath {
    param([string]$Dir, [string]$Label)
    if ([string]::IsNullOrWhiteSpace($Dir)) {
        return
    }
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -like "*$Dir*") {
        return
    }
    if ([string]::IsNullOrEmpty($userPath)) {
        [Environment]::SetEnvironmentVariable("Path", $Dir, "User")
    } else {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$Dir", "User")
    }
    Write-Host "added $Dir to user PATH ($Label)"
    Broadcast-EnvironmentChange
    Refresh-SessionPath
}

function Offer-InstallGit {
    if (Find-GitExe) {
        Ensure-GitOnUserPath
        return
    }
    if (Test-SkipGitPrompt) {
        return
    }
    if (-not (Test-InstallInteractive)) {
        Write-Host "git not found; install Git for tps updates (https://git-scm.com/download/win)" -ForegroundColor Yellow
        return
    }
    $answer = Read-Host "Git is not installed (needed for tps update and git ship). Install Git now? [Y/n]"
    if ($answer -match '^[Nn]') {
        return
    }
    if (Get-Command winget -ErrorAction SilentlyContinue) {
        Write-Host "installing Git.Git via winget..."
        $winget = Start-Process -FilePath winget -ArgumentList @(
            'install', '--id', 'Git.Git', '-e', '--source', 'winget',
            '--accept-package-agreements', '--accept-source-agreements'
        ) -Wait -PassThru -NoNewWindow
        if ($winget.ExitCode -ne 0) {
            Write-Host "winget install failed (exit $($winget.ExitCode)); install manually: https://git-scm.com/download/win" -ForegroundColor Yellow
            return
        }
    } else {
        Write-Host "winget not found; open https://git-scm.com/download/win to install Git" -ForegroundColor Yellow
        Start-Process "https://git-scm.com/download/win"
        return
    }
    if (Find-GitExe) {
        Ensure-GitOnUserPath
        Write-Host "git installed"
    } else {
        Write-Host "git installed; open a new terminal if git is not found yet" -ForegroundColor Yellow
    }
}

function Find-GoExe {
    $cmd = Get-Command go -ErrorAction SilentlyContinue
    if ($cmd) {
        return $cmd.Source
    }
    $candidate = Join-Path ${env:ProgramFiles} "Go\bin\go.exe"
    if (Test-Path -LiteralPath $candidate) {
        return $candidate
    }
    return $null
}

function Ensure-GoOnUserPath {
    $goExe = Find-GoExe
    if (-not $goExe) {
        return
    }
    $goDir = Split-Path -Parent $goExe
    Add-DirToUserPath $goDir "go"
    $goBin = Join-Path $env:USERPROFILE "go\bin"
    try {
        $gopath = & $goExe env GOPATH 2>$null
        if ($gopath) {
            $goBin = Join-Path $gopath.Trim() "bin"
        }
    } catch {
        # use default goBin
    }
    Add-DirToUserPath $goBin "go install bin"
}

function Offer-InstallGo {
    if (Find-GoExe) {
        Ensure-GoOnUserPath
        return
    }
    if (Test-SkipGoPrompt) {
        return
    }
    if (-not (Test-InstallInteractive)) {
        Write-Host "go not found; install Go for tps dev updates (https://go.dev/dl/)" -ForegroundColor Yellow
        return
    }
    $answer = Read-Host "Go is not installed (needed for tps u / go install from the dev repo). Install Go now? [Y/n]"
    if ($answer -match '^[Nn]') {
        return
    }
    if (Get-Command winget -ErrorAction SilentlyContinue) {
        Write-Host "installing Go via winget..."
        $winget = Start-Process -FilePath winget -ArgumentList @(
            'install', '--id', 'GoLang.Go', '-e', '--source', 'winget',
            '--accept-package-agreements', '--accept-source-agreements'
        ) -Wait -PassThru -NoNewWindow
        if ($winget.ExitCode -ne 0) {
            Write-Host "winget install failed (exit $($winget.ExitCode)); install manually: https://go.dev/dl/" -ForegroundColor Yellow
            return
        }
    } else {
        Write-Host "winget not found; open https://go.dev/dl/ to install Go" -ForegroundColor Yellow
        Start-Process "https://go.dev/dl/"
        return
    }
    Refresh-SessionPath
    if (Find-GoExe) {
        Ensure-GoOnUserPath
        Write-Host "go installed"
    } else {
        Write-Host "go installed; open a new terminal if go is not found yet" -ForegroundColor Yellow
    }
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
        Broadcast-EnvironmentChange
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
Offer-InstallGit
Offer-InstallGo
if ($env:TPS_INSTALL_SUBPROCESS -eq '1') {
    Show-ParentSessionPathHint
} else {
    Write-Host "run: tps version"
    Write-Host "terminals opened before install (e.g. inside Cursor) may need that app restarted"
}

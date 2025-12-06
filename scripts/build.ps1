$ErrorActionPreference = "Stop"

# Build script for Windows (PowerShell)
# Steps:
# 1) Build frontend in viewer/viewer-frontend
# 2) Ensure output is placed in viewer/service/router/dist, also copy to troll/dist
# 3) Build viewer cmd program (outputs viewer.exe)
# 4) Build troll main program (outputs troll.exe)

function Ensure-Command {
    param(
        [Parameter(Mandatory=$true)][string]$Name
    )
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Required command '$Name' not found in PATH."
    }
}

$ScriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot = (Resolve-Path (Join-Path $ScriptRoot '..')).Path

$FrontendDir = Join-Path $RepoRoot 'viewer\viewer-frontend'
$RouterDistDir = Join-Path $RepoRoot 'viewer\service\app\dist'
$TrollDistDir = Join-Path $RepoRoot 'troll\dist'

Write-Host "[1/4] Building frontend (viewer/viewer-frontend)" -ForegroundColor Cyan

Push-Location $FrontendDir
try {
    # Prefer pnpm, fallback to npm
    $buildTool = $null
    if (Get-Command pnpm -ErrorAction SilentlyContinue) {
        $buildTool = 'pnpm'
        Write-Host "Using pnpm" -ForegroundColor Green
        pnpm install --frozen-lockfile
        pnpm run build
    } elseif (Get-Command npm -ErrorAction SilentlyContinue) {
        $buildTool = 'npm'
        Write-Host "Using npm" -ForegroundColor Yellow
        npm ci
        npm run build
    } else {
        throw "Neither 'pnpm' nor 'npm' found. Please install one of them."
    }
}
finally {
    Pop-Location
}

Write-Host "[2/4] Sync dist to viewer/service/router and troll/dist" -ForegroundColor Cyan

# Create router dist if not exists
if (-not (Test-Path $RouterDistDir)) { New-Item -ItemType Directory -Path $RouterDistDir | Out-Null }

# If build output already went to router/dist via Vite config, we still ensure copy from frontend/dist as fallback
$FrontendDistDir = Join-Path $FrontendDir 'dist'
if (Test-Path $FrontendDistDir) {
    # Clean target dirs
    if (Test-Path $RouterDistDir) { Remove-Item -Recurse -Force $RouterDistDir }
    New-Item -ItemType Directory -Path $RouterDistDir | Out-Null
    Copy-Item -Path (Join-Path $FrontendDistDir '*') -Destination $RouterDistDir -Recurse -Force
}

# Sync to troll/dist
if (Test-Path $TrollDistDir) { Remove-Item -Recurse -Force $TrollDistDir }
New-Item -ItemType Directory -Path $TrollDistDir | Out-Null
Copy-Item -Path (Join-Path $RouterDistDir '*') -Destination $TrollDistDir -Recurse -Force

Write-Host "[3/4] Building viewer cmd program" -ForegroundColor Cyan

$ViewerDir = Join-Path $RepoRoot 'viewer'
Push-Location $ViewerDir
try {
    Ensure-Command go
    $BinDir = Join-Path $ViewerDir 'bin'
    if (-not (Test-Path $BinDir)) { New-Item -ItemType Directory -Path $BinDir | Out-Null }
    go build -o (Join-Path $BinDir 'viewer.exe') ./cmd
}
finally {
    Pop-Location
}

Write-Host "[4/4] Building troll main program" -ForegroundColor Cyan

$TrollDir = Join-Path $RepoRoot 'troll'
Push-Location $TrollDir
try {
    Ensure-Command go
    go build -o troll.exe
}
finally {
    Pop-Location
}

Write-Host "Build complete." -ForegroundColor Green
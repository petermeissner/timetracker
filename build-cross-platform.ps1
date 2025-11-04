# Cross-platform build script for Timesheet application
# Builds binaries for Windows, Linux, and macOS

Write-Host "🔧 Building Timesheet binaries for multiple platforms..." -ForegroundColor Green

# Create distribution directory if it doesn't exist
$distDir = "distribution"
if (!(Test-Path $distDir)) {
    Write-Host "📁 Creating distribution directory..." -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $distDir | Out-Null
}

# Clean previous builds
Write-Host "🧹 Cleaning previous builds..." -ForegroundColor Yellow
Remove-Item "$distDir\*" -Force -ErrorAction SilentlyContinue

# Build settings
$appName = "timesheet"

Write-Host "📦 Building binaries..." -ForegroundColor Green

# Windows (amd64)
Write-Host "  → Windows (amd64)..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o "$distDir\$appName-windows-amd64.exe" .
if ($LASTEXITCODE -eq 0) {
    Write-Host "    ✅ Windows build successful" -ForegroundColor Green
} else {
    Write-Host "    ❌ Windows build failed" -ForegroundColor Red
}

# Linux (amd64)
Write-Host "  → Linux (amd64)..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o "$distDir\$appName-linux-amd64" .
if ($LASTEXITCODE -eq 0) {
    Write-Host "    ✅ Linux build successful" -ForegroundColor Green
} else {
    Write-Host "    ❌ Linux build failed" -ForegroundColor Red
}

# Linux (arm64) - for Raspberry Pi and ARM servers
Write-Host "  → Linux (arm64)..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -ldflags "-s -w" -o "$distDir\$appName-linux-arm64" .
if ($LASTEXITCODE -eq 0) {
    Write-Host "    ✅ Linux ARM64 build successful" -ForegroundColor Green
} else {
    Write-Host "    ❌ Linux ARM64 build failed" -ForegroundColor Red
}

# macOS (amd64) - Intel Macs
Write-Host "  → macOS (amd64)..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o "$distDir\$appName-macos-amd64" .
if ($LASTEXITCODE -eq 0) {
    Write-Host "    ✅ macOS Intel build successful" -ForegroundColor Green
} else {
    Write-Host "    ❌ macOS Intel build failed" -ForegroundColor Red
}

# macOS (arm64) - Apple Silicon Macs
Write-Host "  → macOS (arm64)..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -ldflags "-s -w" -o "$distDir\$appName-macos-arm64" .
if ($LASTEXITCODE -eq 0) {
    Write-Host "    ✅ macOS Apple Silicon build successful" -ForegroundColor Green
} else {
    Write-Host "    ❌ macOS Apple Silicon build failed" -ForegroundColor Red
}

# Reset environment variables
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "🎉 Build process completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Built binaries:" -ForegroundColor Yellow

# List built files with sizes
Get-ChildItem $distDir | ForEach-Object {
    $sizeKB = [math]::Round($_.Length / 1KB, 1)
    $sizeMB = [math]::Round($_.Length / 1MB, 1)
    if ($sizeMB -gt 1) {
        Write-Host "  📄 $($_.Name) ($sizeMB MB)" -ForegroundColor White
    } else {
        Write-Host "  📄 $($_.Name) ($sizeKB KB)" -ForegroundColor White
    }
}

Write-Host ""
Write-Host "📦 Distribution files are in the '$distDir' directory" -ForegroundColor Green
Write-Host "🚀 Copy the appropriate binary to your target system and run it!" -ForegroundColor Green
# Security scan automation — Tech Challenge Oficina API
# Usage: .\scripts\security-scan.ps1

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$ReportDir = Join-Path $Root "docs\security"

Write-Host "=== Tech Challenge — Security Scan ===" -ForegroundColor Cyan
Write-Host "Project root: $Root"

if (-not (Test-Path $ReportDir)) {
    New-Item -ItemType Directory -Path $ReportDir -Force | Out-Null
}

$goBin = Join-Path $env:USERPROFILE "go\bin"
if ($env:Path -notlike "*$goBin*") {
    $env:Path = "$goBin;$env:Path"
}

Write-Host "`n[1/4] Installing scan tools..." -ForegroundColor Yellow
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

Write-Host "`n[2/4] Running gosec (SAST)..." -ForegroundColor Yellow
Push-Location $Root
try {
    gosec -fmt json -out (Join-Path $ReportDir "gosec-report.json") ./... 2>&1 | Tee-Object -FilePath (Join-Path $ReportDir "gosec-output.txt")
    gosec -fmt html -out (Join-Path $ReportDir "gosec-report.html") ./... 2>&1 | Out-Null
} finally {
    Pop-Location
}

Write-Host "`n[3/4] Running govulncheck (SCA)..." -ForegroundColor Yellow
Push-Location $Root
try {
    govulncheck ./... 2>&1 | Tee-Object -FilePath (Join-Path $ReportDir "govulncheck-output.txt")
    govulncheck -json ./... 2>&1 | Out-File -FilePath (Join-Path $ReportDir "govulncheck-report.json") -Encoding utf8
} finally {
    Pop-Location
}

Write-Host "`n[4/4] Building Docker image and scanning..." -ForegroundColor Yellow
Push-Location $Root
try {
    docker build -t techchalleng-oficina:latest . 2>&1 | Tee-Object -FilePath (Join-Path $ReportDir "docker-build-output.txt")

    $scoutCmd = Get-Command docker -ErrorAction SilentlyContinue
    if ($scoutCmd) {
        docker scout cves techchalleng-oficina:latest 2>&1 | Tee-Object -FilePath (Join-Path $ReportDir "docker-scout-output.txt")
    } else {
        Write-Host "Docker not available — skipping container scan." -ForegroundColor DarkYellow
    }
} catch {
    Write-Host "Docker scan skipped: $_" -ForegroundColor DarkYellow
} finally {
    Pop-Location
}

Write-Host "`n=== Scan complete ===" -ForegroundColor Green
Write-Host "Reports saved to: $ReportDir"

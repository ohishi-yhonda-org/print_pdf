# PowerShell script for development tasks

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "利用可能なコマンド:" -ForegroundColor Green
    Write-Host "  .\dev.ps1 test      - すべてのテストを実行" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1 build     - アプリケーションをビルド" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1 clean     - 生成ファイルを削除" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1 lint      - コードリンティングを実行" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1 fmt       - コードフォーマットを実行" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1 vet       - go vetを実行" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1 coverage  - テストカバレッジを表示" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1 ci        - 全チェック（CI相当）" -ForegroundColor Yellow
}

function Run-Tests {
    Write-Host "テストを実行中..." -ForegroundColor Green
    go test -v ./...
}

function Run-TestsWithRace {
    Write-Host "レース検査付きテストを実行中..." -ForegroundColor Green
    go test -race -short ./...
}

function Run-Coverage {
    Write-Host "テストカバレッジを計算中..." -ForegroundColor Green
    go test -coverprofile=coverage.out ./...
    if (Test-Path coverage.out) {
        go tool cover -html=coverage.out -o coverage.html
        go tool cover -func=coverage.out
        Write-Host "カバレッジレポートが coverage.html に生成されました" -ForegroundColor Green
    }
}

function Build-App {
    Write-Host "アプリケーションをビルド中..." -ForegroundColor Green
    go build -v ./...
}

function Clean-Files {
    Write-Host "生成ファイルを削除中..." -ForegroundColor Green
    go clean
    if (Test-Path coverage.out) { Remove-Item coverage.out }
    if (Test-Path coverage.html) { Remove-Item coverage.html }
    if (Test-Path "*.pdf") { Remove-Item *.pdf }
}

function Run-Lint {
    Write-Host "リンティングを実行中..." -ForegroundColor Green
    if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
        golangci-lint run
    } else {
        Write-Host "golangci-lint がインストールされていません" -ForegroundColor Red
        Write-Host "インストール方法: https://golangci-lint.run/usage/install/" -ForegroundColor Yellow
    }
}

function Format-Code {
    Write-Host "コードをフォーマット中..." -ForegroundColor Green
    go fmt ./...
    if (Get-Command goimports -ErrorAction SilentlyContinue) {
        goimports -w .
    } else {
        Write-Host "goimports がインストールされていません" -ForegroundColor Yellow
        Write-Host "インストール: go install golang.org/x/tools/cmd/goimports@latest" -ForegroundColor Yellow
    }
}

function Run-Vet {
    Write-Host "go vet を実行中..." -ForegroundColor Green
    go vet ./...
}

function Run-CI {
    Write-Host "CI相当のチェックを実行中..." -ForegroundColor Green
    Format-Code
    Run-Vet
    Run-Lint
    Run-Tests
    Run-Coverage
}

# コマンド実行
switch ($Command.ToLower()) {
    "test" { Run-Tests }
    "test-race" { Run-TestsWithRace }
    "build" { Build-App }
    "clean" { Clean-Files }
    "lint" { Run-Lint }
    "fmt" { Format-Code }
    "vet" { Run-Vet }
    "coverage" { Run-Coverage }
    "ci" { Run-CI }
    "help" { Show-Help }
    default { 
        Write-Host "不明なコマンド: $Command" -ForegroundColor Red
        Show-Help 
    }
}

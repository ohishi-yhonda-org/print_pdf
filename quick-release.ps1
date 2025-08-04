# クイックリリーススクリプト
# 使用方法: .\quick-release.ps1 [version]
# 例: .\quick-release.ps1 v1.0.15

param(
    [string]$Version = ""
)

# UTF-8エンコーディング設定
$OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$env:LC_ALL = "C.UTF-8"

Write-Host "🚀 Quick Release Script Starting..." -ForegroundColor Green

# Step 0: Confirm release intention
Write-Host ""
Write-Host "📋 Release Confirmation" -ForegroundColor Cyan
Write-Host "This script will:" -ForegroundColor White
Write-Host "  Default (Enter): Just commit & push changes" -ForegroundColor Green
Write-Host "  'y': Full release workflow with CI" -ForegroundColor White
Write-Host "  'n': Cancel" -ForegroundColor White
Write-Host ""

$releaseConfirm = Read-Host "Action? (Enter=push only, y=release, n=cancel)"

# Default to push-only if Enter pressed
if ([string]::IsNullOrWhiteSpace($releaseConfirm)) {
    Write-Host "📝 Push-only mode (default)" -ForegroundColor Yellow
    
    # Step 1: Commit current changes
    Write-Host "📝 Committing current changes..." -ForegroundColor Yellow
    try {
        git add .
        $commitMessage = Read-Host "Enter commit message (or press Enter for default)"
        if ([string]::IsNullOrWhiteSpace($commitMessage)) {
            $commitMessage = "feat: update $(Get-Date -Format 'yyyy-MM-dd HH:mm')"
        }
        
        # UTF-8エンコーディングでコミット
        $env:LC_ALL = "C.UTF-8"
        git -c core.quotepath=false commit -m $commitMessage
        Write-Host "✅ Changes committed" -ForegroundColor Green
    } catch {
        Write-Host "⚠️  No changes to commit or commit failed" -ForegroundColor Yellow
    }
    
    # Push to main
    git push origin main
    Write-Host "✅ Changes pushed to main" -ForegroundColor Green
    Write-Host "🎯 Push-only complete. No release created." -ForegroundColor Cyan
    exit
}

if ($releaseConfirm -eq "y" -or $releaseConfirm -eq "yes") {
    Write-Host "🚀 Full release mode selected" -ForegroundColor Green
    # Continue to full release workflow below
} elseif ($releaseConfirm -eq "n" -or $releaseConfirm -eq "no") {
    Write-Host "❌ Cancelled by user" -ForegroundColor Red
    exit
} else {
    Write-Host "❌ Invalid option. Use: Enter (push only), y (release), n (cancel)" -ForegroundColor Red
    exit
}

Write-Host "✅ Release confirmed, proceeding..." -ForegroundColor Green

# 最新コミットが既に [release] フラグを持っているかチェック
$latestCommitMessage = git -c core.quotepath=false log -1 --pretty=format:"%s"
if ($latestCommitMessage -match "\[release\]") {
    Write-Host "🏷️  Latest commit already has [release] flag: $latestCommitMessage" -ForegroundColor Cyan
    Write-Host "📤 Proceeding to push for CI auto-tagging" -ForegroundColor Yellow
} else {
    # Step 1: Add release flag to existing commit
    Write-Host "📝 Adding release flag to current commit..." -ForegroundColor Yellow
    try {
        # 変更がある場合は先にコミット
        $hasChanges = (git status --porcelain) -ne $null
        
        if ($hasChanges) {
            git add .
            $commitMessage = Read-Host "コミットメッセージを入力してください（Enterでデフォルト）"
            if ([string]::IsNullOrWhiteSpace($commitMessage)) {
                $commitMessage = "feat: 日本語更新 $(Get-Date -Format 'yyyy-MM-dd HH:mm')"
            }
            # UTF-8エンコーディングで日本語コミット
            $env:GIT_COMMITTER_NAME = "${env:USERNAME}"
            $env:LC_ALL = "C.UTF-8"
            git -c i18n.commitEncoding=utf-8 -c core.quotepath=false commit -m "$commitMessage"
            Write-Host "✅ New changes committed" -ForegroundColor Green
        }
        
        # 最新コミットのメッセージに [release] フラグを追加
        $currentMessage = git -c core.quotepath=false log -1 --pretty=format:"%s"
        $newMessage = "$currentMessage [release]"
        # UTF-8エンコーディングでコミット修正
        $env:LC_ALL = "C.UTF-8"
        git -c core.quotepath=false commit --amend -m $newMessage
        Write-Host "✅ Current commit amended with release flag" -ForegroundColor Green
        Write-Host ("📝 Updated message: {0}" -f $newMessage) -ForegroundColor Cyan
    } catch {
        Write-Host "⚠️  No changes to commit or commit amendment failed" -ForegroundColor Yellow
        Write-Host "❌ Please ensure you have changes to commit before proceeding" -ForegroundColor Red
        exit
    }
}

# Step 2: Push amended commit and trigger CI auto-tagging
Write-Host "🚀 Pushing amended commit to main..." -ForegroundColor Yellow
Write-Host "⚠️  Force push required due to commit amendment" -ForegroundColor Yellow

# Force push the amended commit to main
git push --force-with-lease origin main

Write-Host "✅ Changes pushed to main" -ForegroundColor Green
Write-Host "🤖 CI will automatically create release tag after tests pass" -ForegroundColor Cyan

# Wait a moment for CI to start
Start-Sleep -Seconds 5

# Open GitHub Actions to monitor progress
Write-Host "🌐 Opening GitHub Actions to monitor progress..." -ForegroundColor Yellow
Start-Process "https://github.com/ohishi-yhonda-org/print_pdf/actions"

Write-Host "" 
Write-Host "🎉 Release process initiated!" -ForegroundColor Green
Write-Host "📦 Release will be created automatically by CI" -ForegroundColor Cyan
Write-Host "" 
Write-Host "💡 Future usage:" -ForegroundColor Gray
Write-Host "   Enter = push only, y = full release, n = cancel" -ForegroundColor Gray

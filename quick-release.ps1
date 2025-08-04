# Quick Release Script
# Usage: .\quick-release.ps1 [version]
# Example: .\quick-release.ps1 v1.0.15

param(
    [string]$Version = ""
)

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
        git commit -m $commitMessage
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
$latestCommitMessage = git log -1 --pretty=format:"%s"
if ($latestCommitMessage -match "\[release\]") {
    Write-Host "🏷️  Latest commit already has [release] flag: $latestCommitMessage" -ForegroundColor Cyan
    Write-Host "📤 Proceeding to push for CI auto-tagging" -ForegroundColor Yellow
} else {
    # Step 1: Add release flag to trigger CI auto-tagging
    Write-Host "📝 Adding release flag to trigger CI..." -ForegroundColor Yellow
    try {
        # 変更がある場合はコミット、ない場合は空コミットで [release] フラグを追加
        $hasChanges = (git status --porcelain) -ne $null
        
        if ($hasChanges) {
            git add .
            $commitMessage = Read-Host "Enter commit message (or press Enter for default)"
            if ([string]::IsNullOrWhiteSpace($commitMessage)) {
                $commitMessage = "feat: release update $(Get-Date -Format 'yyyy-MM-dd HH:mm') [release]"
            } else {
                $commitMessage = "$commitMessage [release]"
            }
            git commit -m $commitMessage
            Write-Host "✅ Changes committed with release flag" -ForegroundColor Green
        } else {
            # 変更がない場合は空コミットでリリーストリガー
            git commit --allow-empty -m "trigger: release $(Get-Date -Format 'yyyy-MM-dd HH:mm') [release]"
            Write-Host "✅ Empty commit created with release flag" -ForegroundColor Green
        }
    } catch {
        Write-Host "⚠️  Failed to create release commit" -ForegroundColor Yellow
        exit 1
    }
}

# Step 2: Push changes and trigger CI auto-tagging
Write-Host "🚀 Pushing changes to main..." -ForegroundColor Yellow

# Push changes to main to trigger auto-tagging
git push origin main

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
Write-Host "   Enter = push only | y = full release | n = cancel" -ForegroundColor Gray

# Quick Release Script
# Usage: .\quick-release.ps1 [version]
# Example: .\quick-release.ps1 v1.0.15

param(
    [string]$Version = ""
)

Write-Host "🚀 Quick Release Script Starting..." -ForegroundColor Green

# Step 1: Commit current changes
Write-Host "📝 Committing current changes..." -ForegroundColor Yellow
try {
    git add .
    $commitMessage = Read-Host "Enter commit message (or press Enter for default)"
    if ([string]::IsNullOrWhiteSpace($commitMessage)) {
        $commitMessage = "feat: release update $(Get-Date -Format 'yyyy-MM-dd HH:mm')"
    }
    git commit -m $commitMessage
    Write-Host "✅ Changes committed" -ForegroundColor Green
} catch {
    Write-Host "⚠️  No changes to commit or commit failed" -ForegroundColor Yellow
}

# Step 2: Push changes
Write-Host "📤 Pushing changes to main..." -ForegroundColor Yellow
git push origin main

# Step 3: Trigger release
if ([string]::IsNullOrWhiteSpace($Version)) {
    Write-Host "🔄 Triggering automatic release..." -ForegroundColor Yellow
    Write-Host "GitHub Actions will auto-increment version and create release" -ForegroundColor Cyan
} else {
    Write-Host "🔄 Triggering manual release with version: $Version" -ForegroundColor Yellow
    
    # Create and push tag for specific version
    git tag $Version
    git push origin $Version
    
    Write-Host "✅ Tag $Version created and pushed" -ForegroundColor Green
}

# Step 4: Open GitHub Actions page
Write-Host "🌐 Opening GitHub Actions page..." -ForegroundColor Yellow
Start-Process "https://github.com/ohishi-yhonda-org/print_pdf/actions"

Write-Host ""
Write-Host "🎉 Release process initiated!" -ForegroundColor Green
Write-Host "📊 Check the Actions tab to monitor progress" -ForegroundColor Cyan
Write-Host "📦 Release will be available in 2-5 minutes" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Wait for CI to complete" -ForegroundColor White
Write-Host "  2. Check GitHub Releases page" -ForegroundColor White
Write-Host "  3. Download and test the new release" -ForegroundColor White

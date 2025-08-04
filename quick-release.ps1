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

# Step 3: Generate and push release tag
if ([string]::IsNullOrWhiteSpace($Version)) {
    Write-Host "🔄 Generating next release version..." -ForegroundColor Yellow
    
    # Get latest tag and increment patch version
    try {
        $gitOutput = & git tag --sort=-version:refname 2>$null
        $latestTag = $gitOutput | Where-Object { $_ -match "^v\d+\.\d+\.\d+$" } | Select-Object -First 1
        
        if ($latestTag) {
            Write-Host "Latest release tag found: $latestTag" -ForegroundColor Cyan
        } else {
            Write-Host "No release tags found, checking all tags..." -ForegroundColor Yellow
            $allTags = & git tag --list 2>$null
            Write-Host "All tags: $($allTags -join ', ')" -ForegroundColor Gray
        }
    } catch {
        Write-Host "Git command failed, using fallback" -ForegroundColor Yellow
        $latestTag = $null
    }
    
    if ($latestTag -and $latestTag -match "v(\d+)\.(\d+)\.(\d+)") {
        $major = [int]$matches[1]
        $minor = [int]$matches[2]
        $patch = [int]$matches[3] + 1
        $Version = "v$major.$minor.$patch"
        Write-Host "Generated version: $Version" -ForegroundColor Green
    } else {
        # Fallback version if no tags found - use v1.0.14 based on existing tags
        $Version = "v1.0.14"
        Write-Host "Using fallback version: $Version" -ForegroundColor Yellow
    }
} else {
    Write-Host "Using specified version: $Version" -ForegroundColor Cyan
}

Write-Host "🚀 Creating release tag: $Version" -ForegroundColor Yellow

# Create and push tag for release (non-dev version)
git tag $Version
git push origin $Version

Write-Host "✅ Release tag $Version created and pushed" -ForegroundColor Green

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

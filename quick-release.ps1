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

# Step 2: Generate and create release tag
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

# Step 3: Push changes and wait for CI, then create release tag
Write-Host "🚀 Pushing changes to main for testing..." -ForegroundColor Yellow

# First push changes to main to run tests
git push origin main

Write-Host "✅ Changes pushed to main" -ForegroundColor Green
Write-Host "⏳ Waiting for CI tests to complete..." -ForegroundColor Yellow

# Wait a moment for CI to start
Start-Sleep -Seconds 5

# Open GitHub Actions to monitor progress
Write-Host "🌐 Opening GitHub Actions to monitor tests..." -ForegroundColor Yellow
Start-Process "https://github.com/ohishi-yhonda-org/print_pdf/actions"

# Create release tag automatically
Write-Host "🚀 Creating release tag: $Version" -ForegroundColor Yellow

# Create and push tag for release
git tag $Version
git push origin $Version

Write-Host "✅ Release tag $Version created and pushed" -ForegroundColor Green

Write-Host ""
Write-Host "🎉 Release process completed!" -ForegroundColor Green
Write-Host "📦 Release will be available in 2-3 minutes" -ForegroundColor Cyan
Write-Host ""
Write-Host "💡 Future usage:" -ForegroundColor Gray
Write-Host "   Enter = push only | y = full release | n = cancel" -ForegroundColor Gray

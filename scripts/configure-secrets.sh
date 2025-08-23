#!/bin/bash

# Script to add Docker Hub secrets to both Doppler and GitHub

if [ $# -lt 2 ]; then
    echo "Usage: $0 <dockerhub_username> <dockerhub_token> [doppler_token] [github_repo]"
    echo ""
    echo "Example: $0 johndoe abc123xyz456 dp.pt.123456 owner/repo"
    exit 1
fi

DOCKERHUB_USERNAME=$1
DOCKERHUB_TOKEN=$2
DOPPLER_TOKEN=${3:-""}
GITHUB_REPO=${4:-""}

# Add secrets to Doppler
echo "Adding Docker Hub secrets to Doppler..."
if [ -z "$DOPPLER_TOKEN" ]; then
    # Check if doppler CLI is logged in
    doppler status > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "Doppler CLI not logged in. Please login first:"
        echo "doppler login"
        exit 1
    fi

    # Set secrets using current logged-in credentials
    doppler secrets set DOCKERHUB_USERNAME="$DOCKERHUB_USERNAME"
    doppler secrets set DOCKERHUB_TOKEN="$DOCKERHUB_TOKEN"
else
    # Use provided token for CI/CD
    DOPPLER_TOKEN="$DOPPLER_TOKEN" doppler secrets set DOCKERHUB_USERNAME="$DOCKERHUB_USERNAME" --project devlogai --config development
    DOPPLER_TOKEN="$DOPPLER_TOKEN" doppler secrets set DOCKERHUB_TOKEN="$DOCKERHUB_TOKEN" --project devlogai --config development
fi

echo "✅ Doppler secrets added successfully!"

# Add secrets to GitHub (if github CLI is available and repo is specified)
if command -v gh &> /dev/null && [ -n "$GITHUB_REPO" ]; then
    echo "Adding Docker Hub secrets to GitHub repository..."
    
    # Check if user is logged in to GitHub CLI
    gh auth status > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "GitHub CLI not logged in. Please login first:"
        echo "gh auth login"
        exit 1
    fi

    # Set GitHub secrets
    echo "$DOCKERHUB_USERNAME" | gh secret set DOCKERHUB_USERNAME --repo "$GITHUB_REPO"
    echo "$DOCKERHUB_TOKEN" | gh secret set DOCKERHUB_TOKEN --repo "$GITHUB_REPO"
    echo "✅ GitHub secrets added successfully!"
else
    echo "⚠️ GitHub CLI not available or repository not specified."
    echo "To add secrets to GitHub manually:"
    echo "1. Go to your GitHub repository"
    echo "2. Navigate to Settings > Secrets and variables > Actions"
    echo "3. Add the following secrets:"
    echo "   - DOCKERHUB_USERNAME: $DOCKERHUB_USERNAME"
    echo "   - DOCKERHUB_TOKEN: (your token)"
fi

echo ""
echo "✅ All done! Docker Hub credentials have been configured."

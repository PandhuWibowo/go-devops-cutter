#!/bin/bash
#
# Install git hooks for go-devops-cutter
# This script copies the hooks from scripts/hooks/ to .git/hooks/
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Installing git hooks for go-devops-cutter...${NC}"
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    echo "Please run this script from the project root directory"
    exit 1
fi

# Check if hooks directory exists
if [ ! -d "scripts/hooks" ]; then
    echo -e "${RED}Error: scripts/hooks directory not found${NC}"
    exit 1
fi

# Create .git/hooks directory if it doesn't exist
mkdir -p .git/hooks

# Install pre-commit hook
if [ -f "scripts/hooks/pre-commit" ]; then
    echo -e "${BLUE}Installing pre-commit hook...${NC}"
    cp scripts/hooks/pre-commit .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    echo -e "${GREEN}✓ Pre-commit hook installed${NC}"
else
    echo -e "${YELLOW}⚠️  Pre-commit hook not found${NC}"
fi

# Install pre-push hook
if [ -f "scripts/hooks/pre-push" ]; then
    echo -e "${BLUE}Installing pre-push hook...${NC}"
    cp scripts/hooks/pre-push .git/hooks/pre-push
    chmod +x .git/hooks/pre-push
    echo -e "${GREEN}✓ Pre-push hook installed${NC}"
else
    echo -e "${YELLOW}⚠️  Pre-push hook not found${NC}"
fi

echo ""
echo -e "${GREEN}✅ Git hooks installed successfully!${NC}"
echo ""
echo -e "${BLUE}Hooks installed:${NC}"
echo "  • pre-commit: Runs tests, formatting, and linting before commits"
echo "  • pre-push: Runs tests with coverage before pushing to remote"
echo ""
echo -e "${YELLOW}To skip hooks temporarily:${NC}"
echo "  • Skip pre-commit: git commit --no-verify"
echo "  • Skip pre-push: git push --no-verify"
echo ""
echo -e "${YELLOW}To uninstall hooks:${NC}"
echo "  • Run: rm .git/hooks/pre-commit .git/hooks/pre-push"
echo ""

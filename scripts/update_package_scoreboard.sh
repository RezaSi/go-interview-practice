#!/bin/bash

# Script to manually update the main package scoreboard in README.md
# This aggregates data from all package challenge scoreboards

set -e

echo "🚀 Go Interview Practice - Package Scoreboard Updater"
echo "===================================================="
echo ""

# Check if we're in the right directory
if [ ! -f "README.md" ]; then
    echo "❌ Error: README.md not found. Please run this script from the repository root."
    exit 1
fi

# Check if Python 3 is available
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python 3 is required but not installed."
    exit 1
fi

# Check if the script directory exists
if [ ! -d "scripts" ]; then
    echo "❌ Error: scripts directory not found."
    exit 1
fi

# Check if the generator script exists
if [ ! -f "scripts/generate_package_scoreboard.py" ]; then
    echo "❌ Error: Package scoreboard generator script not found."
    exit 1
fi

# Check if packages directory exists
if [ ! -d "packages" ]; then
    echo "❌ Error: packages directory not found."
    exit 1
fi

echo "📊 Analyzing package challenge scoreboards..."
echo ""

# Count package challenges with submissions
package_challenges_with_submissions=$(find packages -name "SCOREBOARD.md" -exec grep -l -v "^#\|^|\s*Username\|^|\s*---" {} \; 2>/dev/null | wc -l | xargs)
total_package_challenges=$(find packages -name "SCOREBOARD.md" | wc -l | xargs)

# Count packages
total_packages=$(find packages -maxdepth 1 -type d | tail -n +2 | wc -l | xargs)

echo "📈 Found:"
echo "  - Total packages: $total_packages"
echo "  - Total package challenges: $total_package_challenges"
echo "  - Package challenges with submissions: $package_challenges_with_submissions"
echo ""

# List available packages
echo "📦 Available packages:"
for package_dir in packages/*/; do
    if [ -d "$package_dir" ]; then
        package_name=$(basename "$package_dir")
        challenge_count=$(find "$package_dir" -maxdepth 1 -name "challenge-*" -type d | wc -l | xargs)
        echo "  - $package_name ($challenge_count challenges)"
    fi
done
echo ""

# Run the package scoreboard generator
echo "🔄 Generating main package scoreboard..."
python3 scripts/generate_package_scoreboard.py

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Main package scoreboard updated successfully!"
    echo ""
    echo "📝 The package leaderboard in README.md has been updated with:"
    echo "  - Latest package completion statistics"
    echo "  - Updated package achievement badges"
    echo "  - Current package completion rates"
    echo "  - Top 10 package developer rankings"
    echo "  - Per-package progress breakdown"
    echo ""
    echo "🎯 Next steps:"
    echo "  1. Review the updated README.md"
    echo "  2. Commit and push changes if satisfied"
    echo "  3. The package scoreboard will auto-update on future submissions"
    echo ""
else
    echo ""
    echo "❌ Failed to update main package scoreboard!"
    echo "Please check the error messages above and try again."
    exit 1
fi 
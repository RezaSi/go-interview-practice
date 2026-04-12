# Scoreboard Generation Scripts

This directory contains Python scripts to automatically generate and update scoreboards in the main README.md file.

## Scripts Overview

### Individual Scripts

1. **`generate_main_scoreboard.py`** - Updates the classic challenges leaderboard
   - Scans all `challenge-*` directories for SCOREBOARD.md files  
   - Aggregates completion data from all classic challenges
   - Updates the "🏆 Top 10 Leaderboard" section in README.md

2. **`generate_package_scoreboard.py`** - Updates the package challenges leaderboard
   - Scans all package directories in `packages/` for challenge scoreboards
   - Aggregates completion data from package challenges
   - Updates the "🚀 Package Challenges Leaderboard" section in README.md

### Convenience Scripts

3. **`update_all_scoreboards.py`** - Runs both scoreboard generators
   - Executes both classic and package scoreboard scripts in sequence
   - Provides a summary of success/failure status
   - Recommended for updating all scoreboards at once

4. **`test_scoreboards.py`** - Tests the scoreboard scripts
   - Validates that both scripts work together without conflicts
   - Tests running scripts in different orders
   - Checks for proper README.md marker placement
   - Includes backup/restore functionality for safe testing

## Usage

### Quick Update (Recommended)
```bash
# Update both scoreboards at once
python3 scripts/update_all_scoreboards.py
```

### Individual Updates
```bash
# Update only classic challenges leaderboard
python3 scripts/generate_main_scoreboard.py

# Update only package challenges leaderboard  
python3 scripts/generate_package_scoreboard.py
```

### Testing
```bash
# Test that scripts work together properly
python3 scripts/test_scoreboards.py
```

## Script Features

### ✅ **Robust & Independent**
- Scripts can run in **any order** without conflicts
- Each script updates only its specific README section
- **Path-aware**: Works when run from any directory (root or scripts/)
- **Error handling**: Graceful handling of missing files/directories

### ✅ **Safe & Non-Destructive**
- Uses **unique markers** to identify sections:
  - `<!-- END_CLASSIC_LEADERBOARD -->` for classic challenges
  - `<!-- END_PACKAGE_LEADERBOARD -->` for package challenges
- **Preserves existing content** outside of managed sections
- **Idempotent**: Running multiple times produces same result

### ✅ **Automatic Discovery**
- **Classic challenges**: Automatically finds all `challenge-*` directories
- **Package challenges**: Automatically scans `packages/` directory
- **No hardcoded paths** or challenge lists to maintain

## README.md Structure

The scripts maintain this structure in README.md:

```markdown
## 🏆 Top 10 Leaderboard
[Classic challenges leaderboard content]
<!-- END_CLASSIC_LEADERBOARD -->

## 🚀 Package Challenges Leaderboard  
[Package challenges leaderboard content]
<!-- END_PACKAGE_LEADERBOARD -->

## Key Features
[Rest of README content...]
```

## Requirements

- **Python 3.x** 
- **Standard library only** (no external dependencies)
- **SCOREBOARD.md files** in challenge directories with proper format:
  ```
  | Username | Passed Tests | Total Tests | ... |
  |----------|--------------|-------------|-----|
  | user1    | 6           | 6           | ... |
  ```

## How It Works

1. **Data Collection**: Scripts scan challenge directories for SCOREBOARD.md files
2. **Parsing**: Extract usernames and test results using regex parsing
3. **Aggregation**: Count completed challenges per user (only 100% completion counts)
4. **Sorting**: Sort users by completion count, then alphabetically
5. **Formatting**: Generate beautiful GitHub-compatible HTML/Markdown tables
6. **Update**: Replace specific sections in README.md using markers

## Automation

These scripts are designed to be run by:
- **GitHub Actions** (automatic updates when scoreboards change)
- **Local development** (manual updates during testing)
- **CI/CD pipelines** (scheduled or triggered updates)

## Contributing

When adding new challenges:
1. **Classic challenges**: Just create the challenge directory - automatically detected
2. **Package challenges**: Add to package's `learning_path` in `package.json`
3. **Scoreboards**: Follow existing SCOREBOARD.md format
4. **No code changes needed** - scripts will find and process new challenges automatically

---

💡 **Tip**: The `update_all_scoreboards.py` script is the easiest way to keep all scoreboards current! 
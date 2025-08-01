#!/usr/bin/env python3
"""
Wrapper script to update both classic and package scoreboards and generate contributor badges.
This script ensures both scoreboards are updated and profile badges are generated.
"""

import os
import sys
import subprocess
from pathlib import Path


def run_script(script_name, working_dir):
    """Run a scoreboard script and return success status."""
    script_path = Path(__file__).parent / script_name
    
    print(f"\n{'='*60}")
    print(f"🔄 Running {script_name}")
    print(f"{'='*60}")
    
    try:
        result = subprocess.run(
            [sys.executable, str(script_path)], 
            cwd=working_dir,  # Run from the root directory
            text=True
        )
        
        if result.returncode == 0:
            print(f"✅ {script_name} completed successfully!")
            return True
        else:
            print(f"❌ {script_name} failed with exit code {result.returncode}")
            return False
        
    except Exception as e:
        print(f"❌ Error running {script_name}: {e}")
        return False


def main():
    """Main function to update all scoreboards."""
    print("🚀 Updating All Scoreboards")
    print("="*60)
    print("This script will update both classic and package challenge scoreboards.")
    
    # Get the root directory (parent of scripts directory)
    root_dir = Path(__file__).parent.parent
    
    print(f"Working directory: {root_dir}")
    
    scripts = [
        "generate_main_scoreboard.py",
        "generate_package_scoreboard.py",
        "generate_contributor_badges.py"
    ]
    
    success_count = 0
    total_scripts = len(scripts)
    
    for script in scripts:
        if run_script(script, root_dir):
            success_count += 1
    
    print(f"\n{'='*60}")
    print(f"📊 Summary")
    print(f"{'='*60}")
    print(f"Scripts run: {total_scripts}")
    print(f"Successful: {success_count}")
    print(f"Failed: {total_scripts - success_count}")
    
    if success_count == total_scripts:
        print(f"\n🎉 All scoreboards and badges updated successfully!")
        print(f"✅ Classic Challenges Leaderboard updated")
        print(f"✅ Package Challenges Leaderboard updated")
        print(f"✅ Contributor Profile Badges generated")
        print(f"\nREADME.md has been updated with the latest scoreboard data.")
        print(f"🏆 Profile badges are available in the badges/ directory.")
        return 0
    else:
        print(f"\n⚠️  Some scripts failed. Please check the errors above.")
        return 1


if __name__ == "__main__":
    sys.exit(main()) 
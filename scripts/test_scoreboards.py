#!/usr/bin/env python3
"""
Test script to verify that both scoreboard generation scripts work correctly together.
"""

import os
import sys
import subprocess
from pathlib import Path


def run_script(script_name):
    """Run a script and return success status and output."""
    script_path = Path(__file__).parent / script_name
    
    try:
        result = subprocess.run(
            [sys.executable, str(script_path)], 
            capture_output=True, 
            text=True,
            cwd=script_path.parent
        )
        
        print(f"\n{'='*50}")
        print(f"Running {script_name}")
        print(f"{'='*50}")
        print("STDOUT:")
        print(result.stdout)
        if result.stderr:
            print("STDERR:")
            print(result.stderr)
        print(f"Exit code: {result.returncode}")
        
        return result.returncode == 0, result.stdout, result.stderr
        
    except Exception as e:
        print(f"Error running {script_name}: {e}")
        return False, "", str(e)


def check_readme_markers():
    """Check if README.md has the correct markers after updates."""
    
    # Get the root directory (handle running from different locations)
    current_dir = Path('.')
    if 'scripts' in str(current_dir.absolute()):
        current_dir = current_dir.parent
    
    readme_path = current_dir / 'README.md'
    
    try:
        with open(readme_path, 'r') as f:
            content = f.read()
        
        # Check for markers
        classic_start = "## 🏆 Top 10 Leaderboard" in content
        classic_end = "<!-- END_CLASSIC_LEADERBOARD -->" in content
        package_start = "## 🚀 Package Challenges Leaderboard" in content
        package_end = "<!-- END_PACKAGE_LEADERBOARD -->" in content
        
        print(f"\n{'='*50}")
        print("README.md Marker Check")
        print(f"{'='*50}")
        print(f"Classic Leaderboard Start: {'✅' if classic_start else '❌'}")
        print(f"Classic Leaderboard End: {'✅' if classic_end else '❌'}")
        print(f"Package Leaderboard Start: {'✅' if package_start else '❌'}")
        print(f"Package Leaderboard End: {'✅' if package_end else '❌'}")
        
        # Check section order
        if classic_start and package_start:
            classic_pos = content.find("## 🏆 Top 10 Leaderboard")
            package_pos = content.find("## 🚀 Package Challenges Leaderboard")
            
            if classic_pos < package_pos:
                print("Section Order: ✅ Classic before Package")
            else:
                print("Section Order: ❌ Package before Classic")
        
        return all([classic_start, classic_end, package_start, package_end])
        
    except Exception as e:
        print(f"Error reading README.md: {e}")
        return False


def backup_readme():
    """Create a backup of README.md before testing."""
    
    current_dir = Path('.')
    if 'scripts' in str(current_dir.absolute()):
        current_dir = current_dir.parent
    
    readme_path = current_dir / 'README.md'
    backup_path = current_dir / 'README.md.backup'
    
    try:
        with open(readme_path, 'r') as source:
            with open(backup_path, 'w') as backup:
                backup.write(source.read())
        print("✅ README.md backed up")
        return True
    except Exception as e:
        print(f"❌ Failed to backup README.md: {e}")
        return False


def restore_readme():
    """Restore README.md from backup."""
    
    current_dir = Path('.')
    if 'scripts' in str(current_dir.absolute()):
        current_dir = current_dir.parent
    
    readme_path = current_dir / 'README.md'
    backup_path = current_dir / 'README.md.backup'
    
    try:
        if backup_path.exists():
            with open(backup_path, 'r') as backup:
                with open(readme_path, 'w') as target:
                    target.write(backup.read())
            backup_path.unlink()  # Remove backup file
            print("✅ README.md restored from backup")
        return True
    except Exception as e:
        print(f"❌ Failed to restore README.md: {e}")
        return False


def main():
    """Main test function."""
    print("🧪 Testing Scoreboard Scripts")
    print("="*60)
    
    # Backup README.md
    if not backup_readme():
        return 1
    
    try:
        # Test both orders
        test_scenarios = [
            ("Classic First", ["generate_main_scoreboard.py", "generate_package_scoreboard.py"]),
            ("Package First", ["generate_package_scoreboard.py", "generate_main_scoreboard.py"])
        ]
        
        for scenario_name, scripts in test_scenarios:
            print(f"\n🔬 Testing Scenario: {scenario_name}")
            print("="*60)
            
            # Restore README.md to original state for each test
            restore_readme()
            backup_readme()
            
            all_success = True
            
            for script in scripts:
                success, stdout, stderr = run_script(script)
                if not success:
                    print(f"❌ {script} failed!")
                    all_success = False
                else:
                    print(f"✅ {script} succeeded!")
            
            # Check markers after running both scripts
            markers_ok = check_readme_markers()
            
            print(f"\n📊 Scenario Result: {'✅ PASSED' if all_success and markers_ok else '❌ FAILED'}")
        
        # Final test: Run both scripts multiple times to test idempotency
        print(f"\n🔄 Testing Idempotency (running each script twice)")
        print("="*60)
        
        for script in ["generate_main_scoreboard.py", "generate_package_scoreboard.py"]:
            print(f"\nTesting {script} idempotency...")
            
            # Run first time
            success1, _, _ = run_script(script)
            
            # Get content after first run
            current_dir = Path('.')
            if 'scripts' in str(current_dir.absolute()):
                current_dir = current_dir.parent
            
            with open(current_dir / 'README.md', 'r') as f:
                content1 = f.read()
            
            # Run second time
            success2, _, _ = run_script(script)
            
            # Get content after second run
            with open(current_dir / 'README.md', 'r') as f:
                content2 = f.read()
            
            # Check if content is identical
            if success1 and success2 and content1 == content2:
                print(f"✅ {script} is idempotent")
            else:
                print(f"❌ {script} is not idempotent")
        
        print(f"\n🎉 All tests completed!")
        
        # Ask user if they want to keep the updated README
        keep_changes = input("\nDo you want to keep the updated README.md? (y/n): ").lower().strip()
        
        if keep_changes != 'y':
            restore_readme()
            print("✅ README.md restored to original state")
        else:
            # Remove backup
            current_dir = Path('.')
            if 'scripts' in str(current_dir.absolute()):
                current_dir = current_dir.parent
            backup_path = current_dir / 'README.md.backup'
            if backup_path.exists():
                backup_path.unlink()
            print("✅ Updated README.md kept")
        
        return 0
        
    except KeyboardInterrupt:
        print("\n⚠️ Test interrupted by user")
        restore_readme()
        return 1
    except Exception as e:
        print(f"\n❌ Test failed with error: {e}")
        restore_readme()
        return 1


if __name__ == "__main__":
    sys.exit(main()) 
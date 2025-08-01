#!/usr/bin/env python3
"""
Script to generate the package scoreboard for README.md by aggregating data from all package challenge scoreboards.
"""

import os
import re
import sys
from collections import defaultdict
from pathlib import Path


def parse_package_scoreboard_file(filepath):
    """Parse a package challenge SCOREBOARD.md file and extract usernames who completed the challenge (passed ALL tests)."""
    users = set()
    
    try:
        with open(filepath, 'r') as f:
            content = f.read()
        
        # Split into lines and process
        lines = content.strip().split('\n')
        
        for line in lines:
            # Skip header and separator lines
            if not line.strip() or 'Username' in line or '---' in line or line.startswith('#'):
                continue
            
            # Parse table row
            if '|' in line:
                parts = [part.strip() for part in line.split('|')]
                if len(parts) >= 4:  # Username | Passed Tests | Total Tests | (optional extra columns)
                    username = parts[1]
                    passed_tests_str = parts[2]
                    total_tests_str = parts[3]
                    
                    # Skip empty usernames or placeholders
                    if not username or username == '------' or username.isdigit():
                        continue
                    
                    try:
                        # Extract numbers from test counts (handles formats like "6", "6 tests", etc.)
                        passed_tests = int(''.join(filter(str.isdigit, passed_tests_str)))
                        total_tests = int(''.join(filter(str.isdigit, total_tests_str)))
                        
                        # Only count as completed if ALL tests passed
                        if passed_tests > 0 and passed_tests == total_tests:
                            users.add(username)
                            print(f"  ✅ {username}: {passed_tests}/{total_tests} tests passed (COMPLETED)")
                        else:
                            print(f"  ❌ {username}: {passed_tests}/{total_tests} tests passed (incomplete)")
                            
                    except (ValueError, TypeError):
                        # If we can't parse test numbers, skip this entry
                        print(f"  ⚠️  {username}: Could not parse test results")
                        continue
    
    except FileNotFoundError:
        pass
    except Exception as e:
        print(f"Error parsing {filepath}: {e}", file=sys.stderr)
    
    return users


def get_package_challenge_title(challenge_dir):
    """Extract challenge title from README.md or directory name."""
    readme_path = os.path.join(challenge_dir, 'README.md')
    
    try:
        with open(readme_path, 'r') as f:
            content = f.read()
        
        # Look for title in various formats
        title_patterns = [
            r'^#\s+(.+?)$',  # # Title
            r'^\*\*(.+?)\*\*',  # **Title**
            r'Challenge \d+:\s*(.+?)$'  # Challenge N: Title
        ]
        
        lines = content.split('\n')
        for line in lines:
            line = line.strip()
            for pattern in title_patterns:
                match = re.search(pattern, line, re.MULTILINE)
                if match:
                    title = match.group(1).strip()
                    # Clean up title
                    title = re.sub(r'Challenge \d+:\s*', '', title)
                    return title
    
    except FileNotFoundError:
        pass
    
    # Fallback to directory name
    return challenge_dir.replace('challenge-', 'Challenge ')


def generate_package_scoreboard():
    """Generate the package scoreboard by aggregating all package challenge scoreboards."""
    
    # Dictionary to track user completion counts per package
    package_completions = defaultdict(lambda: defaultdict(lambda: {'count': 0, 'challenges': []}))
    overall_completions = defaultdict(lambda: {'count': 0, 'packages': set()})
    
    # Get the root directory (handle running from different locations)
    current_dir = Path('.')
    if 'scripts' in str(current_dir.absolute()):
        # If running from scripts directory, go up one level
        current_dir = current_dir.parent
    
    # Find packages directory
    packages_dir = current_dir / 'packages'
    
    if not packages_dir.exists():
        print("❌ No packages directory found!")
        return ""
    
    package_dirs = sorted([d for d in packages_dir.iterdir() if d.is_dir()])
    
    print(f"Found {len(package_dirs)} package directories")
    
    # Process each package
    for package_dir in package_dirs:
        package_name = package_dir.name
        print(f"\n🔍 Processing package: {package_name}")
        
        # Find challenge directories within this package
        challenge_dirs = sorted([d for d in package_dir.iterdir() 
                               if d.is_dir() and d.name.startswith('challenge-')])
        
        for challenge_dir in challenge_dirs:
            challenge_id = challenge_dir.name
            scoreboard_path = challenge_dir / 'SCOREBOARD.md'
            
            if scoreboard_path.exists():
                users = parse_package_scoreboard_file(scoreboard_path)
                challenge_title = get_package_challenge_title(str(challenge_dir))
                
                print(f"  {challenge_id}: {len(users)} users completed")
                
                for user in users:
                    package_completions[package_name][user]['count'] += 1
                    package_completions[package_name][user]['challenges'].append({
                        'id': challenge_id,
                        'title': challenge_title
                    })
                    overall_completions[user]['count'] += 1
                    overall_completions[user]['packages'].add(package_name)
    
    # Generate the markdown scoreboard
    markdown_lines = [
        "## 🚀 Package Challenges Leaderboard",
        "",
        "Master Go packages through hands-on challenges! Each package offers a structured learning path with real-world scenarios.",
        "",
        "> **Note**: The data below is automatically updated by GitHub Actions when package challenge scoreboards change.",
        "",
    ]
    
    # Generate overall leaderboard across all packages
    sorted_overall = sorted(overall_completions.items(), 
                           key=lambda x: (-x[1]['count'], x[0]))
    
    if sorted_overall:
        # Generate top 10 overall package challenge leaders
        html_table = generate_package_html_leaderboard(sorted_overall[:10], package_completions)
        markdown_lines.append(html_table)
    else:
        markdown_lines.extend([
            "No completed package challenges yet. Be the first to solve a package challenge!",
            ""
        ])
    
    # Add per-package breakdown
    markdown_lines.extend([
        "",
        "### 📦 Per-Package Progress",
        ""
    ])
    
    for package_name in sorted(package_completions.keys()):
        package_users = package_completions[package_name]
        sorted_package_users = sorted(package_users.items(), 
                                    key=lambda x: (-x[1]['count'], x[0]))
        
        if sorted_package_users:
            markdown_lines.extend([
                f"#### {package_name.title()} Package",
                "",
                "| Rank | Developer | Completed | Progress |",
                "|:---:|:---:|:---:|:---|"
            ])
            
            # Find total challenges for this package
            package_path = packages_dir / package_name
            total_challenges = len([d for d in package_path.iterdir() 
                                  if d.is_dir() and d.name.startswith('challenge-')])
            
            for rank, (username, data) in enumerate(sorted_package_users[:5], 1):  # Top 5 per package
                count = data['count']
                progress_bar = generate_progress_bar(count, total_challenges)
                
                rank_emoji = "🥇" if rank == 1 else "🥈" if rank == 2 else "🥉" if rank == 3 else f"{rank}"
                
                markdown_lines.append(
                    f"| {rank_emoji} | **[{username}](https://github.com/{username})** | {count}/{total_challenges} | {progress_bar} |"
                )
            
            markdown_lines.append("")
    
    # Add summary information
    total_package_challenges = sum(len([d for d in package_dir.iterdir() 
                                      if d.is_dir() and d.name.startswith('challenge-')]) 
                                 for package_dir in package_dirs)
    
    markdown_lines.extend([
        "### 📊 Package Challenge Statistics",
        "",
        f"- **Total Package Challenges Available**: {total_package_challenges}",
        f"- **Active Package Learners**: {len(overall_completions)}",
        f"- **Available Packages**: {len(package_dirs)} ({', '.join([d.name for d in package_dirs])})",
        "",
    ])
    
    if sorted_overall:
        top_user = sorted_overall[0]
        markdown_lines.append(f"- **Most Package Challenges Solved**: {top_user[1]['count']} by {top_user[0]}")
        markdown_lines.append("")
    
    markdown_lines.extend([
        "<!-- END_PACKAGE_LEADERBOARD -->",
        ""
    ])
    
    return '\n'.join(markdown_lines)


def generate_progress_bar(completed, total, length=10):
    """Generate a text-based progress bar."""
    if total == 0:
        return "⬜" * length
    
    progress = completed / total
    filled = int(progress * length)
    bar = "🟩" * filled + "⬜" * (length - filled)
    percentage = f"{progress * 100:.0f}%"
    return f"{bar} {percentage}"


def generate_package_html_leaderboard(top_users, package_completions):
    """Generate a beautiful GitHub-compatible package leaderboard table."""
    
    # Start with the table header - simple markdown format
    markdown_lines = [
        '| 🏅 | Developer | Total Solved | Packages | Achievement | Challenge Distribution |',
        '|:---:|:---:|:---:|:---:|:---:|:---|'
    ]
    
    for i, (username, data) in enumerate(top_users, 1):
        total_count = data['count']
        packages_completed = data['packages']
        package_count = len(packages_completed)
        
        # Determine achievement badge
        if total_count >= 15:
            achievement = "🔥 Package Master"
        elif total_count >= 10:
            achievement = "⭐ Package Expert"
        elif total_count >= 5:
            achievement = "💪 Package Advanced"
        elif total_count >= 3:
            achievement = "🚀 Package Intermediate"
        else:
            achievement = "🌱 Package Beginner"
        
        # Rank badge with medals for top 3
        if i == 1:
            rank_badge = "🥇"
        elif i == 2:
            rank_badge = "🥈"
        elif i == 3:
            rank_badge = "🥉"
        else:
            rank_badge = f"{i}"
        
        # Generate package completion breakdown
        package_breakdown = []
        for package_name in sorted(packages_completed):
            if package_name in package_completions:
                user_data = package_completions[package_name].get(username, {'count': 0})
                count = user_data['count']
                package_breakdown.append(f"**{package_name}**: {count}")
        
        breakdown_text = " • ".join(package_breakdown)
        
        # Create profile with GitHub avatar
        profile_cell = f'<img src="https://github.com/{username}.png" width="24" height="24" style="border-radius: 50%;"><br/>**[{username}](https://github.com/{username})**'
        
        # Create solved count
        solved_cell = f"**{total_count}**"
        
        # Create packages count
        packages_cell = f"**{package_count}** pkg{'s' if package_count != 1 else ''}"
        
        # Add row to table
        markdown_lines.append(
            f"| {rank_badge} | {profile_cell} | {solved_cell} | {packages_cell} | {achievement} | {breakdown_text} |"
        )
    
    # Add centered legend
    markdown_lines.extend([
        '',
        '<div align="center">',
        '',
        '🚀 **Package Challenges** - Learn Go packages through practical, real-world scenarios',
        '',
        '</div>'
    ])
    
    return '\n'.join(markdown_lines)


def update_readme_with_package_scoreboard(scoreboard_content):
    """Update README.md with the new package scoreboard content."""
    
    # Get the root directory (handle running from different locations)
    current_dir = Path('.')
    if 'scripts' in str(current_dir.absolute()):
        # If running from scripts directory, go up one level
        current_dir = current_dir.parent
    
    readme_path = current_dir / 'README.md'
    
    try:
        with open(readme_path, 'r') as f:
            content = f.read()
    except FileNotFoundError:
        print("README.md not found!", file=sys.stderr)
        return False
    
    # Define specific markers for the package scoreboard section
    start_marker = "## 🚀 Package Challenges Leaderboard"
    end_marker = "<!-- END_PACKAGE_LEADERBOARD -->"
    
    # Find the positions of markers
    start_pos = content.find(start_marker)
    end_pos = content.find(end_marker)
    
    if start_pos == -1:
        # If package scoreboard doesn't exist, insert after classic leaderboard or before key features
        classic_end = content.find("<!-- END_CLASSIC_LEADERBOARD -->")
        key_features_pos = content.find("## Key Features")
        
        if classic_end != -1:
            # Insert after classic leaderboard
            insertion_point = content.find('\n', classic_end) + 1
        elif key_features_pos != -1:
            # Fallback to before Key Features
            insertion_point = key_features_pos
        else:
            print("Could not find insertion point in README.md", file=sys.stderr)
            return False
        
        # Insert the new package scoreboard
        new_content = (content[:insertion_point] + 
                      scoreboard_content + '\n' + 
                      content[insertion_point:])
    else:
        if end_pos == -1:
            # If start marker exists but no end marker, find next section
            next_section_patterns = [
                "## Key Features",
                "## Getting Started",
                "## Challenge Categories"
            ]
            
            end_pos = len(content)  # Default to end of file
            for pattern in next_section_patterns:
                pattern_pos = content.find(pattern, start_pos + len(start_marker))
                if pattern_pos != -1:
                    end_pos = pattern_pos
                    break
        else:
            # Include the end marker in replacement
            end_pos = content.find('\n', end_pos) + 1
        
        # Replace existing package scoreboard section
        new_content = (content[:start_pos] + 
                      scoreboard_content + 
                      content[end_pos:])
    
    # Write the updated content
    try:
        with open(readme_path, 'w') as f:
            f.write(new_content)
        print("README.md updated successfully with package scoreboard!")
        return True
    except Exception as e:
        print(f"Error writing to README.md: {e}", file=sys.stderr)
        return False


def main():
    """Main function to generate and update the package scoreboard."""
    print("Generating package challenges scoreboard...")
    
    scoreboard_content = generate_package_scoreboard()
    
    if update_readme_with_package_scoreboard(scoreboard_content):
        print("Package scoreboard updated successfully!")
        return 0
    else:
        print("Failed to update package scoreboard!", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main()) 
#!/usr/bin/env python3
"""
Script to generate the main scoreboard for README.md by aggregating data from all challenge scoreboards.
"""

import os
import re
import sys
import requests
from collections import defaultdict
from pathlib import Path


def parse_scoreboard_file(filepath):
    """Parse a SCOREBOARD.md file and extract usernames who completed the challenge (passed ALL tests)."""
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


def load_sponsors():
    """Load sponsor list by scraping the public GitHub sponsors page."""
    sponsors = set()
    
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (compatible; PythonSponsorScraper/1.0)'
        }
        
        response = requests.get(
            'https://github.com/sponsors/RezaSi',
            headers=headers,
            timeout=10
        )
        
        if response.status_code == 200:
            html = response.text
            
            # Extract usernames from alt="@username" attributes
            avatar_pattern = r'alt="@([a-zA-Z0-9][a-zA-Z0-9\-]*)"'
            matches = re.findall(avatar_pattern, html)
            
            for username in matches:
                # Filter out the repository owner from sponsors list
                if username != "RezaSi":
                    sponsors.add(username)
    
    except Exception as e:
        pass  # Silently handle sponsor loading errors
    
    return sponsors


def get_challenge_title(challenge_dir):
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


def generate_main_scoreboard():
    """Generate the main scoreboard by aggregating all challenge scoreboards."""
    
    # Load sponsors from GitHub API
    sponsors = load_sponsors()
    
    # Dictionary to track user completion counts
    user_completions = defaultdict(lambda: {'count': 0, 'challenges': []})
    
    # Get the root directory (handle running from different locations)
    current_dir = Path('.')
    if 'scripts' in str(current_dir.absolute()):
        # If running from scripts directory, go up one level
        current_dir = current_dir.parent
    
    # Find all challenge directories
    challenge_dirs = sorted([d for d in current_dir.iterdir() 
                           if d.is_dir() and d.name.startswith('challenge-')])
    
    print(f"Found {len(challenge_dirs)} challenge directories")
    
    # Process each challenge
    for challenge_dir in challenge_dirs:
        challenge_num = challenge_dir.name.replace('challenge-', '')
        scoreboard_path = challenge_dir / 'SCOREBOARD.md'
        
        if scoreboard_path.exists():
            users = parse_scoreboard_file(scoreboard_path)
            challenge_title = get_challenge_title(str(challenge_dir))
            
            print(f"Challenge {challenge_num}: {len(users)} users completed")
            
            for user in users:
                user_completions[user]['count'] += 1
                user_completions[user]['challenges'].append({
                    'num': int(challenge_num),
                    'title': challenge_title
                })
    
    # Sort users by completion count (descending) and then by username
    sorted_users = sorted(user_completions.items(), 
                         key=lambda x: (-x[1]['count'], x[0]))
    
    # Generate the HTML leaderboard
    markdown_lines = [
        "## 🏆 Top 10 Leaderboard",
        "",
        "Our most accomplished Go developers, ranked by number of challenges completed:",
        "",
        "> **Note**: The data below is automatically updated by GitHub Actions when challenge scoreboards change.",
        "",
    ]
    
    total_challenges = len(challenge_dirs)
    
    if sorted_users:
        # Generate HTML table with styling
        html_table = generate_html_leaderboard(sorted_users[:10], total_challenges, challenge_dirs, sponsors)
        markdown_lines.append(html_table)
    else:
        markdown_lines.extend([
            "No completed challenges yet. Be the first to solve a challenge!",
            ""
        ])
    
    # Add summary information
    markdown_lines.extend([
        "",
        f"*Updated automatically based on {total_challenges} available challenges*",
        "",
        "### Challenge Progress Overview",
        "",
        f"- **Total Challenges Available**: {total_challenges}",
        f"- **Active Developers**: {len(user_completions)}",
        f"- **Most Challenges Solved**: {sorted_users[0][1]['count'] if sorted_users else 0} by {sorted_users[0][0] if sorted_users else 'N/A'}",
        "",
        "<!-- END_CLASSIC_LEADERBOARD -->",
        ""
    ])
    
    return '\n'.join(markdown_lines)


def generate_html_leaderboard(top_users, total_challenges, challenge_dirs, sponsors):
    """Generate a beautiful GitHub-compatible leaderboard table."""
    
    # Get list of challenge numbers for indicators
    challenge_numbers = sorted([int(d.name.replace('challenge-', '')) for d in challenge_dirs])
    
    # Start with the table header - simple markdown format
    markdown_lines = [
        '| 🏅 | Developer | Solved | Rate | Achievement | Progress |',
        '|:---:|:---:|:---:|:---:|:---:|:---|'
    ]
    
    for i, (username, data) in enumerate(top_users, 1):
        count = data['count']
        completion_rate = f"{(count / total_challenges * 100):.1f}%"
        
        # Determine achievement badge
        if count >= 20:
            achievement = "Master"
        elif count >= 15:
            achievement = "Expert"
        elif count >= 10:
            achievement = "Advanced"
        elif count >= 5:
            achievement = "Intermediate"
        else:
            achievement = "Beginner"
        
        # Rank badge with medals for top 3
        if i == 1:
            rank_badge = "🥇"
        elif i == 2:
            rank_badge = "🥈"
        elif i == 3:
            rank_badge = "🥉"
        else:
            rank_badge = f"{i}"
        
        # Generate challenge indicators - show all challenges in two rows
        completed_challenges = {ch['num'] for ch in data['challenges']}
        
        # Split challenges into two rows for better display
        first_half = challenge_numbers[:len(challenge_numbers)//2 + len(challenge_numbers)%2]
        second_half = challenge_numbers[len(challenge_numbers)//2 + len(challenge_numbers)%2:]

        # First row indicators
        first_row = ""
        for ch_num in first_half:
            if ch_num in completed_challenges:
                first_row += "✅"
            else:
                first_row += "⬜"
        
        # Second row indicators  
        second_row = ""
        for ch_num in second_half:
            if ch_num in completed_challenges:
                second_row += "✅"
            else:
                second_row += "⬜"
        
        # Combine both rows with line break
        indicators = f"{first_row}<br/>{second_row}"
        
        # Create simple profile with GitHub avatar - centered
        sponsor_badge = " ❤️" if username in sponsors else ""
        profile_cell = f'<img src="https://github.com/{username}.png" width="24" height="24" style="border-radius: 50%;"><br/>**[{username}](https://github.com/{username})**{sponsor_badge}'
        
        # Create solved count
        solved_cell = f"**{count}**/{total_challenges}"
        
        # Create completion rate
        rate_cell = f"**{completion_rate}**"
        
        # Add row to table
        markdown_lines.append(
            f"| {rank_badge} | {profile_cell} | {solved_cell} | {rate_cell} | {achievement} | {indicators} |"
        )
    
    # Add centered legend
    markdown_lines.extend([
        '',
        '<div align="center">',
        '',
        '✅ Completed • ⬜ Not Completed',
        '',
        f'*All {total_challenges} challenges shown in two rows*',
        '',
        '</div>'
    ])
    
    return '\n'.join(markdown_lines)


def update_readme_with_scoreboard(scoreboard_content):
    """Update README.md with the new scoreboard content."""
    
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
    
    # Define specific markers for the classic leaderboard section
    start_marker = "## 🏆 Top 10 Leaderboard"
    end_marker = "<!-- END_CLASSIC_LEADERBOARD -->"
    
    # Find the positions of markers
    start_pos = content.find(start_marker)
    end_pos = content.find(end_marker)
    
    if start_pos == -1:
        # If classic leaderboard doesn't exist, insert before package challenges or key features
        package_pos = content.find("## 🚀 Package Challenges Leaderboard")
        key_features_pos = content.find("## Key Features")
        
        if package_pos != -1:
            insertion_point = package_pos
        elif key_features_pos != -1:
            insertion_point = key_features_pos
        else:
            print("Could not find insertion point in README.md", file=sys.stderr)
            return False
        
        # Insert the new classic leaderboard
        new_content = (content[:insertion_point] + 
                      scoreboard_content + '\n' + 
                      content[insertion_point:])
    else:
        if end_pos == -1:
            # If start marker exists but no end marker, find next section
            next_section_patterns = [
                "## 🚀 Package Challenges Leaderboard",
                "## Key Features",
                "## Getting Started"
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
        
        # Replace existing classic leaderboard section
        new_content = (content[:start_pos] + 
                      scoreboard_content + 
                      content[end_pos:])
    
    # Write the updated content
    try:
        with open(readme_path, 'w') as f:
            f.write(new_content)
        print("README.md updated successfully with classic leaderboard!")
        return True
    except Exception as e:
        print(f"Error writing to README.md: {e}", file=sys.stderr)
        return False


def main():
    """Main function to generate and update the scoreboard."""
    print("Generating main (classic) scoreboard...")
    
    scoreboard_content = generate_main_scoreboard()
    
    if update_readme_with_scoreboard(scoreboard_content):
        print("Main scoreboard updated successfully!")
        return 0
    else:
        print("Failed to update main scoreboard!", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main()) 
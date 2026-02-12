#!/usr/bin/env python3
"""
Generate contributor profile badges for the Go Interview Practice repository.
This script creates both static and dynamic badges that contributors can use
on their GitHub profiles, LinkedIn, personal websites, etc.
"""

import json
import os
import re
import sys
from pathlib import Path
from typing import Dict, List, Tuple, Optional

class BadgeGenerator:
    def __init__(self):
        # Determine script directory and project root
        script_dir = Path(__file__).parent
        self.project_root = script_dir.parent if script_dir.name == 'scripts' else script_dir
        self.badges_dir = self.project_root / 'badges'
        self.badges_dir.mkdir(exist_ok=True)
        
        # Achievement thresholds (based on your current system)
        self.achievement_levels = {
            'Beginner': {'min_challenges': 1, 'min_rate': 0, 'color': '97ca00', 'emoji': '🌱'},
            'Advanced': {'min_challenges': 10, 'min_rate': 30, 'color': 'orange', 'emoji': '⚡'},
            'Expert': {'min_challenges': 15, 'min_rate': 50, 'color': 'blue', 'emoji': '🎯'},
            'Master': {'min_challenges': 20, 'min_rate': 65, 'color': 'gold', 'emoji': '🏆'}
        }
        
    def get_achievement_level(self, challenges_solved: int, total_challenges: int) -> Tuple[str, str, str]:
        """Determine achievement level based on challenges solved."""
        completion_rate = (challenges_solved / total_challenges * 100) if total_challenges > 0 else 0
        
        for level in ['Master', 'Expert', 'Advanced', 'Beginner']:
            if (challenges_solved >= self.achievement_levels[level]['min_challenges'] and 
                completion_rate >= self.achievement_levels[level]['min_rate']):
                return (level, 
                       self.achievement_levels[level]['color'], 
                       self.achievement_levels[level]['emoji'])
        
        return 'Beginner', self.achievement_levels['Beginner']['color'], self.achievement_levels['Beginner']['emoji']

    def scan_classic_challenges(self) -> Dict[str, int]:
        """Scan all classic challenge directories and count completions per user."""
        user_completions = {}
        total_challenges = 0
        
        # Find all challenge directories
        challenge_dirs = [d for d in self.project_root.iterdir() 
                         if d.is_dir() and d.name.startswith('challenge-')]
        
        total_challenges = len(challenge_dirs)
        
        for challenge_dir in challenge_dirs:
            scoreboard_file = challenge_dir / 'SCOREBOARD.md'
            if not scoreboard_file.exists():
                continue
                
            try:
                with open(scoreboard_file, 'r', encoding='utf-8') as f:
                    content = f.read()
                
                # Extract user data from scoreboard
                lines = content.split('\n')
                for line in lines:
                    if '|' in line and line.count('|') >= 4:
                        parts = [p.strip() for p in line.split('|')]
                        if len(parts) >= 4 and parts[1] and not parts[1].startswith('-'):
                            username = parts[1].strip()
                            try:
                                passed_tests = int(parts[2].strip())
                                total_tests = int(parts[3].strip())
                                
                                # Only count 100% completion
                                if passed_tests == total_tests and total_tests > 0:
                                    user_completions[username] = user_completions.get(username, 0) + 1
                            except (ValueError, IndexError):
                                continue
                                
            except Exception as e:
                print(f"Warning: Could not process {scoreboard_file}: {e}")
                continue
        
        return user_completions, total_challenges

    def scan_package_challenges(self) -> Dict[str, Dict[str, int]]:
        """Scan package challenges and count completions per user per package."""
        packages_dir = self.project_root / 'packages'
        if not packages_dir.exists():
            return {}
            
        user_package_completions = {}
        
        for package_dir in packages_dir.iterdir():
            if not package_dir.is_dir() or package_dir.name == 'README.md':
                continue
                
            package_name = package_dir.name
            challenge_dirs = [d for d in package_dir.iterdir() 
                            if d.is_dir() and d.name.startswith('challenge-')]
            
            for challenge_dir in challenge_dirs:
                scoreboard_file = challenge_dir / 'SCOREBOARD.md'
                if not scoreboard_file.exists():
                    continue
                    
                try:
                    with open(scoreboard_file, 'r', encoding='utf-8') as f:
                        content = f.read()
                    
                    lines = content.split('\n')
                    for line in lines:
                        if '|' in line and line.count('|') >= 4:
                            parts = [p.strip() for p in line.split('|')]
                            if len(parts) >= 4 and parts[1] and not parts[1].startswith('-'):
                                username = parts[1].strip()
                                try:
                                    passed_tests = int(parts[2].strip())
                                    total_tests = int(parts[3].strip())
                                    
                                    if passed_tests == total_tests and total_tests > 0:
                                        if username not in user_package_completions:
                                            user_package_completions[username] = {}
                                        user_package_completions[username][package_name] = \
                                            user_package_completions[username].get(package_name, 0) + 1
                                except (ValueError, IndexError):
                                    continue
                                    
                except Exception as e:
                    print(f"Warning: Could not process package {scoreboard_file}: {e}")
                    continue
        
        return user_package_completions

    def generate_shields_json(self, username: str, challenges_solved: int, total_challenges: int, 
                            package_stats: Optional[Dict[str, int]] = None) -> Dict:
        """Generate shields.io endpoint JSON for dynamic badges."""
        level, color, emoji = self.get_achievement_level(challenges_solved, total_challenges)
        completion_rate = round((challenges_solved / total_challenges * 100), 1) if total_challenges > 0 else 0
        
        # Create the main badge data
        badge_data = {
            "schemaVersion": 1,
            "label": "Go Interview Practice",
            "message": f"{emoji} {level} ({challenges_solved}/{total_challenges})",
            "color": color,
            "style": "for-the-badge"
        }
        
        return badge_data

    def generate_static_badges(self) -> Dict[str, str]:
        """Generate static badge markdown for different achievement levels."""
        badges = {}
        
        # Generic contributor badge
        badges['contributor'] = """[![Go Interview Practice Contributor](https://img.shields.io/badge/Go_Interview_Practice-Contributor-blue?style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)"""
        
        # Achievement level badges
        for level, config in self.achievement_levels.items():
            emoji = config['emoji']
            color = config['color']
            badges[level.lower()] = f"""[![Go Interview Practice {level}](https://img.shields.io/badge/Go_Interview_Practice-{emoji}_{level}-{color}?style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)"""
        
        return badges

    def generate_custom_svg_badge(self, username: str, challenges_solved: int, 
                                total_challenges: int, package_stats: Optional[Dict[str, int]] = None) -> str:
        """Generate a beautiful custom SVG badge for the user."""
        level, color, emoji = self.get_achievement_level(challenges_solved, total_challenges)
        completion_rate = round((challenges_solved / total_challenges * 100), 1) if total_challenges > 0 else 0
        
        # Enhanced color schemes for modern look
        color_schemes = {
            'gold': {'primary': '#FFD700', 'secondary': '#FFA500', 'accent': '#FF8C00'},
            'blue': {'primary': '#4A90E2', 'secondary': '#357ABD', 'accent': '#2E5F87'}, 
            'orange': {'primary': '#FF8C42', 'secondary': '#FF6B1A', 'accent': '#E55A00'},
            '97ca00': {'primary': '#97CA00', 'secondary': '#7BA428', 'accent': '#5F7E1F'}
        }
        
        scheme = color_schemes.get(color, color_schemes['blue'])
        
        # Package challenge stats
        package_count = len(package_stats) if package_stats else 0
        total_package_challenges = sum(package_stats.values()) if package_stats else 0
        
        # Progress calculations
        progress_width = int((challenges_solved / total_challenges) * 140) if total_challenges > 0 else 0
        
        svg_content = f'''<svg width="350" height="120" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <!-- Modern gradients -->
    <linearGradient id="cardGradient" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#f8f9fa;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#e9ecef;stop-opacity:1" />
    </linearGradient>
    
    <linearGradient id="headerGradient" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" style="stop-color:{scheme['primary']};stop-opacity:1" />
      <stop offset="100%" style="stop-color:{scheme['secondary']};stop-opacity:1" />
    </linearGradient>
    
    <linearGradient id="progressGradient" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" style="stop-color:{scheme['accent']};stop-opacity:1" />
      <stop offset="100%" style="stop-color:{scheme['primary']};stop-opacity:1" />
    </linearGradient>
    
    <!-- Shadow filter -->
    <filter id="dropshadow" x="-20%" y="-20%" width="140%" height="140%">
      <feDropShadow dx="2" dy="2" stdDeviation="3" flood-color="#00000020"/>
    </filter>
  </defs>
  
  <!-- Card background with shadow -->
  <rect width="350" height="120" fill="url(#cardGradient)" rx="12" filter="url(#dropshadow)"/>
  
  <!-- Header section -->
  <rect width="350" height="35" fill="url(#headerGradient)" rx="12"/>
  <rect width="350" height="25" fill="url(#headerGradient)"/>
  
  <!-- Repository info -->
  <text x="15" y="15" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" font-weight="600" fill="white" opacity="0.9">GO INTERVIEW PRACTICE</text>
  <text x="15" y="27" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="8" fill="white" opacity="0.8">github.com/RezaSi/go-interview-practice</text>
  
  <!-- Achievement level and emoji -->
  <text x="320" y="22" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="16" text-anchor="middle" fill="white">{emoji}</text>
  
  <!-- User info section -->
  <text x="15" y="58" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="14" font-weight="700" fill="#212529">@{username}</text>
  <text x="15" y="75" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" font-weight="600" fill="{scheme['primary']}">{emoji} {level} Developer</text>
  
  <!-- Classic Challenges Section -->
  <text x="15" y="95" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" font-weight="500" fill="#6c757d">Classic Challenges</text>
  
  <!-- Progress bar background -->
  <rect x="15" y="100" width="140" height="6" fill="#e9ecef" rx="3"/>
  <!-- Progress bar fill -->
  <rect x="15" y="100" width="{progress_width}" height="6" fill="url(#progressGradient)" rx="3"/>
  
  <!-- Progress text -->
  <text x="160" y="106" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="9" font-weight="600" fill="#495057">{challenges_solved}/{total_challenges} ({completion_rate}%)</text>
  
  <!-- Package Challenges Section (if any) -->'''
        
        if package_count > 0:
            svg_content += f'''
  <text x="190" y="58" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" font-weight="500" fill="#6c757d">Package Challenges</text>
  <text x="190" y="72" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" font-weight="600" fill="{scheme['secondary']}">{total_package_challenges} across {package_count} packages</text>
  
  <!-- Package icons -->
  <circle cx="195" cy="82" r="3" fill="{scheme['primary']}" opacity="0.8"/>
  <circle cx="205" cy="82" r="3" fill="{scheme['secondary']}" opacity="0.8"/>
  <circle cx="215" cy="82" r="3" fill="{scheme['accent']}" opacity="0.8"/>'''
        else:
            svg_content += f'''
  <text x="190" y="65" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" font-weight="500" fill="#6c757d">Ready for</text>
  <text x="190" y="78" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" font-weight="600" fill="{scheme['secondary']}">Package Challenges!</text>'''
        
        # Add achievement indicator
        if challenges_solved >= 20:
            svg_content += f'''
  <!-- Achievement indicator -->
  <circle cx="320" cy="85" r="8" fill="{scheme['primary']}" opacity="0.2"/>
  <text x="320" y="89" font-family="SF Pro Display,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" text-anchor="middle" fill="{scheme['primary']}" font-weight="700">★</text>'''
        
        svg_content += '''
</svg>'''
        
        return svg_content

    def generate_compact_badge(self, username: str, challenges_solved: int, 
                             total_challenges: int, package_stats: Optional[Dict[str, int]] = None) -> str:
        """Generate a compact horizontal badge for GitHub README."""
        level, color, emoji = self.get_achievement_level(challenges_solved, total_challenges)
        completion_rate = round((challenges_solved / total_challenges * 100), 1) if total_challenges > 0 else 0
        
        # Enhanced color schemes
        color_schemes = {
            'gold': {'primary': '#FFD700', 'secondary': '#FFA500'},
            'blue': {'primary': '#4A90E2', 'secondary': '#357ABD'}, 
            'orange': {'primary': '#FF8C42', 'secondary': '#FF6B1A'},
            '97ca00': {'primary': '#97CA00', 'secondary': '#7BA428'}
        }
        
        scheme = color_schemes.get(color, color_schemes['blue'])
        
        # Package info
        package_count = len(package_stats) if package_stats else 0
        total_package_challenges = sum(package_stats.values()) if package_stats else 0
        
        # Progress calculation
        progress_width = int((challenges_solved / total_challenges) * 100) if total_challenges > 0 else 0
        
        svg_content = f'''<svg width="400" height="60" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient id="compactGradient" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" style="stop-color:{scheme['primary']};stop-opacity:1" />
      <stop offset="100%" style="stop-color:{scheme['secondary']};stop-opacity:1" />
    </linearGradient>
    
    <linearGradient id="compactBg" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" style="stop-color:#ffffff;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#f8f9fa;stop-opacity:1" />
    </linearGradient>
    
    <filter id="shadow" x="-10%" y="-10%" width="120%" height="120%">
      <feDropShadow dx="1" dy="1" stdDeviation="2" flood-color="#00000015"/>
    </filter>
  </defs>
  
  <!-- Main background -->
  <rect width="400" height="60" fill="url(#compactBg)" rx="8" filter="url(#shadow)"/>
  
  <!-- Left accent -->
  <rect width="6" height="60" fill="url(#compactGradient)" rx="8"/>
  
  <!-- Go logo and title -->
  <text x="20" y="20" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" font-weight="700" fill="#212529">🐹 Go Interview Practice</text>
  
  <!-- Username and level -->
  <text x="20" y="38" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="11" fill="#6c757d">@{username}</text>
  <text x="20" y="52" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" font-weight="600" fill="{scheme['primary']}">{emoji} {level} Developer</text>
  
  <!-- Progress section -->
  <text x="220" y="20" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" font-weight="500" fill="#495057">Classic Progress</text>
  
  <!-- Progress bar -->
  <rect x="220" y="25" width="100" height="4" fill="#e9ecef" rx="2"/>
  <rect x="220" y="25" width="{progress_width}" height="4" fill="url(#compactGradient)" rx="2"/>
  
  <!-- Stats -->
  <text x="220" y="42" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="10" font-weight="600" fill="#495057">{challenges_solved}/{total_challenges} ({completion_rate}%)</text>'''
        
        if package_count > 0:
            svg_content += f'''
  <text x="220" y="54" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="9" fill="{scheme['secondary']}">📦 {total_package_challenges} package challenges</text>'''
        
        # Achievement badge
        if challenges_solved >= 20:
            svg_content += f'''
  <circle cx="365" cy="30" r="12" fill="{scheme['primary']}" opacity="0.1"/>
  <text x="365" y="34" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" text-anchor="middle" fill="{scheme['primary']}">⭐</text>'''
        elif challenges_solved >= 15:
            svg_content += f'''
  <circle cx="365" cy="30" r="12" fill="{scheme['primary']}" opacity="0.1"/>
  <text x="365" y="34" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" text-anchor="middle" fill="{scheme['primary']}">🎯</text>'''
        elif challenges_solved >= 10:
            svg_content += f'''
  <circle cx="365" cy="30" r="12" fill="{scheme['primary']}" opacity="0.1"/>
  <text x="365" y="34" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" text-anchor="middle" fill="{scheme['primary']}">⚡</text>'''
        else:
            svg_content += f'''
  <circle cx="365" cy="30" r="12" fill="{scheme['primary']}" opacity="0.1"/>
  <text x="365" y="34" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif" font-size="12" text-anchor="middle" fill="{scheme['primary']}">🌱</text>'''
        
        svg_content += '''
</svg>'''
        
        return svg_content

    def generate_readme_badges(self, username: str, challenges_solved: int, 
                             total_challenges: int, package_stats: Optional[Dict[str, int]] = None) -> str:
        """Generate a comprehensive badge collection for README files."""
        level, color, emoji = self.get_achievement_level(challenges_solved, total_challenges)
        completion_rate = round((challenges_solved / total_challenges * 100), 1) if total_challenges > 0 else 0
        
        # Dynamic badge (requires the JSON endpoint) - with clickable link
        dynamic_badge = f"""[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/{username}.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)"""
        
        # Static badges - with clickable links
        challenge_badge = f"""[![Challenges Solved](https://img.shields.io/badge/Go_Challenges-{challenges_solved}%2F{total_challenges}-brightgreen?style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)"""
        
        level_badge = f"""[![Achievement Level](https://img.shields.io/badge/Level-{emoji}_{level}-{color}?style=for-the-badge&logo=trophy&logoColor=white)](https://github.com/RezaSi/go-interview-practice)"""
        
        completion_badge = f"""[![Completion Rate](https://img.shields.io/badge/Completion-{completion_rate}%25-{color}?style=for-the-badge&logo=checkmarx&logoColor=white)](https://github.com/RezaSi/go-interview-practice)"""
        
        # Package badges if available - with clickable link
        package_badges = ""
        if package_stats:
            total_package_challenges = sum(package_stats.values())
            package_count = len(package_stats)
            package_badges = f"""[![Package Challenges](https://img.shields.io/badge/Package_Challenges-{total_package_challenges}_across_{package_count}_packages-purple?style=for-the-badge&logo=package&logoColor=white)](https://github.com/RezaSi/go-interview-practice)"""
        
        badges_collection = f"""## 🏆 Go Interview Practice Achievements

### 🎨 Beautiful Custom Badges
*Click any badge to visit the Go Interview Practice repository!*

<!-- Full-size Card Badge - Clickable -->
[![Go Interview Practice Achievement Card](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/{username}.svg)](https://github.com/RezaSi/go-interview-practice)

<!-- Compact Horizontal Badge - Clickable -->
[![Go Interview Practice Compact](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/{username}_compact.svg)](https://github.com/RezaSi/go-interview-practice)

### 🔄 Dynamic Shields.io Badge
<!-- Dynamic Badge (auto-updates) -->
{dynamic_badge}

### 📊 Static Badges Collection
{challenge_badge}
{level_badge}
{completion_badge}
{package_badges}

### 🔗 Repository Link Badge
[![Go Interview Practice Repository](https://img.shields.io/badge/View_Repository-Go_Interview_Practice-blue?style=for-the-badge&logo=github&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

---

### 📈 Your Achievement Summary

**👤 Username:** @{username}  
**🏅 Achievement Level:** {emoji} **{level} Developer**  
**📊 Classic Challenges:** {challenges_solved}/{total_challenges} ({completion_rate}% complete)  
**🔗 Repository:** [Go Interview Practice](https://github.com/RezaSi/go-interview-practice)  
"""
        
        if package_stats:
            badges_collection += f"**Package Challenges:** {sum(package_stats.values())} across {len(package_stats)} packages\n"
        
        return badges_collection

    def run(self):
        """Main execution function."""
        print("🎯 Generating contributor badges...")
        
        # Scan challenges
        print("📊 Scanning classic challenges...")
        classic_completions, total_classic = self.scan_classic_challenges()
        
        print("📦 Scanning package challenges...")
        package_completions = self.scan_package_challenges()
        
        # Generate badges for each contributor
        generated_count = 0
        
        # Combine all contributors
        all_contributors = set(classic_completions.keys()) | set(package_completions.keys())
        
        for username in all_contributors:
            classic_solved = classic_completions.get(username, 0)
            package_stats = package_completions.get(username, {})
            
            # Generate dynamic badge JSON
            badge_json = self.generate_shields_json(username, classic_solved, total_classic, package_stats)
            json_file = self.badges_dir / f"{username}.json"
            
            with open(json_file, 'w', encoding='utf-8') as f:
                json.dump(badge_json, f, indent=2)
            
            # Generate custom SVG (full card)
            svg_content = self.generate_custom_svg_badge(username, classic_solved, total_classic, package_stats)
            svg_file = self.badges_dir / f"{username}.svg"
            
            with open(svg_file, 'w', encoding='utf-8') as f:
                f.write(svg_content)
            
            # Generate compact badge
            compact_svg_content = self.generate_compact_badge(username, classic_solved, total_classic, package_stats)
            compact_svg_file = self.badges_dir / f"{username}_compact.svg"
            
            with open(compact_svg_file, 'w', encoding='utf-8') as f:
                f.write(compact_svg_content)
            
            # Generate README badges collection (optional - comment out to reduce files)
            readme_content = self.generate_readme_badges(username, classic_solved, total_classic, package_stats)
            readme_file = self.badges_dir / f"{username}_badges.md"
            
            with open(readme_file, 'w', encoding='utf-8') as f:
                f.write(readme_content)
            
            generated_count += 1
            print(f"  ✅ Generated badges for {username} ({classic_solved} challenges)")
        
        # Generate static badge templates
        static_badges = self.generate_static_badges()
        static_file = self.badges_dir / "static_badges.md"
        
        with open(static_file, 'w', encoding='utf-8') as f:
            f.write("# Static Badge Templates\n\n")
            f.write("These badges can be used by any contributor:\n\n")
            for badge_name, badge_code in static_badges.items():
                f.write(f"## {badge_name.title()}\n")
                f.write(f"```markdown\n{badge_code}\n```\n")
                f.write(f"{badge_code}\n\n")
        
        # Generate instructions
        instructions = self.generate_instructions()
        instructions_file = self.badges_dir / "README.md"
        
        with open(instructions_file, 'w', encoding='utf-8') as f:
            f.write(instructions)
        
        print(f"\n🎉 Successfully generated badges for {generated_count} contributors!")
        print(f"📁 Badge files saved to: {self.badges_dir}")
        print(f"📖 Instructions: {instructions_file}")

    def generate_instructions(self) -> str:
        """Generate comprehensive instructions for using the badges."""
        return """# 🏆 Go Interview Practice Profile Badges

This directory contains beautiful, modern profile badges that contributors can use to showcase their achievements in the Go Interview Practice repository. All badges feature modern UI/UX design with progress bars, gradients, and comprehensive information display.

## 🎨 Badge Collection Overview

Each contributor gets a complete badge collection:

### 📄 Files Generated For Each User
- `USERNAME.svg` - **Full-size card badge** (350×120px) with complete stats
- `USERNAME_compact.svg` - **Compact horizontal badge** (400×60px) for README headers  
- `USERNAME.json` - **Dynamic badge data** for shields.io integration
- `USERNAME_badges.md` - **Complete collection** with all badge types ready to copy

## 🎯 Badge Types & Usage

### 1. **🎨 Full-Size Card Badge** (Recommended for Profiles)
Beautiful card-style badge with comprehensive information:

```markdown
![Go Interview Practice Achievement](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.svg)
```

**Features:**
- 🏆 Repository branding with your GitHub URL
- 👤 Username and achievement level  
- 📊 Progress bar for classic challenges
- 📦 Package challenges information
- 🎨 Modern gradients and typography
- ⭐ Achievement indicators for high performers

### 2. **⚡ Compact Horizontal Badge** (Great for README Headers)
Sleek horizontal layout perfect for project headers:

```markdown
![Go Interview Practice Compact](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME_compact.svg)
```

**Features:**
- 🚀 Clean, professional design
- 📈 Progress visualization
- 🏅 Achievement level indication
- 📦 Package challenges summary

### 3. **🔄 Dynamic Shields.io Badge** (Auto-Updating)
These badges automatically update based on your current progress:

```markdown
![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.json&style=for-the-badge&logo=go&logoColor=white)
```

**Styles Available:**
- `for-the-badge` (recommended for profiles)
- `flat`, `flat-square`, `plastic`, `social`

### 4. **📊 Static Template Badges**
Ready-to-use badges for any contributor:

```markdown
[![Go Interview Practice Contributor](https://img.shields.io/badge/Go_Interview_Practice-Contributor-blue?style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)
```

## 🏅 Achievement System

Your badges automatically reflect your achievement level with unique colors and emojis:

| Level | Requirements | Badge Color | Visual Theme |
|-------|-------------|-------------|--------------|
| 🌱 **Beginner** | 1+ challenges | Fresh Green | Growth theme |
| ⚡ **Advanced** | 10+ challenges (30%+) | Energy Orange | Power theme |
| 🎯 **Expert** | 15+ challenges (50%+) | Professional Blue | Precision theme |
| 🏆 **Master** | 20+ challenges (65%+) | Golden | Excellence theme |

## 📱 Usage Examples

### 🐙 GitHub Profile README
Perfect for showcasing your coding journey:

```markdown
## 🏆 My Coding Achievements

![Go Interview Practice Achievement](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.svg)

I've been mastering Go programming through structured challenges, solving algorithmic problems and building real-world applications with popular Go frameworks.
```

### 🌐 Personal Website/Portfolio
```html
<div class="achievements">
  <h3>Go Programming Expertise</h3>
  <img src="https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME_compact.svg" 
       alt="Go Interview Practice Achievement" />
  <p>Demonstrated proficiency in Go through completion of coding challenges</p>
</div>
```

### 💼 LinkedIn Profile
1. Screenshot your full-size card badge
2. Upload as a "License & Certification" 
3. Title: "Go Programming - Interview Practice Completion"
4. Link: `https://github.com/RezaSi/go-interview-practice`

### 📋 CV/Resume Section
**Technical Achievements:**
- Go Interview Practice: Completed X/30 challenges ([Your Level] level)
- Demonstrated algorithmic problem-solving skills
- Experience with Go frameworks (Gin, Cobra, GORM) [if applicable]

## 🔄 Auto-Update Features

### When Badges Refresh
Your badges automatically update when:
- ✅ You solve new challenges
- ✅ Your achievement level increases  
- ✅ Package challenges are completed
- ✅ Scoreboards are regenerated
- ✅ New challenges are added

### Update Frequency
- **GitHub Actions**: Automatic regeneration on repository changes
- **CDN Refresh**: Changes appear within 1-5 minutes
- **Manual Trigger**: Repository maintainers can force updates

## 🚀 Getting Your Badges

### Step 1: Contribute
1. Fork the [Go Interview Practice](https://github.com/RezaSi/go-interview-practice) repository
2. Solve at least one challenge
3. Submit your solution via pull request

### Step 2: Automatic Generation
1. Your pull request gets merged
2. GitHub Actions automatically regenerates badges
3. Your badge files appear in this directory within minutes

### Step 3: Use Your Badges
1. Browse to `badges/YOUR_USERNAME_badges.md`
2. Copy the markdown for your preferred badge style
3. Paste into your GitHub profile, website, or portfolio

## 🎨 Design Philosophy

Our badges feature modern UI/UX principles:
- **🎯 Information Hierarchy**: Most important info (achievement level) prominently displayed
- **🌈 Color Psychology**: Colors match achievement levels and create visual progression
- **📱 Responsive Design**: Works perfectly across all devices and platforms
- **⚡ Performance**: Lightweight SVG graphics load instantly
- **🎨 Modern Aesthetics**: Clean typography, subtle shadows, and beautiful gradients

## 📞 Support & Troubleshooting

### Common Issues
- **Badge not found?** → Ensure you've submitted at least one solution
- **Badge not updating?** → Wait 5 minutes for CDN refresh, check GitHub Actions
- **Want different size?** → Use `_compact.svg` for smaller displays

### Get Help
- 📧 **Email**: [rezashiri88@gmail.com](mailto:rezashiri88@gmail.com)
- 🐙 **GitHub Issues**: [Report a problem](https://github.com/RezaSi/go-interview-practice/issues)
- 💬 **Discussions**: [Community support](https://github.com/RezaSi/go-interview-practice/discussions)

---

## 🎉 Badge Gallery

Check out badges from our top contributors:

**🏆 Master Level Examples:**
- [odelbos badges](./odelbos_badges.md) - 28/30 challenges 
- [mick4711 badges](./mick4711_badges.md) - 21/30 challenges

**🎯 Expert Level Examples:**
- [y1hao badges](./y1hao_badges.md) - 21/30 challenges
- [JackDalberg badges](./JackDalberg_badges.md) - 20/30 challenges

---

**Repository**: [Go Interview Practice](https://github.com/RezaSi/go-interview-practice)  
**Start Your Journey**: Fork, solve, and earn your badge today! 🚀

*Beautiful badges that showcase your coding excellence and motivate continuous learning.*
"""

if __name__ == "__main__":
    generator = BadgeGenerator()
    generator.run()
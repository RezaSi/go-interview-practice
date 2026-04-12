# 🏆 Profile Badges for Contributors

Welcome to the **Go Interview Practice** profile badges system! This guide shows you how to showcase your coding achievements using beautiful badges on your GitHub profile, LinkedIn, personal website, or anywhere you want to demonstrate your Go programming skills.

> **✨ All badges are clickable!** When users click your badge, they'll be taken directly to the Go Interview Practice repository, helping you promote the project while showcasing your achievements.

## 📚 Quick Navigation

- [🎯 Quick Start](#-quick-start) - Get your badges in 3 steps
- [🎨 Badge Types & Usage](#-badge-types--usage) - Full overview of available badges
- [✨ Examples](#-badge-examples) - See badges in action
- [📱 Usage Examples](#-usage-examples) - GitHub, LinkedIn, website integration
- [🚀 Getting Your Badges](#-getting-your-badges) - Step-by-step process

## 🎯 Quick Start

### Step 1: Find Your Badge Collection
After contributing to the repository, find your personalized badges:

```
badges/YOUR_USERNAME_badges.md    ← Your complete badge collection (start here!)
badges/YOUR_USERNAME.svg          ← Full-size card badge  
badges/YOUR_USERNAME_compact.svg  ← Compact horizontal badge
```

### Step 2: Copy & Paste
1. Open [`badges/YOUR_USERNAME_badges.md`](../badges/)
2. Copy the markdown for your preferred badge style
3. Paste into your GitHub profile README, website, or portfolio

### Step 3: Show Off Your Skills! 🚀
Your badges automatically update as you solve more challenges - no manual work needed!

## 🎨 Badge Types & Usage

### 1. **Dynamic Badges** ⭐ *Recommended*
These badges automatically update when you solve more challenges:

**Example for user `odelbos`:**
[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/odelbos.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

**Your Dynamic Badge (Clickable):**
```markdown
[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)
```

### 2. **Custom SVG Badges**
Beautiful custom-designed badges with gradients and achievement levels:

## ✨ Badge Examples

**Beautiful Modern Designs with Progress Bars:**

**🏆 Master Level (Gold)** - odelbos with 28/30 challenges:
[![Go Interview Practice Achievement](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/odelbos.svg)](https://github.com/RezaSi/go-interview-practice)

**⚡ Advanced Level (Orange)** - RezaSi with 14/30 challenges:
[![Go Interview Practice Achievement](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/RezaSi.svg)](https://github.com/RezaSi/go-interview-practice)

**🎯 Expert Level (Blue)** - ashwinipatankar with 17/30 challenges:
[![Go Interview Practice Achievement](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/ashwinipatankar.svg)](https://github.com/RezaSi/go-interview-practice)

**⚡ Compact Horizontal Style:**
[![Go Interview Practice Compact](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/odelbos_compact.svg)](https://github.com/RezaSi/go-interview-practice)

**Your SVG Badge (Clickable):**
```markdown
[![Go Interview Practice Achievement](https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.svg)](https://github.com/RezaSi/go-interview-practice)
```

### 3. **Static Badges**
Simple badges that anyone can use regardless of their progress:

[![Go Interview Practice Contributor](https://img.shields.io/badge/Go_Interview_Practice-Contributor-blue?style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

```markdown
[![Go Interview Practice Contributor](https://img.shields.io/badge/Go_Interview_Practice-Contributor-blue?style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)
```

## 🏅 Achievement System

Your badge automatically reflects your achievement level:

| Level | Requirements | Badge Color | Emoji |
|-------|-------------|-------------|-------|
| 🌱 **Beginner** | 1+ challenges | Green | 🌱 |
| ⚡ **Advanced** | 10+ challenges (30%+ completion) | Orange | ⚡ |
| 🎯 **Expert** | 15+ challenges (50%+ completion) | Blue | 🎯 |
| 🏆 **Master** | 20+ challenges (65%+ completion) | Gold | 🏆 |

## 🎨 Customization Options

### Badge Styles
Change the `style` parameter for different looks:
- `for-the-badge` - Large, professional (recommended for profiles)
- `flat` - Minimal, clean
- `flat-square` - Square corners
- `plastic` - Glossy look
- `social` - Social media style

### Colors
Available colors for static badges:
- `brightgreen`, `green`, `yellowgreen`, `yellow`, `orange`, `red`
- `lightgrey`, `blue`, `purple`, `pink`
- Custom hex colors: `#ff69b4`

## 📋 Usage Examples

### GitHub Profile README
```markdown
## 🏆 My Go Interview Practice Journey

[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

I've been practicing Go programming through the [Go Interview Practice](https://github.com/RezaSi/go-interview-practice) repository, where I've solved multiple coding challenges and improved my algorithmic thinking skills.
```

### Personal Portfolio Website
```html
<div class="badges">
  <h3>My Coding Achievements</h3>
  <img src="https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.svg" 
       alt="Go Interview Practice Achievement" />
  <p>Completed multiple Go programming challenges</p>
</div>
```

### LinkedIn Profile
1. Download your SVG badge as an image (screenshot or convert to PNG)
2. Upload as a "License & Certification" or in your summary section
3. Link back to the repository: `https://github.com/RezaSi/go-interview-practice`

### CV/Resume
Include the badge image and mention:
- "Completed X/30 Go programming challenges"
- "Achieved [Your Level] level in algorithmic problem solving"
- "Active contributor to open-source Go learning project"

## 🔄 Auto-Updates

### When Badges Update
Your dynamic badges automatically refresh when:
- ✅ You solve new challenges
- ✅ Your achievement level increases  
- ✅ Scoreboard data is regenerated
- ✅ New challenges are added to the repository

### Update Frequency
- Badges are regenerated when scoreboards change
- GitHub Actions automatically run the badge generator
- Changes appear within minutes of scoreboard updates

## 🎯 Professional Tips

### For Job Applications
```markdown
## Technical Skills Demonstration

I actively practice algorithmic problem-solving through structured challenges:

[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/YOUR_USERNAME.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

**Repository**: [Go Interview Practice](https://github.com/RezaSi/go-interview-practice)
**My Solutions**: [View my submissions](https://github.com/RezaSi/go-interview-practice/tree/main/challenge-*/submissions/YOUR_USERNAME)
```

### For Networking
Use badges in:
- GitHub profile README
- Dev.to articles about your learning journey
- Twitter/LinkedIn posts about your progress
- Personal blog posts about Go programming
- Conference speaker bio slides

## 🚀 Getting Started

### Step 1: Contribute
1. Fork the [Go Interview Practice](https://github.com/RezaSi/go-interview-practice) repository
2. Solve at least one challenge
3. Submit your solution via pull request

### Step 2: Get Your Badge
1. Wait for your solution to be merged
2. Badges are auto-generated after merges
3. Find your badge files in the [`badges/`](../badges/) directory

### Step 3: Show Off Your Achievement
1. Copy the markdown from `badges/YOUR_USERNAME_badges.md`
2. Paste into your GitHub profile README
3. Share your progress with the community!

## 🎊 Badge Examples

Here are some real examples from our top contributors:

### Master Level (20+ challenges)
[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/odelbos.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

### Expert Level (15+ challenges)  
[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/ashwinipatankar.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

### Advanced Level (10+ challenges)
[![Go Interview Practice](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/RezaSi/go-interview-practice/main/badges/RezaSi.json&style=for-the-badge&logo=go&logoColor=white)](https://github.com/RezaSi/go-interview-practice)

## 📞 Need Help?

### Troubleshooting
- **Badge not found?** Make sure you've submitted at least one solution
- **Badge not updating?** Wait a few minutes for GitHub's CDN to refresh
- **Want different style?** Modify the `style` parameter in the URL

### Contact
- 📧 **Email**: [rezashiri88@gmail.com](mailto:rezashiri88@gmail.com)
- 🐙 **GitHub**: [@RezaSi](https://github.com/RezaSi)
- 💬 **Issues**: [Repository Issues](https://github.com/RezaSi/go-interview-practice/issues)

---

**Start your Go journey today!** 🚀  
[**Join Go Interview Practice →**](https://github.com/RezaSi/go-interview-practice)

*Show the world your coding skills with beautiful, automatically-updating achievement badges!*
# Package Challenges - Dynamic System Documentation

This directory contains package-specific coding challenges for popular Go libraries and frameworks. The system is designed to be completely dynamic, allowing easy addition of new packages without requiring code changes.

## Available Packages

### 🗄️ [GORM](./gorm/) - ORM Library
**5 Challenges** | Beginner to Advanced | **6-8 hours**
- Database operations, associations, migrations, advanced queries, and generics API

### 🌐 [Gin](./gin/) - Web Framework  
**5 Challenges** | Beginner to Advanced | **6-8 hours**
- HTTP routing, middleware, authentication, file handling, and testing

### ⚡ [Cobra](./cobra/) - CLI Framework
**4 Challenges** | Beginner to Advanced | **4-6 hours**
- Command-line applications, flags, subcommands, data persistence, and advanced patterns

*More packages coming soon...*

## Directory Structure

```
packages/
├── {package-name}/
│   ├── package.json                    # Package metadata
│   ├── challenge-N-{name}/
│   │   ├── metadata.json              # Challenge metadata (optional)
│   │   ├── README.md                  # Challenge description
│   │   ├── solution-template.go       # Template code for users
│   │   ├── solution-template_test.go  # Test file
│   │   ├── hints.md                   # Hints and tips
│   │   └── submissions/               # User submissions
│   │       └── {username}/
│   │           └── solution.go        # User's solution
│   └── ...
```

## Adding a New Package

### 1. Create Package Directory
Create a new directory with your package name:
```bash
mkdir packages/{package-name}
```

### 2. Create package.json
Define package metadata in `packages/{package-name}/package.json`:

```json
{
  "name": "package-name",
  "display_name": "Package Display Name",
  "description": "Brief description of the package",
  "version": "v1.0.0",
  "github_url": "https://github.com/owner/repo",
  "documentation_url": "https://package-docs.com",
  "stars": 10000,
  "category": "web|cli|database|other",
  "difficulty": "beginner_to_advanced",
  "prerequisites": ["basic_go", "http_concepts"],
  "learning_path": [
    "challenge-1-basic-feature",
    "challenge-2-advanced-feature"
  ],
  "tags": ["tag1", "tag2", "tag3"],
  "estimated_time": "4-6 hours",
  "real_world_usage": [
    "Use case 1",
    "Use case 2"
  ]
}
```

### 3. Create Challenges
For each challenge in your learning path:

#### Challenge Directory
```bash
mkdir packages/{package-name}/challenge-N-{name}
```

#### Required Files

**README.md** - Challenge description and instructions
**solution-template.go** - Starting code template
**solution-template_test.go** - Test cases
**hints.md** - Helpful hints for learners

#### Optional metadata.json
For enhanced challenge information:

```json
{
  "title": "Challenge Title",
  "description": "Detailed description",
  "short_description": "Brief description for cards",
  "difficulty": "Beginner|Intermediate|Advanced",
  "estimated_time": "30-45 min",
  "learning_objectives": [
    "Objective 1",
    "Objective 2"
  ],
  "prerequisites": ["prerequisite1"],
  "tags": ["tag1", "tag2"],
  "real_world_connection": "How this applies in real projects",
  "requirements": [
    "Requirement 1",
    "Requirement 2"
  ],
  "bonus_points": [
    "Bonus task 1"
  ],
  "icon": "bi-icon-name",
  "order": 1
}
```

## How the Dynamic System Works

### 1. Package Discovery
- The system automatically scans the `packages/` directory
- Each subdirectory is treated as a package
- Package metadata is loaded from `package.json`

### 2. Challenge Loading
- Challenges are discovered from the `learning_path` in `package.json`
- Challenge directories are scanned for content
- Metadata is loaded from `metadata.json` if available
- Fallback metadata is generated from directory names and README files

### 3. Template Functions
The system provides dynamic template functions:
- `isPackageActive` - Check if package has available challenges
- `getPackageChallenges` - Get ordered list of challenges
- `getChallengeInfo` - Get metadata for specific challenge
- `getDifficultyBadgeClass` - Get CSS class for difficulty
- `getCategoryIcon` - Get icon for package category
- `isComingSoon` - Check if challenge is not yet available

### 4. Status Management
Challenges automatically have status:
- **available** - Challenge directory exists with content
- **coming-soon** - Challenge listed in learning_path but directory doesn't exist

## Benefits of Dynamic System

1. **Zero Code Changes** - Add packages without modifying application code
2. **Consistent UI** - All packages render with the same templates
3. **Flexible Metadata** - Rich challenge information through JSON
4. **Automatic Discovery** - New packages appear immediately
5. **Fallback Support** - Works with minimal metadata, enhances with more
6. **Easy Maintenance** - Package-specific logic contained in metadata

## Package Categories

- **web** - Web frameworks and HTTP libraries
- **cli** - Command-line tools and frameworks  
- **database** - Database drivers and ORMs
- **other** - General purpose libraries

## Icons

Use Bootstrap Icons for challenges:
- `bi-play-circle` - Basic/intro challenges
- `bi-layers` - Middleware/architecture 
- `bi-shield-check` - Validation/security
- `bi-person-lock` - Authentication
- `bi-cloud-upload` - File handling
- `bi-database` - Database operations
- `bi-terminal` - CLI operations
- `bi-code-slash` - General coding

## Best Practices

1. **Learning Path Order** - Arrange challenges from basic to advanced
2. **Clear Descriptions** - Write helpful, specific challenge descriptions
3. **Good Test Coverage** - Provide comprehensive test cases
4. **Practical Examples** - Use real-world scenarios in challenges
5. **Progressive Difficulty** - Each challenge should build on previous ones
6. **Helpful Hints** - Provide hints without giving away solutions

## Contributing

For detailed contribution guidelines for package challenges, see [CONTRIBUTING.md](../CONTRIBUTING.md#package-challenges-frameworklibrary-focused).

### Quick Guidelines

1. **Follow the Directory Structure** - Use the exact structure shown above
2. **Include All Required Files** - README.md, solution-template.go, tests, and hints
3. **Create Working Solutions** - Include a complete solution in submissions/RezaSi/
4. **Test Thoroughly** - Ensure all tests pass and edge cases are covered
5. **Write Clear Documentation** - Provide comprehensive learning materials
6. **Use Appropriate Difficulty** - Match difficulty to target audience
7. **Ensure Learning Objectives** - Each challenge should have clear educational goals
8. **Follow Package Conventions** - Use consistent naming and structure
9. **Include Dependencies** - Set up proper go.mod with all required packages
10. **Create Executable Scripts** - Provide run_tests.sh for validation

### Template Files Included

The system provides template files for:
- **metadata.json** - Challenge metadata structure
- **go.mod** - Module configuration with dependencies  
- **solution-template.go** - Code template with TODOs
- **solution-template_test.go** - Comprehensive test suite
- **learning.md** - Educational content (400+ lines recommended)
- **hints.md** - Step-by-step guidance
- **run_tests.sh** - Testing and validation script

The system will automatically detect and display your new package challenges! 
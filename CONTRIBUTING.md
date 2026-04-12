# **Contributing to Go Interview Practice**

Thank you for your interest in contributing to the **Go Interview Practice** repository! We welcome contributions from the community to help improve this project. Whether you want to submit solutions, add new challenges, or improve documentation, your efforts are appreciated.

## **Table of Contents**

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
  - [Submitting a Solution](#submitting-a-solution)
  - [Adding a New Challenge](#adding-a-new-challenge)
    - [Classic vs Package Challenges](#classic-vs-package-challenges)
    - [Classic Challenges](#classic-challenges-algorithmdata-structure-focused)
    - [Package Challenges](#package-challenges-frameworklibrary-focused)
- [Style Guidelines](#style-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)
- [Contact](#contact)

---

## **Code of Conduct**

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project, you agree to abide by its terms.

---

## **How to Contribute**

### **Submitting a Solution**

You can submit solutions to both Classic and Package challenges:

#### **For Classic Challenges**

1. **Fork the Repository:**

   - Click the "Fork" button on the repository page.

2. **Clone Your Fork:**

   ```bash
   git clone https://github.com/yourusername/go-interview-practice.git
   ```

3. **Create a New Branch:**

   ```bash
   git checkout -b challenge-[number]-solution
   ```

4. **Set Up Your Submission:**

   - Use the provided script to set up your submission:

     ```bash
     ./create_submission.sh [challenge-number]
     ```

5. **Implement Your Solution:**

   - Edit the `solution-template.go` file in your submission directory.
   - Ensure your code passes all the tests.

6. **Run Tests Locally:**

   - Navigate to the challenge directory and use `run_tests.sh`:

     ```bash
     cd challenge-[number]
     ./run_tests.sh
     ```

7. **Commit and Push:**

   ```bash
   git add challenge-[number]/submissions/yourusername/
   git commit -m "Add Challenge [number] solution by [yourusername]"
   git push origin challenge-[number]-solution
   ```

#### **For Package Challenges**

1. **Fork and Clone** (same as above)

2. **Create a New Branch:**

   ```bash
   git checkout -b package-[package-name]-challenge-[number]-solution
   ```

3. **Navigate to the Package Challenge:**

   ```bash
   cd packages/[package-name]/challenge-[number]-[topic]
   ```

4. **Create Your Submission Directory:**

   ```bash
   mkdir -p submissions/yourusername
   ```

5. **Implement Your Solution:**

   - Copy the `solution-template.go` to your submission directory:

     ```bash
     cp solution-template.go submissions/yourusername/solution.go
     ```

   - Edit `submissions/yourusername/solution.go` and complete all TODOs.
   - Ensure your solution follows the package requirements and passes all tests.

6. **Run Tests Locally:**

   - Use the package challenge test script:

     ```bash
     ./run_tests.sh
     # When prompted, enter your GitHub username
     ```

7. **Commit and Push:**

   ```bash
   git add packages/[package-name]/challenge-[number]-[topic]/submissions/yourusername/
   git commit -m "Add [Package] Challenge [number] solution by [yourusername]"
   git push origin package-[package-name]-challenge-[number]-solution
   ```

#### **General Submission Guidelines**

8. **Create a Pull Request:**

   - Go to your fork on GitHub and open a pull request to the `main` branch.
   - Use a descriptive title and mention which challenge you solved.

9. **Receive Feedback:**

   - The automated tests will run on your pull request.
   - Address any comments or requested changes.
   - Package challenge solutions will be automatically added to the scoreboard upon merge.

### **Adding a New Challenge**

There are two types of challenges you can contribute:

#### **Classic vs Package Challenges**

Before contributing a new challenge, it's important to understand the difference between the two types:

**Classic Challenges:**
- **Focus:** Algorithm and data structure problems
- **Purpose:** Fundamental programming concepts and problem-solving skills
- **Structure:** Single challenge directory with standalone problems
- **Examples:** Binary search, linked list manipulation, dynamic programming
- **Target Audience:** All developers regardless of framework experience
- **Location:** `challenge-[number]/` directories in the root

**Package Challenges:**
- **Focus:** Real-world application development with specific Go packages/frameworks
- **Purpose:** Practical skills for building production applications
- **Structure:** Package-based directory with progressive challenge series
- **Examples:** REST APIs with Gin, CLI tools with Cobra, database operations with GORM
- **Target Audience:** Developers learning specific frameworks or building portfolio projects
- **Location:** `packages/[package-name]/challenge-[number]-[topic]/` directories

**When to Choose Each Type:**

- Choose **Classic Challenges** for:
  - Algorithm problems from coding interviews
  - Data structure implementations
  - Mathematical or logical puzzles
  - Language-agnostic programming concepts

- Choose **Package Challenges** for:
  - Framework-specific tutorials
  - Building complete applications
  - Learning industry-standard libraries
  - Demonstrating real-world development patterns

#### **Classic Challenges (Algorithm/Data Structure Focused)**

For traditional algorithm and data structure challenges:

1. **Create a New Issue:**

   - Open an issue to discuss the new challenge idea.
   - Provide details such as the problem statement and its relevance.

2. **Wait for Approval:**

   - Wait for maintainers or community members to provide feedback.

3. **Create a New Branch:**

   ```bash
   git checkout -b add-challenge-[number]
   ```

4. **Set Up the Challenge Directory:**

   ```
   challenge-[number]/
   ├── README.md
   ├── solution-template.go
   ├── solution-template_test.go
   ├── learning.md
   ├── hints.md
   ├── run_tests.sh
   └── submissions/
   ```

5. **Write the Challenge Description:**

   - Include problem statement, function signature, input/output format, constraints, and sample inputs/outputs in `README.md`.

6. **Create Learning Materials:**

   - In `challenge-[number]/learning.md`, provide:
     - Explanations of relevant Go concepts needed for the challenge
     - Code examples demonstrating these concepts
     - Best practices and efficiency considerations
     - Links to further reading resources

7. **Create the Solution Template:**

   - Provide a skeleton code in `solution-template.go` with appropriate comments.

8. **Write Comprehensive Tests:**

   - Create `solution-template_test.go` with detailed test cases covering various scenarios, including edge cases.

9. **Create Hints:**

   - Provide step-by-step guidance in `hints.md` without giving away the complete solution.

10. **Create Test Script:**

    - Create an executable `run_tests.sh` script for testing submissions.

11. **Update Documentation:**

    - Add the new challenge to the main `README.md`.

#### **Package Challenges (Framework/Library Focused)**

For challenges that focus on specific Go packages/frameworks:

1. **Create a New Issue:**

   - Open an issue to discuss the new package challenge idea.
   - Specify the package/framework (e.g., Gin, Cobra, GORM) and challenge focus.

2. **Wait for Approval:**

   - Wait for maintainers or community members to provide feedback.

3. **Create a New Branch:**

   ```bash
   git checkout -b add-package-[package-name]-challenge-[number]
   ```

4. **Set Up the Package Challenge Directory:**

   ```
   packages/[package-name]/
   ├── package.json                    # Package metadata and learning path
   └── challenge-[number]-[topic]/
       ├── metadata.json               # Challenge-specific metadata
       ├── README.md                   # Challenge description
       ├── solution-template.go        # Template with TODOs
       ├── solution-template_test.go   # Comprehensive tests
       ├── go.mod                      # Module with dependencies
       ├── go.sum                      # Dependency checksums
       ├── learning.md                 # In-depth educational content
       ├── hints.md                    # Step-by-step guidance
       ├── run_tests.sh               # Testing script
       ├── SCOREBOARD.md              # Auto-generated scoreboard
       └── submissions/               # User solutions
           └── [username]/
               └── solution.go        # Complete working solution
   ```

5. **Create Package Metadata (if new package):**

   - Create `packages/[package-name]/package.json` with:
     - Package information (name, description, GitHub repo)
     - Learning path defining challenge progression
     - Categories and difficulty levels

6. **Create Challenge Metadata:**

   - Create `metadata.json` with:
     - Title, description, difficulty, estimated time
     - Learning objectives and prerequisites
     - Requirements and bonus points
     - Tags and real-world connections

7. **Write the Challenge Description:**

   - Include practical problem statement, CLI/API requirements, and testing instructions in `README.md`.

8. **Create Learning Materials:**

   - In `learning.md`, provide comprehensive educational content (400+ lines):
     - Framework fundamentals and core concepts
     - Code examples and patterns
     - Best practices and real-world usage
     - Advanced features and testing strategies

9. **Create the Solution Template:**

   - Provide a structured template in `solution-template.go` with:
     - Proper imports and dependencies
     - Type definitions and structures
     - Function signatures with TODO comments
     - Helper functions and validation logic

10. **Write Comprehensive Tests:**

    - Create `solution-template_test.go` with:
      - Unit tests for all functions and features
      - Integration tests for complete workflows
      - Edge cases and error scenarios
      - Performance and behavior validation

11. **Create Dependencies:**

    - Set up `go.mod` with proper module name and Go version
    - Include all necessary dependencies for the package
    - Run `go mod tidy` to generate `go.sum`

12. **Create Hints:**

    - Provide detailed guidance in `hints.md` with:
      - Step-by-step implementation guidance
      - Code examples and patterns
      - Common pitfalls to avoid
      - Testing and debugging tips

13. **Create Test Script:**

    - Create an executable `run_tests.sh` script that:
      - Tests compilation and functionality
      - Runs unit tests and functional tests
      - Validates flag handling and argument processing
      - Provides detailed feedback and next steps

14. **Create Working Solution:**

    - Implement a complete working solution in `submissions/RezaSi/solution.go`
    - Ensure it passes all tests and demonstrates best practices

15. **Update Documentation:**

    - Update package scoreboard using the package scoreboard scripts
    - Ensure the web UI can discover and display the new challenge

### **General Guidelines for Both Challenge Types**

12. **Commit and Push:**

    ```bash
    # For classic challenges
    git add challenge-[number]/
    git commit -m "Add Challenge [number]: [Challenge Title]"
    
    # For package challenges
    git add packages/[package-name]/
    git commit -m "Add [Package] Challenge [number]: [Challenge Title]"
    
    git push origin [branch-name]
    ```

13. **Create a Pull Request:**

    - Submit the pull request for review.
    - Ensure all tests pass in the CI workflow.
    - Include a detailed description of the challenge and its educational value.

---

## **Style Guidelines**

- **Code Formatting:**

  - Use `gofmt` to format your Go code.
  - Maintain consistent indentation and spacing.

- **Naming Conventions:**

  - Use descriptive variable and function names.
  - Follow Go naming conventions (e.g., `camelCase` for variables and functions).

- **Comments:**

  - Include comments to explain complex logic.
  - Use doc comments (`//`) for exported functions and types.

- **Test Writing:**

  - Write thorough tests covering various input cases.
  - Use subtests (`t.Run()`) to organize test cases.

---

## **Pull Request Process**

1. **Ensure All Tests Pass:**

   - Run tests locally before submitting your pull request.
   - Check that your code does not break existing functionality.

2. **Provide a Clear Description:**

   - Explain what changes you have made and why.
   - Reference any related issues.

3. **One Pull Request per Feature:**

   - Keep your pull requests focused on a single feature or fix.

4. **Wait for Review:**

   - A maintainer will review your pull request.
   - Be responsive to feedback and make necessary changes.

---

## **Reporting Issues**

If you encounter any problems or have suggestions:

- **Open an Issue:**

  - Go to the [Issues](https://github.com/RezaSi/go-interview-practice/issues) tab.
  - Provide a detailed description of the issue or suggestion.
  - Include steps to reproduce the issue if applicable.

---

## **Contact**

For any questions or additional support:

- **Email:** [rezashiri88@gmail.com](mailto:rezashiri88@gmail.com)
- **GitHub:** [RezaSi](https://github.com/RezaSi)

---

Thank you for contributing to the Go Interview Practice repository!

# Go Interview Practice Web UI

This is a web-based user interface for the Go Interview Practice project, providing an interactive environment for solving Go programming challenges.

## Features

- **Challenge Browser**: View all available coding challenges with difficulty indicators.
- **In-browser Code Editor**: Edit and run Go code directly in your browser with syntax highlighting.
- **Test Runner**: Run tests against your solution and see results in real-time.
- **Learning Materials**: Access Go learning materials specific to each challenge to improve your understanding.
- **Scoreboard**: Track your progress and see how you compare to others.
- **Markdown Support**: Challenge descriptions and learning materials rendered with full Markdown support.

## Getting Started

### Prerequisites

- Go 1.16 or later
- Web browser (Chrome, Firefox, Safari, Edge)

### Running the Web UI

1. Navigate to the web-ui directory:
   ```
   cd web-ui
   ```

2. Run the web server:
   ```
   go run main.go
   ```

3. Open your browser and visit:
   ```
   http://localhost:8080
   ```

## Project Structure

```
web-ui/
├── main.go                  # Main server entry point
├── static/                  # Static assets
│   ├── css/                 # CSS stylesheets
│   │   └── style.css        # Custom CSS for the UI
│   └── js/                  # JavaScript files
│       └── main.js          # Common JavaScript utilities
├── templates/               # HTML templates
│   ├── base.html            # Base template with common layout
│   ├── challenge.html       # Challenge page with code editor
│   ├── home.html            # Home page with challenge list
│   └── scoreboard.html      # Scoreboard page for results
└── README.md                # This file
```

## Technical Details

### Templates and HTML Rendering

The web UI uses Go's `html/template` package for server-side rendering, with a base template that defines the common layout and individual content templates for each page type.

### JavaScript Libraries

- **Bootstrap**: For responsive UI components
- **Ace Editor**: For the in-browser code editor
- **Marked**: For Markdown parsing
- **Highlight.js**: For syntax highlighting

### API Endpoints

The web UI exposes the following API endpoints:

- `GET /api/challenges`: Get all challenges
- `GET /api/challenges/{id}`: Get a specific challenge
- `POST /api/run`: Run code for a specific challenge
- `POST /api/submissions`: Submit a solution
- `GET /api/scoreboard/{id}`: Get scoreboard for a challenge

## Development

### Adding New Features

1. If adding new pages, create a new template in the `templates` directory.
2. Add any new API handlers in `main.go`.
3. Add CSS styles to `static/css/style.css`.
4. Add JavaScript utilities to `static/js/main.js`.

### Running in Development Mode

To enable hot-reloading during development, you can use tools like [Air](https://github.com/cosmtrek/air):

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air
```

## Contributing

Contributions to improve the web UI are welcome! Please feel free to submit pull requests or open issues for new features or bug fixes.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 

## Complete Workflow: From Fork to Pull Request

### Prerequisites
1. **Fork the repository** on GitHub (click the "Fork" button)
2. **Clone your fork** locally: `git clone https://github.com/yourusername/go-interview-practice.git`
3. **Start the web UI** as described above

### Solving and Submitting Challenges

1. **Choose a challenge** from the homepage
2. **Write your solution** in the code editor
3. **Test your code** using the "Run Tests" button
4. **Submit your solution** when tests pass

### After Successful Submission

The web UI will guide you through these steps:

1. **Save to Filesystem**: Click the button to save your solution locally
2. **Commit and Push**: 
   ```bash
   git add challenge-X/submissions/yourusername/
   git commit -m "Add solution for Challenge X by yourusername"
   git push origin main
   ```
3. **Create Pull Request**:
   - Go to your fork on GitHub
   - Click "Contribute" → "Open pull request"
   - Add a descriptive title
   - Submit the PR

### What Happens Next

- Your pull request will be reviewed
- Once merged, your solution appears on the public scoreboard
- You get credit for solving the challenge
- Other developers can learn from your approach

This workflow ensures:
- ✅ Your solutions are properly tracked
- ✅ You contribute to the community
- ✅ Your GitHub profile shows your contributions
- ✅ You appear on the public leaderboards

## Submission Process

The web UI provides two ways to submit your solution:

### 1. In-Browser Submission (For Scoreboard Only)

When you click the "Submit Solution" button, your solution will be:
- Tested against the challenge test cases
- Added to the in-memory scoreboard if all tests pass
- Displayed in the challenge scoreboard

This submission is temporary and only exists in the current server session. It doesn't save your solution to the filesystem.

### 2. Filesystem Submission (For Pull Requests)

After successfully submitting a solution that passes all tests, you'll see two options:

#### Option 1: One-Click Filesystem Save

Click the "Save to Filesystem" button to:
- Automatically create a submission directory in your local repository
- Save your solution to `challenge-X/submissions/yourusername/solution-template.go`
- Get a list of Git commands to commit and push your changes

This option creates the actual file structure needed for a GitHub pull request.

#### Option 2: Copy Manual Commands

If you prefer to manage the file creation yourself, you can:
- Click "Copy Commands" to copy shell commands to your clipboard
- Run these commands in your terminal to create the submission files
- Commit and push the changes using the provided Git commands

### Completing Your Submission

After saving your solution to the filesystem (via either method), complete the submission by:
1. Committing your changes
2. Pushing to your fork
3. Creating a pull request to the original repository

This workflow ensures your submission is properly integrated into the project's review system. 
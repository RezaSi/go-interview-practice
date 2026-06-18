# AI Configuration for Real-Time Code Review

## Setup Instructions

### 1. Environment Variables

Create a `.env` file in the project root or set these environment variables:

```bash
# Set your preferred AI provider: gemini, openai, claude, or mock
export AI_PROVIDER=gemini

# API Keys (only set the one you're using)
export GEMINI_API_KEY=your_gemini_api_key_here
export OPENAI_API_KEY=your_openai_api_key_here
export CLAUDE_API_KEY=your_claude_api_key_here

# Optional: Override default models
export AI_MODEL=gemini-pro
```

### 2. Getting API Keys

#### Gemini (Recommended - Free tier available)
1. Go to [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Create a new API key
3. Set `AI_PROVIDER=gemini` and `GEMINI_API_KEY=your_key`

#### OpenAI
1. Go to [OpenAI API Keys](https://platform.openai.com/api-keys)
2. Create a new API key
3. Set `AI_PROVIDER=openai` and `OPENAI_API_KEY=your_key`

#### Claude
1. Go to [Anthropic Console](https://console.anthropic.com/)
2. Create a new API key
3. Set `AI_PROVIDER=claude` and `CLAUDE_API_KEY=your_key`

### 3. Development Mode

For testing without API keys, use mock AI:
```bash
export AI_PROVIDER=mock
```

This provides realistic-looking responses without making external API calls.

### 4. Starting the Server

```bash
cd web-ui
go run main.go
```

The AI features will be available at:
- `POST /api/ai/code-review` - Real-time code analysis
- `POST /api/ai/interviewer-questions` - Generate follow-up questions  
- `POST /api/ai/code-hint` - Context-aware hints

## Features ✅ WORKING

### Real-Time Code Review ✅
- **Overall Score**: 0-100 rating of code quality  
- **Issues Detection**: Bugs, performance, style, logic issues
- **Suggestions**: Optimization and best practice recommendations
- **Complexity Analysis**: Time/space complexity evaluation
- **Interviewer Feedback**: What a real interviewer would say
- **Security**: All content is HTML-escaped for safety

### Dynamic Interview Questions ✅  
- Context-aware questions based on the user's solution
- Progressive difficulty based on user performance
- Go-specific technical probing
- Edge case exploration
- Array of 5 relevant questions per request

### Smart Hints System ✅
- 4 levels of hints (subtle nudge → detailed explanation)
- Context-aware based on current code
- Educational approach that teaches concepts
- Progressive hint buttons (Lv1 → Lv2 → Lv3 → Lv4)

## API Examples

### Code Review
```javascript
POST /api/ai/code-review
{
  "challengeId": 1,
  "code": "func Sum(a, b int) int { return a + b }",
  "context": "Interview started 5 minutes ago"
}
```

### Get Interview Questions
```javascript
POST /api/ai/interviewer-questions
{
  "challengeId": 1, 
  "code": "func Sum(a, b int) int { return a + b }",
  "userProgress": "Completed basic solution"
}
```

### Get Hint
```javascript
POST /api/ai/code-hint
{
  "challengeId": 1,
  "code": "func Sum(a, b int) int { // stuck here }",
  "hintLevel": 2
}
```

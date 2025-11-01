# Challenge 2: Structured Logging & Fields

Take your logging skills to the next level by building a simple **HTTP server with context-aware, structured logging**. This challenge focuses on adding contextual fields to your logs, tracing requests through a system using a correlation ID, and demonstrating custom formatters.

## Challenge Requirements

You will create a web server that uses a **logging middleware**. This middleware will intercept incoming requests, enrich the logging context with request-specific data, and then pass control to the actual request handler.

1.  **Create a Logging Middleware**: This middleware will be the core of the challenge. For every incoming HTTP request, it must:
    * Generate a unique **`request_id`** (using a UUID).
    * Create a logger instance pre-filled with contextual fields: `request_id`, `http_method`, `uri`, and `user_agent`.
    * Log an initial message like "Request received".
    * Pass the enriched logger to the downstream HTTP handler.

2.  **Create an HTTP Handler**: A simple `/hello` handler that:
    * Receives the context-aware logger from the middleware.
    * Uses this logger to log a message (e.g., "Processing hello request"). This message must automatically include the fields added by the middleware.
    * Writes a simple response to the client (e.g., "Hello, world!").

3.  **Main Function**: Sets up the HTTP server, applies the middleware to the handler, and starts listening for requests.

4.  **Custom Formatter Support**:  
    * Configure the logger globally to use a formatter.  
    * Default to `JSONFormatter` for structured logs.  
    * Also demonstrate how to switch to `TextFormatter` (e.g., via a flag or environment variable) to show formatter flexibility.

## Expected Log Output

When you run your server and make a `GET` request to `/hello`, your console output (in **JSON format**) should look similar to this. Notice how the **`request_id` is the same** for both log entries, linking them together.

```json
{"http_method":"GET","level":"info","msg":"Request received","request_id":"a1b2c3d4-e5f6-7890-1234-567890abcdef","time":"...","uri":"/hello","user_agent":"curl/7.79.1"}
{"level":"info","msg":"Processing hello request","request_id":"a1b2c3d4-e5f6-7890-1234-567890abcdef","time":"...","user_id":"user-99"}
```

If you switch to `TextFormatter`, the same logs would be human-readable lines like:

```bash
time="2025-10-02T18:42:00Z" level=info msg="Request received" request_id=a1b2c3d4-e5f6-7890-1234-567890abcdef http_method=GET uri=/hello user_agent="curl/7.79.1"
time="2025-10-02T18:42:00Z" level=info msg="Processing hello request" request_id=a1b2c3d4-e5f6-7890-1234-567890abcdef user_id=user-99
```

> Note: The second log line from the handler includes an extra field, `user_id`, demonstrating how the context can be further enriched

## Implementation Requirements

* **Logger Configuration**:

  * The logger should be configured globally
  * Default to `JSONFormatter`, but allow switching to `TextFormatter` to demonstrate formatter flexibility

* **loggingMiddleware (func)**:

  * Must have the signature `func(http.Handler) http.Handler`
  * Inside, it should create an `http.HandlerFunc`
  * Use a library like `github.com/google/uuid` to generate the `request_id`
  * Create a `logrus.Entry` (a logger with pre-set fields) using `logrus.WithFields()`
  * Use Go's `context` package to pass this `logrus.Entry` to the next handler

* **helloHandler (func)**:

  * Must have the signature `func(http.ResponseWriter, *http.Request)`
  * Retrieve the `logrus.Entry` from the request's context
  * If the logger isn't found in the context, fall back to the global `logrus` logger
  * Add at least one more field to the log (e.g., `user_id`)
  * Write a 200 OK response.

---

### Fields vs Structured Body (Quick comparison)

| Concept | What it is | Example | When to use |
|--------:|:-----------|:--------|:------------|
| Fields | Key-value pairs attached to a log entry (flat). | `request_id=... user_id=user-99 http_method=GET` | Correlating events, simple filtering, fast searches. |
| Structured body | A full structured payload (JSON object) — can include nested objects/arrays. | `{"event":"db.query","sql":"SELECT ...","duration_ms":12}` | Rich context, indexing in log stores, complex queries/visualizations. |
| Combined | Fields + structured body: fields for indexing, body for deep context. | Fields: `request_id=...` Body: `{"user":{"id":"user-99","role":"admin"}}` | Best for production: quick filters + full context when needed. |

---
## Testing Requirements

Your solution must pass tests that verify:

* The middleware correctly adds the required fields (`request_id`, `http_method`, so on...)
* The `request_id` is a valid UUID and is consistent for logs within the same request
* The handler successfully retrieves and uses the context-aware logger
* The final log output from the handler contains both the middleware fields and the handler-specific fields
* The logger supports switching between JSON and Text formatters
---
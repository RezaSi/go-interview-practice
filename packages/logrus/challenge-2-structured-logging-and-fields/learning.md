# Learning: Structured Logging & Fields

This lesson builds on basic logging concepts and shows how to use **structured fields, context-aware logging, custom formatters, and correlation IDs** to trace requests in Go applications.

---

## **Structured Logging with Fields**

Plain logs like:

```
Task failed
```

aren’t enough. Instead, use structured logs with **fields**:

```go
logrus.WithFields(logrus.Fields{
    "task": "ProcessPayment",
    "user_id": "user-456",
    "status": "failed",
}).Error("Task failed")
```

Output:

```json
{"level":"error","msg":"Task failed","task":"ProcessPayment","user_id":"user-456","status":"failed"}
```

Fields make logs queryable and easier to filter by user, component, or operation

---
## Structured logging: fields vs structured body
|                       Concept | Description                                                                                                         | Example (JSON)                                                                                                                  | Pros                                                                              | Cons                                                                                    | Notes / Tooling                                                                                                        |
| ----------------------------: | :------------------------------------------------------------------------------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------ | :-------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------- |
|       Fields (flat key/value) | Flat key-value fields attached to each log entry (log metadata). Easy to index & filter.                            | `{"level":"info","msg":"Request received","request_id":"...","http_method":"GET"}`                                              | Fast to search and aggregate; good for correlation IDs, status codes, user IDs.   | Can be limiting for deeply nested context.                                              | Use for primary searchable dimensions (request ids, user ids, status, latency). `logrus.WithFields()` is a common API. |
| Structured body (nested JSON) | A richer JSON payload inside the log entry—can include nested objects, arrays, and full event contexts.             | `{"event":"db.query","sql":"SELECT ...","params":{"user_id":"user-99"},"duration_ms":12}`                                       | Stores deep context for debugging and replay; ideal for log analytics and traces. | Bigger logs, potentially more storage and noise; needs schemas for consistent querying. | Use for complex events (DB queries, error stacks, HTTP payloads). Index selectively to avoid cost.                     |
|             Combined approach | Use fields for the frequently-filtered keys and a structured body for full context.                                 | `{"level":"info","msg":"Processed order","request_id":"...","user_id":"user-99","body":{"order":{"id":"o-123","items":[...]}}}` | Best of both worlds: quick filters + full context on demand.                      | More work to design consistent schemas.                                                 | Recommended pattern in production; allows fast dashboards and deep investigations.                                     |
|               Correlation IDs | A dedicated field (eg `request_id`) propagated across services and logs so all related entries can be linked.       | `{"request_id":"...","msg":"..."} `                                                                                             | Enables tracing across services, useful for distributed debugging.                | Needs generation & propagation discipline (middleware, headers).                        | Use UUIDs, include in headers (e.g., `X-Request-ID`) and log fields.                                                   |
|                 Practical tip | Keep a small set of indexed fields (IDs, status, latency buckets) and push large payloads to body only when needed. | —                                                                                                                               | Reduces storage/cost while keeping searchability.                                 | —                                                                                       | Consider log retention/rotation and PII policies.                                                                      |


---
## **Context-Aware Logging**

Logs are more powerful when tied to a specific request or operation
> Note: in the challenge and tests we use the field name http_method (not method) and prefer r.URL.Path or r.RequestURI for the uri field - keep names consistent between README, tests, and learning content

### Middleware Setup

```go
// typed key prevents collisions
type loggerKeyType string

const loggerKey loggerKeyType = "logger"

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        logger := logrus.WithFields(logrus.Fields{
            "request_id":  requestID,
            "http_method": r.Method,
            "uri":         r.URL.Path,
            "user_agent":  r.UserAgent(),
        })

        // put the *logrus.Entry into the context
        ctx := context.WithValue(r.Context(), loggerKey, logger)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Using in Handlers

```go
// safe retrieval: fall back to global logger if missing
logger, ok := r.Context().Value(loggerKey).(*logrus.Entry)
if !ok || logger == nil {
    logger = logrus.NewEntry(logrus.StandardLogger())
}

logger = logger.WithField("user_id", "user-99")
logger.Info("Processing request")
```

Every handler log now carries the same `request_id`, `method`, and `uri` fields

---

## **3. Custom Formatters**

Logrus supports multiple output formats.

### JSON Formatter (default for structured logs)

```go
logrus.SetFormatter(&logrus.JSONFormatter{})
```

### Text Formatter (human-friendly)

```go
logrus.SetFormatter(&logrus.TextFormatter{
    FullTimestamp: true,
})
```

You can swap formatters depending on your environment (JSON for production, text for local debugging)

---

## **4. Request Tracing & Correlation IDs**

A **Correlation ID** ties logs across different components into one trace

```go
import "github.com/google/uuid"

requestID := uuid.New().String()
logger := logrus.WithField("request_id", requestID)
```

When passed via headers (e.g., `X-Request-ID`), correlation IDs can trace:

* API calls across microservices
* Database queries linked to a request
* Error paths back to a single user action

---

## **Best Practices**

1. **Always log with fields** — structured > plain text
2. **Add request_id** early in the request lifecycle
3. **Pick the right formatter** (JSON for machines, text for humans)
4. **Propagate correlation IDs** across service boundaries
5. Keep messages clear, avoid secrets

---

## **Resources**

* [Logrus – Fields & Formatters](https://github.com/sirupsen/logrus)
* [Go Docs – context](https://pkg.go.dev/context)
* [Correlation IDs in Go](https://blog.golang.org/context)

---
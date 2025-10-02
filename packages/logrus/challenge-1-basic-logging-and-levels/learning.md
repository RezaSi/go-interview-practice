# Learning: Basic Logging in Go with Logrus

## **What is Logging?**

Logging is how your application records what’s happening while it runs. It’s your “black box.”
When something fails in production, logs are usually the only way to see what went wrong.

**Why log?**

* **Debugging** – track down issues without a debugger
* **Monitoring** – check health/performance
* **Auditing** – record important events
* **Clarity** – know what your code was doing at a given time

---

## **Go’s Built-in Logging**

Go ships with the `log` package:

```go
import "log"

func main() {
    log.Println("Hello from standard log")
}
```

Output looks like:

```
2025/10/02 18:00:00 Hello from standard log
```

It’s nice for basics, but:

* No levels (everything is just a line)
* No JSON/structured output
* Hard to configure

---

## **Logrus: A Better Logger**

[Logrus](https://github.com/sirupsen/logrus) is the most common logging library for Go.
It’s compatible with `log` but adds:

* **Levels**: Debug, Info, Warn, Error, Fatal, Panic
* **Formatters**: Text or JSON
* **Configurable Output**: Console, file, or anything that implements `io.Writer`

---

## **Core Concepts**

### **1. Setup Logger**

```go
import (
    "os"
    "github.com/sirupsen/logrus"
)

func setupLogger() {
    logrus.SetOutput(os.Stdout)                 // send logs to console
    logrus.SetFormatter(&logrus.JSONFormatter{}) // use JSON format
    logrus.SetLevel(logrus.InfoLevel)           // default: Info and above
}
```

---

### **2. Log Levels**

Levels control which logs are shown.

From lowest → highest severity:

* `Debug` → development details
* `Info` → normal operations
* `Warn` → something’s off but still running
* `Error` → operation failed
* `Fatal` → critical, app exits
* `Panic` → logs + panics

Example:

```go
logrus.Debug("Checking system status...")
logrus.Info("Logbook app starting")
logrus.Warn("Disk space is low")
logrus.Error("Failed to connect to backup service")
logrus.Fatal("Critical config file missing") // exits
logrus.Panic("Database connection issue")    // panics
```

---

### **3. Choosing Levels at Runtime**

You can set the level dynamically from input:

```go
func setupLogger(level string) {
    logrus.SetOutput(os.Stdout)
    logrus.SetFormatter(&logrus.JSONFormatter{})

    lvl, err := logrus.ParseLevel(level)
    if err != nil {
        lvl = logrus.InfoLevel
    }
    logrus.SetLevel(lvl)
}
```

---

## **Building a Simple Logbook App**

```go
func runLogbookOperations() {
    logrus.Debug("Checking system status...")
    logrus.Info("Logbook application starting up.")
    logrus.Warn("Disk space is running low.")
    logrus.Error("Failed to connect to remote backup service.")
    logrus.Fatal("Critical configuration file 'config.yml' not found.")
    logrus.Panic("Unhandled database connection issue.")
}

func main() {
    logLevel := "info"
    if len(os.Args) > 1 {
        logLevel = os.Args[1]
    }

    setupLogger(logLevel)
    logrus.Infof("Log level set to '%s'", logrus.GetLevel().String())
    runLogbookOperations()
}
```

---

## **Best Practices**

1. Pick the right level: `Debug` for dev, `Info` for normal ops, `Error` when something breaks.
2. Default to JSON formatter — easy to parse later.
3. Don’t log secrets (passwords, keys).
4. Keep messages clear: “Failed to connect to DB” > “error 17”.

---

## 🔗 **Resources**

* [Logrus GitHub](https://github.com/sirupsen/logrus)
* [Go by Example – Logging](https://gobyexample.com/logging)

---
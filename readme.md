# Go Task Worker API

A simple Golang HTTP service that accepts task requests, and processes them asynchronously using worker goroutines.

This project demonstrates key Go concepts:
- Concurrency with goroutines and channels
- HTTP server & JSON handling
- Environment-based configuration (`GOMAXPROCS`)
- Graceful task queue handling

---

# Features

- REST API for submitting background tasks
- Worker goroutine that processes queued tasks
- Configurable CPU concurrency via `GOMAXPROCS`
- Basic health check endpoint
- Channel-based in-memory queue with rate limiting (returns **429 Too Many Requests** when full)

---




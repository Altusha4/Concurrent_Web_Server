# Concurrent Web Server — Go (net/http)

A concurrent RESTful web server built using Go’s net/http package.
The application demonstrates safe concurrent request handling, shared state protection with mutexes, background workers, graceful shutdown, and a web dashboard interface

---

## Technologies Used

- ![](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
- ![](https://img.shields.io/badge/net/http-000000?style=for-the-badge&logo=go&logoColor=white)
- ![](https://img.shields.io/badge/Goroutines-00ADD8?style=for-the-badge&logo=go&logoColor=white)
- ![](https://img.shields.io/badge/Mutex-4B5563?style=for-the-badge)
- ![](https://img.shields.io/badge/Context-2563EB?style=for-the-badge)
- ![](https://img.shields.io/badge/JSON-F59E0B?style=for-the-badge)


---

## Features

- RESTful API built with Go `net/http`
- Concurrent request handling using goroutines
- Thread-safe storage with sync.RWMutex protection
- Background worker with periodic status logging
- Graceful shutdown using context and signals
- Web dashboard interface with real-time statistics
- JSON-based request/response format
- Clean architecture with separation of concerns

## API Endpoints

| Method | Endpoint | Description |
|------|--------|------------|
| POST | `/data` | Store key-value pair	{"key": "string", "value": "string"} |
| GET | `/data` | Delete data by key |
| DELETE | `/data/{key}` | Delete entry by key (ID) |
| GET | `/stats` | Get server statistics |

---

## Project Structure

```
assignment2/
├── cmd/
│   └── webserver/
│       └── main.go              # Server entry point & graceful shutdown
├── internal/
│   ├── handler/
│   │   └── handler.go           # HTTP request handlers
│   ├── models/
│   │   └── models.go            # Data models (KeyValue)
│   ├── service/
│   │   └── service.go           # Business logic layer
│   ├── storage/
│   │   └── storage.go           # Thread-safe in-memory storage
│   └── worker/
│       └── worker.go            # Background worker
├── pkg/
│   └── frontend/
│       └── frontend.go          # Web dashboard serving
│       └── index.html 
│       └── style.css  
│       └── app.js 
│
├── go.mod                       # Go module definition
└── README.md                    # This file
```
---

## Concurrency & Thread Safety

- net/http automatically handles each request in a separate goroutine
- Shared resources are protected using:
    - sync.RWMutex for the in-memory map
    - sync.Mutex for request counters
- Prevents race conditions during concurrent access

---

## Background Worker

- Runs in a separate goroutine
- Logs server status every 5 seconds
- Uses time.Ticker and select statement
- Stops cleanly when shutdown signal is received


---

## Graceful Shutdown

- Implemented using context.Context
- OS signals (Ctrl+C, SIGTERM) are captured
- Server shuts down without interrupting active requests
- Background worker stops gracefully

This follows industry-standard practices.

---

# How to Run the Project

## Install Go
Make sure Go is installed:
```
go version
```
Run the server
```
go run .
```
Server address
```
http://localhost:8080
```
Example Usage (curl)
Add data
```
curl -X POST http://localhost:8080/data \
  -H "Content-Type: application/json" \
  -d '{"id":"SE-2416","subject":"ADP1","day":"Wednesday","time":"14:00","room":"C1.1.239","teacher":"Nurlybek"}'
```
Get all data
```
curl http://localhost:8080/data
```
Delete data
```
curl -X DELETE http://localhost:8080/data/CS101
```
Get statistics
```
curl http://localhost:8080/stats
```

# **How It Works**

## **Request Processing Flow**

### **Complete Request Cycle:**
```
Client → HTTP Request → net/http → Goroutine → Handler → Service → Storage → JSON Response
                             │                                           │
                        Auto-spawned                               Mutex-protected
                        per request                               thread-safe access
```

---

## **Key Components**

### **1. Concurrency (Goroutines)**
- Each HTTP request runs in separate goroutine
- Automatically managed by `net/http`
- Supports thousands of concurrent connections

```go
// net/http handles goroutines automatically
http.HandleFunc("/api/data", handler) // Each call = new goroutine
```

### **2. Thread Safety (Mutex)**
- **Multiple readers** allowed simultaneously (RWMutex)
- **Single writer** with exclusive access
- Prevents race conditions in shared memory

```go
// Storage layer protection
func (s *Storage) Get(key string) {
    s.mu.RLock()    //  Multiple readers OK
    defer s.mu.RUnlock()
    // Read data
}

func (s *Storage) Set(key, value string) {
    s.mu.Lock()     //  One writer at a time
    defer s.mu.Unlock()
    // Write data
}
```

### **3. Background Worker**
- Runs independently every 5 seconds
- Monitors: Request count & Database size
- Clean shutdown via context cancellation

```go
// Worker with context control
for {
    select {
    case <-ticker.C:        // Every 5 seconds
        logStatus()
    case <-ctx.Done():      // Shutdown signal
        return              // Clean exit
    }
}
```

### **4. Graceful Shutdown**
1. `Ctrl+C` signal received
2. Context cancelled
3. Worker stops
4. Server stops accepting new requests
5. Active requests complete
6. Clean exit

```go
signal.Notify(ch, os.Interrupt)  // Listen for Ctrl+C
<-ch                              // Wait for signal
server.Shutdown(ctx)              // Graceful stop
```

---

## **Performance Features**

- **Auto-scaling**: Each request = new goroutine
- **Thread-safe**: Mutex-protected data access
- **Non-blocking**: Background worker doesn't block requests
- **Resource-efficient**: Minimal memory footprint

---

## **Real-time Monitoring**

**Terminal Output:**
```
[Worker] Status - Requests: 15, Database size: 3
[Worker] Status - Requests: 17, Database size: 4
```

**Dashboard Updates:**
- Live request counter
- Real-time database size
- Interactive API testing
- Visual feedback

---

## **Why It Works**

| Feature | Benefit |
|---------|---------|
| **Goroutines** | High concurrency, low overhead |
| **Mutex** | Thread-safe without data races |
| **Context** | Clean shutdown and cancellation |
| **Atomic** | Fast, lock-free counters |
| **Embed** | Single binary deployment |

## Author

Altynay Yertay  
Software Engineering Student  
Astana IT University


package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Request struct {
	Status   string `json:"status"`
	Duration int    `json:"duration"` // required by our logic, but stdlib won't enforce it
}

var taskChan = make(chan Request, 10)

func main() {
	v := os.Getenv("GOMAXPROCS")
	if v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			runtime.GOMAXPROCS(n)
			fmt.Printf("set GOMAXPROCS to %d\n", n)
		} else {
			fmt.Printf("invalid GOMAXPROCS=%q: %v\n", v, err)
		}
	} else {
		fmt.Printf("GOMAXPROCS is %d\n", runtime.GOMAXPROCS(0))
	}

	go worker()

	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/healthCheck", healthHandler)

	fmt.Println("listening on port 8000")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func worker() {
	for t := range taskChan {
		fmt.Printf("request status %q and duration %d\n", t.Status, t.Duration)
		time.Sleep(time.Duration(t.Duration) * time.Second)
		fmt.Printf("Task completed: status=%s, duration=%d\n", t.Status, t.Duration)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Duration <= 0 {
		http.Error(w, "duration must be > 0", http.StatusBadRequest)
		return
	}

	// Normalize status – we could also keep req.Status if you prefer
	resp := Request{
		Status:   "Completed",
		Duration: req.Duration,
	}

	// Try to enqueue first, *then* write headers/body
	select {
	case taskChan <- resp:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			// at this point headers are sent; just log
			log.Printf("error writing response: %v", err)
		}
	default:
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

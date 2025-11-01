package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Status   string `json:"status"`
	Duration int    `json:"duration" binding:"required"`
}

var taskChan = make(chan Request, 10)

func main() {
	go worker()
	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/healthCheck", healthHandler)

	fmt.Println("listening on port 8000")

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func worker() {
	for t := range taskChan {
		fmt.Printf("request status %v and duration %d\n", t.Status, t.Duration)
		time.Sleep(time.Duration(t.Duration) * time.Second)
		fmt.Printf("Task completed: status=%s, duration=%d\n", t.Status, t.Duration)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var req Request

	//Decode the input into req

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Input JSON ", http.StatusBadRequest)
		return
	}

	if req.Status == "" || req.Duration == 0 {
		http.Error(w, "Missing Input values ", http.StatusBadRequest)
		return

	}

	//modify the status

	resp := Request{
		Status:   "Completed",
		Duration: req.Duration,
	}

	// set the headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	select {
	case taskChan <- resp:
		{
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				http.Error(w, "Invalid Input JSON ", http.StatusBadRequest)
				return
			}
		}
	default:

		http.Error(w, "Too Many requests", http.StatusTooManyRequests)
	}

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

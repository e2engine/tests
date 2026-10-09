package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	e2enginehttp "github.com/e2engine/instrumentation-go/http"
)

func main() {
	var wait int

	if len(os.Args) > 1 {
		var err error
		wait, err = strconv.Atoi(os.Args[1]) // Just to simulate some processing based on input
		if err != nil {
			log.Fatalf("Invalid gateway http path argument: %v", err)
		}
	}

	client := &http.Client{
		Transport: e2enginehttp.Transport(),
	}

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodGet,
			"http://127.0.0.1:8083/healthz",
			http.NoBody,
		)
		if err != nil {
			http.Error(w, "failed to create dependency request", http.StatusInternalServerError)
			return
		}

		const testExecutionIDHeader = "E2Engine-Test-Execution-ID"
		req.Header.Set(
			testExecutionIDHeader,
			r.Header.Get(testExecutionIDHeader),
		)

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "dependency request failed", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "dependency returned non-200 status", http.StatusBadGateway)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "failed to read dependency response", http.StatusBadGateway)
			return
		}

		time.Sleep(time.Duration(wait) * time.Second) // Simulate a long processing time

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(body)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(time.Duration(wait) * time.Second) // Simulate a long processing time

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_, _ = w.Write([]byte(`{"error": "simulated error response"}`))
	})

	if err := http.ListenAndServe("127.0.0.1:9000", nil); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	e2enginehttp "github.com/e2engine/instrumentation-go/http"
)

func main() {
	client := &http.Client{
		Transport: e2enginehttp.Transport(),
	}

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodPost,
			"http://127.0.0.1:8083/healthz",
			bytes.NewBufferString(`{"customer_id":"customer-1"}`),
		)
		if err != nil {
			http.Error(w, "failed to create dependency request", http.StatusBadGateway)
			return
		}

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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(body)
	})

	if err := http.ListenAndServe(
		"127.0.0.1:9000",
		e2enginehttp.Handler(),
	); err != nil {
		log.Fatal(err)
	}
}

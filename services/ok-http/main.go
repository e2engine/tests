package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ok", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}

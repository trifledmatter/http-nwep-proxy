package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/usenwep/nwfetch-go"
)

func main() {
	nwfetch.Init()
	defer nwfetch.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("url")
		if target == "" {
			http.Error(w, "missing ?url= parameter", http.StatusBadRequest)
			return
		}

		resp, err := nwfetch.Get(target)
		if err != nil {
			http.Error(w, fmt.Sprintf("nwfetch error: %v", err), http.StatusBadGateway)
			return
		}

		if err := resp.StatusError(); err != nil {
			http.Error(w, fmt.Sprintf("upstream error: %s — %s", resp.Status, resp.StatusDetails), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(resp.Body)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

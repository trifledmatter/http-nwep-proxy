package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/usenwep/nwfetch-go"
)

var client *nwfetch.Client

func main() {
	var err error
	client, err = nwfetch.NewClient(nwfetch.WithTimeout(3 * time.Second))
	if err != nil {
		log.Fatalf("failed to create nwfetch client: %v", err)
	}
	defer client.Close()

	http.HandleFunc("/raw", handleRaw)
	http.HandleFunc("/", handleIframe)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleRaw(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("addr")
	if target == "" {
		http.Error(w, "missing ?addr= parameter", http.StatusBadRequest)
		return
	}

	resp, err := client.Get(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to reach %s", target), http.StatusBadGateway)
		return
	}

	if err := resp.StatusError(); err != nil {
		http.Error(w, fmt.Sprintf("upstream error: %s — %s", resp.Status, resp.StatusDetails), http.StatusBadGateway)
		return
	}

	contentType := "text/plain; charset=utf-8"
	if ct, ok := resp.Header("content-type"); ok && ct != "" {
		contentType = ct
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(resp.Body)
}

func handleIframe(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("addr")
	if target == "" {
		http.Error(w, "missing ?addr= parameter", http.StatusBadRequest)
		return
	}

	escaped := html.EscapeString("/raw?addr=" + target)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>nwep proxy</title>
<style>*{margin:0;padding:0}iframe{width:100%%;height:100vh;border:none}</style>
</head>
<body><iframe src="%s"></iframe></body>
</html>`, escaped)
}

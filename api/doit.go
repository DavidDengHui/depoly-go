package handler

import (
	"fmt"
	"net/http"
)
 
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
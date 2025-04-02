package handler

import (
	"fmt"
	"net/http"
)

// 関数名は "Handler" 固定
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Go + Vercel Serverless!")
}

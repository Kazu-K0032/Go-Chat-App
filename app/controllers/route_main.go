package controllers

import (
	"net/http"
)

// '/'へのアクセス
func top(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "layout", "navbar", "top")
}

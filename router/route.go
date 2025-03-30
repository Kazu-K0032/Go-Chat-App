package router

import (
	"net/http"
)

// '/'へのアクセス
func top(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "layout", "header", "top", "footer") 
}

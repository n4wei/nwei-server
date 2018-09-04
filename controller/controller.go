package controller

import (
	"io"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "success\n")
}

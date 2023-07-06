package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// Greet godoc
// @Summary      Greeter
// @Description  Show "Hello, World" and current time
// @Success      200
// @Router       / [get]
func Greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

package main

import (
	"net/http"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message + ". I am Latest\n"
	w.Write([]byte(message))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Just fine"))
}

func main() {
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":7080", nil); err != nil {
		panic(err)
	}
}
